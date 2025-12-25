// tests/unit/repository/user_repository_test.go
package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"library-management-api/internal/models"
	"library-management-api/internal/repository"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewUserRepository(gormDB)

	user := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Role:         models.RoleMember,
		IsActive:     true,
	}

	// Mock INSERT query
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(
			user.Username,
			user.Email,
			user.PasswordHash,
			user.Role,
			user.IsActive,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // id
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.Create(user)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewUserRepository(gormDB)

	// Mock SELECT query
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "role", "is_active", "created_at", "updated_at"}).
		AddRow(1, "testuser", "test@example.com", "hashed_password", "member", true, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewUserRepository(gormDB)

	// Mock empty result
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1`).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByID(999)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewUserRepository(gormDB)

	rows := sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(1, "john_doe", "john@example.com")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1`).
		WithArgs("john_doe").
		WillReturnRows(rows)

	user, err := repo.FindByUsername("john_doe")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john_doe", user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewUserRepository(gormDB)

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "users"`).
		WillReturnRows(countRows)

	// Mock SELECT with pagination
	rows := sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(1, "user1", "user1@example.com").
		AddRow(2, "user2", "user2@example.com")

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(rows)

	users, total, err := repo.List(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1", users[0].Username)
	assert.Equal(t, "user2", users[1].Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}
