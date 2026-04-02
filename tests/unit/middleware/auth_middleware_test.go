// tests/unit/middleware/auth_middleware_test.go
package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alpardfm/library-management-api/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/alpardfm/library-management-api/internal/middleware"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// Test middleware
	router.Use(middleware.AuthMiddleware("test-secret"))

	router.GET("/protected", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if exists {
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "no auth"})
		}
	})

	t.Run("No Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Bearer Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Expired Token", func(t *testing.T) {
		token, err := auth.GenerateToken(1, "testuser", "member", "test-secret", -1)
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, false, response["success"])
		assert.Equal(t, "token has expired", response["message"])
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, false, response["success"])
		assert.Equal(t, "invalid token", response["message"])
	})
}

func TestRoleMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// Add auth context first
	router.Use(func(c *gin.Context) {
		c.Set("role", "member")
		c.Next()
	})

	// Test admin-only middleware
	router.Use(middleware.RoleMiddleware("admin"))

	router.GET("/admin-only", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "admin access"})
	})

	t.Run("Insufficient Permissions", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin-only", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Has Required Role", func(t *testing.T) {
		// Create new router with admin role
		adminRouter := gin.New()
		adminRouter.Use(func(c *gin.Context) {
			c.Set("role", "admin")
			c.Next()
		})
		adminRouter.Use(middleware.RoleMiddleware("admin"))

		adminRouter.GET("/admin-only", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "admin access"})
		})

		req := httptest.NewRequest("GET", "/admin-only", nil)
		w := httptest.NewRecorder()

		adminRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
