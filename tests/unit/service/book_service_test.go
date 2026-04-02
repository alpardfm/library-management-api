package service_test

import (
	"errors"
	"testing"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookService_CreateBook(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	req := dto.CreateBookRequest{
		ISBN:        "9781234567897",
		Title:       "Test Book",
		Author:      "Test Author",
		TotalCopies: 5,
	}

	mockRepo.On("FindByISBN", req.ISBN).
		Return((*models.Book)(nil), errors.New("not found")).
		Once()

	mockRepo.On("Create", mock.AnythingOfType("*models.Book")).
		Run(func(args mock.Arguments) {
			book := args.Get(0).(*models.Book)
			book.ID = 1
			assert.Equal(t, req.ISBN, book.ISBN)
			assert.Equal(t, req.Title, book.Title)
			assert.Equal(t, req.Author, book.Author)
			assert.Equal(t, 5, book.TotalCopies)
			assert.Equal(t, 5, book.AvailableCopies)
		}).
		Return(nil).
		Once()

	book, err := bookService.CreateBook(req)

	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, uint(1), book.ID)
	mockRepo.AssertExpectations(t)
}

func TestBookService_CreateBook_DuplicateISBN(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	req := dto.CreateBookRequest{
		ISBN:        "9781234567897",
		Title:       "Test Book",
		Author:      "Test Author",
		TotalCopies: 5,
	}

	existingBook := &models.Book{
		ID:    1,
		ISBN:  req.ISBN,
		Title: "Existing Book",
	}

	mockRepo.On("FindByISBN", req.ISBN).Return(existingBook, nil).Once()

	book, err := bookService.CreateBook(req)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Equal(t, "book with this ISBN already exists", err.Error())
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestBookService_GetBookByID(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	expectedBook := &models.Book{
		ID:     1,
		ISBN:   "9781234567897",
		Title:  "Test Book",
		Author: "Test Author",
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedBook, nil).Once()

	book, err := bookService.GetBookByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedBook, book)
	mockRepo.AssertExpectations(t)
}

func TestBookService_GetBookByID_NotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	mockRepo.On("FindByID", uint(999)).
		Return((*models.Book)(nil), errors.New("record not found")).
		Once()

	book, err := bookService.GetBookByID(999)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "book not found")
	mockRepo.AssertExpectations(t)
}

func TestBookService_UpdateBook(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	existingBook := &models.Book{
		ID:              1,
		ISBN:            "9781234567897",
		Title:           "Old Title",
		Author:          "Old Author",
		TotalCopies:     5,
		AvailableCopies: 3,
	}

	req := dto.UpdateBookRequest{
		Title:       "New Title",
		Author:      "New Author",
		TotalCopies: 10,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingBook, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.Book")).
		Run(func(args mock.Arguments) {
			book := args.Get(0).(*models.Book)
			assert.Equal(t, "New Title", book.Title)
			assert.Equal(t, "New Author", book.Author)
			assert.Equal(t, 10, book.TotalCopies)
			assert.Equal(t, 8, book.AvailableCopies)
		}).
		Return(nil).
		Once()

	book, err := bookService.UpdateBook(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, book)
	mockRepo.AssertExpectations(t)
}

func TestBookService_UpdateBook_RejectsTotalCopiesBelowBorrowedCopies(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	existingBook := &models.Book{
		ID:              1,
		ISBN:            "9781234567897",
		Title:           "Distributed Systems",
		Author:          "Author",
		TotalCopies:     5,
		AvailableCopies: 1,
	}

	req := dto.UpdateBookRequest{TotalCopies: 3}

	mockRepo.On("FindByID", uint(1)).Return(existingBook, nil).Once()

	book, err := bookService.UpdateBook(1, req)

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Equal(t, "total copies cannot be less than borrowed copies", err.Error())
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestBookService_UpdateBook_RejectsInconsistentExistingStock(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	existingBook := &models.Book{
		ID:              1,
		ISBN:            "9781234567897",
		Title:           "Distributed Systems",
		Author:          "Author",
		TotalCopies:     2,
		AvailableCopies: 3,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingBook, nil).Once()

	book, err := bookService.UpdateBook(1, dto.UpdateBookRequest{Title: "New Title"})

	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Equal(t, "book stock is inconsistent", err.Error())
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestBookService_DeleteBook(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	existingBook := &models.Book{
		ID:              1,
		Title:           "Test Book",
		TotalCopies:     5,
		AvailableCopies: 5,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingBook, nil).Once()
	mockRepo.On("Delete", uint(1)).Return(nil).Once()

	err := bookService.DeleteBook(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookService_DeleteBook_WithActiveBorrows(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	existingBook := &models.Book{
		ID:              1,
		Title:           "Test Book",
		TotalCopies:     5,
		AvailableCopies: 3,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingBook, nil).Once()

	err := bookService.DeleteBook(1)

	assert.Error(t, err)
	assert.Equal(t, "cannot delete book with active borrows", err.Error())
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
}

func TestBookService_ListBooks(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	expectedBooks := []models.Book{
		{ID: 1, Title: "Book 1"},
		{ID: 2, Title: "Book 2"},
	}

	mockRepo.On("List", 1, 10, "test", "created_at_desc").
		Return(expectedBooks, int64(2), nil).
		Once()

	books, total, err := bookService.ListBooks(1, 10, "test", "created_at_desc")

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, books, 2)
	mockRepo.AssertExpectations(t)
}

func TestBookService_CheckAvailability(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	book := &models.Book{
		ID:              1,
		Title:           "Test Book",
		TotalCopies:     5,
		AvailableCopies: 3,
	}

	mockRepo.On("FindByID", uint(1)).Return(book, nil).Once()

	available, err := bookService.CheckAvailability(1)

	assert.NoError(t, err)
	assert.True(t, available)
	mockRepo.AssertExpectations(t)
}

func TestBookService_CheckAvailability_NotAvailable(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	book := &models.Book{
		ID:              1,
		Title:           "Test Book",
		TotalCopies:     5,
		AvailableCopies: 0,
	}

	mockRepo.On("FindByID", uint(1)).Return(book, nil).Once()

	available, err := bookService.CheckAvailability(1)

	assert.NoError(t, err)
	assert.False(t, available)
	mockRepo.AssertExpectations(t)
}

func TestBookService_CheckAvailability_RejectsInconsistentStock(t *testing.T) {
	mockRepo := new(MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	book := &models.Book{
		ID:              1,
		Title:           "Broken Book",
		TotalCopies:     2,
		AvailableCopies: 3,
	}

	mockRepo.On("FindByID", uint(1)).Return(book, nil).Once()

	available, err := bookService.CheckAvailability(1)

	assert.Error(t, err)
	assert.False(t, available)
	assert.Equal(t, "book stock is inconsistent", err.Error())
	mockRepo.AssertExpectations(t)
}
