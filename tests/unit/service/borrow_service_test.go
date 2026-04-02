package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/internal/repository"
	"github.com/alpardfm/library-management-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) WithTx(tx *gorm.DB) repository.BookRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.BookRepository)
}

func (m *MockBookRepository) Create(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) FindByID(id uint) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepository) FindByIDForUpdate(id uint) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepository) FindByISBN(isbn string) (*models.Book, error) {
	args := m.Called(isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepository) Update(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookRepository) List(page, limit int, search, sort string) ([]models.Book, int64, error) {
	args := m.Called(page, limit, search, sort)
	return args.Get(0).([]models.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRepository) UpdateAvailableCopies(id uint, change int) error {
	args := m.Called(id, change)
	return args.Error(0)
}

type MockBorrowRepository struct {
	mock.Mock
}

func (m *MockBorrowRepository) WithTx(tx *gorm.DB) repository.BorrowRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.BorrowRepository)
}

func (m *MockBorrowRepository) Create(record *models.BorrowRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockBorrowRepository) FindByID(id uint) (*models.BorrowRecord, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BorrowRecord), args.Error(1)
}

func (m *MockBorrowRepository) FindByIDForUpdate(id uint) (*models.BorrowRecord, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BorrowRecord), args.Error(1)
}

func (m *MockBorrowRepository) FindActiveByUserAndBook(userID, bookID uint) (*models.BorrowRecord, error) {
	args := m.Called(userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BorrowRecord), args.Error(1)
}

func (m *MockBorrowRepository) Update(record *models.BorrowRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockBorrowRepository) ListByUser(userID uint, page, limit int, sort string) ([]models.BorrowRecord, int64, error) {
	args := m.Called(userID, page, limit, sort)
	return args.Get(0).([]models.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRepository) ListActive(page, limit int, sort string) ([]models.BorrowRecord, int64, error) {
	args := m.Called(page, limit, sort)
	return args.Get(0).([]models.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRepository) ListOverdue(page, limit int, sort string) ([]models.BorrowRecord, int64, error) {
	args := m.Called(page, limit, sort)
	return args.Get(0).([]models.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowRepository) CountActiveByUser(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) WithTx(tx *gorm.DB) repository.UserRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.UserRepository)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByIDForUpdate(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List(page, limit int) ([]models.User, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mockDB, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mockDB
}

func newBorrowService(t *testing.T) (*MockBorrowRepository, *MockBookRepository, *MockUserRepository, sqlmock.Sqlmock, service.BorrowService) {
	t.Helper()

	mockBorrowRepo := new(MockBorrowRepository)
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	gormDB, mockDB := newMockDB(t)

	svc := service.NewBorrowService(gormDB, mockBorrowRepo, mockBookRepo, mockUserRepo, service.BorrowServiceConfig{
		MaxBooksPerUser: 5,
		BorrowDays:      7,
		FinePerDay:      1000,
	})

	return mockBorrowRepo, mockBookRepo, mockUserRepo, mockDB, svc
}

func TestBorrowService_BorrowBook_Success(t *testing.T) {
	mockBorrowRepo, mockBookRepo, mockUserRepo, sqlMock, borrowService := newBorrowService(t)

	userID := uint(1)
	req := dto.BorrowBookRequest{BookID: 1}

	user := &models.User{ID: userID, IsActive: true}
	book := &models.Book{ID: 1, TotalCopies: 5, AvailableCopies: 3}

	sqlMock.ExpectBegin()
	mockUserRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockUserRepo).Once()
	mockBookRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookRepo).Once()
	mockBorrowRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBorrowRepo).Once()
	mockUserRepo.On("FindByIDForUpdate", userID).Return(user, nil).Once()
	mockBorrowRepo.On("CountActiveByUser", userID).Return(int64(1), nil).Once()
	mockBookRepo.On("FindByIDForUpdate", uint(1)).Return(book, nil).Once()
	mockBorrowRepo.On("FindActiveByUserAndBook", userID, uint(1)).Return((*models.BorrowRecord)(nil), gorm.ErrRecordNotFound).Once()
	mockBookRepo.On("Update", mock.AnythingOfType("*models.Book")).
		Run(func(args mock.Arguments) {
			updatedBook := args.Get(0).(*models.Book)
			assert.Equal(t, 2, updatedBook.AvailableCopies)
		}).
		Return(nil).
		Once()
	mockBorrowRepo.On("Create", mock.AnythingOfType("*models.BorrowRecord")).
		Run(func(args mock.Arguments) {
			record := args.Get(0).(*models.BorrowRecord)
			assert.Equal(t, userID, record.UserID)
			assert.Equal(t, uint(1), record.BookID)
			assert.False(t, record.DueDate.IsZero())
		}).
		Return(nil).
		Once()
	sqlMock.ExpectCommit()

	borrowRecord, err := borrowService.BorrowBook(userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, borrowRecord)
	assert.Equal(t, uint(1), borrowRecord.BookID)
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
	mockBorrowRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBorrowService_ReturnBook_Success(t *testing.T) {
	mockBorrowRepo, mockBookRepo, _, sqlMock, borrowService := newBorrowService(t)

	userID := uint(1)
	req := dto.ReturnBookRequest{BorrowRecordID: 1}
	now := time.Now()
	borrowRecord := &models.BorrowRecord{
		ID:         1,
		UserID:     userID,
		BookID:     1,
		BorrowDate: now.Add(-10 * 24 * time.Hour),
		DueDate:    now.Add(-(72*time.Hour + time.Minute)),
		Status:     models.StatusBorrowed,
	}
	book := &models.Book{
		ID:              1,
		TotalCopies:     5,
		AvailableCopies: 2,
	}

	sqlMock.ExpectBegin()
	mockBookRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookRepo).Once()
	mockBorrowRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBorrowRepo).Once()
	mockBorrowRepo.On("FindByIDForUpdate", uint(1)).Return(borrowRecord, nil).Once()
	mockBookRepo.On("FindByIDForUpdate", uint(1)).Return(book, nil).Once()
	mockBookRepo.On("Update", mock.AnythingOfType("*models.Book")).
		Run(func(args mock.Arguments) {
			updatedBook := args.Get(0).(*models.Book)
			assert.Equal(t, 3, updatedBook.AvailableCopies)
		}).
		Return(nil).
		Once()
	mockBorrowRepo.On("Update", mock.AnythingOfType("*models.BorrowRecord")).
		Run(func(args mock.Arguments) {
			updatedRecord := args.Get(0).(*models.BorrowRecord)
			assert.NotNil(t, updatedRecord.ReturnDate)
			assert.Equal(t, models.StatusReturned, updatedRecord.Status)
		}).
		Return(nil).
		Once()
	sqlMock.ExpectCommit()

	returnedRecord, fine, err := borrowService.ReturnBook(userID, "member", req)

	assert.NoError(t, err)
	assert.NotNil(t, returnedRecord)
	assert.Equal(t, 3000, fine)
	mockBookRepo.AssertExpectations(t)
	mockBorrowRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBorrowService_ReturnBook_NotOwner_RollsBack(t *testing.T) {
	mockBorrowRepo, mockBookRepo, _, sqlMock, borrowService := newBorrowService(t)

	req := dto.ReturnBookRequest{BorrowRecordID: 1}
	borrowRecord := &models.BorrowRecord{
		ID:     1,
		UserID: 99,
		BookID: 1,
		Status: models.StatusBorrowed,
	}

	sqlMock.ExpectBegin()
	mockBookRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookRepo).Once()
	mockBorrowRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBorrowRepo).Once()
	mockBorrowRepo.On("FindByIDForUpdate", uint(1)).Return(borrowRecord, nil).Once()
	sqlMock.ExpectRollback()

	returnedRecord, fine, err := borrowService.ReturnBook(1, "member", req)

	assert.Error(t, err)
	assert.Nil(t, returnedRecord)
	assert.Equal(t, 0, fine)
	assert.Equal(t, "not authorized to return this book", err.Error())
	mockBookRepo.AssertExpectations(t)
	mockBorrowRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBorrowService_ReturnBook_AdminCanReturnOtherUserBook(t *testing.T) {
	mockBorrowRepo, mockBookRepo, _, sqlMock, borrowService := newBorrowService(t)

	req := dto.ReturnBookRequest{BorrowRecordID: 1}
	now := time.Now()
	borrowRecord := &models.BorrowRecord{
		ID:         1,
		UserID:     99,
		BookID:     1,
		BorrowDate: now.Add(-2 * 24 * time.Hour),
		DueDate:    now.Add(5 * 24 * time.Hour),
		Status:     models.StatusBorrowed,
	}
	book := &models.Book{ID: 1, TotalCopies: 5, AvailableCopies: 2}

	sqlMock.ExpectBegin()
	mockBookRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookRepo).Once()
	mockBorrowRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBorrowRepo).Once()
	mockBorrowRepo.On("FindByIDForUpdate", uint(1)).Return(borrowRecord, nil).Once()
	mockBookRepo.On("FindByIDForUpdate", uint(1)).Return(book, nil).Once()
	mockBookRepo.On("Update", mock.AnythingOfType("*models.Book")).Return(nil).Once()
	mockBorrowRepo.On("Update", mock.AnythingOfType("*models.BorrowRecord")).Return(nil).Once()
	sqlMock.ExpectCommit()

	returnedRecord, fine, err := borrowService.ReturnBook(1, "admin", req)

	assert.NoError(t, err)
	assert.NotNil(t, returnedRecord)
	assert.Equal(t, 0, fine)
	mockBookRepo.AssertExpectations(t)
	mockBorrowRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBorrowService_ReturnBook_StockAlreadyFull_RollsBack(t *testing.T) {
	mockBorrowRepo, mockBookRepo, _, sqlMock, borrowService := newBorrowService(t)

	req := dto.ReturnBookRequest{BorrowRecordID: 1}
	borrowRecord := &models.BorrowRecord{
		ID:     1,
		UserID: 1,
		BookID: 1,
		Status: models.StatusBorrowed,
	}
	book := &models.Book{
		ID:              1,
		TotalCopies:     3,
		AvailableCopies: 3,
	}

	sqlMock.ExpectBegin()
	mockBookRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookRepo).Once()
	mockBorrowRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBorrowRepo).Once()
	mockBorrowRepo.On("FindByIDForUpdate", uint(1)).Return(borrowRecord, nil).Once()
	mockBookRepo.On("FindByIDForUpdate", uint(1)).Return(book, nil).Once()
	sqlMock.ExpectRollback()

	returnedRecord, fine, err := borrowService.ReturnBook(1, "member", req)

	assert.Error(t, err)
	assert.Nil(t, returnedRecord)
	assert.Equal(t, 0, fine)
	assert.Equal(t, "book stock is already full, cannot process return", err.Error())
	mockBookRepo.AssertExpectations(t)
	mockBorrowRepo.AssertExpectations(t)
	mockBookRepo.AssertNotCalled(t, "Update", mock.Anything)
	mockBorrowRepo.AssertNotCalled(t, "Update", mock.Anything)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBorrowService_BorrowBook_TransactionError_RollsBack(t *testing.T) {
	mockBorrowRepo, mockBookRepo, mockUserRepo, sqlMock, borrowService := newBorrowService(t)

	userID := uint(1)
	req := dto.BorrowBookRequest{BookID: 1}
	user := &models.User{ID: userID, IsActive: true}
	book := &models.Book{ID: 1, TotalCopies: 5, AvailableCopies: 1}

	sqlMock.ExpectBegin()
	mockUserRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockUserRepo).Once()
	mockBookRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookRepo).Once()
	mockBorrowRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBorrowRepo).Once()
	mockUserRepo.On("FindByIDForUpdate", userID).Return(user, nil).Once()
	mockBorrowRepo.On("CountActiveByUser", userID).Return(int64(0), nil).Once()
	mockBookRepo.On("FindByIDForUpdate", uint(1)).Return(book, nil).Once()
	mockBorrowRepo.On("FindActiveByUserAndBook", userID, uint(1)).Return((*models.BorrowRecord)(nil), gorm.ErrRecordNotFound).Once()
	mockBookRepo.On("Update", mock.AnythingOfType("*models.Book")).Return(nil).Once()
	mockBorrowRepo.On("Create", mock.AnythingOfType("*models.BorrowRecord")).Return(errors.New("insert failed")).Once()
	sqlMock.ExpectRollback()

	borrowRecord, err := borrowService.BorrowBook(userID, req)

	assert.Error(t, err)
	assert.Nil(t, borrowRecord)
	assert.Equal(t, "failed to create borrow record", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
	mockBorrowRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBorrowService_BorrowBook_DuplicateActiveBorrow_RollsBack(t *testing.T) {
	mockBorrowRepo, mockBookRepo, mockUserRepo, sqlMock, borrowService := newBorrowService(t)

	userID := uint(1)
	req := dto.BorrowBookRequest{BookID: 1}
	user := &models.User{ID: userID, IsActive: true}
	book := &models.Book{ID: 1, TotalCopies: 5, AvailableCopies: 2}
	existingBorrow := &models.BorrowRecord{ID: 99, UserID: userID, BookID: req.BookID, Status: models.StatusBorrowed}

	sqlMock.ExpectBegin()
	mockUserRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockUserRepo).Once()
	mockBookRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookRepo).Once()
	mockBorrowRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBorrowRepo).Once()
	mockUserRepo.On("FindByIDForUpdate", userID).Return(user, nil).Once()
	mockBorrowRepo.On("CountActiveByUser", userID).Return(int64(1), nil).Once()
	mockBookRepo.On("FindByIDForUpdate", req.BookID).Return(book, nil).Once()
	mockBorrowRepo.On("FindActiveByUserAndBook", userID, req.BookID).Return(existingBorrow, nil).Once()
	sqlMock.ExpectRollback()

	borrowRecord, err := borrowService.BorrowBook(userID, req)

	assert.Error(t, err)
	assert.Nil(t, borrowRecord)
	assert.Equal(t, "user has already borrowed this book", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
	mockBorrowRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBorrowService_BorrowBook_MaxLimitReached_RollsBack(t *testing.T) {
	mockBorrowRepo, mockBookRepo, mockUserRepo, sqlMock, borrowService := newBorrowService(t)

	userID := uint(1)
	req := dto.BorrowBookRequest{BookID: 1}
	user := &models.User{ID: userID, IsActive: true}

	sqlMock.ExpectBegin()
	mockUserRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockUserRepo).Once()
	mockBookRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookRepo).Once()
	mockBorrowRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBorrowRepo).Once()
	mockUserRepo.On("FindByIDForUpdate", userID).Return(user, nil).Once()
	mockBorrowRepo.On("CountActiveByUser", userID).Return(int64(5), nil).Once()
	sqlMock.ExpectRollback()

	borrowRecord, err := borrowService.BorrowBook(userID, req)

	assert.Error(t, err)
	assert.Nil(t, borrowRecord)
	assert.Contains(t, err.Error(), "maximum borrow limit")
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
	mockBorrowRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBorrowService_BorrowBook_BookUnavailable_RollsBack(t *testing.T) {
	mockBorrowRepo, mockBookRepo, mockUserRepo, sqlMock, borrowService := newBorrowService(t)

	userID := uint(1)
	req := dto.BorrowBookRequest{BookID: 1}
	user := &models.User{ID: userID, IsActive: true}
	book := &models.Book{ID: 1, TotalCopies: 1, AvailableCopies: 0}

	sqlMock.ExpectBegin()
	mockUserRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockUserRepo).Once()
	mockBookRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookRepo).Once()
	mockBorrowRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBorrowRepo).Once()
	mockUserRepo.On("FindByIDForUpdate", userID).Return(user, nil).Once()
	mockBorrowRepo.On("CountActiveByUser", userID).Return(int64(0), nil).Once()
	mockBookRepo.On("FindByIDForUpdate", req.BookID).Return(book, nil).Once()
	sqlMock.ExpectRollback()

	borrowRecord, err := borrowService.BorrowBook(userID, req)

	assert.Error(t, err)
	assert.Nil(t, borrowRecord)
	assert.Equal(t, "book is not available for borrowing", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
	mockBorrowRepo.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
