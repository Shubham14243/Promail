package configs

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	connStr := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Disable server-side prepared statements.
	// Required when running behind PgBouncer in transaction-pooling mode
	// (or any pooler that rotates physical connections between transactions).
	// Without this, lib/pq caches prepared statement names on a connection
	// and fails with: pq: unnamed prepared statement does not exist (26000)
	// when the connection is rotated to a different backend.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	DB = db
}
