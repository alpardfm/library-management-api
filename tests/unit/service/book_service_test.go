// // tests/unit/service/book_service_test.go
package service_test

// import (
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"

// 	"library-management-api/internal/dto"
// 	"library-management-api/internal/models"
// 	"library-management-api/internal/service"
// )

// // Mock Book Repository
// type MockBookRepository struct {
// 	mock.Mock
// }

// func (m *MockBookRepository) Create(book *models.Book) error {
// 	args := m.Called(book)
// 	return args.Error(0)
// }

// func (m *MockBookRepository) FindByID(id uint) (*models.Book, error) {
// 	args := m.Called(id)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*models.Book), args.Error(1)
// }

// func (m *MockBookRepository) FindByISBN(isbn string) (*models.Book, error) {
// 	args := m.Called(isbn)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*models.Book), args.Error(1)
// }

// func (m *MockBookRepository) Update(book *models.Book) error {
// 	args := m.Called(book)
// 	return args.Error(0)
// }

// func (m *MockBookRepository) Delete(id uint) error {
// 	args := m.Called(id)
// 	return args.Error(0)
// }

// func (m *MockBookRepository) List(page, limit int, search string) ([]models.Book, int64, error) {
// 	args := m.Called(page, limit, search)
// 	return args.Get(0).([]models.Book), args.Get(1).(int64), args.Error(2)
// }

// func (m *MockBookRepository) UpdateAvailableCopies(id uint, change int) error {
// 	args := m.Called(id, change)
// 	return args.Error(0)
// }

// func TestBookService_CreateBook(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	req := dto.CreateBookRequest{
// 		ISBN:        "9781234567897",
// 		Title:       "Test Book",
// 		Author:      "Test Author",
// 		TotalCopies: 5,
// 	}

// 	// Mock ISBN check
// 	mockRepo.On("FindByISBN", "9781234567897").
// 		Return((*models.Book)(nil), errors.New("not found")).
// 		Once()

// 	// Mock create
// 	mockRepo.On("Create", mock.AnythingOfType("*models.Book")).
// 		Run(func(args mock.Arguments) {
// 			book := args.Get(0).(*models.Book)
// 			book.ID = 1
// 			assert.Equal(t, "9781234567897", book.ISBN)
// 			assert.Equal(t, "Test Book", book.Title)
// 			assert.Equal(t, 5, book.TotalCopies)
// 			assert.Equal(t, 5, book.AvailableCopies)
// 		}).
// 		Return(nil).
// 		Once()

// 	book, err := bookService.CreateBook(req)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, book)
// 	assert.Equal(t, uint(1), book.ID)
// 	mockRepo.AssertExpectations(t)
// }

// func TestBookService_CreateBook_DuplicateISBN(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	req := dto.CreateBookRequest{
// 		ISBN:        "9781234567897",
// 		Title:       "Test Book",
// 		Author:      "Test Author",
// 		TotalCopies: 5,
// 	}

// 	existingBook := &models.Book{
// 		ID:    1,
// 		ISBN:  "9781234567897",
// 		Title: "Existing Book",
// 	}

// 	mockRepo.On("FindByISBN", "9781234567897").
// 		Return(existingBook, nil).
// 		Once()

// 	book, err := bookService.CreateBook(req)

// 	assert.Error(t, err)
// 	assert.Nil(t, book)
// 	assert.Equal(t, "book with this ISBN already exists", err.Error())
// 	mockRepo.AssertExpectations(t)
// }

// func TestBookService_GetBookByID(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	expectedBook := &models.Book{
// 		ID:     1,
// 		ISBN:   "9781234567897",
// 		Title:  "Test Book",
// 		Author: "Test Author",
// 	}

// 	mockRepo.On("FindByID", uint(1)).
// 		Return(expectedBook, nil).
// 		Once()

// 	book, err := bookService.GetBookByID(1)

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedBook, book)
// 	mockRepo.AssertExpectations(t)
// }

// func TestBookService_GetBookByID_NotFound(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	mockRepo.On("FindByID", uint(999)).
// 		Return((*models.Book)(nil), errors.New("record not found")).
// 		Once()

// 	book, err := bookService.GetBookByID(999)

// 	assert.Error(t, err)
// 	assert.Nil(t, book)
// 	assert.Contains(t, err.Error(), "book not found")
// 	mockRepo.AssertExpectations(t)
// }

// func TestBookService_UpdateBook(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	existingBook := &models.Book{
// 		ID:              1,
// 		ISBN:            "9781234567897",
// 		Title:           "Old Title",
// 		Author:          "Old Author",
// 		TotalCopies:     5,
// 		AvailableCopies: 3,
// 	}

// 	req := dto.UpdateBookRequest{
// 		Title:       "New Title",
// 		Author:      "New Author",
// 		TotalCopies: 10,
// 	}

// 	mockRepo.On("FindByID", uint(1)).
// 		Return(existingBook, nil).
// 		Once()

// 	mockRepo.On("Update", mock.AnythingOfType("*models.Book")).
// 		Run(func(args mock.Arguments) {
// 			book := args.Get(0).(*models.Book)
// 			assert.Equal(t, "New Title", book.Title)
// 			assert.Equal(t, "New Author", book.Author)
// 			assert.Equal(t, 10, book.TotalCopies)
// 			assert.Equal(t, 8, book.AvailableCopies) // 3 + (10-5) = 8
// 		}).
// 		Return(nil).
// 		Once()

// 	book, err := bookService.UpdateBook(1, req)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, book)
// 	mockRepo.AssertExpectations(t)
// }

// func TestBookService_DeleteBook(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	existingBook := &models.Book{
// 		ID:              1,
// 		Title:           "Test Book",
// 		TotalCopies:     5,
// 		AvailableCopies: 5, // All copies available
// 	}

// 	mockRepo.On("FindByID", uint(1)).
// 		Return(existingBook, nil).
// 		Once()

// 	mockRepo.On("Delete", uint(1)).
// 		Return(nil).
// 		Once()

// 	err := bookService.DeleteBook(1)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestBookService_DeleteBook_WithActiveBorrows(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	existingBook := &models.Book{
// 		ID:              1,
// 		Title:           "Test Book",
// 		TotalCopies:     5,
// 		AvailableCopies: 3, // 2 copies borrowed
// 	}

// 	mockRepo.On("FindByID", uint(1)).
// 		Return(existingBook, nil).
// 		Once()

// 	err := bookService.DeleteBook(1)

// 	assert.Error(t, err)
// 	assert.Equal(t, "cannot delete book with active borrows", err.Error())
// 	mockRepo.AssertExpectations(t)
// 	mockRepo.AssertNotCalled(t, "Delete")
// }

// func TestBookService_ListBooks(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	expectedBooks := []models.Book{
// 		{ID: 1, Title: "Book 1"},
// 		{ID: 2, Title: "Book 2"},
// 	}

// 	mockRepo.On("List", 1, 10, "test").
// 		Return(expectedBooks, int64(2), nil).
// 		Once()

// 	books, total, err := bookService.ListBooks(1, 10, "test")

// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(2), total)
// 	assert.Len(t, books, 2)
// 	mockRepo.AssertExpectations(t)
// }

// func TestBookService_ListBooks_DefaultPagination(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	expectedBooks := []models.Book{
// 		{ID: 1, Title: "Book 1"},
// 	}

// 	mockRepo.On("List", 1, 10, "").
// 		Return(expectedBooks, int64(1), nil).
// 		Once()

// 	books, total, err := bookService.ListBooks(0, 0, "") // Invalid values

// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(1), total)
// 	assert.Len(t, books, 1)
// 	mockRepo.AssertExpectations(t)
// }

// func TestBookService_CheckAvailability(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	book := &models.Book{
// 		ID:              1,
// 		Title:           "Test Book",
// 		TotalCopies:     5,
// 		AvailableCopies: 3,
// 	}

// 	mockRepo.On("FindByID", uint(1)).
// 		Return(book, nil).
// 		Once()

// 	available, err := bookService.CheckAvailability(1)

// 	assert.NoError(t, err)
// 	assert.True(t, available)
// 	mockRepo.AssertExpectations(t)
// }

// func TestBookService_CheckAvailability_NotAvailable(t *testing.T) {
// 	mockRepo := new(MockBookRepository)
// 	bookService := service.NewBookService(mockRepo)

// 	book := &models.Book{
// 		ID:              1,
// 		Title:           "Test Book",
// 		TotalCopies:     5,
// 		AvailableCopies: 0,
// 	}

// 	mockRepo.On("FindByID", uint(1)).
// 		Return(book, nil).
// 		Once()

// 	available, err := bookService.CheckAvailability(1)

// 	assert.NoError(t, err)
// 	assert.False(t, available)
// 	mockRepo.AssertExpectations(t)
// }
