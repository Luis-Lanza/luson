package ports

import (
	"context"
	"time"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context, filter UserFilter) ([]domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
}

// UserFilter defines the filter options for listing users.
type UserFilter struct {
	Role     *string
	BranchID *uuid.UUID
	Active   *bool
	Limit    int
	Offset   int
}

// BranchRepository defines the interface for branch data access.
type BranchRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error)
	List(ctx context.Context, filter BranchFilter) ([]domain.Branch, error)
	Create(ctx context.Context, branch *domain.Branch) error
	Update(ctx context.Context, branch *domain.Branch) error
}

// BranchFilter defines the filter options for listing branches.
type BranchFilter struct {
	Active *bool
	Limit  int
	Offset int
}

// SupplierRepository defines the interface for supplier data access.
type SupplierRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error)
	List(ctx context.Context, filter SupplierFilter) ([]domain.Supplier, error)
	Create(ctx context.Context, supplier *domain.Supplier) error
	Update(ctx context.Context, supplier *domain.Supplier) error
}

// SupplierFilter defines the filter options for listing suppliers.
type SupplierFilter struct {
	Active *bool
	Limit  int
	Offset int
}

// TokenRepository defines the interface for refresh token storage.
type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	DeleteUserTokens(ctx context.Context, userID uuid.UUID) error
}

// RefreshToken represents a stored refresh token.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// ProductRepository defines the interface for product data access.
type ProductRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	FindByName(ctx context.Context, name string) (*domain.Product, error)
	List(ctx context.Context, filter ProductFilter) ([]domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, product *domain.Product) error
}

// ProductFilter defines the filter options for listing products.
type ProductFilter struct {
	ProductType *string
	Brand       *string
	VehicleType *string
	BatteryType *string
	Active      *bool
	Search      *string
	Limit       int
	Offset      int
}

// StockRepository defines the interface for stock data access.
type StockRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Stock, error)
	FindByProductAndLocation(ctx context.Context, productID uuid.UUID, locationType string, locationID uuid.UUID) (*domain.Stock, error)
	ListByLocation(ctx context.Context, locationType string, locationID uuid.UUID, filter StockFilter) ([]domain.Stock, error)
	ListByProduct(ctx context.Context, productID uuid.UUID) ([]domain.Stock, error)
	Create(ctx context.Context, stock *domain.Stock) error
	Update(ctx context.Context, stock *domain.Stock) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// StockFilter defines the filter options for listing stock.
type StockFilter struct {
	LowStockOnly *bool
	Limit        int
	Offset       int
}

// PurchaseBatchRepository defines the interface for purchase batch data access.
type PurchaseBatchRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.PurchaseBatch, error)
	List(ctx context.Context, filter PurchaseBatchFilter) ([]domain.PurchaseBatch, error)
	Create(ctx context.Context, batch *domain.PurchaseBatch) error
	MarkAsReceived(ctx context.Context, id uuid.UUID, receivedBy uuid.UUID) error
}

// PurchaseBatchFilter defines the filter options for listing purchase batches.
type PurchaseBatchFilter struct {
	SupplierID *uuid.UUID
	Received   *bool
	FromDate   *time.Time
	ToDate     *time.Time
	Limit      int
	Offset     int
}

// TransferRepository defines the interface for transfer data access.
type TransferRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Transfer, error)
	FindWithDetails(ctx context.Context, id uuid.UUID) (*domain.TransferWithDetails, error)
	List(ctx context.Context, filter TransferFilter) ([]domain.Transfer, error)
	Create(ctx context.Context, transfer *domain.Transfer, details []domain.TransferDetail) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransferStatus, userID *uuid.UUID) error
}

// TransferFilter defines the filter options for listing transfers.
type TransferFilter struct {
	OriginType      *string
	OriginID        *uuid.UUID
	DestinationType *string
	DestinationID   *uuid.UUID
	Status          *string
	RequestedBy     *uuid.UUID
	Limit           int
	Offset          int
}
