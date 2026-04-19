package domain

import "errors"

// Stock errors
var (
	ErrStockInvalidProductID      = errors.New("product_id is required")
	ErrStockInvalidLocationID     = errors.New("location_id is required")
	ErrStockInvalidLocationType   = errors.New("location_type is required")
	ErrStockNegativeQuantity      = errors.New("quantity cannot be negative")
	ErrStockNegativeMinStockAlert = errors.New("min_stock_alert cannot be negative")
	ErrStockNotFound              = errors.New("stock not found")
	ErrStockInsufficientQuantity  = errors.New("insufficient stock quantity")
	ErrStockDuplicate             = errors.New("stock entry already exists for this product and location")
)

// Transfer errors
var (
	ErrTransferInvalidOrigin       = errors.New("origin location is required")
	ErrTransferInvalidDestination  = errors.New("destination location is required")
	ErrTransferInvalidLocationType = errors.New("location type is required")
	ErrTransferSameLocation        = errors.New("origin and destination cannot be the same")
	ErrTransferInvalidStatus       = errors.New("invalid transfer status")
	ErrTransferInvalidTransition   = errors.New("invalid status transition")
	ErrTransferNoItems             = errors.New("transfer must have at least one item")
	ErrTransferInvalidQuantity     = errors.New("item quantity must be positive")
	ErrTransferNotFound            = errors.New("transfer not found")
)

// Product errors
var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductDuplicateName = errors.New("product name already exists")
)

// Purchase Batch errors
var (
	ErrPurchaseBatchNotFound        = errors.New("purchase batch not found")
	ErrPurchaseBatchNoItems         = errors.New("purchase batch must have at least one item")
	ErrPurchaseBatchInvalidCost     = errors.New("unit cost must be positive")
	ErrPurchaseBatchInvalidQuantity = errors.New("quantity must be positive")
	ErrPurchaseBatchAlreadyReceived = errors.New("purchase batch already received")
)
