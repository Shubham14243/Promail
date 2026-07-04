package repositories

import (
	"database/sql"
	"promail/models"
)

type RefreshTokenRepository struct {
	DB *sql.DB
}

func (r *RefreshTokenRepository) ValidateRefreshToken(token string) (*models.RefreshToken, error) {

	var rt models.RefreshToken

	err := r.DB.QueryRow(`
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1 AND expires_at > NOW()
	`, token).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &rt, nil
}

func (r *RefreshTokenRepository) CreateRefreshToken(rt models.RefreshTokenCreate) error {

	_, err := r.DB.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET
			token = EXCLUDED.token,
			expires_at = EXCLUDED.expires_at
	`,
		rt.UserID,
		rt.Token,
		rt.ExpiresAt,
	)

	return err
}

func (r *RefreshTokenRepository) DeleteRefreshToken(userID int64) error {

	_, err := r.DB.Exec(`
		DELETE FROM refresh_tokens
		WHERE user_id=$1
	`, userID)

	return err
}
