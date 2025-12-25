// tests/unit/repository/borrow_repository_test.go
package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"library-management-api/internal/models"
	"library-management-api/internal/repository"
)

func TestBorrowRepository_FindActiveByUserAndBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewBorrowRepository(gormDB)

	rows := sqlmock.NewRows([]string{"id", "user_id", "book_id", "status"}).
		AddRow(1, 1, 1, "borrowed")

	mock.ExpectQuery(`SELECT \* FROM "borrow_records" WHERE user_id = \$1 AND book_id = \$2 AND status = \$3`).
		WithArgs(1, 1, "borrowed").
		WillReturnRows(rows)

	record, err := repo.FindActiveByUserAndBook(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, uint(1), record.ID)
	assert.Equal(t, models.StatusBorrowed, record.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBorrowRepository_CountActiveByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewBorrowRepository(gormDB)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "borrow_records" WHERE user_id = \$1 AND status = \$2`).
		WithArgs(1, "borrowed").
		WillReturnRows(countRows)

	count, err := repo.CountActiveByUser(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBorrowRepository_ListOverdue(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewBorrowRepository(gormDB)

	// Mock count
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "borrow_records" WHERE status = \$1 OR \(return_date IS NULL AND due_date < \$2\)`).
		WillReturnRows(countRows)

	// Mock SELECT with preload
	rows := sqlmock.NewRows([]string{"id", "user_id", "book_id", "status"}).
		AddRow(1, 1, 1, "overdue").
		AddRow(2, 2, 2, "borrowed")

	mock.ExpectQuery(`SELECT \* FROM "borrow_records" WHERE status = \$1 OR \(return_date IS NULL AND due_date < \$2\) ORDER BY due_date ASC`).
		WillReturnRows(rows)

	// Mock user preload
	userRows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(1, "user1").
		AddRow(2, "user2")
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" IN \(\$1,\$2\)`).
		WillReturnRows(userRows)

	// Mock book preload
	bookRows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "Book 1").
		AddRow(2, "Book 2")
	mock.ExpectQuery(`SELECT \* FROM "books" WHERE "books"."id" IN \(\$1,\$2\)`).
		WillReturnRows(bookRows)

	records, total, err := repo.ListOverdue(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, records, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
