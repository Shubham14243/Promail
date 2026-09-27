package models

import (
	"time"
)

type RefreshToken struct {
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

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type PasswordResetToken struct {
	UserID        int64     `json:"user_id"`
	PasswordToken string    `json:"password_token"`
	ExpiresAt     time.Time `json:"expires_at"`
}

type SetNewPassword struct {
	PasswordToken string `json:"password_token"`
	NewPassword   string `json:"new_password"`
}
