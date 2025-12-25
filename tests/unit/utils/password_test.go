// tests/unit/utils/password_test.go
package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alpardfm/library-management-api/pkg/utils"
)

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"

	hashed, err := utils.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecurePassword123"

	hashed, err := utils.HashPassword(password)
	assert.NoError(t, err)

	// Test correct password
	valid := utils.CheckPasswordHash(password, hashed)
	assert.True(t, valid)

	// Test incorrect password
	invalid := utils.CheckPasswordHash("wrongPassword", hashed)
	assert.False(t, invalid)

	// Test empty password
	emptyValid := utils.CheckPasswordHash("", hashed)
	assert.False(t, emptyValid)
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	hashed, err := utils.HashPassword("")

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	valid := utils.CheckPasswordHash("password", "invalid-hash-format")
	assert.False(t, valid)
}
