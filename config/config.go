package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort int

	DBHost    string
	DBPort    int
	DBName    string
	DBUser    string
	DBPassword string
	DBSSLMode string

	LogLevel string
}

func Load() (AppConfig, error) {
	getString := func(key, def string) string {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
		return def
	}

	getInt := func(key string, def int) (int, error) {
		v := strings.TrimSpace(os.Getenv(key))
		if v == "" {
			return def, nil
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("%s must be int: %w", key, err)
		}
		return i, nil
	}

	appPort, err := getInt("APP_PORT", 8080)
	if err != nil {
		return AppConfig{}, err
	}
	if appPort <= 0 || appPort > 65535 {
		return AppConfig{}, fmt.Errorf("APP_PORT out of range")
	}

	dbPort, err := getInt("DB_PORT", 5432)
	if err != nil {
		return AppConfig{}, err
	}
	if dbPort <= 0 || dbPort > 65535 {
		return AppConfig{}, fmt.Errorf("DB_PORT out of range")
	}

	cfg := AppConfig{
		AppName:     getString("APP_NAME", "gobanking-v2"),
		AppEnv:      getString("APP_ENV", "development"),
		AppPort:     appPort,
		DBHost:      getString("DB_HOST", "localhost"),
		DBPort:      dbPort,
		DBName:      getString("DB_NAME", "gobanking"),
		DBUser:      getString("DB_USER", "postgres"),
		DBPassword:  getString("DB_PASSWORD", "postgres"),
		DBSSLMode:   getString("DB_SSLMODE", "disable"),
		LogLevel:    strings.ToLower(getString("LOG_LEVEL", "info")),
	}

	// minimal validation for required envs in production
	if cfg.AppEnv == "production" {
		if strings.TrimSpace(cfg.DBPassword) == "" {
			return AppConfig{}, fmt.Errorf("DB_PASSWORD is required in production")
		}
	}

	return cfg, nil
}
