package handlers

import (
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/dto"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SupplierHandler handles supplier management HTTP requests.
type SupplierHandler struct {
	supplierService ports.SupplierService
}

// NewSupplierHandler creates a new supplier handler.
func NewSupplierHandler(supplierService ports.SupplierService) *SupplierHandler {
	return &SupplierHandler{
		supplierService: supplierService,
	}
}

// List handles GET /api/suppliers.
func (h *SupplierHandler) List(c *gin.Context) {
	var query dto.ListQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	filter := ports.SupplierFilter{
		Limit:  query.Limit,
		Offset: query.Offset,
		Active: query.Active,
	}

	suppliers, err := h.supplierService.List(c.Request.Context(), filter)
	if err != nil {
		dto.InternalError(c, err.Error())
		return
	}

	total := len(suppliers)
	dto.SuccessWithMeta(c, suppliers, &dto.Meta{
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	})
}

// GetByID handles GET /api/suppliers/:id.
func (h *SupplierHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid supplier ID")
		return
	}

	supplier, err := h.supplierService.GetByID(c.Request.Context(), id)
	if err != nil {
		dto.NotFound(c, "Supplier not found")
		return
	}

	dto.Success(c, supplier)
}

// Create handles POST /api/suppliers.
func (h *SupplierHandler) Create(c *gin.Context) {
	var req dto.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	serviceReq := ports.CreateSupplierRequest{
		Name:    req.Name,
		Contact: req.Contact,
		Address: req.Address,
	}

	supplier, err := h.supplierService.Create(c.Request.Context(), serviceReq)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Created(c, supplier)
}

// Update handles PUT /api/suppliers/:id.
func (h *SupplierHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid supplier ID")
		return
	}

	var req dto.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	serviceReq := ports.UpdateSupplierRequest{
		Name:    req.Name,
		Contact: req.Contact,
		Address: req.Address,
		Active:  req.Active,
	}

	supplier, err := h.supplierService.Update(c.Request.Context(), id, serviceReq)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, supplier)
}

// Deactivate handles DELETE /api/suppliers/:id.
func (h *SupplierHandler) Deactivate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid supplier ID")
		return
	}

	if err := h.supplierService.Deactivate(c.Request.Context(), id); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Supplier deactivated successfully"})
}
