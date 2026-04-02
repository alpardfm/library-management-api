// tests/e2e/library_flow_test.go
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/alpardfm/library-management-api/internal/dto"
)

type LibraryE2ETestSuite struct {
	suite.Suite
	baseURL string
	client  *http.Client
	token   string
	userID  uint
}

func (suite *LibraryE2ETestSuite) SetupSuite() {
	baseURL := os.Getenv("E2E_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	suite.baseURL = baseURL
	suite.client = &http.Client{Timeout: 10 * time.Second}

	// Wait for server to be ready
	suite.waitForServer()
}

func (suite *LibraryE2ETestSuite) waitForServer() {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := suite.client.Get(suite.baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		time.Sleep(1 * time.Second)
	}
	suite.T().Skipf("skipping e2e test: server not ready at %s/health after %d seconds", suite.baseURL, maxRetries)
}

func (suite *LibraryE2ETestSuite) TestCompleteLibraryFlow() {
	suffix := time.Now().UnixNano()
	username := fmt.Sprintf("e2euser-%d", suffix)
	email := fmt.Sprintf("e2e-%d@example.com", suffix)
	apiBaseURL := suite.baseURL + "/api/v1"

	suite.T().Run("1. Register User", func(t *testing.T) {
		reqBody := dto.RegisterRequest{
			Username: username,
			Email:    email,
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)
		resp, err := suite.client.Post(apiBaseURL+"/auth/register", "application/json", bytes.NewBuffer(jsonBody))

		assert.NoError(t, err)
		require.NotNil(t, resp)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))

		assert.Equal(t, true, response["success"])
		assert.Equal(t, "User registered successfully", response["message"])
		assert.NotNil(t, response["data"])
	})

	suite.T().Run("2. Login", func(t *testing.T) {
		reqBody := dto.LoginRequest{
			Username: username,
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(reqBody)
		resp, err := suite.client.Post(apiBaseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonBody))

		assert.NoError(t, err)
		require.NotNil(t, resp)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))

		assert.Equal(t, true, response["success"])
		data := response["data"].(map[string]interface{})
		assert.NotNil(t, data["token"])
		suite.token = data["token"].(string)

		userData := data["user"].(map[string]interface{})
		suite.userID = uint(userData["id"].(float64))
	})

	suite.T().Run("3. List Books", func(t *testing.T) {
		req, _ := http.NewRequest("GET", apiBaseURL+"/books", nil)
		req.Header.Set("Authorization", "Bearer "+suite.token)

		resp, err := suite.client.Do(req)

		assert.NoError(t, err)
		require.NotNil(t, resp)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))

		assert.Equal(t, true, response["success"])
		assert.NotNil(t, response["data"])
		assert.NotNil(t, response["meta"])
		meta := response["meta"].(map[string]interface{})
		assert.NotNil(t, meta["page"])
		assert.NotNil(t, meta["limit"])
		assert.NotNil(t, meta["total"])
		assert.NotNil(t, meta["total_pages"])
	})
}

func TestLibraryE2ETestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}
	suite.Run(t, new(LibraryE2ETestSuite))
}
