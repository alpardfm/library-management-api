// internal/middleware/logger.go
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log after request is processed
		latencyMs := time.Since(start).Milliseconds()
		requestID := c.GetString(RequestIDKey)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		event := log.Info()
		if c.Writer.Status() >= 500 {
			event = log.Error()
		}

		event.
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Int64("latency_ms", latencyMs).
			Fields(requestActorFields(c)).
			Msg("request completed")
	}
}

func requestActorFields(c *gin.Context) map[string]interface{} {
	fields := map[string]interface{}{}

	if userID, exists := c.Get("user_id"); exists {
		fields["user_id"] = userID
	}

	if role, exists := c.Get("role"); exists {
		if roleStr, ok := role.(string); ok {
			fields["role"] = roleStr
		}
	}

	return fields
}
