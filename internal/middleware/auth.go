// internal/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/alpardfm/library-management-api/pkg/apperror"
	"github.com/alpardfm/library-management-api/pkg/auth"
	httpresponse "github.com/alpardfm/library-management-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httpresponse.Error(c, apperror.Unauthorized("authorization header required"))
			c.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			httpresponse.Error(c, apperror.Unauthorized("invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := auth.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			if err == auth.ErrExpiredToken {
				httpresponse.Error(c, apperror.Unauthorized("token has expired"))
			} else {
				httpresponse.Error(c, apperror.Unauthorized("invalid token"))
			}
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user role not found"})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role type"})
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}
