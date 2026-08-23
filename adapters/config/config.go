package config

import (
	"fmt"
	"os"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/joho/godotenv"
)

// loadDotEnv loads variables from a .env file if present. It only warns on failure.
func loadDotEnv() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}
}

func DBConfig() db.DBConfig {
	return db.DBConfig{
		MaxConns:          20, // Increased for production load
		MinConns:          5,  // Increased minimum connections
		ConnTimeout:       15 * time.Second,
		MaxConnLifetime:   30 * time.Minute,
		MaxConnIdleTime:   5 * time.Minute,
		HealthCheckPeriod: time.Minute,
	}
}

type EnvConfiguration struct {
	HTTP_PORT    string
	DATABASE_URL string

	// Paystack Service
	PAYSTACK_BASE_URL   string
	PAYSTACK_SECRET_KEY string
	PAYSTACK_PUBLIC_KEY string
	// Redis Cache
	REDIS_URL string
}

// LoadEnvironmentVariables loads configuration from env and .env for the MEDIVUE service.
func LoadEnvironmentVariables() (*EnvConfiguration, error) {
	loadDotEnv()

	cfg := &EnvConfiguration{
		HTTP_PORT:           os.Getenv("HTTP_PORT"),
		DATABASE_URL:        os.Getenv("DATABASE_URL"),
		PAYSTACK_BASE_URL:   os.Getenv("PAYSTACK_BASE_URL"),
		PAYSTACK_PUBLIC_KEY: os.Getenv("PAYSTACK_PUBLIC_KEY"),
		PAYSTACK_SECRET_KEY: os.Getenv("PAYSTACK_SECRET_KEY"),
		REDIS_URL:           os.Getenv("REDIS_URL"),
	}

	// Validate required fields
	if cfg.HTTP_PORT == "" {
		return nil, fmt.Errorf("missing required environment variables: PORT")
	}

	return cfg, nil
}
