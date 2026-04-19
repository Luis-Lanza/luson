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

// purchaseBatchService implements the ports.PurchaseBatchService interface.
type purchaseBatchService struct {
	batchRepo   ports.PurchaseBatchRepository
	productRepo ports.ProductRepository
	stockRepo   ports.StockRepository
}

// NewPurchaseBatchService creates a new instance of PurchaseBatchService.
func NewPurchaseBatchService(batchRepo ports.PurchaseBatchRepository, productRepo ports.ProductRepository, stockRepo ports.StockRepository) ports.PurchaseBatchService {
	return &purchaseBatchService{
		batchRepo:   batchRepo,
		productRepo: productRepo,
		stockRepo:   stockRepo,
	}
}

// Create creates a new purchase batch with items.
func (s *purchaseBatchService) Create(ctx context.Context, req ports.CreatePurchaseBatchRequest) (*domain.PurchaseBatch, error) {
	// Validate input
	if req.CreatedBy == uuid.Nil {
		return nil, errors.New("created_by is required")
	}

	if len(req.Items) == 0 {
		return nil, domain.ErrPurchaseBatchNoItems
	}

	// Validate all products exist and calculate total
	var totalCost float64
	for i, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("item %d: %w", i, domain.ErrPurchaseBatchInvalidQuantity)
		}
		if item.UnitCost <= 0 {
			return nil, fmt.Errorf("item %d: %w", i, domain.ErrPurchaseBatchInvalidCost)
		}

		_, err := s.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("item %d: product not found: %w", i, err)
		}

		totalCost += float64(item.Quantity) * item.UnitCost
	}

	// Create batch entity
	batch := &domain.PurchaseBatch{
		ID:            uuid.New(),
		SupplierID:    req.SupplierID,
		InvoiceNumber: req.InvoiceNumber,
		PurchaseDate:  req.PurchaseDate,
		Notes:         req.Notes,
		TotalCost:     totalCost,
		Received:      false,
		CreatedBy:     req.CreatedBy,
		CreatedAt:     time.Now(),
	}

	// Validate entity
	if err := batch.IsValid(); err != nil {
		return nil, err
	}

	// Save to repository (details are handled separately in a real implementation with transaction)
	if err := s.batchRepo.Create(ctx, batch); err != nil {
		return nil, fmt.Errorf("failed to create purchase batch: %w", err)
	}

	return batch, nil
}

// GetByID retrieves a purchase batch by ID with its details.
// Note: This is a simplified implementation. In production, you'd add FindWithDetails to the interface.
func (s *purchaseBatchService) GetByID(ctx context.Context, id uuid.UUID) (*domain.PurchaseBatchWithDetails, error) {
	batch, err := s.batchRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("purchase batch not found: %w", err)
	}

	// Return batch without details for now (details would need a separate interface method)
	return &domain.PurchaseBatchWithDetails{
		Batch:   *batch,
		Details: []domain.PurchaseBatchDetail{}, // Empty for now
	}, nil
}

// List retrieves purchase batches with filtering.
func (s *purchaseBatchService) List(ctx context.Context, filter ports.PurchaseBatchFilter) ([]domain.PurchaseBatch, error) {
	batches, err := s.batchRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list purchase batches: %w", err)
	}

	return batches, nil
}

// Receive marks a purchase batch as received and updates stock.
func (s *purchaseBatchService) Receive(ctx context.Context, id uuid.UUID, receivedBy uuid.UUID) error {
	// Find the batch
	batch, err := s.batchRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase batch not found: %w", err)
	}

	// Check if already received
	if batch.IsReceived() {
		return domain.ErrPurchaseBatchAlreadyReceived
	}

	// Mark as received
	if err := s.batchRepo.MarkAsReceived(ctx, id, receivedBy); err != nil {
		return fmt.Errorf("failed to mark purchase batch as received: %w", err)
	}

	// TODO: Update stock quantities (this would be done in a real implementation)
	// For now, we just mark the batch as received

	return nil
}
