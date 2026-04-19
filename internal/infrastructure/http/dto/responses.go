package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response wrapper.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta contains metadata for paginated responses.
type Meta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// LoginResponse represents a successful login response.
type LoginResponse struct {
	User         interface{} `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

// PaginatedResponse represents a paginated list response.
type PaginatedResponse struct {
	Items  interface{} `json:"items"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

// Success returns a successful API response.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta returns a successful API response with metadata.
func SuccessWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Created returns a 201 Created response.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Error returns an error response with the given status code.
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

// BadRequest returns a 400 Bad Request response.
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized returns a 401 Unauthorized response.
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden returns a 403 Forbidden response.
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	Error(c, http.StatusForbidden, message)
}

// NotFound returns a 404 Not Found response.
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Not found"
	}
	Error(c, http.StatusNotFound, message)
}

// InternalError returns a 500 Internal Server Error response.
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	Error(c, http.StatusInternalServerError, message)
}

// ValidationError returns a 400 response with validation errors.
func ValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"errors":  errors,
	})
}
