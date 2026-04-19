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
