// internal/service/book_service.go
package service

import (
	"errors"
	"fmt"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/internal/repository"
)

type BookService interface {
	CreateBook(req dto.CreateBookRequest) (*models.Book, error)
	GetBookByID(id uint) (*models.Book, error)
	UpdateBook(id uint, req dto.UpdateBookRequest) (*models.Book, error)
	DeleteBook(id uint) error
	ListBooks(page, limit int, search string) ([]models.Book, int64, error)
	CheckAvailability(id uint) (bool, error)
}

type bookService struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{bookRepo: bookRepo}
}

func (s *bookService) CreateBook(req dto.CreateBookRequest) (*models.Book, error) {
	// Check if ISBN exists
	existingBook, _ := s.bookRepo.FindByISBN(req.ISBN)
	if existingBook != nil {
		return nil, errors.New("book with this ISBN already exists")
	}

	// Create book
	book := &models.Book{
		ISBN:            req.ISBN,
		Title:           req.Title,
		Author:          req.Author,
		Publisher:       req.Publisher,
		PublicationYear: req.PublicationYear,
		Genre:           req.Genre,
		Description:     req.Description,
		TotalCopies:     req.TotalCopies,
		AvailableCopies: req.TotalCopies,
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) GetBookByID(id uint) (*models.Book, error) {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("book not found: %w", err)
	}
	return book, nil
}

func (s *bookService) UpdateBook(id uint, req dto.UpdateBookRequest) (*models.Book, error) {
	// Get existing book
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("book not found: %w", err)
	}

	// Update fields if provided
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Publisher != "" {
		book.Publisher = req.Publisher
	}
	if req.PublicationYear > 0 {
		book.PublicationYear = req.PublicationYear
	}
	if req.Genre != "" {
		book.Genre = req.Genre
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.TotalCopies > 0 {
		// Adjust available copies
		diff := req.TotalCopies - book.TotalCopies
		book.TotalCopies = req.TotalCopies
		book.AvailableCopies += diff

		// Ensure available copies doesn't go negative
		if book.AvailableCopies < 0 {
			book.AvailableCopies = 0
		}
	}

	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) DeleteBook(id uint) error {
	// Check if book exists
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("book not found: %w", err)
	}

	// Check if book has active borrows
	if book.AvailableCopies != book.TotalCopies {
		return errors.New("cannot delete book with active borrows")
	}

	return s.bookRepo.Delete(id)
}

func (s *bookService) ListBooks(page, limit int, search string) ([]models.Book, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.bookRepo.List(page, limit, search)
}

func (s *bookService) CheckAvailability(id uint) (bool, error) {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return false, fmt.Errorf("book not found: %w", err)
	}

	return book.CanBorrow(), nil
}
