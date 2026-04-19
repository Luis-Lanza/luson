package ports

import (
	"context"

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
