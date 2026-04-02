package database

import (
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockPostgresDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestPostgresMigrations_AreIdempotent(t *testing.T) {
	migrations := postgresMigrations()
	require.NotEmpty(t, migrations)

	for _, migration := range migrations {
		sql := normalizeSQL(migration.statement)

		assert.NotEmpty(t, migration.name)
		assert.NotEmpty(t, sql)

		if strings.HasPrefix(sql, "DO $$") {
			assert.Contains(t, sql, "IF NOT EXISTS", migration.name)
			continue
		}

		assert.Contains(t, sql, "IF NOT EXISTS", migration.name)
	}
}

func TestApplyPostgresMigrations_ExecutesAllStatements(t *testing.T) {
	gormDB, mock := newMockPostgresDB(t)

	for _, migration := range postgresMigrations() {
		mock.ExpectExec(regexp.QuoteMeta(normalizeSQL(migration.statement))).
			WillReturnResult(sqlmock.NewResult(0, 0))
	}

	err := applyPostgresMigrations(gormDB)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
