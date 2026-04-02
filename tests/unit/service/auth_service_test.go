package service_test

import (
	"testing"
	"time"

	"github.com/alpardfm/library-management-api/internal/models"
	"github.com/alpardfm/library-management-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_GenerateToken_UsesConfiguredJWTExpiry(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockUserRepo, "test-secret", 2*time.Hour)

	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     models.RoleMember,
	}

	token, err := authService.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := authService.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	remaining := time.Until(claims.ExpiresAt.Time)
	assert.Greater(t, remaining, 90*time.Minute)
	assert.LessOrEqual(t, remaining, 2*time.Hour)
	mockUserRepo.AssertExpectations(t)
}
