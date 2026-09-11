package repositories

import (
	"database/sql"
	"encoding/json"
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
			t.variables,
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
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	var templates []models.TemplateData

	for rows.Next() {
		var template models.TemplateData
		var variables []byte
		if err := rows.Scan(&template.ID, &template.Name, &template.Slug, &template.Subject, &template.Type, &template.Content, &variables, &template.Status, &template.CreatedAt, &template.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(variables, &template.Variables); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	return templates, nil
}

func (r *TemplateRepository) GetAppTemplateSingle(templateID int64, userID int64) (*models.TemplateData, error) {

	var template models.TemplateData
	var variables []byte

	err := r.DB.QueryRow(`
		SELECT t.id, t.app_id, t.name, t.slug, t.subject, t.type, t.content, t.variables, t.status, t.created_at, t.updated_at
        FROM templates t
		JOIN apps a ON a.id = t.app_id
		WHERE
			t.id = $1
			AND a.user_id = $2
	`,
		templateID, userID).Scan(&template.ID, &template.AppID, &template.Name, &template.Slug, &template.Subject, &template.Type, &template.Content, &variables, &template.Status, &template.CreatedAt, &template.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(variables, &template.Variables); err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *TemplateRepository) GetAppTemplateBySlug(slug string, userID int64) (*models.TemplateData, error) {

	var template models.TemplateData
	var variables []byte

	err := r.DB.QueryRow(`
        SELECT t.id, t.name, t.slug, t.subject, t.type, t.content, t.variables, t.status, t.created_at, t.updated_at
        FROM templates t
		JOIN apps a ON a.id = t.app_id
		WHERE
			t.slug = $1
			AND a.user_id = $2
	`,
		slug, userID).Scan(&template.ID, &template.Name, &template.Slug, &template.Subject, &template.Type, &template.Content, &variables, &template.Status, &template.CreatedAt, &template.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(variables, &template.Variables); err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *TemplateRepository) CreateTemplate(template models.TemplateCreate) error {

	variables, err := json.Marshal(template.Variables)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(`INSERT INTO templates(app_id, name, slug, subject, type, content, variables) values($1, $2, $3, $4, $5, $6, $7)`, template.AppID, template.Name, template.Slug, template.Subject, template.Type, template.Content, string(variables))

	return err
}

func (r *TemplateRepository) UpdateTemplate(templateID int64, template models.TemplateUpdate) error {

	_, err := r.DB.Exec(`UPDATE templates SET name=$1, slug=$2, subject=$3, status=$4 WHERE id=$5`, template.Name, template.Slug, template.Subject, template.Status, templateID)

	return err
}

func (r *TemplateRepository) UpdateTemplateVariables(templateID int64, variables map[string]string) error {

	variableData, err := json.Marshal(variables)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(`UPDATE templates SET variables=$1 WHERE id=$2`, string(variableData), templateID)

	return err
}

func (r *TemplateRepository) UpdateTemplateContent(templateID int64, templateContent models.TemplateContent) error {

	variables, err := json.Marshal(templateContent.Variables)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(`UPDATE templates SET type=$1, content=$2, variables=$3 WHERE id=$4`, templateContent.Type, templateContent.Content, string(variables), templateID)

	return err
}

func (r *TemplateRepository) DeleteTemplate(templateID int64, userID int64) error {

	_, err := r.DB.Exec(`DELETE FROM templates WHERE id=$1 AND app_id IN (SELECT id FROM apps WHERE user_id=$2)`, templateID, userID)

	return err
}
