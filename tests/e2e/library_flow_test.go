// tests/e2e/library_flow_test.go
package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"library-management-api/internal/dto"
)

type LibraryE2ETestSuite struct {
	suite.Suite
	baseURL string
	token   string
	userID  uint
	bookID  uint
}

func (suite *LibraryE2ETestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8080/api/v1"

	// Wait for server to be ready
	suite.waitForServer()
}

func (suite *LibraryE2ETestSuite) waitForServer() {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(suite.baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		time.Sleep(1 * time.Second)
	}
	suite.T().Fatalf("Server not ready after %d seconds", maxRetries)
}

func (suite *LibraryE2ETestSuite) TestCompleteLibraryFlow() {
	suite.T().Run("1. Register User", func(t *testing.T) {
		reqBody := dto.RegisterRequest{
			Username: "e2euser",
			Email:    "e2e@example.com",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)
		resp, err := http.Post(suite.baseURL+"/auth/register", "application/json", bytes.NewBuffer(jsonBody))

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()

		assert.Equal(t, "User registered successfully", response["message"])
	})

	suite.T().Run("2. Login", func(t *testing.T) {
		reqBody := dto.LoginRequest{
			Username: "e2euser",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)
		resp, err := http.Post(suite.baseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonBody))

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()

		assert.NotNil(t, response["token"])
		suite.token = response["token"].(string)

		userData := response["user"].(map[string]interface{})
		suite.userID = uint(userData["id"].(float64))
	})

	suite.T().Run("3. List Books", func(t *testing.T) {
		req, _ := http.NewRequest("GET", suite.baseURL+"/books", nil)
		req.Header.Set("Authorization", "Bearer "+suite.token)

		client := &http.Client{}
		resp, err := client.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()

		assert.NotNil(t, response["data"])
	})

	// Note: For admin operations, we would need admin token
	// This shows the complete flow structure
}

func TestLibraryE2ETestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}
	suite.Run(t, new(LibraryE2ETestSuite))
}
