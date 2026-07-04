package repositories

import (
	"database/sql"
	"promail/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {

	rows, err := r.DB.Query(`SELECT id, uuid, name, email, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {

		var user models.User

		rows.Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		users = append(users, user)
	}

	return users, nil

}

func (r *UserRepository) GetUserByID(id int64) (*models.User, error) {

	var user models.User

	err := r.DB.QueryRow(`SELECT id, uuid, name, email, created_at, updated_at FROM users where id=$1`, id).Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {

	var user models.User

	err := r.DB.QueryRow(`SELECT id, email, password_hash FROM users where email=$1`, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *UserRepository) UserExists(email string) (bool, error) {

	var exists bool

	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM users WHERE email = $1
		)
	`, email).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) CreateUser(user models.User) error {

	_, err := r.DB.Exec(`
		INSERT INTO users(uuid, name, email, password_hash)
		VALUES($1,$2,$3,$4)
	`,
		user.UUID,
		user.Name,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (r *UserRepository) UpdateUser(id int64, user models.UserUpdateRequest) error {

	_, err := r.DB.Exec(`
		UPDATE users
		SET name=$1,email=$2
		WHERE id=$3
	`,
		user.Name,
		user.Email,
		id,
	)

	return err
}

func (r *UserRepository) DeleteUser(id int64) error {

	_, err := r.DB.Exec(`
		DELETE FROM users
		WHERE id=$1
	`, id)

	return err
}
