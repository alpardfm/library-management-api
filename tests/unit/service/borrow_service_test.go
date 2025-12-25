// // tests/unit/service/borrow_service_test.go
package service_test

// import (
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"

// 	"library-management-api/internal/dto"
// 	"library-management-api/internal/models"
// 	"library-management-api/internal/service"
// )

// // Mock Repositories
// type MockBorrowRepository struct {
// 	mock.Mock
// }

// func (m *MockBorrowRepository) Create(record *models.BorrowRecord) error {
// 	args := m.Called(record)
// 	return args.Error(0)
// }

// func (m *MockBorrowRepository) FindByID(id uint) (*models.BorrowRecord, error) {
// 	args := m.Called(id)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*models.BorrowRecord), args.Error(1)
// }

// func (m *MockBorrowRepository) FindActiveByUserAndBook(userID, bookID uint) (*models.BorrowRecord, error) {
// 	args := m.Called(userID, bookID)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*models.BorrowRecord), args.Error(1)
// }

// func (m *MockBorrowRepository) Update(record *models.BorrowRecord) error {
// 	args := m.Called(record)
// 	return args.Error(0)
// }

// func (m *MockBorrowRepository) ListByUser(userID uint, page, limit int) ([]models.BorrowRecord, int64, error) {
// 	args := m.Called(userID, page, limit)
// 	return args.Get(0).([]models.BorrowRecord), args.Get(1).(int64), args.Error(2)
// }

// func (m *MockBorrowRepository) ListActive(page, limit int) ([]models.BorrowRecord, int64, error) {
// 	args := m.Called(page, limit)
// 	return args.Get(0).([]models.BorrowRecord), args.Get(1).(int64), args.Error(2)
// }

// func (m *MockBorrowRepository) ListOverdue(page, limit int) ([]models.BorrowRecord, int64, error) {
// 	args := m.Called(page, limit)
// 	return args.Get(0).([]models.BorrowRecord), args.Get(1).(int64), args.Error(2)
// }

// func (m *MockBorrowRepository) CountActiveByUser(userID uint) (int64, error) {
// 	args := m.Called(userID)
// 	return args.Get(0).(int64), args.Error(1)
// }

// type MockUserRepository struct {
// 	mock.Mock
// }

// func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
// 	args := m.Called(id)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*models.User), args.Error(1)
// }

// func TestBorrowService_BorrowBook_Success(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	userID := uint(1)
// 	req := dto.BorrowBookRequest{
// 		BookID: 1,
// 	}

// 	// Mock user
// 	user := &models.User{
// 		ID:       userID,
// 		Username: "testuser",
// 		IsActive: true,
// 	}
// 	mockUserRepo.On("FindByID", userID).Return(user, nil).Once()

// 	// Mock book
// 	book := &models.Book{
// 		ID:              1,
// 		Title:           "Test Book",
// 		TotalCopies:     5,
// 		AvailableCopies: 3,
// 	}
// 	mockBookRepo.On("FindByID", uint(1)).Return(book, nil).Once()

// 	// Mock active count
// 	mockBorrowRepo.On("CountActiveByUser", userID).Return(int64(2), nil).Once()

// 	// Mock check if already borrowed
// 	mockBorrowRepo.On("FindActiveByUserAndBook", userID, uint(1)).
// 		Return((*models.BorrowRecord)(nil), errors.New("not found")).
// 		Once()

// 	// Mock book update
// 	mockBookRepo.On("Update", mock.AnythingOfType("*models.Book")).
// 		Run(func(args mock.Arguments) {
// 			book := args.Get(0).(*models.Book)
// 			assert.Equal(t, 2, book.AvailableCopies) // Reduced by 1
// 		}).
// 		Return(nil).
// 		Once()

// 	// Mock borrow record creation
// 	mockBorrowRepo.On("Create", mock.AnythingOfType("*models.BorrowRecord")).
// 		Run(func(args mock.Arguments) {
// 			record := args.Get(0).(*models.BorrowRecord)
// 			assert.Equal(t, userID, record.UserID)
// 			assert.Equal(t, uint(1), record.BookID)
// 			assert.Equal(t, models.StatusBorrowed, record.Status)
// 		}).
// 		Return(nil).
// 		Once()

// 	borrowRecord, err := borrowService.BorrowBook(userID, req)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, borrowRecord)
// 	mockUserRepo.AssertExpectations(t)
// 	mockBookRepo.AssertExpectations(t)
// 	mockBorrowRepo.AssertExpectations(t)
// }

// func TestBorrowService_BorrowBook_UserNotFound(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	userID := uint(999)
// 	req := dto.BorrowBookRequest{
// 		BookID: 1,
// 	}

// 	mockUserRepo.On("FindByID", userID).
// 		Return((*models.User)(nil), errors.New("user not found")).
// 		Once()

// 	borrowRecord, err := borrowService.BorrowBook(userID, req)

// 	assert.Error(t, err)
// 	assert.Nil(t, borrowRecord)
// 	assert.Contains(t, err.Error(), "user not found")

// 	mockUserRepo.AssertExpectations(t)
// 	mockBookRepo.AssertNotCalled(t, "FindByID")
// 	mockBorrowRepo.AssertNotCalled(t, "CountActiveByUser")
// }

// func TestBorrowService_BorrowBook_UserInactive(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	userID := uint(1)
// 	req := dto.BorrowBookRequest{
// 		BookID: 1,
// 	}

// 	user := &models.User{
// 		ID:       userID,
// 		Username: "testuser",
// 		IsActive: false,
// 	}
// 	mockUserRepo.On("FindByID", userID).Return(user, nil).Once()

// 	borrowRecord, err := borrowService.BorrowBook(userID, req)

// 	assert.Error(t, err)
// 	assert.Nil(t, borrowRecord)
// 	assert.Equal(t, "user account is deactivated", err.Error())

// 	mockUserRepo.AssertExpectations(t)
// 	mockBookRepo.AssertNotCalled(t, "FindByID")
// }

// func TestBorrowService_BorrowBook_BookNotAvailable(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	userID := uint(1)
// 	req := dto.BorrowBookRequest{
// 		BookID: 1,
// 	}

// 	user := &models.User{
// 		ID:       userID,
// 		Username: "testuser",
// 		IsActive: true,
// 	}
// 	mockUserRepo.On("FindByID", userID).Return(user, nil).Once()

// 	book := &models.Book{
// 		ID:              1,
// 		Title:           "Test Book",
// 		TotalCopies:     5,
// 		AvailableCopies: 0, // No copies available
// 	}
// 	mockBookRepo.On("FindByID", uint(1)).Return(book, nil).Once()

// 	borrowRecord, err := borrowService.BorrowBook(userID, req)

// 	assert.Error(t, err)
// 	assert.Nil(t, borrowRecord)
// 	assert.Equal(t, "book is not available for borrowing", err.Error())

// 	mockUserRepo.AssertExpectations(t)
// 	mockBookRepo.AssertExpectations(t)
// 	mockBorrowRepo.AssertNotCalled(t, "CountActiveByUser")
// }

// func TestBorrowService_BorrowBook_ExceedLimit(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	userID := uint(1)
// 	req := dto.BorrowBookRequest{
// 		BookID: 1,
// 	}

// 	user := &models.User{
// 		ID:       userID,
// 		Username: "testuser",
// 		IsActive: true,
// 	}
// 	mockUserRepo.On("FindByID", userID).Return(user, nil).Once()

// 	book := &models.Book{
// 		ID:              1,
// 		Title:           "Test Book",
// 		TotalCopies:     5,
// 		AvailableCopies: 3,
// 	}
// 	mockBookRepo.On("FindByID", uint(1)).Return(book, nil).Once()

// 	// User already has 5 active borrows (max limit)
// 	mockBorrowRepo.On("CountActiveByUser", userID).Return(int64(5), nil).Once()

// 	borrowRecord, err := borrowService.BorrowBook(userID, req)

// 	assert.Error(t, err)
// 	assert.Nil(t, borrowRecord)
// 	assert.Contains(t, err.Error(), "maximum borrow limit")

// 	mockUserRepo.AssertExpectations(t)
// 	mockBookRepo.AssertExpectations(t)
// 	mockBorrowRepo.AssertExpectations(t)
// 	mockBookRepo.AssertNotCalled(t, "Update")
// }

// func TestBorrowService_ReturnBook_Success(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	userID := uint(1)
// 	req := dto.ReturnBookRequest{
// 		BorrowRecordID: 1,
// 	}

// 	borrowDate := time.Now().Add(-10 * 24 * time.Hour)
// 	dueDate := time.Now().Add(-3 * 24 * time.Hour) // 3 days overdue

// 	borrowRecord := &models.BorrowRecord{
// 		ID:         1,
// 		UserID:     userID,
// 		BookID:     1,
// 		BorrowDate: borrowDate,
// 		DueDate:    dueDate,
// 		ReturnDate: nil,
// 		Status:     models.StatusBorrowed,
// 	}

// 	mockBorrowRepo.On("FindByID", uint(1)).Return(borrowRecord, nil).Once()

// 	book := &models.Book{
// 		ID:              1,
// 		Title:           "Test Book",
// 		TotalCopies:     5,
// 		AvailableCopies: 2,
// 	}
// 	mockBookRepo.On("FindByID", uint(1)).Return(book, nil).Once()

// 	// Mock book update (increase available copies)
// 	mockBookRepo.On("Update", mock.AnythingOfType("*models.Book")).
// 		Run(func(args mock.Arguments) {
// 			book := args.Get(0).(*models.Book)
// 			assert.Equal(t, 3, book.AvailableCopies) // Increased by 1
// 		}).
// 		Return(nil).
// 		Once()

// 	// Mock borrow record update
// 	mockBorrowRepo.On("Update", mock.AnythingOfType("*models.BorrowRecord")).
// 		Run(func(args mock.Arguments) {
// 			record := args.Get(0).(*models.BorrowRecord)
// 			assert.NotNil(t, record.ReturnDate)
// 			assert.Equal(t, models.StatusReturned, record.Status)
// 		}).
// 		Return(nil).
// 		Once()

// 	returnedRecord, fine, err := borrowService.ReturnBook(userID, req)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, returnedRecord)
// 	assert.Equal(t, 3000, fine) // 3 days overdue * 1000 per day
// 	mockBorrowRepo.AssertExpectations(t)
// 	mockBookRepo.AssertExpectations(t)
// }

// func TestBorrowService_ReturnBook_NotOwner(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	userID := uint(2) // Different user
// 	req := dto.ReturnBookRequest{
// 		BorrowRecordID: 1,
// 	}

// 	borrowRecord := &models.BorrowRecord{
// 		ID:         1,
// 		UserID:     1, // Owned by user 1
// 		BookID:     1,
// 		ReturnDate: nil,
// 	}

// 	mockBorrowRepo.On("FindByID", uint(1)).Return(borrowRecord, nil).Once()

// 	returnedRecord, fine, err := borrowService.ReturnBook(userID, req)

// 	assert.Error(t, err)
// 	assert.Nil(t, returnedRecord)
// 	assert.Equal(t, 0, fine)
// 	assert.Equal(t, "not authorized to return this book", err.Error())

// 	mockBorrowRepo.AssertExpectations(t)
// 	mockBookRepo.AssertNotCalled(t, "FindByID")
// }

// func TestBorrowService_ReturnBook_AlreadyReturned(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	userID := uint(1)
// 	req := dto.ReturnBookRequest{
// 		BorrowRecordID: 1,
// 	}

// 	returnDate := time.Now().Add(-1 * time.Hour)
// 	borrowRecord := &models.BorrowRecord{
// 		ID:         1,
// 		UserID:     userID,
// 		BookID:     1,
// 		ReturnDate: &returnDate,
// 		Status:     models.StatusReturned,
// 	}

// 	mockBorrowRepo.On("FindByID", uint(1)).Return(borrowRecord, nil).Once()

// 	returnedRecord, fine, err := borrowService.ReturnBook(userID, req)

// 	assert.Error(t, err)
// 	assert.Nil(t, returnedRecord)
// 	assert.Equal(t, 0, fine)
// 	assert.Equal(t, "book already returned", err.Error())

// 	mockBorrowRepo.AssertExpectations(t)
// 	mockBookRepo.AssertNotCalled(t, "FindByID")
// }

// func TestBorrowService_CalculateFine(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	borrowDate := time.Now().Add(-20 * 24 * time.Hour)
// 	dueDate := time.Now().Add(-10 * 24 * time.Hour) // 10 days overdue

// 	borrowRecord := &models.BorrowRecord{
// 		ID:         1,
// 		UserID:     1,
// 		BookID:     1,
// 		BorrowDate: borrowDate,
// 		DueDate:    dueDate,
// 		ReturnDate: nil,
// 		Status:     models.StatusOverdue,
// 	}

// 	mockBorrowRepo.On("FindByID", uint(1)).Return(borrowRecord, nil).Once()

// 	fine, err := borrowService.CalculateFine(1)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 10000, fine) // 10 days * 1000
// 	mockBorrowRepo.AssertExpectations(t)
// }

// func TestBorrowService_GetUserBorrows(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	userID := uint(1)
// 	expectedRecords := []models.BorrowRecord{
// 		{ID: 1, UserID: userID},
// 		{ID: 2, UserID: userID},
// 	}

// 	mockBorrowRepo.On("ListByUser", userID, 1, 10).
// 		Return(expectedRecords, int64(2), nil).
// 		Once()

// 	records, total, err := borrowService.GetUserBorrows(userID, 1, 10)

// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(2), total)
// 	assert.Len(t, records, 2)
// 	mockBorrowRepo.AssertExpectations(t)
// }

// func TestBorrowService_GetActiveBorrows(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	expectedRecords := []models.BorrowRecord{
// 		{ID: 1, Status: models.StatusBorrowed},
// 		{ID: 2, Status: models.StatusBorrowed},
// 	}

// 	mockBorrowRepo.On("ListActive", 1, 10).
// 		Return(expectedRecords, int64(2), nil).
// 		Once()

// 	records, total, err := borrowService.GetActiveBorrows(1, 10)

// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(2), total)
// 	assert.Len(t, records, 2)
// 	mockBorrowRepo.AssertExpectations(t)
// }

// func TestBorrowService_GetOverdueBorrows(t *testing.T) {
// 	mockBorrowRepo := new(MockBorrowRepository)
// 	mockBookRepo := new(MockBookRepository)
// 	mockUserRepo := new(MockUserRepository)

// 	borrowService := service.NewBorrowService(mockBorrowRepo, mockBookRepo, mockUserRepo)

// 	expectedRecords := []models.BorrowRecord{
// 		{ID: 1, Status: models.StatusOverdue},
// 		{ID: 2, Status: models.StatusOverdue},
// 	}

// 	mockBorrowRepo.On("ListOverdue", 1, 10).
// 		Return(expectedRecords, int64(2), nil).
// 		Once()

// 	records, total, err := borrowService.GetOverdueBorrows(1, 10)

// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(2), total)
// 	assert.Len(t, records, 2)
// 	mockBorrowRepo.AssertExpectations(t)
// }
