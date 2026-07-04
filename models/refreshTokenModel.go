package models

import (
	"time"
)

type RefreshToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"uuid"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenCreate struct {
	UserID    int64     `json:"user_id"`
	Token     string    `json:"uuid"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AuthToken string `json:"auth_token"`
}
