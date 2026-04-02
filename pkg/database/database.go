// pkg/database/database.go
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alpardfm/library-management-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfig() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "library_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func Connect() (*gorm.DB, error) {
	config := NewConfig()

	// Custom logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Enable connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	models := []interface{}{
		&models.User{},
		&models.Book{},
		&models.BorrowRecord{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	if db.Dialector.Name() != "postgres" {
		return nil
	}

	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conname = 'available_copies_non_negative'
			) THEN
				ALTER TABLE books
				ADD CONSTRAINT available_copies_non_negative
				CHECK (available_copies >= 0);
			END IF;
		END $$;
	`).Error; err != nil {
		return fmt.Errorf("failed to create non-negative stock constraint: %w", err)
	}

	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conname = 'available_copies_not_exceed_total'
			) THEN
				ALTER TABLE books
				ADD CONSTRAINT available_copies_not_exceed_total
				CHECK (available_copies <= total_copies);
			END IF;
		END $$;
	`).Error; err != nil {
		return fmt.Errorf("failed to create max stock constraint: %w", err)
	}

	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_borrow_records_active_user_book
		ON borrow_records (user_id, book_id)
		WHERE return_date IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to create active borrow unique index: %w", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
