package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranch_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		branch  Branch
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid branch",
			branch: Branch{
				ID:               uuid.New(),
				Name:             "Sucursal Central",
				Address:          "Av. Principal 123",
				PettyCashBalance: 1000.50,
				Active:           true,
				CreatedAt:        time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid branch without address",
			branch: Branch{
				ID:               uuid.New(),
				Name:             "Sucursal Norte",
				Address:          "",
				PettyCashBalance: 0,
				Active:           true,
				CreatedAt:        time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid - empty name",
			branch: Branch{
				ID:               uuid.New(),
				Name:             "",
				Address:          "Av. Principal 123",
				PettyCashBalance: 0,
				Active:           true,
				CreatedAt:        time.Now(),
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "inactive branch",
			branch: Branch{
				ID:               uuid.New(),
				Name:             "Sucursal Cerrada",
				Address:          "Av. Vieja 456",
				PettyCashBalance: 500.00,
				Active:           false,
				CreatedAt:        time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.branch.IsValid()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBranch_ActiveStatus(t *testing.T) {
	t.Run("active branch", func(t *testing.T) {
		b := Branch{Active: true}
		assert.True(t, b.Active)
	})

	t.Run("inactive branch", func(t *testing.T) {
		b := Branch{Active: false}
		assert.False(t, b.Active)
	})
}

func TestBranch_PettyCashBalance(t *testing.T) {
	t.Run("default balance is zero", func(t *testing.T) {
		b := Branch{
			ID:   uuid.New(),
			Name: "Test Branch",
		}
		assert.Equal(t, 0.0, b.PettyCashBalance)
	})

	t.Run("can have positive balance", func(t *testing.T) {
		b := Branch{
			ID:               uuid.New(),
			Name:             "Test Branch",
			PettyCashBalance: 1500.75,
		}
		assert.Equal(t, 1500.75, b.PettyCashBalance)
	})
}

func TestBranch_UpdatePettyCash(t *testing.T) {
	b := Branch{
		ID:               uuid.New(),
		Name:             "Test Branch",
		PettyCashBalance: 1000.00,
		Active:           true,
	}

	// Simulate updating petty cash
	b.PettyCashBalance = 1500.00
	assert.Equal(t, 1500.00, b.PettyCashBalance)
}
