package repositories

import (
	"database/sql"
	"promail/models"
)

type AppConfigRepository struct {
	DB *sql.DB
}

func (r *AppConfigRepository) AppConfigExistsByAppID(appID int64, userID int64) (bool, error) {

	var exists bool

	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM app_configs c
			JOIN apps a ON a.id = c.app_id
			WHERE c.app_id = $1
			  AND a.user_id = $2
		)
	`, appID, userID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AppConfigRepository) AppConfigExistsByID(configID int64, userID int64) (bool, error) {

	var exists bool

	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM app_configs c
			JOIN apps a ON a.id = c.app_id
			WHERE c.id = $1
			  AND a.user_id = $2
		)
	`, configID, userID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AppConfigRepository) GetAppConfigs(appID int64, userID int64) (*models.AppConfigData, error) {

	var appConfig models.AppConfigData

	err := r.DB.QueryRow(
		`SELECT
			c.id,
			c.app_id,
			c.host,
			c.port,
			c.name,
			c.username,
			c.password,
			c.open_track,
			c.click_track,
			c.auto_retry,
			c.retry_max_count,
			c.created_at,
			c.updated_at
		FROM app_configs c
		JOIN apps a ON a.id = c.app_id
		WHERE a.id = $1
		  AND a.user_id = $2`,
		appID, userID).Scan(&appConfig.ID, &appConfig.AppID, &appConfig.SMTPHost, &appConfig.SMTPPort, &appConfig.SMTPName, &appConfig.SMTPUsername, &appConfig.SMTPPassword, &appConfig.OpenTrack, &appConfig.ClickTrack, &appConfig.AutoRetry, &appConfig.RetryMaxCount, &appConfig.CreatedAt, &appConfig.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}

func (r *AppConfigRepository) CreateAppConfig(config models.AppConfigCreate) error {

	_, err := r.DB.Exec(`INSERT INTO app_configs(app_id, host, port, name, username, password, open_track, click_track, auto_retry, retry_max_count) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, config.AppID, config.SMTPHost, config.SMTPPort, config.SMTPName, config.SMTPUsername, config.SMTPPassword, config.OpenTrack, config.ClickTrack, config.AutoRetry, config.RetryMaxCount)

	return err
}

func (r *AppConfigRepository) UpdateAppConfig(configID int64, config models.AppConfigUpdate) error {

	_, err := r.DB.Exec(`UPDATE app_configs SET host=$1, port=$2, name=$3, username=$4, password=$5, open_track=$6, click_track=$7, auto_retry=$8, retry_max_count=$9 WHERE id=$10`, config.SMTPHost, config.SMTPPort, config.SMTPName, config.SMTPUsername, config.SMTPPassword, config.OpenTrack, config.ClickTrack, config.AutoRetry, config.RetryMaxCount, configID)

	return err
}
