package handlers

import (
	"time"

	"github.com/Luis-Lanza/luson/internal/infrastructure/http/dto"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PurchaseBatchHandler handles purchase batch management HTTP requests.
type PurchaseBatchHandler struct {
	batchService ports.PurchaseBatchService
}

// NewPurchaseBatchHandler creates a new purchase batch handler.
func NewPurchaseBatchHandler(batchService ports.PurchaseBatchService) *PurchaseBatchHandler {
	return &PurchaseBatchHandler{
		batchService: batchService,
	}
}

// List handles GET /api/purchase-batches.
func (h *PurchaseBatchHandler) List(c *gin.Context) {
	var query struct {
		Limit    int        `form:"limit,default=20" binding:"max=100"`
		Offset   int        `form:"offset,default=0"`
		Received *bool      `form:"received,omitempty"`
		Supplier *uuid.UUID `form:"supplier_id,omitempty"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	filter := ports.PurchaseBatchFilter{
		Limit:      query.Limit,
		Offset:     query.Offset,
		Received:   query.Received,
		SupplierID: query.Supplier,
	}

	batches, err := h.batchService.List(c.Request.Context(), filter)
	if err != nil {
		dto.InternalError(c, err.Error())
		return
	}

	total := len(batches)
	dto.SuccessWithMeta(c, batches, &dto.Meta{
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	})
}

// GetByID handles GET /api/purchase-batches/:id.
func (h *PurchaseBatchHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid purchase batch ID")
		return
	}

	batch, err := h.batchService.GetByID(c.Request.Context(), id)
	if err != nil {
		dto.NotFound(c, "Purchase batch not found")
		return
	}

	dto.Success(c, batch)
}

// Create handles POST /api/purchase-batches.
func (h *PurchaseBatchHandler) Create(c *gin.Context) {
	var req dto.CreatePurchaseBatchRequest
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

	// Parse date if provided
	var purchaseDate time.Time
	if req.PurchaseDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.PurchaseDate)
		if err != nil {
			dto.BadRequest(c, "Invalid purchase_date format (expected RFC3339)")
			return
		}
		purchaseDate = parsed
	} else {
		purchaseDate = time.Now()
	}

	// Convert items
	items := make([]ports.PurchaseBatchItemRequest, len(req.Items))
	for i, item := range req.Items {
		items[i] = ports.PurchaseBatchItemRequest{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitCost:  item.UnitCost,
		}
	}

	serviceReq := ports.CreatePurchaseBatchRequest{
		SupplierID:    req.SupplierID,
		InvoiceNumber: req.InvoiceNumber,
		PurchaseDate:  purchaseDate,
		Notes:         req.Notes,
		Items:         items,
		CreatedBy:     userID.(uuid.UUID),
	}

	batch, err := h.batchService.Create(c.Request.Context(), serviceReq)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Created(c, batch)
}

// Receive handles POST /api/purchase-batches/:id/receive.
func (h *PurchaseBatchHandler) Receive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid purchase batch ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		dto.Unauthorized(c, "User not authenticated")
		return
	}

	if err := h.batchService.Receive(c.Request.Context(), id, userID.(uuid.UUID)); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Purchase batch marked as received"})
}
