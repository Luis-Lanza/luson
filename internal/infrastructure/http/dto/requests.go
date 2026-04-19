package dto

import "github.com/google/uuid"

// LoginRequest represents a login request.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest represents a token refresh request.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// CreateUserRequest represents a request to create a user.
type CreateUserRequest struct {
	Username string     `json:"username" binding:"required"`
	Password string     `json:"password" binding:"required,min=6"`
	Role     string     `json:"role" binding:"required,oneof=admin encargado_almacen cajero"`
	BranchID *uuid.UUID `json:"branch_id,omitempty"`
}

// UpdateUserRequest represents a request to update a user.
type UpdateUserRequest struct {
	Role     *string    `json:"role,omitempty" binding:"omitempty,oneof=admin encargado_almacen cajero"`
	BranchID *uuid.UUID `json:"branch_id,omitempty"`
	Active   *bool      `json:"active,omitempty"`
}

// CreateBranchRequest represents a request to create a branch.
type CreateBranchRequest struct {
	Name             string  `json:"name" binding:"required"`
	Address          string  `json:"address,omitempty"`
	PettyCashBalance float64 `json:"petty_cash_balance,omitempty"`
}

// UpdateBranchRequest represents a request to update a branch.
type UpdateBranchRequest struct {
	Name             *string  `json:"name,omitempty"`
	Address          *string  `json:"address,omitempty"`
	PettyCashBalance *float64 `json:"petty_cash_balance,omitempty"`
	Active           *bool    `json:"active,omitempty"`
}

// CreateSupplierRequest represents a request to create a supplier.
type CreateSupplierRequest struct {
	Name    string  `json:"name" binding:"required"`
	Contact *string `json:"contact,omitempty"`
	Address *string `json:"address,omitempty"`
}

// UpdateSupplierRequest represents a request to update a supplier.
type UpdateSupplierRequest struct {
	Name    *string `json:"name,omitempty"`
	Contact *string `json:"contact,omitempty"`
	Address *string `json:"address,omitempty"`
	Active  *bool   `json:"active,omitempty"`
}

// ListQueryParams represents common query parameters for list endpoints.
type ListQueryParams struct {
	Limit  int   `form:"limit,default=20" binding:"max=100"`
	Offset int   `form:"offset,default=0"`
	Active *bool `form:"active,omitempty"`
}
