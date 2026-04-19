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
	Config               *config.Config
	DB                   interface{}
	AuthService          ports.AuthService
	UserService          ports.UserService
	BranchService        ports.BranchService
	SupplierService      ports.SupplierService
	ProductService       ports.ProductService
	StockService         ports.StockService
	PurchaseBatchService ports.PurchaseBatchService
	TransferService      ports.TransferService
}

// NewRouter creates and configures the Gin router with all routes and middleware.
func NewRouter(cfg *config.Config, db interface{}) *gin.Engine {
	// This is a simplified version - in production, you'd wire all dependencies
	return setupRouter(RouterConfig{
		Config: cfg,
		DB:     db,
	})
}

// NewRouterWithServices creates a router with all services wired up.
func NewRouterWithServices(config RouterConfig) *gin.Engine {
	return setupRouter(config)
}

func setupRouter(config RouterConfig) *gin.Engine {
	cfg := config.Config
	db := config.DB
	authService := config.AuthService
	userService := config.UserService
	branchService := config.BranchService
	supplierService := config.SupplierService
	productService := config.ProductService
	stockService := config.StockService
	purchaseBatchService := config.PurchaseBatchService
	transferService := config.TransferService

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

	var productHandler *handlers.ProductHandler
	if productService != nil {
		productHandler = handlers.NewProductHandler(productService)
	}

	var stockHandler *handlers.StockHandler
	if stockService != nil {
		stockHandler = handlers.NewStockHandler(stockService)
	}

	var purchaseBatchHandler *handlers.PurchaseBatchHandler
	if purchaseBatchService != nil {
		purchaseBatchHandler = handlers.NewPurchaseBatchHandler(purchaseBatchService)
	}

	var transferHandler *handlers.TransferHandler
	if transferService != nil {
		transferHandler = handlers.NewTransferHandler(transferService)
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

			// Product routes (accessible to all authenticated users)
			if productHandler != nil {
				products := authorized.Group("/products")
				products.GET("", productHandler.List)
				products.GET("/:id", productHandler.GetByID)
				products.POST("", productHandler.Create)
				products.PUT("/:id", productHandler.Update)
				products.DELETE("/:id", productHandler.Deactivate)
			}

			// Stock routes (accessible to all authenticated users)
			if stockHandler != nil {
				stock := authorized.Group("/stock")
				stock.GET("", stockHandler.GetByProductAndLocation)
				stock.GET("/product/:productId", stockHandler.ListByProduct)
				stock.GET("/location/:locationType/:locationId", stockHandler.ListByLocation)
				stock.PUT("/:id/min-alert", stockHandler.SetMinStockAlert)
			}

			// Purchase batch routes (accessible to all authenticated users)
			if purchaseBatchHandler != nil {
				batches := authorized.Group("/purchase-batches")
				batches.GET("", purchaseBatchHandler.List)
				batches.GET("/:id", purchaseBatchHandler.GetByID)
				batches.POST("", purchaseBatchHandler.Create)
				batches.POST("/:id/receive", purchaseBatchHandler.Receive)
			}

			// Transfer routes (accessible to all authenticated users)
			if transferHandler != nil {
				transfers := authorized.Group("/transfers")
				transfers.GET("", transferHandler.List)
				transfers.GET("/:id", transferHandler.GetByID)
				transfers.POST("", transferHandler.Create)
				transfers.POST("/:id/approve", transferHandler.Approve)
				transfers.POST("/:id/reject", transferHandler.Reject)
				transfers.POST("/:id/send", transferHandler.MarkAsSent)
				transfers.POST("/:id/receive", transferHandler.MarkAsReceived)
				transfers.DELETE("/:id", transferHandler.Cancel)
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
