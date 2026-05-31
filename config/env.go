package config

import (
	"github.com/joho/godotenv"
)

// LoadDotEnv loads .env (if present) to simplify local development.
// In production, environment variables should be injected externally (e.g., Kubernetes, Docker secrets).
func LoadDotEnv() error {
	// Ignore if missing; for CI/production this file should not be required.
	return godotenv.Load()
}
