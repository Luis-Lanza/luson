package handlers

import (
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/dto"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StockHandler handles stock management HTTP requests.
type StockHandler struct {
	stockService ports.StockService
}

// NewStockHandler creates a new stock handler.
func NewStockHandler(stockService ports.StockService) *StockHandler {
	return &StockHandler{
		stockService: stockService,
	}
}

// ListByLocation handles GET /api/stock/location/:locationType/:locationId.
func (h *StockHandler) ListByLocation(c *gin.Context) {
	locationType := c.Param("locationType")
	locationID, err := uuid.Parse(c.Param("locationId"))
	if err != nil {
		dto.BadRequest(c, "Invalid location ID")
		return
	}

	var query struct {
		Limit        int  `form:"limit,default=20" binding:"max=100"`
		Offset       int  `form:"offset,default=0"`
		LowStockOnly bool `form:"low_stock_only,default=false"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	filter := ports.StockFilter{
		Limit:        query.Limit,
		Offset:       query.Offset,
		LowStockOnly: &query.LowStockOnly,
	}

	stock, err := h.stockService.ListByLocation(c.Request.Context(), locationType, locationID, filter)
	if err != nil {
		dto.InternalError(c, err.Error())
		return
	}

	total := len(stock)
	dto.SuccessWithMeta(c, stock, &dto.Meta{
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	})
}

// ListByProduct handles GET /api/stock/product/:productId.
func (h *StockHandler) ListByProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		dto.BadRequest(c, "Invalid product ID")
		return
	}

	stock, err := h.stockService.ListByProduct(c.Request.Context(), productID)
	if err != nil {
		dto.InternalError(c, err.Error())
		return
	}

	dto.Success(c, stock)
}

// GetByID handles GET /api/stock/:id.
func (h *StockHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid stock ID")
		return
	}

	stock, err := h.stockService.GetByID(c.Request.Context(), id)
	if err != nil {
		dto.NotFound(c, "Stock entry not found")
		return
	}

	dto.Success(c, stock)
}

// GetByProductAndLocation handles GET /api/stock.
func (h *StockHandler) GetByProductAndLocation(c *gin.Context) {
	var query struct {
		ProductID    uuid.UUID `form:"product_id" binding:"required"`
		LocationType string    `form:"location_type" binding:"required"`
		LocationID   uuid.UUID `form:"location_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	stock, err := h.stockService.GetByProductAndLocation(c.Request.Context(), query.ProductID, query.LocationType, query.LocationID)
	if err != nil {
		dto.NotFound(c, "Stock entry not found")
		return
	}

	dto.Success(c, stock)
}

// SetMinStockAlert handles PUT /api/stock/:id/min-alert.
func (h *StockHandler) SetMinStockAlert(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid stock ID")
		return
	}

	var req struct {
		MinAlert int `json:"min_alert" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	stock, err := h.stockService.SetMinStockAlert(c.Request.Context(), id, req.MinAlert)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, stock)
}
