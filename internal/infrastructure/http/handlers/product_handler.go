package handlers

import (
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/dto"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProductHandler handles product management HTTP requests.
type ProductHandler struct {
	productService ports.ProductService
}

// NewProductHandler creates a new product handler.
func NewProductHandler(productService ports.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// List handles GET /api/products.
func (h *ProductHandler) List(c *gin.Context) {
	var query dto.ListQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	filter := ports.ProductFilter{
		Limit:  query.Limit,
		Offset: query.Offset,
		Active: query.Active,
	}

	// Optional filters from query params
	if productType := c.Query("product_type"); productType != "" {
		filter.ProductType = &productType
	}
	if brand := c.Query("brand"); brand != "" {
		filter.Brand = &brand
	}
	if vehicleType := c.Query("vehicle_type"); vehicleType != "" {
		filter.VehicleType = &vehicleType
	}
	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	products, err := h.productService.List(c.Request.Context(), filter)
	if err != nil {
		dto.InternalError(c, err.Error())
		return
	}

	total := len(products)
	dto.SuccessWithMeta(c, products, &dto.Meta{
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	})
}

// GetByID handles GET /api/products/:id.
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid product ID")
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		dto.NotFound(c, "Product not found")
		return
	}

	dto.Success(c, product)
}

// Create handles POST /api/products.
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		dto.Unauthorized(c, "User not authenticated")
		return
	}

	serviceReq := ports.CreateProductRequest{
		Name:         req.Name,
		Description:  req.Description,
		ProductType:  req.ProductType,
		Brand:        req.Brand,
		Model:        req.Model,
		Voltage:      req.Voltage,
		Amperage:     req.Amperage,
		BatteryType:  req.BatteryType,
		Polarity:     req.Polarity,
		AcidLiters:   req.AcidLiters,
		VehicleType:  req.VehicleType,
		MinSalePrice: req.MinSalePrice,
		CreatedBy:    userID.(uuid.UUID),
	}

	product, err := h.productService.Create(c.Request.Context(), serviceReq)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Created(c, product)
}

// Update handles PUT /api/products/:id.
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid product ID")
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	serviceReq := ports.UpdateProductRequest{
		Name:         req.Name,
		Description:  req.Description,
		Brand:        req.Brand,
		Model:        req.Model,
		Voltage:      req.Voltage,
		Amperage:     req.Amperage,
		BatteryType:  req.BatteryType,
		Polarity:     req.Polarity,
		AcidLiters:   req.AcidLiters,
		VehicleType:  req.VehicleType,
		MinSalePrice: req.MinSalePrice,
		Active:       req.Active,
	}

	product, err := h.productService.Update(c.Request.Context(), id, serviceReq)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, product)
}

// Deactivate handles DELETE /api/products/:id.
func (h *ProductHandler) Deactivate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid product ID")
		return
	}

	if err := h.productService.Deactivate(c.Request.Context(), id); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Product deactivated successfully"})
}
