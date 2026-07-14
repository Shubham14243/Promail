package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateApp struct {
	UserId      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MailKey     uuid.UUID `json:"mail_key"`
	Status      string    `json:"status"`
}

type UpdateApp struct {
	UserId      int64  `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type AppData struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AppMailKey struct {
	ID      int64     `json:"id"`
	MailKey uuid.UUID `json:"mail_key"`
}
