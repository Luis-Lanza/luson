package handlers

import (
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/dto"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user management HTTP requests.
type UserHandler struct {
	userService ports.UserService
}

// NewUserHandler creates a new user handler.
func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// List handles GET /api/users.
func (h *UserHandler) List(c *gin.Context) {
	var query dto.ListQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.BadRequest(c, "Invalid query parameters: "+err.Error())
		return
	}

	filter := ports.UserFilter{
		Limit:  query.Limit,
		Offset: query.Offset,
		Active: query.Active,
	}

	// Role filter from query
	if role := c.Query("role"); role != "" {
		filter.Role = &role
	}

	users, err := h.userService.List(c.Request.Context(), filter)
	if err != nil {
		dto.InternalError(c, err.Error())
		return
	}

	total := len(users) // Note: In production, you'd want a count query
	dto.SuccessWithMeta(c, users, &dto.Meta{
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	})
}

// GetByID handles GET /api/users/:id.
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		dto.NotFound(c, "User not found")
		return
	}

	dto.Success(c, user)
}

// Create handles POST /api/users.
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	serviceReq := ports.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
		BranchID: req.BranchID,
	}

	user, err := h.userService.Create(c.Request.Context(), serviceReq)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Created(c, user)
}

// Update handles PUT /api/users/:id.
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	serviceReq := ports.UpdateUserRequest{
		Role:     req.Role,
		BranchID: req.BranchID,
		Active:   req.Active,
	}

	user, err := h.userService.Update(c.Request.Context(), id, serviceReq)
	if err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, user)
}

// UpdatePassword handles PUT /api/users/:id/password.
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid user ID")
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	if err := h.userService.UpdatePassword(c.Request.Context(), id, req.Password); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Password updated successfully"})
}

// Deactivate handles DELETE /api/users/:id.
func (h *UserHandler) Deactivate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid user ID")
		return
	}

	if err := h.userService.Deactivate(c.Request.Context(), id); err != nil {
		dto.BadRequest(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "User deactivated successfully"})
}
