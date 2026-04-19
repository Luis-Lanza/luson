package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Pinger defines the interface for checking database connectivity.
type Pinger interface {
	Ping(ctx context.Context) error
}

// HealthHandler handles health check requests.
type HealthHandler struct {
	db Pinger
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(db interface{}) *HealthHandler {
	if p, ok := db.(Pinger); ok {
		return &HealthHandler{db: p}
	}
	return &HealthHandler{db: nil}
}

// Check handles GET /api/health requests.
// Returns {"status":"ok"} if database is connected,
// or {"status":"degraded"} if database is unreachable.
func (h *HealthHandler) Check(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "degraded",
			"error":  "database unreachable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
