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
		user_id BIGINT NOT NULL UNIQUE,
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
		user_id BIGINT NOT NULL,
		name VARCHAR(50) NOT NULL UNIQUE,
		description TEXT,
		mail_key UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
		status VARCHAR(20) NOT NULL DEFAULT 'active',
			CHECK (status IN ('active', 'inactive')),
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		CONSTRAINT fk_app_user
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS templates (
		id BIGSERIAL PRIMARY KEY,
		app_id BIGINT NOT NULL,
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

	CREATE TABLE IF NOT EXISTS app_configs (
		id SERIAL PRIMARY KEY,
		app_id BIGINT NOT NULL UNIQUE,
		host VARCHAR(100) NOT NULL,
		port INT NOT NULL DEFAULT 587,
		name VARCHAR(100),
		username VARCHAR(100) NOT NULL,
		password TEXT NOT NULL,
		open_track VARCHAR(10) NOT NULL DEFAULT 'active'
        	CHECK (open_track IN ('active', 'inactive')),
		click_track VARCHAR(10) NOT NULL DEFAULT 'active'
        	CHECK (click_track IN ('active', 'inactive')),
		auto_retry VARCHAR(10) NOT NULL DEFAULT 'active'
        	CHECK (auto_retry IN ('active', 'inactive')),
		retry_max_count INT NOT NULL DEFAULT 3,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		CONSTRAINT fk_smtp_app
			FOREIGN KEY (app_id)
			REFERENCES apps(id)
			ON DELETE CASCADE,
		CONSTRAINT uq_config_app
        	UNIQUE (app_id, id)
	);

	CREATE TABLE IF NOT EXISTS email_logs (
		id BIGSERIAL PRIMARY KEY,
		uuid UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
		user_id BIGINT NOT NULL,
		app_id BIGINT NOT NULL,
		template_id BIGINT,
		to_email VARCHAR(100) NOT NULL,
		subject TEXT NOT NULL,
		variable_data TEXT NOT NULL,
		rendered_body TEXT NOT NULL DEFAULT '',
		status VARCHAR(20) NOT NULL DEFAULT 'queued'
			CHECK (status IN ('queued', 'processing', 'sent', 'opened', 'failed')),
		error_message TEXT,
		sent_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		CONSTRAINT fk_email_app
			FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
		CONSTRAINT fk_email_template
			FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS email_queue (
		id BIGSERIAL PRIMARY KEY,
		email_log_id BIGINT NOT NULL UNIQUE,
		attempts INT NOT NULL DEFAULT 0,
		last_attempted_at TIMESTAMPTZ,
		CONSTRAINT fk_queue_email
			FOREIGN KEY (email_log_id) REFERENCES email_logs(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS email_analytics (
		id BIGSERIAL PRIMARY KEY,
		email_log_id BIGINT NOT NULL,
		type VARCHAR(10) NOT NULL CHECK (type IN ('open', 'click')),
		original_url TEXT,
		tracking_token UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
		opened_at TIMESTAMPTZ,
		clicked_at TIMESTAMPTZ,
		CONSTRAINT fk_click_email
			FOREIGN KEY (email_log_id) REFERENCES email_logs(id) ON DELETE CASCADE
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
