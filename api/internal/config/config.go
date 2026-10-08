package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv, ServerPort, DatabaseURL, JWTSecret, AIAPIKey, AIModel string
	JWTExpiration                                                 time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load()
	expiration, err := time.ParseDuration(value("JWT_EXPIRATION", "24h"))
	if err != nil {
		return Config{}, err
	}
	c := Config{AppEnv: value("APP_ENV", "development"), ServerPort: value("SERVER_PORT", "8080"), DatabaseURL: os.Getenv("DATABASE_URL"), JWTSecret: os.Getenv("JWT_SECRET"), AIAPIKey: os.Getenv("AI_API_KEY"), AIModel: os.Getenv("AI_MODEL"), JWTExpiration: expiration}
	if c.DatabaseURL == "" || c.JWTSecret == "" {
		return Config{}, errors.New("DATABASE_URL and JWT_SECRET are required")
	}
	return c, nil
}

func value(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
