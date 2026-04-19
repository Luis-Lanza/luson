package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupplierRepository_FindByID(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "suppliers")

	repo := NewSupplierRepository(db)
	ctx := context.Background()

	t.Run("finds supplier by ID", func(t *testing.T) {
		contact := "+591 77777777"
		address := "Av. Test 123"
		supplier := &domain.Supplier{
			ID:        uuid.New(),
			Name:      "Test Supplier",
			Contact:   &contact,
			Address:   &address,
			Active:    true,
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, supplier)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, supplier.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Test Supplier", found.Name)
		assert.Equal(t, &contact, found.Contact)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestSupplierRepository_Create(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "suppliers")

	repo := NewSupplierRepository(db)
	ctx := context.Background()

	t.Run("creates supplier with all fields", func(t *testing.T) {
		contact := "contact@test.com"
		address := "Test Address"
		supplier := &domain.Supplier{
			ID:        uuid.New(),
			Name:      "New Supplier",
			Contact:   &contact,
			Address:   &address,
			Active:    true,
			CreatedAt: time.Now(),
		}

		err := repo.Create(ctx, supplier)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, supplier.ID)
		require.NoError(t, err)
		assert.Equal(t, "New Supplier", found.Name)
	})

	t.Run("creates supplier with only required fields", func(t *testing.T) {
		supplier := &domain.Supplier{
			ID:        uuid.New(),
			Name:      "Minimal Supplier",
			Active:    true,
			CreatedAt: time.Now(),
		}

		err := repo.Create(ctx, supplier)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, supplier.ID)
		require.NoError(t, err)
		assert.Equal(t, "Minimal Supplier", found.Name)
		assert.Nil(t, found.Contact)
		assert.Nil(t, found.Address)
	})
}

func TestSupplierRepository_Update(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "suppliers")

	repo := NewSupplierRepository(db)
	ctx := context.Background()

	t.Run("updates supplier fields", func(t *testing.T) {
		supplier := &domain.Supplier{
			ID:        uuid.New(),
			Name:      "Original Name",
			Active:    true,
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, supplier)
		require.NoError(t, err)

		// Update fields
		supplier.Name = "Updated Name"
		supplier.Active = false
		newContact := "new@contact.com"
		supplier.Contact = &newContact

		err = repo.Update(ctx, supplier)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, supplier.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
		assert.False(t, found.Active)
		assert.Equal(t, &newContact, found.Contact)
	})
}

func TestSupplierRepository_List(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "suppliers")

	repo := NewSupplierRepository(db)
	ctx := context.Background()

	// Create test suppliers
	suppliers := []domain.Supplier{
		{ID: uuid.New(), Name: "Supplier 1", Active: true, CreatedAt: time.Now()},
		{ID: uuid.New(), Name: "Supplier 2", Active: false, CreatedAt: time.Now()},
		{ID: uuid.New(), Name: "Supplier 3", Active: true, CreatedAt: time.Now()},
	}

	for i := range suppliers {
		err := repo.Create(ctx, &suppliers[i])
		require.NoError(t, err)
	}

	t.Run("lists all suppliers", func(t *testing.T) {
		result, err := repo.List(ctx, ports.SupplierFilter{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 3)
	})

	t.Run("filters by active status", func(t *testing.T) {
		active := true
		result, err := repo.List(ctx, ports.SupplierFilter{Active: &active})
		require.NoError(t, err)
		for _, s := range result {
			assert.True(t, s.Active)
		}
	})

	t.Run("respects pagination", func(t *testing.T) {
		result, err := repo.List(ctx, ports.SupplierFilter{Limit: 2, Offset: 0})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(result), 2)
	})
}
