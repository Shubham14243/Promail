package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	"promail/models" // wherever AppConfigData / SendEmail / PrepareEmailBody live
)

const (
	pollInterval    = 5 * time.Second
	batchSize       = 20
	sendConcurrency = 5 // max simultaneous SMTP sends per tick
)

type EmailJob struct {
	QueueID    int64
	EmailLogID int64
	Attempts   int

	ToEmail string
	Subject string
	Body    string
	Type    string // "html" | "text"

	AppID         int64
	AutoRetry     string
	RetryMaxCount int

	Conf models.AppConfigData
}

type Worker struct {
	DB *sql.DB
}

func NewWorker(db *sql.DB) *Worker {
	return &Worker{DB: db}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.tick(ctx); err != nil {
				log.Printf("worker tick error: %v", err)
			}
		}
	}
}

func (w *Worker) tick(ctx context.Context) error {
	jobs, err := w.claimBatch(ctx)
	if err != nil {
		return err
	}
	if len(jobs) == 0 {
		return nil
	}

	sem := make(chan struct{}, sendConcurrency)
	done := make(chan struct{})
	for _, job := range jobs {
		job := job
		sem <- struct{}{}
		go func() {
			defer func() { <-sem; done <- struct{}{} }()
			w.processJob(ctx, job)
		}()
	}
	for range jobs {
		<-done
	}
	return nil
}

// claimBatch locks a batch of due jobs, flips email_logs to 'processing' so
// no other worker instance grabs them, and returns everything needed to send.
func (w *Worker) claimBatch(ctx context.Context) ([]EmailJob, error) {
	tx, err := w.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
		SELECT
			eq.id, eq.email_log_id, eq.attempts,
			el.to_email, el.subject, el.rendered_body,
			t.type,
			ac.host, ac.port, ac.username, ac.password, ac.name,
			ac.auto_retry, ac.retry_max_count,
			a.id
		FROM email_queue eq
		JOIN email_logs el   ON el.id = eq.email_log_id
		JOIN apps a          ON a.id = el.app_id
		JOIN app_configs ac  ON ac.app_id = a.id
		LEFT JOIN templates t ON t.id = el.template_id
		WHERE a.status = 'active'
		  AND el.status IN ('queued')
		  AND (
		    eq.last_attempted_at IS NULL
		    OR eq.last_attempted_at < NOW() - (INTERVAL '20 seconds' * POWER(2, eq.attempts))
		  )
		ORDER BY eq.id
		LIMIT $1
		FOR UPDATE OF eq SKIP LOCKED
	`, batchSize)
	if err != nil {
		return nil, err
	}

	var jobs []EmailJob
	var claimedLogIDs []int64

	for rows.Next() {
		var j EmailJob
		var tType sql.NullString
		var confName sql.NullString

		if err := rows.Scan(
			&j.QueueID, &j.EmailLogID, &j.Attempts,
			&j.ToEmail, &j.Subject, &j.Body,
			&tType,
			&j.Conf.SMTPHost, &j.Conf.SMTPPort, &j.Conf.SMTPUsername, &j.Conf.SMTPPassword, &confName,
			&j.AutoRetry, &j.RetryMaxCount,
			&j.AppID,
		); err != nil {
			rows.Close()
			return nil, err
		}

		if tType.Valid {
			j.Type = tType.String
		} else {
			j.Type = "html" // ad-hoc emails: default, or store type on email_logs too
		}
		if confName.Valid {
			j.Conf.SMTPName = confName.String
		}

		j.Conf.SMTPPassword, err = Decrypt(j.Conf.SMTPPassword)
		if err != nil {
			log.Printf("failed to decrypt SMTP password. \n%v", err)
		}

		jobs = append(jobs, j)
		claimedLogIDs = append(claimedLogIDs, j.EmailLogID)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		return nil, tx.Commit()
	}

	// Mark claimed logs as processing so they don't get re-selected while in flight.
	if _, err := tx.ExecContext(ctx, `
		UPDATE email_logs SET status = 'processing', updated_at = NOW()
		WHERE id = ANY($1)
	`, pqArray(claimedLogIDs)); err != nil {
		return nil, err
	}
	// Bump attempts + last_attempted_at now, before sending — guarantees
	// a crash mid-send doesn't leave the job retry-able forever without cooldown.
	for _, j := range jobs {
		if _, err := tx.ExecContext(ctx, `
			UPDATE email_queue SET attempts = attempts + 1, last_attempted_at = NOW()
			WHERE id = $1
		`, j.QueueID); err != nil {
			return nil, err
		}
	}

	return jobs, tx.Commit()
}

func (w *Worker) processJob(ctx context.Context, job EmailJob) {
	err := SendEmail(&job.Conf, job.ToEmail, job.Subject, job.Body, job.Type)
	if err != nil {
		w.handleFailure(ctx, job, err)
		return
	}
	w.handleSuccess(ctx, job)
}

func (w *Worker) handleSuccess(ctx context.Context, job EmailJob) {
	_, err := w.DB.ExecContext(ctx, `
		UPDATE email_logs SET status = 'sent', sent_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, job.EmailLogID)
	if err != nil {
		log.Printf("failed to mark email_log %d sent: %v", job.EmailLogID, err)
		return
	}
	if _, err := w.DB.ExecContext(ctx, `DELETE FROM email_queue WHERE id = $1`, job.QueueID); err != nil {
		log.Printf("failed to delete queue row %d: %v", job.QueueID, err)
	}
}

func (w *Worker) handleFailure(ctx context.Context, job EmailJob, sendErr error) {
	// attempts was already incremented in claimBatch (this send attempt counts).
	currentAttempts := job.Attempts + 1

	maxAttempts := job.RetryMaxCount
	if maxAttempts <= 0 {
		maxAttempts = 1
	}

	giveUp := job.AutoRetry != "active" || currentAttempts >= maxAttempts

	if giveUp {
		_, err := w.DB.ExecContext(ctx, `
			UPDATE email_logs
			SET status = 'failed', error_message = $2, updated_at = NOW()
			WHERE id = $1
		`, job.EmailLogID, sendErr.Error())
		if err != nil {
			log.Printf("failed to mark email_log %d failed: %v", job.EmailLogID, err)
			return
		}
		if _, err := w.DB.ExecContext(ctx, `DELETE FROM email_queue WHERE id = $1`, job.QueueID); err != nil {
			log.Printf("failed to delete queue row %d: %v", job.QueueID, err)
		}
		return
	}

	// Leave it in the queue for retry; revert status to 'queued' so claimBatch
	// picks it up again once the backoff window passes.
	_, err := w.DB.ExecContext(ctx, `
		UPDATE email_logs
		SET status = 'queued', error_message = $2, updated_at = NOW()
		WHERE id = $1
	`, job.EmailLogID, sendErr.Error())
	if err != nil {
		log.Printf("failed to requeue email_log %d: %v", job.EmailLogID, err)
	}
}

// pqArray is a placeholder — use github.com/lib/pq.Array(ids) in real code
// if you're on lib/pq, or pgx's equivalent if on pgx.
func pqArray(ids []int64) interface{} {
	return ids
}
