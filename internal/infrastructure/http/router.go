package http

import (
	"github.com/Luis-Lanza/luson/internal/config"
	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/handlers"
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/middleware"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/gin-gonic/gin"
)

// RouterConfig holds all dependencies needed to configure the router.
type RouterConfig struct {
	Config          *config.Config
	DB              interface{}
	AuthService     ports.AuthService
	UserService     ports.UserService
	BranchService   ports.BranchService
	SupplierService ports.SupplierService
}

// NewRouter creates and configures the Gin router with all routes and middleware.
func NewRouter(cfg *config.Config, db interface{}) *gin.Engine {
	// This is a simplified version - in production, you'd wire all dependencies
	return setupRouter(cfg, db, nil, nil, nil, nil)
}

// NewRouterWithServices creates a router with all services wired up.
func NewRouterWithServices(config RouterConfig) *gin.Engine {
	return setupRouter(config.Config, config.DB, config.AuthService, config.UserService,
		config.BranchService, config.SupplierService)
}

func setupRouter(cfg *config.Config, db interface{}, authService ports.AuthService,
	userService ports.UserService, branchService ports.BranchService,
	supplierService ports.SupplierService) *gin.Engine {

	router := gin.New()

	// Global middleware
	router.Use(middleware.Logging())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Health check endpoint (public)
	healthHandler := handlers.NewHealthHandler(db)
	router.GET("/api/health", healthHandler.Check)

	// Initialize JWT service for auth middleware
	var jwtService middleware.JWTService
	if cfg != nil && cfg.JWTSecret != "" {
		jwtService = middleware.NewJWTService(cfg.JWTSecret)
	}

	// Initialize handlers
	var authHandler *handlers.AuthHandler
	if authService != nil {
		authHandler = handlers.NewAuthHandler(authService)
	}

	var userHandler *handlers.UserHandler
	if userService != nil {
		userHandler = handlers.NewUserHandler(userService)
	}

	var branchHandler *handlers.BranchHandler
	if branchService != nil {
		branchHandler = handlers.NewBranchHandler(branchService)
	}

	var supplierHandler *handlers.SupplierHandler
	if supplierService != nil {
		supplierHandler = handlers.NewSupplierHandler(supplierService)
	}

	// Public routes (no authentication required)
	if authHandler != nil {
		router.POST("/api/auth/login", authHandler.Login)
		router.POST("/api/auth/refresh", authHandler.Refresh)
	}

	// Protected routes (authentication required)
	if jwtService != nil {
		authorized := router.Group("/api")
		authorized.Use(middleware.Auth(jwtService))
		{
			// Auth routes (authenticated users)
			if authHandler != nil {
				authorized.GET("/auth/me", authHandler.Me)
				authorized.POST("/auth/logout", authHandler.Logout)
			}

			// Supplier routes (accessible to all authenticated users)
			if supplierHandler != nil {
				suppliers := authorized.Group("/suppliers")
				suppliers.GET("", supplierHandler.List)
				suppliers.GET("/:id", supplierHandler.GetByID)
				suppliers.POST("", supplierHandler.Create)
				suppliers.PUT("/:id", supplierHandler.Update)
				suppliers.DELETE("/:id", supplierHandler.Deactivate)
			}
		}

		// Admin-only routes
		adminOnly := authorized.Group("")
		adminOnly.Use(middleware.RequireRole(domain.UserRoleAdmin))
		{
			// Branch routes (admin only)
			if branchHandler != nil {
				branches := adminOnly.Group("/branches")
				branches.GET("", branchHandler.List)
				branches.GET("/:id", branchHandler.GetByID)
				branches.POST("", branchHandler.Create)
				branches.PUT("/:id", branchHandler.Update)
				branches.DELETE("/:id", branchHandler.Deactivate)
			}

			// User routes (admin only)
			if userHandler != nil {
				users := adminOnly.Group("/users")
				users.GET("", userHandler.List)
				users.GET("/:id", userHandler.GetByID)
				users.POST("", userHandler.Create)
				users.PUT("/:id", userHandler.Update)
				users.PUT("/:id/password", userHandler.UpdatePassword)
				users.DELETE("/:id", userHandler.Deactivate)
			}
		}
	}

	return router
}
