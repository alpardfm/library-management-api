package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alpardfm/library-management-api/configs"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupSystemRoutesTest(t *testing.T) (*gin.Engine, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	require.NoError(t, err)

	router := gin.New()
	registerSystemRoutes(router, gormDB, &configs.Config{
		AppName:    "Library Management API",
		AppVersion: "1.0.0",
		AppEnv:     "test",
	})

	return router, mock
}

func TestSystemRoutes_Health(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _ := setupSystemRoutesTest(t)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "healthy", body["status"])
	assert.Equal(t, "Library Management API", body["app"])
}

func TestSystemRoutes_Ready_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock := setupSystemRoutesTest(t)
	mock.ExpectPing()

	req := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ready", body["status"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSystemRoutes_Ready_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock := setupSystemRoutesTest(t)
	mock.ExpectPing().WillReturnError(assert.AnError)

	req := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var body map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "not_ready", body["status"])
	assert.NoError(t, mock.ExpectationsWereMet())
}
