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
		SELECT user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1 AND expires_at > NOW()
	`, token).Scan(
		&rt.UserID,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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

func (r *RefreshTokenRepository) CreatePasswordResetToken(rt models.PasswordResetToken) error {

	_, err := r.DB.Exec(`
		INSERT INTO refresh_tokens (user_id, password_token, password_expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET
			password_token = EXCLUDED.password_token,
			password_expires_at = EXCLUDED.password_expires_at
	`,
		rt.UserID,
		rt.PasswordToken,
		rt.ExpiresAt,
	)

	return err
}

func (r *RefreshTokenRepository) GetPasswordUserByToken(password_token string) (int64, error) {

	var userID int64

	err := r.DB.QueryRow(`SELECT user_id FROM refresh_tokens WHERE password_token=$1 AND password_expires_at > NOW()`, password_token).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	_, err = r.DB.Exec(`
		UPDATE refresh_tokens
		SET password_token = NULL, password_expires_at = NULL
		WHERE user_id = $1 AND password_token = $2
	`, userID, password_token)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
