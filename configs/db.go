package configs

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDB() {
	connStr := os.Getenv("DATABASE_URL")

	if connStr != "" && !containsQueryExecMode(connStr) {
		sep := "?"
		if hasQuery(connStr) {
			sep = "&"
		}
		connStr += sep + "default_query_exec_mode=simple_protocol"
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed to open postgres connection: %v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(90 * time.Second)

	if err = db.Ping(); err != nil {
		log.Fatalf("postgres ping failed: %v", err)
	}

	DB = db
}

func containsQueryExecMode(s string) bool {
	for i := 0; i+len("default_query_exec_mode") <= len(s); i++ {
		if s[i:i+len("default_query_exec_mode")] == "default_query_exec_mode" {
			return true
		}
	}
	return false
}

func hasQuery(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '?' {
			return true
		}
	}
	return false
}

func ParseIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}
