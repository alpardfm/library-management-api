package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/handler"
	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) CreateBook(req dto.CreateBookRequest) (*models.Book, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) GetBookByID(id uint) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) UpdateBook(id uint, req dto.UpdateBookRequest) (*models.Book, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) DeleteBook(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookService) ListBooks(page, limit int, search, sort string) ([]models.Book, int64, error) {
	args := m.Called(page, limit, search, sort)
	return args.Get(0).([]models.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookService) CheckAvailability(id uint) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

type MockBorrowService struct {
	mock.Mock
}

func (m *MockBorrowService) BorrowBook(userID uint, req dto.BorrowBookRequest) (*models.BorrowRecord, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BorrowRecord), args.Error(1)
}

func (m *MockBorrowService) ReturnBook(userID uint, role string, req dto.ReturnBookRequest) (*models.BorrowRecord, int, error) {
	args := m.Called(userID, role, req)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).(*models.BorrowRecord), args.Int(1), args.Error(2)
}

func (m *MockBorrowService) GetUserBorrows(userID uint, page, limit int, sort string) ([]models.BorrowRecord, int64, error) {
	args := m.Called(userID, page, limit, sort)
	return args.Get(0).([]models.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowService) GetActiveBorrows(page, limit int, sort string) ([]models.BorrowRecord, int64, error) {
	args := m.Called(page, limit, sort)
	return args.Get(0).([]models.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowService) GetOverdueBorrows(page, limit int, sort string) ([]models.BorrowRecord, int64, error) {
	args := m.Called(page, limit, sort)
	return args.Get(0).([]models.BorrowRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockBorrowService) CalculateFine(borrowID uint) (int, error) {
	args := m.Called(borrowID)
	return args.Int(0), args.Error(1)
}

func TestBookHandler_ListBooks_UsesParsedQueryAndMeta(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockBookService)
	bookHandler := handler.NewBookHandler(mockService)

	router := gin.New()
	router.GET("/books", bookHandler.ListBooks)

	expectedBooks := []models.Book{
		{ID: 1, Title: "Clean Code"},
		{ID: 2, Title: "Domain-Driven Design"},
	}

	mockService.On("ListBooks", 2, 100, "golang", "title_asc").
		Return(expectedBooks, int64(201), nil).
		Once()

	req := httptest.NewRequest("GET", "/books?page=2&limit=999&search=golang&sort=title_asc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])

	meta := response["meta"].(map[string]any)
	assert.Equal(t, float64(2), meta["page"])
	assert.Equal(t, float64(100), meta["limit"])
	assert.Equal(t, float64(201), meta["total"])
	assert.Equal(t, float64(3), meta["total_pages"])
	assert.Equal(t, "title_asc", meta["sort"])
	assert.Equal(t, "golang", meta["search"])

	mockService.AssertExpectations(t)
}

func TestBookHandler_ListBooks_InvalidSortReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockBookService)
	bookHandler := handler.NewBookHandler(mockService)

	router := gin.New()
	router.GET("/books", bookHandler.ListBooks)

	req := httptest.NewRequest("GET", "/books?sort=drop_table", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "invalid sort query", response["message"])

	mockService.AssertNotCalled(t, "ListBooks")
}

func TestBorrowHandler_GetMyBorrows_InvalidPageReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockBorrowService)
	borrowHandler := handler.NewBorrowHandler(mockService)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uint(7))
		c.Set("claims", &auth.Claims{UserID: 7})
		c.Next()
	})
	router.GET("/borrow/my-books", borrowHandler.GetMyBorrows)

	req := httptest.NewRequest("GET", "/borrow/my-books?page=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "invalid page query", response["message"])

	mockService.AssertNotCalled(t, "GetUserBorrows")
}
