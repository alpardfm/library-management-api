// internal/repository/book_repository.go
package repository

import (
	"library-management-api/internal/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *models.Book) error
	FindByID(id uint) (*models.Book, error)
	FindByISBN(isbn string) (*models.Book, error)
	Update(book *models.Book) error
	Delete(id uint) error
	List(page, limit int, search string) ([]models.Book, int64, error)
	UpdateAvailableCopies(id uint, change int) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) FindByID(id uint) (*models.Book, error) {
	var book models.Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) FindByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	err := r.db.Where("isbn = ?", isbn).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&models.Book{}, id).Error
}

func (r *bookRepository) List(page, limit int, search string) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.Book{})

	// Add search if provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR author ILIKE ? OR isbn ILIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	// Count total
	query.Count(&total)

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&books).Error

	return books, total, err
}

func (r *bookRepository) UpdateAvailableCopies(id uint, change int) error {
	return r.db.Model(&models.Book{}).
		Where("id = ?", id).
		Update("available_copies", gorm.Expr("available_copies + ?", change)).
		Error
}
