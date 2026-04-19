package middleware

import (
	"strings"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/dto"
	"github.com/Luis-Lanza/luson/internal/infrastructure/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Context keys for storing user information in Gin context.
const (
	ContextUserIDKey   = "user_id"
	ContextUserRoleKey = "user_role"
	ContextBranchIDKey = "branch_id"
)

// JWTService defines the interface for JWT operations needed by middleware.
type JWTService interface {
	ValidateAccessToken(token string) (*jwt.AccessTokenClaims, error)
}

// jwtService implements JWTService using the jwt package.
type jwtService struct {
	secret string
}

// NewJWTService creates a new JWT service.
func NewJWTService(secret string) JWTService {
	return &jwtService{secret: secret}
}

// ValidateAccessToken validates an access token.
func (s *jwtService) ValidateAccessToken(token string) (*jwt.AccessTokenClaims, error) {
	return jwt.ValidateAccessToken(token, s.secret)
}

// Auth returns a middleware that validates JWT tokens and sets user info in context.
func Auth(jwtSvc JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			dto.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			dto.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwtSvc.ValidateAccessToken(tokenString)
		if err != nil {
			dto.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Parse user ID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			dto.Unauthorized(c, "Invalid user ID in token")
			c.Abort()
			return
		}

		// Parse branch ID if present
		var branchID *uuid.UUID
		if claims.BranchID != nil && *claims.BranchID != "" {
			parsed, err := uuid.Parse(*claims.BranchID)
			if err == nil {
				branchID = &parsed
			}
		}

		// Set user info in context
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextUserRoleKey, domain.UserRole(claims.Role))
		if branchID != nil {
			c.Set(ContextBranchIDKey, *branchID)
		}

		c.Next()
	}
}

// RequireRole returns a middleware that checks if the user has one of the required roles.
func RequireRole(roles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		roleVal, exists := c.Get(ContextUserRoleKey)
		if !exists {
			dto.Forbidden(c, "User role not found")
			c.Abort()
			return
		}

		userRole, ok := roleVal.(domain.UserRole)
		if !ok {
			dto.Forbidden(c, "Invalid user role")
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			dto.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin returns a middleware that requires admin role.
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(domain.UserRoleAdmin)
}

// GetUserID retrieves the user ID from the Gin context.
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get(ContextUserIDKey)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := userID.(uuid.UUID)
	return id, ok
}

// GetUserRole retrieves the user role from the Gin context.
func GetUserRole(c *gin.Context) (domain.UserRole, bool) {
	role, exists := c.Get(ContextUserRoleKey)
	if !exists {
		return "", false
	}
	r, ok := role.(domain.UserRole)
	return r, ok
}

// GetBranchID retrieves the branch ID from the Gin context.
func GetBranchID(c *gin.Context) (uuid.UUID, bool) {
	branchID, exists := c.Get(ContextBranchIDKey)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := branchID.(uuid.UUID)
	return id, ok
}
