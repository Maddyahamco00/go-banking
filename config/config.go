package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
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

func Load() (Config, error) {
	getEnv := func(key string) (string, bool) {
		v, ok := os.LookupEnv(key)
		return v, ok
	}

	require := func(key string) (string, error) {
		v, ok := getEnv(key)
		if !ok || strings.TrimSpace(v) == "" {
			return "", fmt.Errorf("missing required env var %s", key)
		}
		return v, nil
	}

	// App
	appName, err := require("APP_NAME")
	if err != nil {
		return Config{}, err
	}
	appEnv, err := require("APP_ENV")
	if err != nil {
		return Config{}, err
	}

	appPortStr, ok := getEnv("APP_PORT")
	if !ok || strings.TrimSpace(appPortStr) == "" {
		appPortStr = "8080"
	}
	appPort, err := strconv.Atoi(appPortStr)
	if err != nil {
		return Config{}, errors.New("APP_PORT must be an integer")
	}

	// DB
	dbHost, err := require("DB_HOST")
	if err != nil {
		return Config{}, err
	}
	// Optional: keep DB_PORT default for local/dev.
	dbPortStr, ok := getEnv("DB_PORT")
	if !ok || strings.TrimSpace(dbPortStr) == "" {
		dbPortStr = "5432"
	}
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return Config{}, errors.New("DB_PORT must be an integer")
	}

	dbName, err := require("DB_NAME")
	if err != nil {
		return Config{}, err
	}
	dbUser, err := require("DB_USER")
	if err != nil {
		return Config{}, err
	}
	dbPassword, err := require("DB_PASSWORD")
	if err != nil {
		return Config{}, err
	}
	dbSSLMode, ok := getEnv("DB_SSLMODE")
	if !ok || strings.TrimSpace(dbSSLMode) == "" {
		dbSSLMode = "disable"
	}

	// Logging
	logLevel, ok := getEnv("LOG_LEVEL")
	if !ok || strings.TrimSpace(logLevel) == "" {
		logLevel = "info"
	}

	return Config{
		AppName: appName,
		AppEnv:  appEnv,
		AppPort: appPort,

		DBHost:     dbHost,
		DBPort:     dbPort,
		DBName:     dbName,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBSSLMode:  dbSSLMode,

		LogLevel: logLevel,
	}, nil
}

