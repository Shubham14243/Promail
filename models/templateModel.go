package models

type TemplateCreate struct {
	AppID     int64             `json:"app_id"`
	Name      string            `json:"name"`
	Slug      string            `json:"slug"`
	Subject   string            `json:"subject"`
	Type      string            `json:"type"`
	Content   string            `json:"content"`
	Variables map[string]string `json:"variables"`
	Status    string            `json:"status"`
}

type TemplateUpdate struct {
	Name      string            `json:"name"`
	Slug      string            `json:"slug"`
	Subject   string            `json:"subject"`
	Variables map[string]string `json:"variables"`
	Status    string            `json:"status"`
}

type TemplateContent struct {
	Type      string            `json:"type"`
	Content   string            `json:"content"`
	Variables map[string]string `json:"variables"`
}

type TemplateData struct {
	ID        int64             `json:"template_id"`
	AppID     int64             `json:"app_id"`
	Name      string            `json:"name"`
	Slug      string            `json:"slug"`
	Subject   string            `json:"subject"`
	Type      string            `json:"type"`
	Content   string            `json:"content"`
	Variables map[string]string `json:"variables"`
	Status    string            `json:"status"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}
