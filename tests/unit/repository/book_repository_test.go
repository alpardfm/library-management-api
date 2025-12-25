// tests/unit/repository/book_repository_test.go
package repository_test

import (
	"testing"

	"github.com/alpardfm/library-management-api/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/alpardfm/library-management-api/internal/repository"
)

func TestBookRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewBookRepository(gormDB)

	book := &models.Book{
		ISBN:            "9781234567897",
		Title:           "Test Book",
		Author:          "Test Author",
		TotalCopies:     5,
		AvailableCopies: 5,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "books"`).
		WithArgs(
			book.ISBN,
			book.Title,
			book.Author,
			book.Publisher,
			book.PublicationYear,
			book.Genre,
			book.Description,
			book.TotalCopies,
			book.AvailableCopies,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // id
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.Create(book)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), book.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_UpdateAvailableCopies(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewBookRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "books" SET "available_copies" = available_copies \+ \$1 WHERE id = \$2`).
		WithArgs(-1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.UpdateAvailableCopies(1, -1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_ListWithSearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewBookRepository(gormDB)

	// Mock count with search
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "books" WHERE \(title ILIKE \$1 OR author ILIKE \$2 OR isbn ILIKE \$3\)`).
		WithArgs("%test%", "%test%", "%test%").
		WillReturnRows(countRows)

	// Mock SELECT with search and pagination
	rows := sqlmock.NewRows([]string{"id", "isbn", "title", "author"}).
		AddRow(1, "9781234567897", "Test Book", "Test Author")

	mock.ExpectQuery(`SELECT \* FROM "books" WHERE \(title ILIKE \$1 OR author ILIKE \$2 OR isbn ILIKE \$3\) ORDER BY created_at DESC`).
		WithArgs("%test%", "%test%", "%test%").
		WillReturnRows(rows)

	books, total, err := repo.List(1, 10, "test")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, books, 1)
	assert.Equal(t, "Test Book", books[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}
