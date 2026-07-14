package repositories

import (
	"database/sql"
)

type EmailRepository struct {
	DB *sql.DB
}
