package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStock_IsValid(t *testing.T) {
	productID := uuid.New()
	locationID := uuid.New()

	tests := []struct {
		name    string
		stock   Stock
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid stock",
			stock: Stock{
				ID:            uuid.New(),
				ProductID:     productID,
				ProductType:   ProductTypeBateria,
				LocationType:  "branch",
				LocationID:    locationID,
				Quantity:      10,
				MinStockAlert: 5,
				UpdatedAt:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid stock - zero quantity",
			stock: Stock{
				ID:            uuid.New(),
				ProductID:     productID,
				ProductType:   ProductTypeAccesorio,
				LocationType:  "warehouse",
				LocationID:    locationID,
				Quantity:      0,
				MinStockAlert: 2,
				UpdatedAt:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid - nil product_id",
			stock: Stock{
				ID:            uuid.New(),
				ProductID:     uuid.Nil,
				ProductType:   ProductTypeBateria,
				LocationType:  "branch",
				LocationID:    locationID,
				Quantity:      10,
				MinStockAlert: 5,
				UpdatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "product_id is required",
		},
		{
			name: "invalid - nil location_id",
			stock: Stock{
				ID:            uuid.New(),
				ProductID:     productID,
				ProductType:   ProductTypeBateria,
				LocationType:  "branch",
				LocationID:    uuid.Nil,
				Quantity:      10,
				MinStockAlert: 5,
				UpdatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "location_id is required",
		},
		{
			name: "invalid - empty location_type",
			stock: Stock{
				ID:            uuid.New(),
				ProductID:     productID,
				ProductType:   ProductTypeBateria,
				LocationType:  "",
				LocationID:    locationID,
				Quantity:      10,
				MinStockAlert: 5,
				UpdatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "location_type is required",
		},
		{
			name: "invalid - negative quantity",
			stock: Stock{
				ID:            uuid.New(),
				ProductID:     productID,
				ProductType:   ProductTypeBateria,
				LocationType:  "branch",
				LocationID:    locationID,
				Quantity:      -5,
				MinStockAlert: 5,
				UpdatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "quantity cannot be negative",
		},
		{
			name: "invalid - negative min_stock_alert",
			stock: Stock{
				ID:            uuid.New(),
				ProductID:     productID,
				ProductType:   ProductTypeBateria,
				LocationType:  "branch",
				LocationID:    locationID,
				Quantity:      10,
				MinStockAlert: -1,
				UpdatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "min_stock_alert cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.stock.IsValid()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStock_IsLowStock(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		minAlert int
		expected bool
	}{
		{"stock above alert", 10, 5, false},
		{"stock at alert level", 5, 5, true},
		{"stock below alert", 3, 5, true},
		{"zero stock", 0, 5, true},
		{"zero alert with stock", 1, 0, false},
		{"zero alert zero stock", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Stock{
				Quantity:      tt.quantity,
				MinStockAlert: tt.minAlert,
			}
			assert.Equal(t, tt.expected, s.IsLowStock())
		})
	}
}

func TestStock_CanFulfill(t *testing.T) {
	s := Stock{
		Quantity: 10,
	}

	tests := []struct {
		name      string
		requested int
		expected  bool
	}{
		{"can fulfill exact", 10, true},
		{"can fulfill less", 5, true},
		{"cannot fulfill more", 11, false},
		{"can fulfill zero", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, s.CanFulfill(tt.requested))
		})
	}
}
