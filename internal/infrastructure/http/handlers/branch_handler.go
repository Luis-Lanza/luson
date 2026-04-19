package handlers

import (
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/dto"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BranchHandler handles branch management HTTP requests.
type BranchHandler struct {
	branchService ports.BranchService
}

// NewBranchHandler creates a new branch handler.
func NewBranchHandler(branchService ports.BranchService) *BranchHandler {
	return &BranchHandler{
		branchService: branchService,
	}
}

// List handles GET /api/branches.
func (h *BranchHandler) List(c *gin.Context) {
	var query dto.ListQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	filter := ports.BranchFilter{
		Limit:  query.Limit,
		Offset: query.Offset,
		Active: query.Active,
	}

	branches, err := h.branchService.List(c.Request.Context(), filter)
	if err != nil {
		dto.InternalError(c, err.Error())
		return
	}

	total := len(branches)
	dto.SuccessWithMeta(c, branches, &dto.Meta{
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	})
}

// GetByID handles GET /api/branches/:id.
func (h *BranchHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid branch ID")
		return
	}

	branch, err := h.branchService.GetByID(c.Request.Context(), id)
	if err != nil {
		dto.NotFound(c, "Branch not found")
		return
	}

	dto.Success(c, branch)
}

// Create handles POST /api/branches.
func (h *BranchHandler) Create(c *gin.Context) {
	var req dto.CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	serviceReq := ports.CreateBranchRequest{
		Name:             req.Name,
		Address:          req.Address,
		PettyCashBalance: req.PettyCashBalance,
	}

	branch, err := h.branchService.Create(c.Request.Context(), serviceReq)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Created(c, branch)
}

// Update handles PUT /api/branches/:id.
func (h *BranchHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid branch ID")
		return
	}

	var req dto.UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	serviceReq := ports.UpdateBranchRequest{
		Name:             req.Name,
		Address:          req.Address,
		PettyCashBalance: req.PettyCashBalance,
		Active:           req.Active,
	}

	branch, err := h.branchService.Update(c.Request.Context(), id, serviceReq)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, branch)
}

// Deactivate handles DELETE /api/branches/:id.
func (h *BranchHandler) Deactivate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid branch ID")
		return
	}

	if err := h.branchService.Deactivate(c.Request.Context(), id); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Branch deactivated successfully"})
}
