// tests/integration/api_integration_test.go
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/pkg/database"
)

type APIIntegrationTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
	token  string
	userID uint
}

func (suite *APIIntegrationTestSuite) SetupSuite() {
	// Setup test database
	os.Setenv("DB_NAME", "library_test")
	os.Setenv("DB_SSLMODE", "disable")

	var err error
	suite.db, err = database.Connect()
	if err != nil {
		suite.T().Fatalf("Failed to connect to database: %v", err)
	}

	// Clean database
	suite.cleanDatabase()

	// Run migrations
	if err := database.AutoMigrate(suite.db); err != nil {
		suite.T().Fatalf("Failed to migrate database: %v", err)
	}

	// Setup router (simplified version)
	gin.SetMode(gin.TestMode)
	suite.router = suite.setupTestRouter()
}

func (suite *APIIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

func (suite *APIIntegrationTestSuite) SetupTest() {
	suite.cleanDatabase()
	suite.setupTestData()
}

func (suite *APIIntegrationTestSuite) cleanDatabase() {
	suite.db.Exec("DELETE FROM borrow_records")
	suite.db.Exec("DELETE FROM books")
	suite.db.Exec("DELETE FROM users")
}

func (suite *APIIntegrationTestSuite) setupTestData() {
	// Create test user
	user := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW", // "password123"
		Role:         models.RoleMember,
		IsActive:     true,
	}
	suite.db.Create(user)
	suite.userID = user.ID

	// Create test book
	book := &models.Book{
		ISBN:            "9781234567897",
		Title:           "Test Book",
		Author:          "Test Author",
		TotalCopies:     5,
		AvailableCopies: 5,
	}
	suite.db.Create(book)
}

func (suite *APIIntegrationTestSuite) setupTestRouter() *gin.Engine {
	// Simplified router setup for integration tests
	router := gin.New()

	// Mock auth endpoint
	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Simple mock login
		if req.Username == "testuser" && req.Password == "password123" {
			c.JSON(http.StatusOK, gin.H{
				"token": "mock-jwt-token",
				"user": gin.H{
					"id":       suite.userID,
					"username": "testuser",
					"email":    "test@example.com",
					"role":     "member",
				},
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		}
	})

	// Mock books endpoint
	router.GET("/api/v1/books", func(c *gin.Context) {
		var books []models.Book
		suite.db.Find(&books)

		c.JSON(http.StatusOK, gin.H{
			"data": books,
			"meta": gin.H{
				"total": len(books),
			},
		})
	})

	return router
}

func (suite *APIIntegrationTestSuite) TestLogin_Success() {
	reqBody := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotNil(suite.T(), response["token"])
	assert.NotNil(suite.T(), response["user"])

	suite.token = response["token"].(string)
}

func (suite *APIIntegrationTestSuite) TestLogin_InvalidCredentials() {
	reqBody := dto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotNil(suite.T(), response["error"])
}

func (suite *APIIntegrationTestSuite) TestListBooks() {
	req := httptest.NewRequest("GET", "/api/v1/books", nil)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotNil(suite.T(), response["data"])
	books := response["data"].([]interface{})
	assert.Greater(suite.T(), len(books), 0)
}

func TestAPIIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(APIIntegrationTestSuite))
}
