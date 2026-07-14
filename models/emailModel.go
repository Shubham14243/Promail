package models

type EmailSend struct {
	AppID   int64  `json:"app_id"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
