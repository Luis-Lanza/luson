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

// productService implements the ports.ProductService interface.
type productService struct {
	productRepo ports.ProductRepository
}

// NewProductService creates a new instance of ProductService.
func NewProductService(productRepo ports.ProductRepository) ports.ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

// Create creates a new product.
func (s *productService) Create(ctx context.Context, req ports.CreateProductRequest) (*domain.Product, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.MinSalePrice <= 0 {
		return nil, errors.New("min_sale_price must be positive")
	}

	// Validate product type
	var productType domain.ProductType
	switch req.ProductType {
	case "bateria":
		productType = domain.ProductTypeBateria
	case "accesorio":
		productType = domain.ProductTypeAccesorio
	default:
		return nil, fmt.Errorf("invalid product_type: %s", req.ProductType)
	}

	// Check for duplicate name
	existing, err := s.productRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check for duplicate name: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrProductDuplicateName
	}

	// Parse optional fields
	var batteryType *domain.BatteryType
	if req.BatteryType != nil {
		bt := domain.BatteryType(*req.BatteryType)
		batteryType = &bt
	}

	var polarity *domain.Polarity
	if req.Polarity != nil {
		p := domain.Polarity(*req.Polarity)
		polarity = &p
	}

	var vehicleType *domain.VehicleType
	if req.VehicleType != nil {
		vt := domain.VehicleType(*req.VehicleType)
		vehicleType = &vt
	}

	// Create product entity
	product := &domain.Product{
		ID:           uuid.New(),
		Name:         req.Name,
		Description:  req.Description,
		ProductType:  productType,
		Brand:        req.Brand,
		Model:        req.Model,
		Voltage:      req.Voltage,
		Amperage:     req.Amperage,
		BatteryType:  batteryType,
		Polarity:     polarity,
		AcidLiters:   req.AcidLiters,
		VehicleType:  vehicleType,
		MinSalePrice: req.MinSalePrice,
		Active:       true,
		CreatedAt:    time.Now(),
		CreatedBy:    req.CreatedBy,
	}

	// Validate entity
	if err := product.IsValid(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// GetByID retrieves a product by ID.
func (s *productService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return product, nil
}

// List retrieves products with filtering.
func (s *productService) List(ctx context.Context, filter ports.ProductFilter) ([]domain.Product, error) {
	products, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return products, nil
}

// Update updates a product's information.
func (s *productService) Update(ctx context.Context, id uuid.UUID, req ports.UpdateProductRequest) (*domain.Product, error) {
	// Find existing product
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Apply updates
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.Brand != nil {
		product.Brand = req.Brand
	}
	if req.Model != nil {
		product.Model = req.Model
	}
	if req.Voltage != nil {
		product.Voltage = req.Voltage
	}
	if req.Amperage != nil {
		product.Amperage = req.Amperage
	}
	if req.BatteryType != nil {
		bt := domain.BatteryType(*req.BatteryType)
		product.BatteryType = &bt
	}
	if req.Polarity != nil {
		p := domain.Polarity(*req.Polarity)
		product.Polarity = &p
	}
	if req.AcidLiters != nil {
		product.AcidLiters = req.AcidLiters
	}
	if req.VehicleType != nil {
		vt := domain.VehicleType(*req.VehicleType)
		product.VehicleType = &vt
	}
	if req.MinSalePrice != nil {
		// Store current price as previous before updating
		product.PreviousPrice = &product.MinSalePrice
		product.EffectiveDate = ptr(time.Now())
		product.MinSalePrice = *req.MinSalePrice
	}
	if req.Active != nil {
		product.Active = *req.Active
	}

	// Validate entity
	if err := product.IsValid(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// Deactivate deactivates a product.
func (s *productService) Deactivate(ctx context.Context, id uuid.UUID) error {
	// Find existing product
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	product.Active = false

	// Save to repository
	if err := s.productRepo.Update(ctx, product); err != nil {
		return fmt.Errorf("failed to deactivate product: %w", err)
	}

	return nil
}

// Helper function to create a pointer to time.Time
func ptr(t time.Time) *time.Time {
	return &t
}
