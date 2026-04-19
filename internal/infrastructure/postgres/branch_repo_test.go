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

func TestBranchRepository_FindByID(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "branches")

	repo := NewBranchRepository(db)
	ctx := context.Background()

	t.Run("finds branch by ID", func(t *testing.T) {
		branch := &domain.Branch{
			ID:               uuid.New(),
			Name:             "Test Branch",
			Address:          "Test Address 123",
			PettyCashBalance: 1000.00,
			Active:           true,
			CreatedAt:        time.Now(),
		}
		err := repo.Create(ctx, branch)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, branch.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Test Branch", found.Name)
		assert.Equal(t, 1000.00, found.PettyCashBalance)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestBranchRepository_Create(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "branches")

	repo := NewBranchRepository(db)
	ctx := context.Background()

	t.Run("creates branch with all fields", func(t *testing.T) {
		branch := &domain.Branch{
			ID:               uuid.New(),
			Name:             "New Branch",
			Address:          "New Address",
			PettyCashBalance: 500.50,
			Active:           true,
			CreatedAt:        time.Now(),
		}

		err := repo.Create(ctx, branch)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, branch.ID)
		require.NoError(t, err)
		assert.Equal(t, "New Branch", found.Name)
	})
}

func TestBranchRepository_Update(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "branches")

	repo := NewBranchRepository(db)
	ctx := context.Background()

	t.Run("updates branch fields", func(t *testing.T) {
		branch := &domain.Branch{
			ID:               uuid.New(),
			Name:             "Original Name",
			Address:          "Original Address",
			PettyCashBalance: 1000.00,
			Active:           true,
			CreatedAt:        time.Now(),
		}
		err := repo.Create(ctx, branch)
		require.NoError(t, err)

		// Update fields
		branch.Name = "Updated Name"
		branch.PettyCashBalance = 2000.00
		branch.Active = false

		err = repo.Update(ctx, branch)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, branch.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
		assert.Equal(t, 2000.00, found.PettyCashBalance)
		assert.False(t, found.Active)
	})
}

func TestBranchRepository_List(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "branches")

	repo := NewBranchRepository(db)
	ctx := context.Background()

	// Create test branches
	branches := []domain.Branch{
		{ID: uuid.New(), Name: "Branch 1", Address: "Addr 1", PettyCashBalance: 100, Active: true, CreatedAt: time.Now()},
		{ID: uuid.New(), Name: "Branch 2", Address: "Addr 2", PettyCashBalance: 200, Active: false, CreatedAt: time.Now()},
		{ID: uuid.New(), Name: "Branch 3", Address: "Addr 3", PettyCashBalance: 300, Active: true, CreatedAt: time.Now()},
	}

	for i := range branches {
		err := repo.Create(ctx, &branches[i])
		require.NoError(t, err)
	}

	t.Run("lists all branches", func(t *testing.T) {
		result, err := repo.List(ctx, ports.BranchFilter{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 3)
	})

	t.Run("filters by active status", func(t *testing.T) {
		active := true
		result, err := repo.List(ctx, ports.BranchFilter{Active: &active})
		require.NoError(t, err)
		for _, b := range result {
			assert.True(t, b.Active)
		}
	})

	t.Run("respects pagination", func(t *testing.T) {
		result, err := repo.List(ctx, ports.BranchFilter{Limit: 2, Offset: 0})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(result), 2)
	})
}
