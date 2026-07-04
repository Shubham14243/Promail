package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           int64     `json:"id"`
	UUID         uuid.UUID `json:"uuid"`
	Name         string    `json:"name"`
	Email        string    `json:"Email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID           int64     `json:"id"`
	UUID         uuid.UUID `json:"uuid"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	AuthToken    string    `json:"auth_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MeResponse struct {
	ID        int64     `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
