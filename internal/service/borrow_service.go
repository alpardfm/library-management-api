// internal/service/borrow_service.go
package service

import (
	"errors"
	"fmt"
	"time"

	"library-management-api/internal/dto"
	"library-management-api/internal/models"
	"library-management-api/internal/repository"
)

type BorrowService interface {
	BorrowBook(userID uint, req dto.BorrowBookRequest) (*models.BorrowRecord, error)
	ReturnBook(userID uint, req dto.ReturnBookRequest) (*models.BorrowRecord, int, error)
	GetUserBorrows(userID uint, page, limit int) ([]models.BorrowRecord, int64, error)
	GetActiveBorrows(page, limit int) ([]models.BorrowRecord, int64, error)
	GetOverdueBorrows(page, limit int) ([]models.BorrowRecord, int64, error)
	CalculateFine(borrowID uint) (int, error)
}

type borrowService struct {
	borrowRepo      repository.BorrowRepository
	bookRepo        repository.BookRepository
	userRepo        repository.UserRepository
	maxBooksPerUser int
}

func NewBorrowService(
	borrowRepo repository.BorrowRepository,
	bookRepo repository.BookRepository,
	userRepo repository.UserRepository,
) BorrowService {
	return &borrowService{
		borrowRepo:      borrowRepo,
		bookRepo:        bookRepo,
		userRepo:        userRepo,
		maxBooksPerUser: 5, // Default max 5 books per user
	}
}

func (s *borrowService) BorrowBook(userID uint, req dto.BorrowBookRequest) (*models.BorrowRecord, error) {
	// Check if user exists and is active
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Check if book exists
	book, err := s.bookRepo.FindByID(req.BookID)
	if err != nil {
		return nil, fmt.Errorf("book not found: %w", err)
	}

	// Check if book is available
	if !book.CanBorrow() {
		return nil, errors.New("book is not available for borrowing")
	}

	// Check user's active borrow count
	activeCount, err := s.borrowRepo.CountActiveByUser(userID)
	if err != nil {
		return nil, err
	}
	if activeCount >= int64(s.maxBooksPerUser) {
		return nil, fmt.Errorf("user has reached maximum borrow limit of %d books", s.maxBooksPerUser)
	}

	// Check if user already borrowed this book and hasn't returned it
	existingBorrow, _ := s.borrowRepo.FindActiveByUserAndBook(userID, req.BookID)
	if existingBorrow != nil {
		return nil, errors.New("user has already borrowed this book")
	}

	// Create borrow record
	borrowRecord := &models.BorrowRecord{
		UserID: userID,
		BookID: req.BookID,
	}

	// Set custom due date if provided
	if !req.DueDate.IsZero() {
		borrowRecord.DueDate = req.DueDate
	}

	// Use transaction to ensure data consistency
	// (In real implementation, use DB transaction)

	// Update book available copies
	if err := book.Borrow(); err != nil {
		return nil, err
	}
	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	// Create borrow record
	if err := s.borrowRepo.Create(borrowRecord); err != nil {
		// Rollback book update
		book.Return()
		s.bookRepo.Update(book)
		return nil, err
	}

	return borrowRecord, nil
}

func (s *borrowService) ReturnBook(userID uint, req dto.ReturnBookRequest) (*models.BorrowRecord, int, error) {
	// Get borrow record
	borrowRecord, err := s.borrowRepo.FindByID(req.BorrowRecordID)
	if err != nil {
		return nil, 0, fmt.Errorf("borrow record not found: %w", err)
	}

	// Check if user owns this borrow record or is admin/librarian
	if borrowRecord.UserID != userID {
		// In real app, check user role
		// For now, only allow user to return their own books
		return nil, 0, errors.New("not authorized to return this book")
	}

	// Check if already returned
	if borrowRecord.ReturnDate != nil {
		return nil, 0, errors.New("book already returned")
	}

	// Get book
	book, err := s.bookRepo.FindByID(borrowRecord.BookID)
	if err != nil {
		return nil, 0, fmt.Errorf("book not found: %w", err)
	}

	// Calculate fine
	fine := borrowRecord.CalculateFine()

	// Update book available copies
	book.Return()
	if err := s.bookRepo.Update(book); err != nil {
		return nil, 0, err
	}

	// Update borrow record
	now := time.Now()
	borrowRecord.ReturnDate = &now
	borrowRecord.Status = models.StatusReturned

	if err := s.borrowRepo.Update(borrowRecord); err != nil {
		// Rollback book update
		book.Borrow()
		s.bookRepo.Update(book)
		return nil, 0, err
	}

	return borrowRecord, fine, nil
}

func (s *borrowService) GetUserBorrows(userID uint, page, limit int) ([]models.BorrowRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	return s.borrowRepo.ListByUser(userID, page, limit)
}

func (s *borrowService) GetActiveBorrows(page, limit int) ([]models.BorrowRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	return s.borrowRepo.ListActive(page, limit)
}

func (s *borrowService) GetOverdueBorrows(page, limit int) ([]models.BorrowRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	return s.borrowRepo.ListOverdue(page, limit)
}

func (s *borrowService) CalculateFine(borrowID uint) (int, error) {
	borrowRecord, err := s.borrowRepo.FindByID(borrowID)
	if err != nil {
		return 0, fmt.Errorf("borrow record not found: %w", err)
	}

	return borrowRecord.CalculateFine(), nil
}
