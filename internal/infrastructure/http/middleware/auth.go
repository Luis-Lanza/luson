package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// Auth returns a middleware that handles authentication.
// Currently a placeholder that passes through all requests.
// TODO: Implement JWT validation.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request for debugging
		slog.Debug("Auth middleware",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)

		// Placeholder: pass through all requests
		// In production, validate JWT token here
		c.Next()
	}
}
