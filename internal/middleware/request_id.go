package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
	maxRequestIDLen = 128
)

var fallbackRequestIDCounter uint64

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" || len(requestID) > maxRequestIDLen {
			requestID = generateRequestID()
		}

		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set(RequestIDHeader, requestID)

		c.Next()
	}
}

func generateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		counter := atomic.AddUint64(&fallbackRequestIDCounter, 1)
		return fmt.Sprintf("fallback-%d-%d", time.Now().UnixNano(), counter)
	}
	return hex.EncodeToString(bytes)
}
