package domain

import (
	"time"

	"github.com/google/uuid"
)

// Stock represents the inventory of a product at a specific location.
type Stock struct {
	ID            uuid.UUID   `json:"id"`
	ProductID     uuid.UUID   `json:"product_id"`
	ProductType   ProductType `json:"product_type"`
	LocationType  string      `json:"location_type"`
	LocationID    uuid.UUID   `json:"location_id"`
	Quantity      int         `json:"quantity"`
	MinStockAlert int         `json:"min_stock_alert"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// IsValid validates the stock entity according to business rules.
func (s Stock) IsValid() error {
	if s.ProductID == uuid.Nil {
		return ErrStockInvalidProductID
	}
	if s.LocationID == uuid.Nil {
		return ErrStockInvalidLocationID
	}
	if s.LocationType == "" {
		return ErrStockInvalidLocationType
	}
	if s.Quantity < 0 {
		return ErrStockNegativeQuantity
	}
	if s.MinStockAlert < 0 {
		return ErrStockNegativeMinStockAlert
	}
	return nil
}

// IsLowStock returns true if the quantity is at or below the minimum stock alert level.
func (s Stock) IsLowStock() bool {
	return s.Quantity <= s.MinStockAlert
}

// CanFulfill returns true if the stock can fulfill the requested quantity.
func (s Stock) CanFulfill(requested int) bool {
	return s.Quantity >= requested
}
