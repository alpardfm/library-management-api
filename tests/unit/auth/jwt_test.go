// tests/unit/auth/jwt_test.go
package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"library-management-api/pkg/auth"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key"
	userID := uint(1)
	username := "testuser"
	role := "member"
	expiry := 1 * time.Hour

	// Generate token
	token, err := auth.GenerateToken(userID, username, role, secret, expiry)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	claims, err := auth.ValidateToken(token, secret)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestValidateToken_InvalidToken(t *testing.T) {
	secret := "test-secret-key"

	// Invalid token format
	claims, err := auth.ValidateToken("invalid.token.here", secret)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, auth.ErrInvalidToken, err)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	secret1 := "secret-key-1"
	secret2 := "secret-key-2"

	token, err := auth.GenerateToken(1, "testuser", "member", secret1, time.Hour)
	assert.NoError(t, err)

	claims, err := auth.ValidateToken(token, secret2)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, auth.ErrInvalidToken, err)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	secret := "test-secret"

	// Generate token with past expiry
	token, err := auth.GenerateToken(1, "testuser", "member", secret, -1*time.Hour)
	assert.NoError(t, err)

	claims, err := auth.ValidateToken(token, secret)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, auth.ErrExpiredToken, err)
}
