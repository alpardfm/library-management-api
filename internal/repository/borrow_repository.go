// internal/repository/borrow_repository.go
package repository

import (
	"time"

	"library-management-api/internal/models"

	"gorm.io/gorm"
)

type BorrowRepository interface {
	Create(record *models.BorrowRecord) error
	FindByID(id uint) (*models.BorrowRecord, error)
	FindActiveByUserAndBook(userID, bookID uint) (*models.BorrowRecord, error)
	Update(record *models.BorrowRecord) error
	ListByUser(userID uint, page, limit int) ([]models.BorrowRecord, int64, error)
	ListActive(page, limit int) ([]models.BorrowRecord, int64, error)
	ListOverdue(page, limit int) ([]models.BorrowRecord, int64, error)
	CountActiveByUser(userID uint) (int64, error)
}

type borrowRepository struct {
	db *gorm.DB
}

func NewBorrowRepository(db *gorm.DB) BorrowRepository {
	return &borrowRepository{db: db}
}

func (r *borrowRepository) Create(record *models.BorrowRecord) error {
	return r.db.Create(record).Error
}

func (r *borrowRepository) FindByID(id uint) (*models.BorrowRecord, error) {
	var record models.BorrowRecord
	err := r.db.Preload("User").Preload("Book").First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *borrowRepository) FindActiveByUserAndBook(userID, bookID uint) (*models.BorrowRecord, error) {
	var record models.BorrowRecord
	err := r.db.Where("user_id = ? AND book_id = ? AND status = ?",
		userID, bookID, models.StatusBorrowed).
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *borrowRepository) Update(record *models.BorrowRecord) error {
	return r.db.Save(record).Error
}

func (r *borrowRepository) ListByUser(userID uint, page, limit int) ([]models.BorrowRecord, int64, error) {
	var records []models.BorrowRecord
	var total int64

	offset := (page - 1) * limit

	query := r.db.Preload("Book").Where("user_id = ?", userID)
	query.Model(&models.BorrowRecord{}).Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&records).Error

	return records, total, err
}

func (r *borrowRepository) ListActive(page, limit int) ([]models.BorrowRecord, int64, error) {
	var records []models.BorrowRecord
	var total int64

	offset := (page - 1) * limit

	query := r.db.Preload("User").Preload("Book").
		Where("status = ?", models.StatusBorrowed)

	query.Model(&models.BorrowRecord{}).Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("due_date ASC").
		Find(&records).Error

	return records, total, err
}

func (r *borrowRepository) ListOverdue(page, limit int) ([]models.BorrowRecord, int64, error) {
	var records []models.BorrowRecord
	var total int64

	offset := (page - 1) * limit
	now := time.Now()

	query := r.db.Preload("User").Preload("Book").
		Where("status = ? OR (return_date IS NULL AND due_date < ?)",
			models.StatusOverdue, now)

	query.Model(&models.BorrowRecord{}).Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("due_date ASC").
		Find(&records).Error

	return records, total, err
}

func (r *borrowRepository) CountActiveByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.BorrowRecord{}).
		Where("user_id = ? AND status = ?", userID, models.StatusBorrowed).
		Count(&count).Error
	return count, err
}
