package ports

import (
	"context"
	"time"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/google/uuid"
)

// AuthService defines the interface for authentication use cases.
type AuthService interface {
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	GetCurrentUser(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	Logout(ctx context.Context, refreshToken string) error
}

// LoginResult contains the result of a successful login.
type LoginResult struct {
	User         domain.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

// TokenPair contains a new pair of access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserService defines the interface for user management use cases.
type UserService interface {
	Create(ctx context.Context, req CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context, filter UserFilter) ([]domain.User, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*domain.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	Deactivate(ctx context.Context, id uuid.UUID) error
}

// CreateUserRequest contains the data needed to create a user.
type CreateUserRequest struct {
	Username string
	Password string
	Role     string
	BranchID *uuid.UUID
}

// UpdateUserRequest contains the data for partial user update.
type UpdateUserRequest struct {
	Role     *string
	BranchID *uuid.UUID
	Active   *bool
}

// BranchService defines the interface for branch management use cases.
type BranchService interface {
	Create(ctx context.Context, req CreateBranchRequest) (*domain.Branch, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error)
	List(ctx context.Context, filter BranchFilter) ([]domain.Branch, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateBranchRequest) (*domain.Branch, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
}

// CreateBranchRequest contains the data needed to create a branch.
type CreateBranchRequest struct {
	Name             string
	Address          string
	PettyCashBalance float64
}

// UpdateBranchRequest contains the data for partial branch update.
type UpdateBranchRequest struct {
	Name             *string
	Address          *string
	PettyCashBalance *float64
	Active           *bool
}

// SupplierService defines the interface for supplier management use cases.
type SupplierService interface {
	Create(ctx context.Context, req CreateSupplierRequest) (*domain.Supplier, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error)
	List(ctx context.Context, filter SupplierFilter) ([]domain.Supplier, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateSupplierRequest) (*domain.Supplier, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
}

// CreateSupplierRequest contains the data needed to create a supplier.
type CreateSupplierRequest struct {
	Name    string
	Contact *string
	Address *string
}

// UpdateSupplierRequest contains the data for partial supplier update.
type UpdateSupplierRequest struct {
	Name    *string
	Contact *string
	Address *string
	Active  *bool
}

// ProductService defines the interface for product management use cases.
type ProductService interface {
	Create(ctx context.Context, req CreateProductRequest) (*domain.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	List(ctx context.Context, filter ProductFilter) ([]domain.Product, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*domain.Product, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
}

// CreateProductRequest contains the data needed to create a product.
type CreateProductRequest struct {
	Name         string
	Description  *string
	ProductType  string
	Brand        *string
	Model        *string
	Voltage      *float64
	Amperage     *float64
	BatteryType  *string
	Polarity     *string
	AcidLiters   *float64
	VehicleType  *string
	MinSalePrice float64
	CreatedBy    uuid.UUID
}

// UpdateProductRequest contains the data for partial product update.
type UpdateProductRequest struct {
	Name         *string
	Description  *string
	Brand        *string
	Model        *string
	Voltage      *float64
	Amperage     *float64
	BatteryType  *string
	Polarity     *string
	AcidLiters   *float64
	VehicleType  *string
	MinSalePrice *float64
	Active       *bool
}

// StockService defines the interface for stock management use cases.
type StockService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Stock, error)
	GetByProductAndLocation(ctx context.Context, productID uuid.UUID, locationType string, locationID uuid.UUID) (*domain.Stock, error)
	ListByLocation(ctx context.Context, locationType string, locationID uuid.UUID, filter StockFilter) ([]domain.Stock, error)
	ListByProduct(ctx context.Context, productID uuid.UUID) ([]domain.Stock, error)
	SetMinStockAlert(ctx context.Context, id uuid.UUID, minAlert int) (*domain.Stock, error)
}

// PurchaseBatchService defines the interface for purchase batch management use cases.
type PurchaseBatchService interface {
	Create(ctx context.Context, req CreatePurchaseBatchRequest) (*domain.PurchaseBatch, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PurchaseBatchWithDetails, error)
	List(ctx context.Context, filter PurchaseBatchFilter) ([]domain.PurchaseBatch, error)
	Receive(ctx context.Context, id uuid.UUID, receivedBy uuid.UUID) error
}

// CreatePurchaseBatchRequest contains the data needed to create a purchase batch.
type CreatePurchaseBatchRequest struct {
	SupplierID    *uuid.UUID
	InvoiceNumber *string
	PurchaseDate  time.Time
	Notes         *string
	Items         []PurchaseBatchItemRequest
	CreatedBy     uuid.UUID
}

// PurchaseBatchItemRequest represents an item in a purchase batch request.
type PurchaseBatchItemRequest struct {
	ProductID uuid.UUID
	Quantity  int
	UnitCost  float64
}

// TransferService defines the interface for transfer management use cases.
type TransferService interface {
	Create(ctx context.Context, req CreateTransferRequest) (*domain.Transfer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TransferWithDetails, error)
	List(ctx context.Context, filter TransferFilter) ([]domain.Transfer, error)
	Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error
	Reject(ctx context.Context, id uuid.UUID, rejectedBy uuid.UUID, reason string) error
	MarkAsSent(ctx context.Context, id uuid.UUID, sentBy uuid.UUID) error
	MarkAsReceived(ctx context.Context, id uuid.UUID, receivedBy uuid.UUID) error
	Cancel(ctx context.Context, id uuid.UUID) error
}

// CreateTransferRequest contains the data needed to create a transfer.
type CreateTransferRequest struct {
	OriginType      string
	OriginID        uuid.UUID
	DestinationType string
	DestinationID   uuid.UUID
	Notes           *string
	Items           []TransferItemRequest
	RequestedBy     uuid.UUID
}

// TransferItemRequest represents an item in a transfer request.
type TransferItemRequest struct {
	ProductID uuid.UUID
	Quantity  int
}
