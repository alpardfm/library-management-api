package integration

import (
	"fmt"
	"testing"

	"github.com/alpardfm/library-management-api/configs"
	"github.com/alpardfm/library-management-api/internal/handler"
	"github.com/alpardfm/library-management-api/internal/repository"
	"github.com/alpardfm/library-management-api/internal/service"
	"github.com/alpardfm/library-management-api/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupIntegrationTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}

	t.Setenv("DB_NAME", "library_test")
	t.Setenv("DB_SSLMODE", "disable")

	db, err := database.Connect()
	if err != nil {
		t.Skipf("skipping integration test: test database is unavailable: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		sqlDB, sqlErr := db.DB()
		if sqlErr == nil {
			_ = sqlDB.Close()
		}
		t.Fatalf("failed to migrate integration test database: %v", err)
	}

	require.NoError(t, resetIntegrationTestDB(db))

	cleanup := func() {
		if err := resetIntegrationTestDB(db); err != nil {
			t.Fatalf("failed to clean integration test database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("failed to access integration sql db handle: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Fatalf("failed to close integration test database: %v", err)
		}
	}

	return db, cleanup
}

func resetIntegrationTestDB(db *gorm.DB) error {
	if err := db.Exec("TRUNCATE TABLE borrow_records, books, users RESTART IDENTITY CASCADE").Error; err != nil {
		return fmt.Errorf("truncate integration tables: %w", err)
	}
	return nil
}

func setupIntegrationRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	cfg := configs.Load()

	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	borrowRepo := repository.NewBorrowRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	bookService := service.NewBookService(bookRepo)
	borrowService := service.NewBorrowService(db, borrowRepo, bookRepo, userRepo, service.BorrowServiceConfig{
		MaxBooksPerUser: cfg.MaxBooksPerUser,
		BorrowDays:      cfg.BorrowDays,
		FinePerDay:      cfg.FinePerDay,
	})

	authHandler := handler.NewAuthHandler(authService)
	_ = handler.NewBookHandler(bookService)
	_ = handler.NewBorrowHandler(borrowService)

	router := gin.New()
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
	}

	return router
}
