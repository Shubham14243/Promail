package configs

import (
	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Load .env file if it exists, but don't fail if it doesn't.
	// In production, environment variables are set via the system/container.
	_ = godotenv.Load()
}
