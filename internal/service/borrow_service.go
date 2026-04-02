package service

import (
	"errors"
	"time"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/internal/repository"
	"github.com/alpardfm/library-management-api/pkg/apperror"
	"gorm.io/gorm"
)

type BorrowService interface {
	BorrowBook(userID uint, req dto.BorrowBookRequest) (*models.BorrowRecord, error)
	ReturnBook(userID uint, role string, req dto.ReturnBookRequest) (*models.BorrowRecord, int, error)
	GetUserBorrows(userID uint, page, limit int, sort string) ([]models.BorrowRecord, int64, error)
	GetActiveBorrows(page, limit int, sort string) ([]models.BorrowRecord, int64, error)
	GetOverdueBorrows(page, limit int, sort string) ([]models.BorrowRecord, int64, error)
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
	var borrowRecord *models.BorrowRecord

	err := s.db.Transaction(func(tx *gorm.DB) error {
		userRepoTx := s.userRepo.WithTx(tx)
		bookRepoTx := s.bookRepo.WithTx(tx)
		borrowRepoTx := s.borrowRepo.WithTx(tx)

		user, err := userRepoTx.FindByIDForUpdate(userID)
		if err != nil {
			return apperror.NotFound("user")
		}
		if !user.IsActive {
			return apperror.Forbidden("user account is deactivated")
		}

		activeCount, err := borrowRepoTx.CountActiveByUser(userID)
		if err != nil {
			return apperror.Internal("failed to count active borrows", err)
		}
		if activeCount >= int64(s.config.MaxBooksPerUser) {
			return apperror.Conflict("user has reached maximum borrow limit")
		}

		book, err := bookRepoTx.FindByIDForUpdate(req.BookID)
		if err != nil {
			return apperror.NotFound("book")
		}
		if err := validateBookStock(book); err != nil {
			return err
		}

		if !book.CanBorrow() {
			return apperror.Conflict("book is not available for borrowing")
		}

		existingBorrow, err := borrowRepoTx.FindActiveByUserAndBook(userID, req.BookID)
		if err == nil && existingBorrow != nil {
			return apperror.Conflict("user has already borrowed this book")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.Internal("failed to check active borrow", err)
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
			return apperror.Conflict(err.Error())
		}
		if err := bookRepoTx.Update(book); err != nil {
			return apperror.Internal("failed to update book", err)
		}

		if err := borrowRepoTx.Create(borrowRecord); err != nil {
			return apperror.Internal("failed to create borrow record", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return borrowRecord, nil
}

func (s *borrowService) ReturnBook(userID uint, role string, req dto.ReturnBookRequest) (*models.BorrowRecord, int, error) {
	var borrowRecord *models.BorrowRecord
	var fine int
	var err error
	err = s.db.Transaction(func(tx *gorm.DB) error {
		bookRepoTx := s.bookRepo.WithTx(tx)
		borrowRepoTx := s.borrowRepo.WithTx(tx)

		borrowRecord, err = borrowRepoTx.FindByIDForUpdate(req.BorrowRecordID)
		if err != nil {
			return apperror.NotFound("borrow record")
		}

		if role == "" {
			role = string(models.RoleMember)
		}

		if !canManageBorrowReturn(role) && borrowRecord.UserID != userID {
			return apperror.Forbidden("not authorized to return this book")
		}

		if borrowRecord.ReturnDate != nil {
			return apperror.Conflict("book already returned")
		}

		book, err := bookRepoTx.FindByIDForUpdate(borrowRecord.BookID)
		if err != nil {
			return apperror.NotFound("book")
		}
		if err := validateBookStock(book); err != nil {
			return err
		}
		if book.AvailableCopies >= book.TotalCopies {
			return apperror.Conflict("book stock is already full, cannot process return")
		}

		fine = borrowRecord.CalculateFine(s.config.FinePerDay)

		book.Return()
		if err := bookRepoTx.Update(book); err != nil {
			return apperror.Internal("failed to update book", err)
		}

		now := time.Now()
		borrowRecord.ReturnDate = &now
		borrowRecord.Status = models.StatusReturned

		if err := borrowRepoTx.Update(borrowRecord); err != nil {
			return apperror.Internal("failed to update borrow record", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return borrowRecord, fine, nil
}

func canManageBorrowReturn(role string) bool {
	return role == string(models.RoleAdmin) || role == string(models.RoleLibrarian)
}

func (s *borrowService) GetUserBorrows(userID uint, page, limit int, sort string) ([]models.BorrowRecord, int64, error) {
	return s.borrowRepo.ListByUser(userID, page, limit, sort)
}

func (s *borrowService) GetActiveBorrows(page, limit int, sort string) ([]models.BorrowRecord, int64, error) {
	return s.borrowRepo.ListActive(page, limit, sort)
}

func (s *borrowService) GetOverdueBorrows(page, limit int, sort string) ([]models.BorrowRecord, int64, error) {
	return s.borrowRepo.ListOverdue(page, limit, sort)
}

func (s *borrowService) CalculateFine(borrowID uint) (int, error) {
	borrowRecord, err := s.borrowRepo.FindByID(borrowID)
	if err != nil {
		return 0, apperror.NotFound("borrow record")
	}

	return borrowRecord.CalculateFine(s.config.FinePerDay), nil
}
