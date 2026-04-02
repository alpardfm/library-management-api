package integration

import (
	"sync"
	"testing"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/internal/repository"
	"github.com/alpardfm/library-management-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupBorrowConcurrencyTest(t *testing.T, cfg service.BorrowServiceConfig) (*gorm.DB, service.BorrowService) {
	t.Helper()

	db, cleanup := setupIntegrationTestDB(t)
	t.Cleanup(cleanup)

	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	borrowRepo := repository.NewBorrowRepository(db)
	borrowService := service.NewBorrowService(db, borrowRepo, bookRepo, userRepo, cfg)

	return db, borrowService
}

func TestBorrowBook_ConcurrentDuplicateBorrow_OnlyOneSucceeds(t *testing.T) {
	db, borrowService := setupBorrowConcurrencyTest(t, service.BorrowServiceConfig{
		MaxBooksPerUser: 5,
		BorrowDays:      7,
		FinePerDay:      1000,
	})

	user := &models.User{
		Username:     "race-user",
		Email:        "race-user@example.com",
		PasswordHash: "hashed-password",
		Role:         models.RoleMember,
		IsActive:     true,
	}
	book := &models.Book{
		ISBN:            "9781234567811",
		Title:           "Race Condition Book",
		Author:          "Tester",
		TotalCopies:     2,
		AvailableCopies: 2,
	}

	require.NoError(t, db.Create(user).Error)
	require.NoError(t, db.Create(book).Error)

	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup

	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := borrowService.BorrowBook(user.ID, dto.BorrowBookRequest{BookID: book.ID})
			results <- err
		}()
	}

	close(start)
	wg.Wait()
	close(results)

	successCount := 0
	failedCount := 0
	for err := range results {
		if err == nil {
			successCount++
		} else {
			failedCount++
		}
	}

	assert.Equal(t, 1, successCount)
	assert.Equal(t, 1, failedCount)

	var activeBorrowCount int64
	require.NoError(t, db.Model(&models.BorrowRecord{}).
		Where("user_id = ? AND book_id = ? AND return_date IS NULL", user.ID, book.ID).
		Count(&activeBorrowCount).Error)
	assert.Equal(t, int64(1), activeBorrowCount)

	var updatedBook models.Book
	require.NoError(t, db.First(&updatedBook, book.ID).Error)
	assert.Equal(t, 1, updatedBook.AvailableCopies)
}
