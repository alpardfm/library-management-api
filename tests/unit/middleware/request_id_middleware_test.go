package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alpardfm/library-management-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware_GeneratesHeaderAndStoresContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"request_id": c.GetString(middleware.RequestIDKey),
		})
	})

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	headerID := w.Header().Get(middleware.RequestIDHeader)
	assert.NotEmpty(t, headerID)

	var body map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, headerID, body["request_id"])
}

func TestRequestIDMiddleware_UsesIncomingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())
	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/ping", nil)
	req.Header.Set(middleware.RequestIDHeader, "req-123")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "req-123", w.Header().Get(middleware.RequestIDHeader))
}

func TestRequestIDMiddleware_RegeneratesOverlongIncomingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())
	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/ping", nil)
	req.Header.Set(middleware.RequestIDHeader, string(bytes.Repeat([]byte("a"), 256)))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	requestID := w.Header().Get(middleware.RequestIDHeader)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, requestID)
	assert.NotEqual(t, string(bytes.Repeat([]byte("a"), 256)), requestID)
	assert.LessOrEqual(t, len(requestID), 128)
}

func TestLoggerMiddleware_IncludesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	originalLogger := log.Logger
	log.Logger = zerolog.New(&buf)
	defer func() {
		log.Logger = originalLogger
	}()

	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/ping", nil)
	req.Header.Set(middleware.RequestIDHeader, "trace-abc")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	output := buf.String()
	assert.Contains(t, output, `"request_id":"trace-abc"`)
	assert.Contains(t, output, `"path":"/ping"`)
	assert.Contains(t, output, `"status":200`)
	assert.Contains(t, output, `"latency_ms":`)
}

func TestLoggerMiddleware_IncludesActorFieldsWithoutSensitiveData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	originalLogger := log.Logger
	log.Logger = zerolog.New(&buf)
	defer func() {
		log.Logger = originalLogger
	}()

	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uint(7))
		c.Set("role", "member")
		c.Next()
	})
	router.Use(middleware.LoggerMiddleware())
	router.GET("/secure", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	output := buf.String()
	assert.Contains(t, output, `"request_id":"`)
	assert.Contains(t, output, `"user_id":7`)
	assert.Contains(t, output, `"role":"member"`)
	assert.NotContains(t, output, "secret-token")
	assert.NotContains(t, output, "Authorization")
}
