package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/internal/repository"
	"gorm.io/gorm"
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
	db         *gorm.DB
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
	db *gorm.DB,
	borrowRepo repository.BorrowRepository,
	bookRepo repository.BookRepository,
	userRepo repository.UserRepository,
	config BorrowServiceConfig,
) BorrowService {
	return &borrowService{
		db:         db,
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

	activeCount, err := s.borrowRepo.CountActiveByUser(userID)
	if err != nil {
		return nil, err
	}
	if activeCount >= int64(s.config.MaxBooksPerUser) {
		return nil, fmt.Errorf("user has reached maximum borrow limit of %d books", s.config.MaxBooksPerUser)
	}

	var borrowRecord *models.BorrowRecord

	err = s.db.Transaction(func(tx *gorm.DB) error {
		bookRepoTx := s.bookRepo.WithTx(tx)
		borrowRepoTx := s.borrowRepo.WithTx(tx)

		book, err := bookRepoTx.FindByIDForUpdate(req.BookID)
		if err != nil {
			return fmt.Errorf("book not found: %w", err)
		}

		if !book.CanBorrow() {
			return errors.New("book is not available for borrowing")
		}

		existingBorrow, err := borrowRepoTx.FindActiveByUserAndBook(userID, req.BookID)
		if err == nil && existingBorrow != nil {
			return errors.New("user has already borrowed this book")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		borrowRecord = &models.BorrowRecord{
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
			return err
		}
		if err := bookRepoTx.Update(book); err != nil {
			return err
		}

		if err := borrowRepoTx.Create(borrowRecord); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return borrowRecord, nil
}

func (s *borrowService) ReturnBook(userID uint, req dto.ReturnBookRequest) (*models.BorrowRecord, int, error) {
	var borrowRecord *models.BorrowRecord
	var fine int
	var err error
	err = s.db.Transaction(func(tx *gorm.DB) error {
		bookRepoTx := s.bookRepo.WithTx(tx)
		borrowRepoTx := s.borrowRepo.WithTx(tx)

		borrowRecord, err = borrowRepoTx.FindByIDForUpdate(req.BorrowRecordID)
		if err != nil {
			return fmt.Errorf("borrow record not found: %w", err)
		}

		if borrowRecord.UserID != userID {
			return errors.New("not authorized to return this book")
		}

		if borrowRecord.ReturnDate != nil {
			return errors.New("book already returned")
		}

		book, err := bookRepoTx.FindByIDForUpdate(borrowRecord.BookID)
		if err != nil {
			return fmt.Errorf("book not found: %w", err)
		}

		fine = borrowRecord.CalculateFine(s.config.FinePerDay)

		book.Return()
		if err := bookRepoTx.Update(book); err != nil {
			return fmt.Errorf("failed to update book: %w", err)
		}

		now := time.Now()
		borrowRecord.ReturnDate = &now
		borrowRecord.Status = models.StatusReturned

		if err := borrowRepoTx.Update(borrowRecord); err != nil {
			return fmt.Errorf("failed to update borrow record: %w", err)
		}

		return nil
	})
	if err != nil {
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
