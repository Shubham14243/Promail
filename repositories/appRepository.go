package repositories

import (
	"database/sql"
	"promail/models"

	"github.com/google/uuid"
)

type AppRepository struct {
	DB *sql.DB
}

func (r *AppRepository) AppExists(name string, userID int64) (bool, error) {

	var exists bool

	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM apps WHERE name = $1 AND user_id = $2
		)
	`, name, userID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AppRepository) AppExistsByID(id int64, userID int64) (bool, error) {

	var exists bool

	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM apps WHERE id = $1 AND user_id = $2
		)
	`, id, userID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AppRepository) GetUserApps(userID int64, limit int, offset int) ([]models.AppData, error) {

	rows, err := r.DB.Query(
		`SELECT id, name, description, status, created_at, updated_at
		FROM apps
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var apps []models.AppData

	for rows.Next() {
		var app models.AppData
		rows.Scan(&app.ID, &app.Name, &app.Description, &app.Status, &app.CreatedAt, &app.UpdatedAt)
		apps = append(apps, app)
	}

	return apps, nil
}

func (r *AppRepository) GetUserAppSingle(appID int64, userID int64) (*models.AppData, error) {

	var app models.AppData

	err := r.DB.QueryRow(
		`SELECT id, name, description, status, created_at, updated_at FROM apps WHERE id=$1 AND user_id=$2`, appID, userID).Scan(&app.ID, &app.Name, &app.Description, &app.Status, &app.CreatedAt, &app.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (r *AppRepository) GetUserAppKey(appID int64, userID int64) (*models.AppMailKey, error) {

	var app models.AppMailKey

	err := r.DB.QueryRow(`SELECT id, mail_key from apps where id=$1 AND user_id=$2`, appID, userID).Scan(&app.ID, &app.MailKey)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (r *AppRepository) CreateApp(app models.CreateApp) error {

	_, err := r.DB.Exec(`INSERT INTO apps(user_id, name, description, mail_Key, status) values($1, $2, $3, $4, $5)`, app.UserId, app.Name, app.Description, app.MailKey, app.Status)

	return err
}

func (r *AppRepository) UpdateApp(appID int64, app models.UpdateApp) error {

	_, err := r.DB.Exec(`UPDATE apps SET name=$1, description=$2, status=$3 WHERE user_id=$4 AND id=$5`, app.Name, app.Description, app.Status, app.UserId, appID)

	return err
}

func (r *AppRepository) RefreshMailKey(appID int64, userID int64, mailKey uuid.UUID) error {

	_, err := r.DB.Exec(`UPDATE apps SET mail_key=$1 WHERE id=$2 AND user_id=$3`, mailKey, appID, userID)

	return err
}

func (r *AppRepository) DeleteApp(appID int64, userID int64) error {

	_, err := r.DB.Exec(`DELETE FROM apps WHERE id=$1 AND user_id=$2`, appID, userID)

	return err
}
