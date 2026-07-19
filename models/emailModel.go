package models

import (
	"github.com/google/uuid"
)

type EmailSendTest struct {
	AppID   int64  `json:"app_id"`
	MailKey string `json:"mail_key"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailSend struct {
	AppID        int64             `json:"app_id"`
	TemplateSlug string            `json:"template_slug"`
	MailKey      string            `json:"mail_key"`
	To           string            `json:"to"`
	Variables    map[string]string `json:"variables"`
}

type EmailLogDataCreate struct {
	UUID         uuid.UUID `json:"uuid"`
	UserID       int64     `json:"user_id"`
	AppID        int64     `json:"app_id"`
	TemplateID   int64     `json:"template_id"`
	ToEmail      string    `json:"to_email"`
	Subject      string    `json:"subject"`
	VariableData string    `json:"variable_data"`
	Body         string    `json:"body"`
	Status       string    `json:"status"`
}

type LogResponse struct {
	AckID uuid.UUID `json:"ack_id"`
	LogID int64     `json:"log_id"`
}

type EmailLogData struct {
	ID           int64          `json:"id"`
	LogUUID      uuid.UUID      `json:"uuid"`
	UserID       int64          `json:"user_id"`
	AppID        int64          `json:"app_id"`
	TemplateID   *int64         `json:"template_id"`
	ToEmail      string         `json:"to_email"`
	Subject      string         `json:"subject"`
	Variables    *string        `json:"variables"`
	Body         string         `json:"body"`
	Status       string         `json:"status"`
	ErrorMessage *string        `json:"error_message"`
	SentAt       *string        `json:"sentAt"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
	Tracking     []TrackingData `json:"tracking"`
}

type TrackingData struct {
	ID            int64   `json:"id"`
	EmailLogID    int64   `json:"email_log_id"`
	Type          string  `json:"type"`
	Url           *string `json:"url"`
	TrackingToken string  `json:"trackingtoken"`
	OpenedAt      *string `json:"opened_at"`
	ClickedAt     *string `json:"clicked_at"`
}

type EmailLogFilter struct {
	AppID         *int64
	TemplateID    *int64
	ToEmail       string
	StartDateTime string
	EndDateTime   string
}

type ClickTracking struct {
	Token       uuid.UUID
	OriginalURL string
}
