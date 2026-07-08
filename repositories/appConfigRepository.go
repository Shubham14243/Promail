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
			c.username,
			c.password,
			c.open_track,
			c.click_track,
			c.created_at,
			c.updated_at
		FROM app_configs c
		JOIN apps a ON a.id = c.app_id
		WHERE a.id = $1
		  AND a.user_id = $2`,
		appID, userID).Scan(&appConfig.ID, &appConfig.AppID, &appConfig.SMTPHost, &appConfig.SMTPPort, &appConfig.SMTPUsername, &appConfig.SMTPPassword, &appConfig.OpenTrack, &appConfig.ClickTrack, &appConfig.CreatedAt, &appConfig.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}

func (r *AppConfigRepository) CreateAppConfig(config models.AppConfigCreate) error {

	_, err := r.DB.Exec(`INSERT INTO app_configs(app_id, host, port, username, password, open_track, click_track) VALUES($1, $2, $3, $4, $5, $6, $7)`, config.AppID, config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword, config.OpenTrack, config.ClickTrack)

	return err
}

func (r *AppConfigRepository) UpdateAppConfig(configID int64, config models.AppConfigUpdate) error {

	_, err := r.DB.Exec(`UPDATE app_configs SET host=$1, port=$2, username=$3, password=$4, open_track=$5, click_track=$6 WHERE id=$7`, config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword, config.OpenTrack, config.ClickTrack, configID)

	return err
}
