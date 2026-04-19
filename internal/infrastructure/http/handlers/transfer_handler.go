package handlers

import (
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/dto"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TransferHandler handles transfer management HTTP requests.
type TransferHandler struct {
	transferService ports.TransferService
}

// NewTransferHandler creates a new transfer handler.
func NewTransferHandler(transferService ports.TransferService) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
	}
}

// List handles GET /api/transfers.
func (h *TransferHandler) List(c *gin.Context) {
	var query struct {
		Limit           int        `form:"limit,default=20" binding:"max=100"`
		Offset          int        `form:"offset,default=0"`
		Status          *string    `form:"status,omitempty"`
		OriginID        *uuid.UUID `form:"origin_id,omitempty"`
		OriginType      *string    `form:"origin_type,omitempty"`
		DestinationID   *uuid.UUID `form:"destination_id,omitempty"`
		DestinationType *string    `form:"destination_type,omitempty"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	filter := ports.TransferFilter{
		Limit:           query.Limit,
		Offset:          query.Offset,
		Status:          query.Status,
		OriginID:        query.OriginID,
		OriginType:      query.OriginType,
		DestinationID:   query.DestinationID,
		DestinationType: query.DestinationType,
	}

	transfers, err := h.transferService.List(c.Request.Context(), filter)
	if err != nil {
		dto.InternalError(c, err.Error())
		return
	}

	total := len(transfers)
	dto.SuccessWithMeta(c, transfers, &dto.Meta{
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	})
}

// GetByID handles GET /api/transfers/:id.
func (h *TransferHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid transfer ID")
		return
	}

	transfer, err := h.transferService.GetByID(c.Request.Context(), id)
	if err != nil {
		dto.NotFound(c, "Transfer not found")
		return
	}

	dto.Success(c, transfer)
}

// Create handles POST /api/transfers.
func (h *TransferHandler) Create(c *gin.Context) {
	var req dto.CreateTransferRequest
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

	// Convert items
	items := make([]ports.TransferItemRequest, len(req.Items))
	for i, item := range req.Items {
		items[i] = ports.TransferItemRequest{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	serviceReq := ports.CreateTransferRequest{
		OriginType:      req.OriginType,
		OriginID:        req.OriginID,
		DestinationType: req.DestinationType,
		DestinationID:   req.DestinationID,
		Notes:           req.Notes,
		Items:           items,
		RequestedBy:     userID.(uuid.UUID),
	}

	transfer, err := h.transferService.Create(c.Request.Context(), serviceReq)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Created(c, transfer)
}

// Approve handles POST /api/transfers/:id/approve.
func (h *TransferHandler) Approve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid transfer ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		dto.Unauthorized(c, "User not authenticated")
		return
	}

	if err := h.transferService.Approve(c.Request.Context(), id, userID.(uuid.UUID)); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Transfer approved successfully"})
}

// Reject handles POST /api/transfers/:id/reject.
func (h *TransferHandler) Reject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid transfer ID")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
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

	if err := h.transferService.Reject(c.Request.Context(), id, userID.(uuid.UUID), req.Reason); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Transfer rejected successfully"})
}

// MarkAsSent handles POST /api/transfers/:id/send.
func (h *TransferHandler) MarkAsSent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid transfer ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		dto.Unauthorized(c, "User not authenticated")
		return
	}

	if err := h.transferService.MarkAsSent(c.Request.Context(), id, userID.(uuid.UUID)); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Transfer marked as sent"})
}

// MarkAsReceived handles POST /api/transfers/:id/receive.
func (h *TransferHandler) MarkAsReceived(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid transfer ID")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		dto.Unauthorized(c, "User not authenticated")
		return
	}

	if err := h.transferService.MarkAsReceived(c.Request.Context(), id, userID.(uuid.UUID)); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Transfer marked as received"})
}

// Cancel handles DELETE /api/transfers/:id.
func (h *TransferHandler) Cancel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid transfer ID")
		return
	}

	if err := h.transferService.Cancel(c.Request.Context(), id); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Transfer cancelled successfully"})
}
