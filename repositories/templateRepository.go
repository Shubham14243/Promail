package repositories

import (
	"database/sql"
	"promail/models"
)

type TemplateRepository struct {
	DB *sql.DB
}

func (r *TemplateRepository) TemplateExistsBySlug(slug string, userID int64) (bool, error) {

	var exists bool

	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM templates t
			JOIN apps a ON a.id = t.app_id
			WHERE t.slug = $1
			  AND a.user_id = $2
		)
	`, slug, userID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TemplateRepository) TemplateExistsByID(templateID int64, userID int64) (bool, error) {

	var exists bool

	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM templates t
			JOIN apps a ON a.id = t.app_id
			WHERE t.id = $1
			  AND a.user_id = $2
		)
	`, templateID, userID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TemplateRepository) GetAppTemplates(appID int64, userID int64, limit int, offset int) ([]models.TemplateData, error) {

	rows, err := r.DB.Query(
		`SELECT
			t.id,
			t.name,
			t.slug,
			t.subject,
			t.type,
			t.content,
			t.status,
			t.created_at,
			t.updated_at
		FROM templates t
		JOIN apps a ON a.id = t.app_id
		WHERE a.id = $1
		AND a.user_id = $2
		ORDER BY t.created_at DESC
		LIMIT $3 OFFSET $4`,
		appID, userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var templates []models.TemplateData

	for rows.Next() {
		var template models.TemplateData
		rows.Scan(&template.ID, &template.Name, &template.Slug, &template.Subject, &template.Type, &template.Content, &template.Status, &template.CreatedAt, &template.UpdatedAt)
		templates = append(templates, template)
	}

	return templates, nil
}

func (r *TemplateRepository) GetAppTemplateSingle(templateID int64, userID int64) (*models.TemplateData, error) {

	var template models.TemplateData

	err := r.DB.QueryRow(`
        SELECT t.id, t.name, t.slug, t.subject, t.type, t.content, t.status, t.created_at, t.updated_at
        FROM templates t
		JOIN apps a ON a.id = t.app_id
		WHERE
			t.id = $1
			AND a.user_id = $2
	`,
		templateID, userID).Scan(&template.ID, &template.Name, &template.Slug, &template.Subject, &template.Type, &template.Content, &template.Status, &template.CreatedAt, &template.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *TemplateRepository) CreateTemplate(template models.TemplateCreate) error {

	_, err := r.DB.Exec(`INSERT INTO templates(app_id, name, slug, subject, type, content) values($1, $2, $3, $4, $5, $6)`, template.AppID, template.Name, template.Slug, template.Subject, template.Type, template.Content)

	return err
}

func (r *TemplateRepository) UpdateTemplate(templateID int64, template models.TemplateUpdate) error {

	_, err := r.DB.Exec(`UPDATE templates SET name=$1, slug=$2, subject=$3, status=$4 WHERE id=$5`, template.Name, template.Slug, template.Subject, template.Status, templateID)

	return err
}

func (r *TemplateRepository) UpdateTemplateContent(templateID int64, templateContent models.TemplateContent) error {

	_, err := r.DB.Exec(`UPDATE templates SET type=$1, content=$2 WHERE id=$3`, templateContent.Type, templateContent.Content, templateID)

	return err
}

func (r *TemplateRepository) DeleteTemplate(templateID int64, userID int64) error {

	_, err := r.DB.Exec(`DELETE FROM templates WHERE id=$1 AND app_id IN (SELECT id FROM apps WHERE user_id=$2)`, templateID, userID)

	return err
}
