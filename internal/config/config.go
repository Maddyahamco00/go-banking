package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	HMAC     HMACConfig     `yaml:"hmac"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

// ServerConfig holds server settings
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// DatabaseConfig holds PostgreSQL settings
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

// RedisConfig holds Redis settings
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// JWTConfig holds JWT settings
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpiryHours int    `yaml:"expiry_hours"`
}

// HMACConfig holds HMAC settings
type HMACConfig struct {
	Secret string `yaml:"secret"`
}

// RateLimitConfig holds rate limiting settings
type RateLimitConfig struct {
	RequestsPerMin int `yaml:"requests_per_min"`
}

// DSN returns the PostgreSQL connection string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Name)
}

// Addr returns the Redis address
func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// Load reads configuration from config.yaml and .env
func Load(path string) (*Config, error) {
	// Load .env if it exists
	_ = godotenv.Load()

	cfg := &Config{}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(file, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Override with env vars if set
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}
	if hmacSecret := os.Getenv("HMAC_SECRET"); hmacSecret != "" {
		cfg.HMAC.Secret = hmacSecret
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}

	return cfg, nil
}