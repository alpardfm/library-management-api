// tests/unit/handler/auth_handler_test.go
package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alpardfm/library-management-api/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/handler"
	"github.com/alpardfm/library-management-api/internal/models"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req dto.RegisterRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockAuthService) GenerateToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Claims), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockService)

	router := gin.New()
	router.POST("/register", authHandler.Register)

	reqBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     models.RoleMember,
	}

	// Test Case 1: Success
	t.Run("Success", func(t *testing.T) {
		mockService.On("Register", reqBody).Return(expectedUser, nil).Once()

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "User registered successfully", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	// Test Case 2: Invalid request
	t.Run("Invalid Request", func(t *testing.T) {
		invalidReq := map[string]interface{}{
			"username": "test",
			"email":    "invalid-email",
			// password missing
		}

		jsonBody, _ := json.Marshal(invalidReq)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertNotCalled(t, "Register")
	})

	// Test Case 3: Service error
	t.Run("Service Error", func(t *testing.T) {
		mockService.On("Register", reqBody).
			Return((*models.User)(nil), assert.AnError).
			Once()

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertExpectations(t)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthService)
	authHandler := handler.NewAuthHandler(mockService)

	router := gin.New()
	router.POST("/login", authHandler.Login)

	reqBody := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	expectedResponse := &dto.LoginResponse{
		Token: "jwt-token-123",
		User: struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "member",
		},
	}

	// Test Case 1: Success
	t.Run("Success", func(t *testing.T) {
		mockService.On("Login", reqBody).Return(expectedResponse, nil).Once()

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "jwt-token-123", response["token"])
		assert.NotNil(t, response["user"])

		mockService.AssertExpectations(t)
	})

	// Test Case 2: Invalid credentials
	t.Run("Invalid Credentials", func(t *testing.T) {
		mockService.On("Login", reqBody).
			Return((*dto.LoginResponse)(nil), errors.New("invalid credentials")).
			Once()

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.NotNil(t, response["error"])

		mockService.AssertExpectations(t)
	})
}
