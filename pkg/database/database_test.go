package database

import (
	"errors"
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
	migrations := allPostgresMigrations()
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

func TestApplyPostgresMigrations_ExecutesBaseStatements(t *testing.T) {
	gormDB, mock := newMockPostgresDB(t)
	t.Setenv("ENABLE_PG_TRGM", "false")

	for _, migration := range basePostgresMigrations() {
		mock.ExpectExec(regexp.QuoteMeta(normalizeSQL(migration.statement))).
			WillReturnResult(sqlmock.NewResult(0, 0))
	}

	err := applyPostgresMigrations(gormDB)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplyPostgresMigrations_ExecutesTrigramStatementsWhenEnabled(t *testing.T) {
	gormDB, mock := newMockPostgresDB(t)
	t.Setenv("ENABLE_PG_TRGM", "true")

	for _, migration := range basePostgresMigrations() {
		mock.ExpectExec(regexp.QuoteMeta(normalizeSQL(migration.statement))).
			WillReturnResult(sqlmock.NewResult(0, 0))
	}

	mock.ExpectExec(regexp.QuoteMeta(normalizeSQL(pgTrgmExtensionMigration().statement))).
		WillReturnResult(sqlmock.NewResult(0, 0))

	for _, migration := range pgTrgmIndexMigrations() {
		mock.ExpectExec(regexp.QuoteMeta(normalizeSQL(migration.statement))).
			WillReturnResult(sqlmock.NewResult(0, 0))
	}

	err := applyPostgresMigrations(gormDB)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplyPostgresMigrations_GracefullySkipsTrigramIndexesWhenExtensionFails(t *testing.T) {
	gormDB, mock := newMockPostgresDB(t)
	t.Setenv("ENABLE_PG_TRGM", "true")

	for _, migration := range basePostgresMigrations() {
		mock.ExpectExec(regexp.QuoteMeta(normalizeSQL(migration.statement))).
			WillReturnResult(sqlmock.NewResult(0, 0))
	}

	mock.ExpectExec(regexp.QuoteMeta(normalizeSQL(pgTrgmExtensionMigration().statement))).
		WillReturnError(errors.New("permission denied to create extension"))

	err := applyPostgresMigrations(gormDB)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
