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

// supplierService implements the ports.SupplierService interface.
type supplierService struct {
	supplierRepo ports.SupplierRepository
}

// NewSupplierService creates a new instance of SupplierService.
func NewSupplierService(supplierRepo ports.SupplierRepository) ports.SupplierService {
	return &supplierService{
		supplierRepo: supplierRepo,
	}
}

// Create creates a new supplier.
func (s *supplierService) Create(ctx context.Context, req ports.CreateSupplierRequest) (*domain.Supplier, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	// Create supplier entity
	supplier := &domain.Supplier{
		ID:        uuid.New(),
		Name:      req.Name,
		Contact:   req.Contact,
		Address:   req.Address,
		Active:    true,
		CreatedAt: time.Now(),
	}

	// Validate entity
	if err := supplier.IsValid(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.supplierRepo.Create(ctx, supplier); err != nil {
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	return supplier, nil
}

// GetByID retrieves a supplier by ID.
func (s *supplierService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	supplier, err := s.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("supplier not found: %w", err)
	}

	return supplier, nil
}

// List retrieves suppliers with filtering.
func (s *supplierService) List(ctx context.Context, filter ports.SupplierFilter) ([]domain.Supplier, error) {
	suppliers, err := s.supplierRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list suppliers: %w", err)
	}

	return suppliers, nil
}

// Update updates a supplier's information.
func (s *supplierService) Update(ctx context.Context, id uuid.UUID, req ports.UpdateSupplierRequest) (*domain.Supplier, error) {
	// Find existing supplier
	supplier, err := s.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("supplier not found: %w", err)
	}

	// Apply updates
	if req.Name != nil {
		supplier.Name = *req.Name
	}
	if req.Contact != nil {
		supplier.Contact = req.Contact
	}
	if req.Address != nil {
		supplier.Address = req.Address
	}
	if req.Active != nil {
		supplier.Active = *req.Active
	}

	// Validate entity
	if err := supplier.IsValid(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.supplierRepo.Update(ctx, supplier); err != nil {
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}

	return supplier, nil
}

// Deactivate deactivates a supplier.
func (s *supplierService) Deactivate(ctx context.Context, id uuid.UUID) error {
	// Find existing supplier
	supplier, err := s.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("supplier not found: %w", err)
	}

	supplier.Active = false

	// Save to repository
	if err := s.supplierRepo.Update(ctx, supplier); err != nil {
		return fmt.Errorf("failed to deactivate supplier: %w", err)
	}

	return nil
}
