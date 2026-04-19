package handlers

import (
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/dto"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService ports.AuthService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles POST /api/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		dto.Unauthorized(c, err.Error())
		return
	}

	response := dto.LoginResponse{
		User:         result.User,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	dto.Success(c, response)
}

// Refresh handles POST /api/auth/refresh.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	tokenPair, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		dto.Unauthorized(c, err.Error())
		return
	}

	dto.Success(c, tokenPair)
}

// Me handles GET /api/auth/me.
func (h *AuthHandler) Me(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		dto.Unauthorized(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		dto.Unauthorized(c, "Invalid user ID")
		return
	}

	user, err := h.authService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		dto.NotFound(c, "User not found")
		return
	}

	dto.Success(c, user)
}

// Logout handles POST /api/auth/logout.
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		dto.InternalError(c, err.Error())
		return
	}

	dto.Success(c, gin.H{"message": "Logged out successfully"})
}
