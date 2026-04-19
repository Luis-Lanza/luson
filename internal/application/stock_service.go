package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/google/uuid"
)

// stockService implements the ports.StockService interface.
type stockService struct {
	stockRepo   ports.StockRepository
	productRepo ports.ProductRepository
}

// NewStockService creates a new instance of StockService.
func NewStockService(stockRepo ports.StockRepository, productRepo ports.ProductRepository) ports.StockService {
	return &stockService{
		stockRepo:   stockRepo,
		productRepo: productRepo,
	}
}

// GetByID retrieves a stock entry by ID.
func (s *stockService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
	stock, err := s.stockRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %w", err)
	}

	return stock, nil
}

// GetByProductAndLocation retrieves stock for a specific product at a specific location.
func (s *stockService) GetByProductAndLocation(ctx context.Context, productID uuid.UUID, locationType string, locationID uuid.UUID) (*domain.Stock, error) {
	// Validate product exists
	_, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	stock, err := s.stockRepo.FindByProductAndLocation(ctx, productID, locationType, locationID)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %w", err)
	}

	if stock == nil {
		return nil, domain.ErrStockNotFound
	}

	return stock, nil
}

// ListByLocation retrieves stock entries for a specific location.
func (s *stockService) ListByLocation(ctx context.Context, locationType string, locationID uuid.UUID, filter ports.StockFilter) ([]domain.Stock, error) {
	stock, err := s.stockRepo.ListByLocation(ctx, locationType, locationID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock: %w", err)
	}

	return stock, nil
}

// ListByProduct retrieves stock entries for a specific product across all locations.
func (s *stockService) ListByProduct(ctx context.Context, productID uuid.UUID) ([]domain.Stock, error) {
	// Validate product exists
	_, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	stock, err := s.stockRepo.ListByProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock: %w", err)
	}

	return stock, nil
}

// SetMinStockAlert updates the minimum stock alert level for a stock entry.
func (s *stockService) SetMinStockAlert(ctx context.Context, id uuid.UUID, minAlert int) (*domain.Stock, error) {
	if minAlert < 0 {
		return nil, errors.New("min_stock_alert cannot be negative")
	}

	// Find existing stock
	stock, err := s.stockRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("stock not found: %w", err)
	}

	stock.MinStockAlert = minAlert

	// Save to repository
	if err := s.stockRepo.Update(ctx, stock); err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	return stock, nil
}
