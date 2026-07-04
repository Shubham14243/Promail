package configs

import (
	"log"
)

func Migrate() error {

	query := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		uuid UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL UNIQUE,
		token TEXT NOT NULL UNIQUE,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS apps (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		name VARCHAR(50) NOT NULL UNIQUE,
		description TEXT,
		sender_name VARCHAR(50) NOT NULL,
		sender_email VARCHAR(100) NOT NULL,
		mail_key UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
		status VARCHAR(20) NOT NULL DEFAULT 'active',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		CONSTRAINT fk_app_user
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS templates (
    id BIGSERIAL PRIMARY KEY,
    app_id INT NOT NULL,

    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,

    subject TEXT NOT NULL,

    type VARCHAR(10) NOT NULL
        CHECK (type IN ('html', 'text')),

    status VARCHAR(10) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'inactive')),

    content TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_template_app
        FOREIGN KEY (app_id)
        REFERENCES apps(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_template_slug
        UNIQUE (app_id, slug)
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Println("Migration failed:", err)
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}
