package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/infrastructure/jwt"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/google/uuid"
)

// PasswordHasher defines the interface for password operations.
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
}

// userService implements the ports.UserService interface.
type userService struct {
	userRepo       ports.UserRepository
	passwordHasher PasswordHasher
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo ports.UserRepository, passwordHasher PasswordHasher) ports.UserService {
	return &userService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

// Create creates a new user.
func (s *userService) Create(ctx context.Context, req ports.CreateUserRequest) (*domain.User, error) {
	// Validate input
	if req.Username == "" {
		return nil, errors.New("username is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	// Validate role
	role := domain.UserRole(req.Role)
	switch role {
	case domain.UserRoleAdmin, domain.UserRoleEncargadoAlmacen, domain.UserRoleCajero:
		// valid roles
	default:
		return nil, fmt.Errorf("invalid role: %s", req.Role)
	}

	// Cajero must have a branch
	if role == domain.UserRoleCajero && req.BranchID == nil {
		return nil, errors.New("cajero must have a branch_id")
	}

	// Hash password
	passwordHash, err := s.passwordHasher.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user := &domain.User{
		ID:           uuid.New(),
		Username:     req.Username,
		PasswordHash: passwordHash,
		Role:         role,
		BranchID:     req.BranchID,
		Active:       true,
		CreatedAt:    time.Now(),
	}

	// Validate entity
	if err := user.IsValid(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return user without password hash
	return &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		BranchID:  user.BranchID,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
	}, nil
}

// GetByID retrieves a user by ID.
func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Return user without password hash
	return &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		BranchID:  user.BranchID,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
	}, nil
}

// List retrieves users with filtering.
func (s *userService) List(ctx context.Context, filter ports.UserFilter) ([]domain.User, error) {
	users, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Return users without password hashes
	result := make([]domain.User, len(users))
	for i, user := range users {
		result[i] = domain.User{
			ID:        user.ID,
			Username:  user.Username,
			Role:      user.Role,
			BranchID:  user.BranchID,
			Active:    user.Active,
			CreatedAt: user.CreatedAt,
		}
	}

	return result, nil
}

// Update updates a user's information.
func (s *userService) Update(ctx context.Context, id uuid.UUID, req ports.UpdateUserRequest) (*domain.User, error) {
	// Find existing user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Apply updates
	if req.Role != nil {
		role := domain.UserRole(*req.Role)
		switch role {
		case domain.UserRoleAdmin, domain.UserRoleEncargadoAlmacen, domain.UserRoleCajero:
			user.Role = role
		default:
			return nil, fmt.Errorf("invalid role: %s", *req.Role)
		}
	}

	if req.BranchID != nil {
		user.BranchID = req.BranchID
	}

	if req.Active != nil {
		user.Active = *req.Active
	}

	// Validate that cajero has branch_id
	if user.Role == domain.UserRoleCajero && user.BranchID == nil {
		return nil, errors.New("cajero must have a branch_id")
	}

	// Validate entity
	if err := user.IsValid(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Return user without password hash
	return &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		BranchID:  user.BranchID,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
	}, nil
}

// UpdatePassword updates a user's password.
func (s *userService) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	// Find existing user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Hash new password
	passwordHash, err := s.passwordHasher.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = passwordHash

	// Save to repository
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// Deactivate deactivates a user.
func (s *userService) Deactivate(ctx context.Context, id uuid.UUID) error {
	// Find existing user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.Active = false

	// Save to repository
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// Ensure jwt.PasswordHasher implements PasswordHasher interface
var _ PasswordHasher = (*jwt.PasswordHasher)(nil)
