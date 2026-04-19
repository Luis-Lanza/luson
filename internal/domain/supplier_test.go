package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupplier_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		supplier Supplier
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid supplier with all fields",
			supplier: Supplier{
				ID:        uuid.New(),
				Name:      "Baterías Bolivia S.A.",
				Contact:   ptrString("contacto@baterias.com"),
				Address:   ptrString("Av. Industrial 456, Zona Sur"),
				Active:    true,
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid supplier with only name",
			supplier: Supplier{
				ID:        uuid.New(),
				Name:      "Proveedor Minimo",
				Contact:   nil,
				Address:   nil,
				Active:    true,
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid - empty name",
			supplier: Supplier{
				ID:        uuid.New(),
				Name:      "",
				Contact:   ptrString("contacto@test.com"),
				Active:    true,
				CreatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "inactive supplier",
			supplier: Supplier{
				ID:        uuid.New(),
				Name:      "Proveedor Inactivo",
				Contact:   ptrString("viejo@proveedor.com"),
				Active:    false,
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.supplier.IsValid()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSupplier_ActiveStatus(t *testing.T) {
	t.Run("active supplier", func(t *testing.T) {
		s := Supplier{Active: true}
		assert.True(t, s.Active)
	})

	t.Run("inactive supplier", func(t *testing.T) {
		s := Supplier{Active: false}
		assert.False(t, s.Active)
	})
}

func TestSupplier_OptionalFields(t *testing.T) {
	t.Run("contact is optional", func(t *testing.T) {
		s := Supplier{
			ID:      uuid.New(),
			Name:    "Test Supplier",
			Contact: nil,
			Active:  true,
		}
		assert.Nil(t, s.Contact)
		assert.NoError(t, s.IsValid())
	})

	t.Run("address is optional", func(t *testing.T) {
		s := Supplier{
			ID:      uuid.New(),
			Name:    "Test Supplier",
			Address: nil,
			Active:  true,
		}
		assert.Nil(t, s.Address)
		assert.NoError(t, s.IsValid())
	})

	t.Run("with contact and address", func(t *testing.T) {
		contact := "+591 77777777"
		address := "Av. Test 123"
		s := Supplier{
			ID:      uuid.New(),
			Name:    "Test Supplier",
			Contact: &contact,
			Address: &address,
			Active:  true,
		}
		assert.Equal(t, "+591 77777777", *s.Contact)
		assert.Equal(t, "Av. Test 123", *s.Address)
	})
}

// Helper function to create pointer to string
func ptrString(s string) *string {
	return &s
}
