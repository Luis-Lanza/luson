package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/google/uuid"
)

// branchService implements the ports.BranchService interface.
type branchService struct {
	branchRepo ports.BranchRepository
}

// NewBranchService creates a new instance of BranchService.
func NewBranchService(branchRepo ports.BranchRepository) ports.BranchService {
	return &branchService{
		branchRepo: branchRepo,
	}
}

// Create creates a new branch.
func (s *branchService) Create(ctx context.Context, req ports.CreateBranchRequest) (*domain.Branch, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	// Create branch entity
	branch := &domain.Branch{
		ID:               uuid.New(),
		Name:             req.Name,
		Address:          req.Address,
		PettyCashBalance: req.PettyCashBalance,
		Active:           true,
		CreatedAt:        time.Now(),
	}

	// Validate entity
	if err := branch.IsValid(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.branchRepo.Create(ctx, branch); err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	return branch, nil
}

// GetByID retrieves a branch by ID.
func (s *branchService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error) {
	branch, err := s.branchRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("branch not found: %w", err)
	}

	return branch, nil
}

// List retrieves branches with filtering.
func (s *branchService) List(ctx context.Context, filter ports.BranchFilter) ([]domain.Branch, error) {
	branches, err := s.branchRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	return branches, nil
}

// Update updates a branch's information.
func (s *branchService) Update(ctx context.Context, id uuid.UUID, req ports.UpdateBranchRequest) (*domain.Branch, error) {
	// Find existing branch
	branch, err := s.branchRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("branch not found: %w", err)
	}

	// Apply updates
	if req.Name != nil {
		branch.Name = *req.Name
	}
	if req.Address != nil {
		branch.Address = *req.Address
	}
	if req.PettyCashBalance != nil {
		branch.PettyCashBalance = *req.PettyCashBalance
	}
	if req.Active != nil {
		branch.Active = *req.Active
	}

	// Validate entity
	if err := branch.IsValid(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.branchRepo.Update(ctx, branch); err != nil {
		return nil, fmt.Errorf("failed to update branch: %w", err)
	}

	return branch, nil
}

// Deactivate deactivates a branch.
func (s *branchService) Deactivate(ctx context.Context, id uuid.UUID) error {
	// Find existing branch
	branch, err := s.branchRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("branch not found: %w", err)
	}

	branch.Active = false

	// Save to repository
	if err := s.branchRepo.Update(ctx, branch); err != nil {
		return fmt.Errorf("failed to deactivate branch: %w", err)
	}

	return nil
}
