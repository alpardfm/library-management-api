// tests/integration/auth_integration_test.go
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alpardfm/library-management-api/pkg/database"

	"github.com/alpardfm/library-management-api/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAuthIntegration(t *testing.T) {
	// Setup test database
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer cleanupTestDB(db)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := setupTestRouter(db)

	t.Run("Register new user", func(t *testing.T) {
		reqBody := dto.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "User registered successfully", response["message"])
		assert.NotNil(t, response["data"])
	})

	t.Run("Login with registered user", func(t *testing.T) {
		reqBody := dto.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.NotEmpty(t, response["token"])
		assert.NotNil(t, response["user"])
	})
}

func setupTestDB() (*gorm.DB, error) {
	// Use SQLite for testing
	// Or connect to test PostgreSQL instance
	return database.Connect()
}

func cleanupTestDB(db *gorm.DB) {
	// Clean up test data
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM books")
	db.Exec("DELETE FROM borrow_records")
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	// Initialize app with test database
	// Similar to main() but with test config
	return nil // Return configured router
}
