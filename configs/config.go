// configs/config.go
package configs

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server
	AppName    string
	AppEnv     string
	AppPort    string
	AppVersion string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret string
	JWTExpiry time.Duration

	// Server Timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Application
	MaxBooksPerUser int
	BorrowDays      int
	FinePerDay      int
}

func Load() *Config {
	return &Config{
		// Server
		AppName:    getEnv("APP_NAME", "Library Management API"),
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8080"),
		AppVersion: getEnv("APP_VERSION", "1.0.0"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "library_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTExpiry: parseDuration(getEnv("JWT_EXPIRY", "24h")),

		// Server Timeouts
		ReadTimeout:  parseDuration(getEnv("READ_TIMEOUT", "10s")),
		WriteTimeout: parseDuration(getEnv("WRITE_TIMEOUT", "10s")),
		IdleTimeout:  parseDuration(getEnv("IDLE_TIMEOUT", "60s")),

		// Application
		MaxBooksPerUser: parseInt(getEnv("MAX_BOOKS_PER_USER", "5")),
		BorrowDays:      parseInt(getEnv("BORROW_DAYS", "14")),
		FinePerDay:      parseInt(getEnv("FINE_PER_DAY", "1000")),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(value string) int {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	return 0
}

func parseDuration(value string) time.Duration {
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}
	return 0
}
