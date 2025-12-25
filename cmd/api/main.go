// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"library-management-api/configs"
	"library-management-api/internal/handler"
	"library-management-api/internal/middleware"
	"library-management-api/internal/repository"
	"library-management-api/internal/service"
	"library-management-api/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := configs.Load()

	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	borrowRepo := repository.NewBorrowRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	bookService := service.NewBookService(bookRepo)
	borrowService := service.NewBorrowService(borrowRepo, bookRepo, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	bookHandler := handler.NewBookHandler(bookService)
	borrowHandler := handler.NewBorrowHandler(borrowService)

	// Setup router
	router := gin.New()

	// Global middleware
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"app":     cfg.AppName,
			"version": cfg.AppVersion,
			"env":     cfg.AppEnv,
		})
	})

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Books
		books := protected.Group("/books")
		{
			books.GET("", bookHandler.ListBooks)
			books.GET("/:id", bookHandler.GetBook)

			// Admin/Librarian only
			books.POST("", middleware.RoleMiddleware("admin", "librarian"), bookHandler.CreateBook)
			books.PUT("/:id", middleware.RoleMiddleware("admin", "librarian"), bookHandler.UpdateBook)
			books.DELETE("/:id", middleware.RoleMiddleware("admin", "librarian"), bookHandler.DeleteBook)
		}

		// Borrow
		borrow := protected.Group("/borrow")
		{
			borrow.POST("", borrowHandler.BorrowBook)
			borrow.POST("/return", borrowHandler.ReturnBook)
			borrow.GET("/my-books", borrowHandler.GetMyBorrows)

			// Admin/Librarian only
			borrow.GET("/active", middleware.RoleMiddleware("admin", "librarian"), borrowHandler.GetActiveBorrows)
			borrow.GET("/overdue", middleware.RoleMiddleware("admin", "librarian"), borrowHandler.GetOverdueBorrows)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
