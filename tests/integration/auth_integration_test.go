// tests/integration/auth_integration_test.go
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alpardfm/library-management-api/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestAuthIntegration(t *testing.T) {
	db, cleanup := setupIntegrationTestDB(t)
	defer cleanup()

	router := setupIntegrationRouter(db)

	t.Run("Register new user", func(t *testing.T) {
		assert.NoError(t, resetIntegrationTestDB(db))

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
		assert.NoError(t, resetIntegrationTestDB(db))

		registerBody := dto.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		registerJSON, _ := json.Marshal(registerBody)
		registerReq := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(registerJSON))
		registerReq.Header.Set("Content-Type", "application/json")

		registerW := httptest.NewRecorder()
		router.ServeHTTP(registerW, registerReq)
		assert.Equal(t, http.StatusCreated, registerW.Code)

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

		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
		assert.NotNil(t, data["user"])
	})
}
