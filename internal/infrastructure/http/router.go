package http

import (
	"github.com/Luis-Lanza/luson/internal/config"
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/handlers"
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the Gin router with all routes and middleware.
func NewRouter(cfg *config.Config, db interface{}) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(middleware.Logging())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Health check endpoint
	healthHandler := handlers.NewHealthHandler(db)
	router.GET("/api/health", healthHandler.Check)

	// API routes group
	api := router.Group("/api")
	api.Use(middleware.Auth())
	{
		// Auth routes (placeholder)
		api.POST("/auth/login", func(c *gin.Context) {
			c.JSON(404, gin.H{"error": "not implemented"})
		})
		api.POST("/auth/refresh", func(c *gin.Context) {
			c.JSON(404, gin.H{"error": "not implemented"})
		})
	}

	return router
}
