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

// CreateProductRequest represents a request to create a product.
type CreateProductRequest struct {
	Name         string   `json:"name" binding:"required"`
	Description  *string  `json:"description,omitempty"`
	ProductType  string   `json:"product_type" binding:"required,oneof=bateria accesorio"`
	Brand        *string  `json:"brand,omitempty"`
	Model        *string  `json:"model,omitempty"`
	Voltage      *float64 `json:"voltage,omitempty"`
	Amperage     *float64 `json:"amperage,omitempty"`
	BatteryType  *string  `json:"battery_type,omitempty" binding:"omitempty,oneof=seca liquida"`
	Polarity     *string  `json:"polarity,omitempty" binding:"omitempty,oneof=izquierda derecha"`
	AcidLiters   *float64 `json:"acid_liters,omitempty"`
	VehicleType  *string  `json:"vehicle_type,omitempty" binding:"omitempty,oneof=auto moto otro"`
	MinSalePrice float64  `json:"min_sale_price" binding:"required,gt=0"`
}

// UpdateProductRequest represents a request to update a product.
type UpdateProductRequest struct {
	Name         *string  `json:"name,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Brand        *string  `json:"brand,omitempty"`
	Model        *string  `json:"model,omitempty"`
	Voltage      *float64 `json:"voltage,omitempty"`
	Amperage     *float64 `json:"amperage,omitempty"`
	BatteryType  *string  `json:"battery_type,omitempty" binding:"omitempty,oneof=seca liquida"`
	Polarity     *string  `json:"polarity,omitempty" binding:"omitempty,oneof=izquierda derecha"`
	AcidLiters   *float64 `json:"acid_liters,omitempty"`
	VehicleType  *string  `json:"vehicle_type,omitempty" binding:"omitempty,oneof=auto moto otro"`
	MinSalePrice *float64 `json:"min_sale_price,omitempty" binding:"omitempty,gt=0"`
	Active       *bool    `json:"active,omitempty"`
}

// CreatePurchaseBatchRequest represents a request to create a purchase batch.
type CreatePurchaseBatchRequest struct {
	SupplierID    *uuid.UUID                 `json:"supplier_id,omitempty"`
	InvoiceNumber *string                    `json:"invoice_number,omitempty"`
	PurchaseDate  string                     `json:"purchase_date,omitempty"`
	Notes         *string                    `json:"notes,omitempty"`
	Items         []PurchaseBatchItemRequest `json:"items" binding:"required,min=1"`
}

// PurchaseBatchItemRequest represents an item in a purchase batch request.
type PurchaseBatchItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
	UnitCost  float64   `json:"unit_cost" binding:"required,gt=0"`
}

// CreateTransferRequest represents a request to create a transfer.
type CreateTransferRequest struct {
	OriginType      string                `json:"origin_type" binding:"required,oneof=branch warehouse"`
	OriginID        uuid.UUID             `json:"origin_id" binding:"required"`
	DestinationType string                `json:"destination_type" binding:"required,oneof=branch warehouse"`
	DestinationID   uuid.UUID             `json:"destination_id" binding:"required"`
	Notes           *string               `json:"notes,omitempty"`
	Items           []TransferItemRequest `json:"items" binding:"required,min=1"`
}

// TransferItemRequest represents an item in a transfer request.
type TransferItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}
