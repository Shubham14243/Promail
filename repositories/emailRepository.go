package repositories

import (
	"database/sql"
	"promail/models"
)

type EmailRepository struct {
	DB *sql.DB
}

func (r *EmailRepository) AddEmailLogAndQueue(emailLog models.EmailLogDataCreate) (models.LogResponse, error) {

	var logRes models.LogResponse

	tx, err := r.DB.Begin()
	if err != nil {
		return logRes, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
		INSERT INTO email_logs (
			uuid, user_id, app_id, template_id, to_email, subject, variable_data, rendered_body
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING uuid, id
	`,
		emailLog.UUID,
		emailLog.UserID,
		emailLog.AppID,
		emailLog.TemplateID,
		emailLog.ToEmail,
		emailLog.Subject,
		emailLog.VariableData,
		emailLog.Body,
	).Scan(&logRes.AckID, &logRes.LogID)
	if err != nil {
		return logRes, err
	}

	_, err = tx.Exec(`
		INSERT INTO email_queue (email_log_id)
		VALUES ($1)
	`, logRes.LogID)
	if err != nil {
		return logRes, err
	}

	if err = tx.Commit(); err != nil {
		return logRes, err
	}

	return logRes, nil
}

func (r *EmailRepository) AddOpenTracking(logData models.LogResponse) error {

	_, err := r.DB.Exec(`INSERT INTO email_analytics(email_log_id, type, tracking_token) VALUES($1, $2, $3)`, logData.LogID, "open", logData.AckID)

	return err
}

func (r *EmailRepository) UpdateOpenTracking(emailLogID int64) error {
	_, err := r.DB.Exec(`UPDATE email_analytics SET opened_at = NOW() WHERE email_log_id = $1 AND type = 'open'`, emailLogID)

	return err
}

func (r *EmailRepository) GetOpenWithUUID(logUUID string) (*models.TrackingData, error) {

	var trackData models.TrackingData

	err := r.DB.QueryRow(
		`SELECT
			id,
			email_log_id,
			type,
			original_url,
			tracking_token,
			opened_at,
			clicked_at
			FROM email_analytics
			WHERE tracking_token = $1
			`,
		logUUID).Scan(&trackData.ID, &trackData.EmailLogID, &trackData.Type, &trackData.Url, &trackData.TrackingToken, &trackData.OpenedAt, &trackData.ClickedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &trackData, nil
}

func (r *EmailRepository) GetLogWithUUID(logUUID string) (*models.EmailLogData, error) {

	var logData models.EmailLogData

	err := r.DB.QueryRow(
		`SELECT
			id,
			uuid,
			user_id,
			app_id,
			template_id,
			to_email,
			subject,
			variable_data,
			rendered_body,
			status,
			error_message,
			sent_at,
			created_at,
			updated_at
			FROM email_logs
			WHERE uuid = $1
			`,
		logUUID).Scan(&logData.ID, &logData.LogUUID, &logData.UserID, &logData.AppID, &logData.TemplateID, &logData.ToEmail, &logData.Subject, &logData.Variables, &logData.Body, &logData.Status, &logData.ErrorMessage, &logData.SentAt, &logData.CreatedAt, &logData.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &logData, nil
}
