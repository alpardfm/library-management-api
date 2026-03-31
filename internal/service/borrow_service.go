package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/internal/repository"
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
	borrowRepo repository.BorrowRepository
	bookRepo   repository.BookRepository
	userRepo   repository.UserRepository
	config     BorrowServiceConfig
}

type BorrowServiceConfig struct {
	MaxBooksPerUser int
	BorrowDays      int
	FinePerDay      int
}

func NewBorrowService(
	borrowRepo repository.BorrowRepository,
	bookRepo repository.BookRepository,
	userRepo repository.UserRepository,
	config BorrowServiceConfig,
) BorrowService {
	return &borrowService{
		borrowRepo: borrowRepo,
		bookRepo:   bookRepo,
		userRepo:   userRepo,
		config:     config,
	}
}

func (s *borrowService) BorrowBook(userID uint, req dto.BorrowBookRequest) (*models.BorrowRecord, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	book, err := s.bookRepo.FindByID(req.BookID)
	if err != nil {
		return nil, fmt.Errorf("book not found: %w", err)
	}

	if !book.CanBorrow() {
		return nil, errors.New("book is not available for borrowing")
	}

	activeCount, err := s.borrowRepo.CountActiveByUser(userID)
	if err != nil {
		return nil, err
	}
	if activeCount >= int64(s.config.MaxBooksPerUser) {
		return nil, fmt.Errorf("user has reached maximum borrow limit of %d books", s.config.MaxBooksPerUser)
	}

	existingBorrow, _ := s.borrowRepo.FindActiveByUserAndBook(userID, req.BookID)
	if existingBorrow != nil {
		return nil, errors.New("user has already borrowed this book")
	}

	borrowRecord := &models.BorrowRecord{
		UserID:     userID,
		BookID:     req.BookID,
		BorrowDate: time.Now(),
	}

	if !req.DueDate.IsZero() {
		borrowRecord.DueDate = req.DueDate
	} else {
		borrowRecord.DueDate = borrowRecord.BorrowDate.Add(time.Duration(s.config.BorrowDays) * 24 * time.Hour)
	}
	if err := book.Borrow(); err != nil {
		return nil, err
	}
	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	if err := s.borrowRepo.Create(borrowRecord); err != nil {
		book.Return()
		s.bookRepo.Update(book)
		return nil, err
	}

	return borrowRecord, nil
}

func (s *borrowService) ReturnBook(userID uint, req dto.ReturnBookRequest) (*models.BorrowRecord, int, error) {
	borrowRecord, err := s.borrowRepo.FindByID(req.BorrowRecordID)
	if err != nil {
		return nil, 0, fmt.Errorf("borrow record not found: %w", err)
	}

	if borrowRecord.UserID != userID {
		return nil, 0, errors.New("not authorized to return this book")
	}

	if borrowRecord.ReturnDate != nil {
		return nil, 0, errors.New("book already returned")
	}

	book, err := s.bookRepo.FindByID(borrowRecord.BookID)
	if err != nil {
		return nil, 0, fmt.Errorf("book not found: %w", err)
	}

	fine := borrowRecord.CalculateFine(s.config.FinePerDay)

	book.Return()
	if err := s.bookRepo.Update(book); err != nil {
		return nil, 0, err
	}

	now := time.Now()
	borrowRecord.ReturnDate = &now
	borrowRecord.Status = models.StatusReturned

	if err := s.borrowRepo.Update(borrowRecord); err != nil {
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

	return borrowRecord.CalculateFine(s.config.FinePerDay), nil
}
