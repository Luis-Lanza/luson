package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	ServerPort           string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	JWTExpiration        time.Duration
	JWTRefreshSecret     string
	JWTRefreshExpiration time.Duration
}

// Load reads configuration from .env file and environment variables.
// It returns an error if required variables are missing.
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.local")

	cfg := &Config{
		ServerPort:           getEnv("PORT", "8080"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "battery_pos"),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTExpiration:        parseDuration(getEnv("JWT_EXPIRATION", "15m")),
		JWTRefreshSecret:     getEnv("JWT_REFRESH_SECRET", ""),
		JWTRefreshExpiration: parseDuration(getEnv("JWT_REFRESH_EXPIRATION", "168h")),
	}

	// Validate required fields
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.JWTRefreshSecret == "" {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET is required")
	}

	return cfg, nil
}

// DatabaseURL returns the PostgreSQL connection string.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}
