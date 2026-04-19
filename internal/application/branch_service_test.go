package application

import (
	"context"
	"errors"
	"testing"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockBranchRepository is a mock implementation of ports.BranchRepository
type MockBranchRepository struct {
	mock.Mock
}

func (m *MockBranchRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}

func (m *MockBranchRepository) List(ctx context.Context, filter ports.BranchFilter) ([]domain.Branch, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]domain.Branch), args.Error(1)
}

func (m *MockBranchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	args := m.Called(ctx, branch)
	return args.Error(0)
}

func (m *MockBranchRepository) Update(ctx context.Context, branch *domain.Branch) error {
	args := m.Called(ctx, branch)
	return args.Error(0)
}

func TestBranchService_Create(t *testing.T) {
	branchRepo := new(MockBranchRepository)
	service := NewBranchService(branchRepo)
	ctx := context.Background()

	t.Run("successful branch creation", func(t *testing.T) {
		req := ports.CreateBranchRequest{
			Name:             "Main Branch",
			Address:          "123 Main St",
			PettyCashBalance: 1000.00,
		}

		branchRepo.On("Create", ctx, mock.AnythingOfType("*domain.Branch")).Return(nil).Once()

		branch, err := service.Create(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, branch)
		assert.Equal(t, req.Name, branch.Name)
		assert.Equal(t, req.Address, branch.Address)
		assert.Equal(t, req.PettyCashBalance, branch.PettyCashBalance)
		assert.True(t, branch.Active)
		branchRepo.AssertExpectations(t)
	})

	t.Run("fails with empty name", func(t *testing.T) {
		req := ports.CreateBranchRequest{
			Name:    "",
			Address: "123 Main St",
		}

		branch, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, branch)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("fails when repository returns error", func(t *testing.T) {
		req := ports.CreateBranchRequest{
			Name: "Test Branch",
		}

		branchRepo.On("Create", ctx, mock.AnythingOfType("*domain.Branch")).Return(errors.New("db error")).Once()

		branch, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, branch)
		branchRepo.AssertExpectations(t)
	})
}

func TestBranchService_GetByID(t *testing.T) {
	branchRepo := new(MockBranchRepository)
	service := NewBranchService(branchRepo)
	ctx := context.Background()

	t.Run("returns branch by id", func(t *testing.T) {
		branchID := uuid.New()
		branch := &domain.Branch{
			ID:     branchID,
			Name:   "Test Branch",
			Active: true,
		}

		branchRepo.On("FindByID", ctx, branchID).Return(branch, nil).Once()

		result, err := service.GetByID(ctx, branchID)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, branchID, result.ID)
		assert.Equal(t, "Test Branch", result.Name)
		branchRepo.AssertExpectations(t)
	})

	t.Run("returns error for non-existent branch", func(t *testing.T) {
		branchID := uuid.New()

		branchRepo.On("FindByID", ctx, branchID).Return(nil, errors.New("branch not found")).Once()

		result, err := service.GetByID(ctx, branchID)

		assert.Error(t, err)
		assert.Nil(t, result)
		branchRepo.AssertExpectations(t)
	})
}

func TestBranchService_List(t *testing.T) {
	branchRepo := new(MockBranchRepository)
	service := NewBranchService(branchRepo)
	ctx := context.Background()

	t.Run("returns list of branches", func(t *testing.T) {
		active := true
		filter := ports.BranchFilter{
			Active: &active,
			Limit:  10,
			Offset: 0,
		}

		branches := []domain.Branch{
			{ID: uuid.New(), Name: "Branch 1", Active: true},
			{ID: uuid.New(), Name: "Branch 2", Active: true},
		}

		branchRepo.On("List", ctx, filter).Return(branches, nil).Once()

		result, err := service.List(ctx, filter)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Branch 1", result[0].Name)
		branchRepo.AssertExpectations(t)
	})

	t.Run("returns empty list when no branches", func(t *testing.T) {
		filter := ports.BranchFilter{}

		branchRepo.On("List", ctx, filter).Return([]domain.Branch{}, nil).Once()

		result, err := service.List(ctx, filter)

		require.NoError(t, err)
		assert.Empty(t, result)
		branchRepo.AssertExpectations(t)
	})
}

func TestBranchService_Update(t *testing.T) {
	branchRepo := new(MockBranchRepository)
	service := NewBranchService(branchRepo)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		branchID := uuid.New()
		newName := "Updated Name"
		newAddress := "Updated Address"
		newBalance := 2000.00
		newActive := false

		existingBranch := &domain.Branch{
			ID:               branchID,
			Name:             "Original Name",
			Address:          "Original Address",
			PettyCashBalance: 1000.00,
			Active:           true,
		}

		req := ports.UpdateBranchRequest{
			Name:             &newName,
			Address:          &newAddress,
			PettyCashBalance: &newBalance,
			Active:           &newActive,
		}

		branchRepo.On("FindByID", ctx, branchID).Return(existingBranch, nil).Once()
		branchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Branch")).Return(nil).Once()

		result, err := service.Update(ctx, branchID, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, newName, result.Name)
		assert.Equal(t, newAddress, result.Address)
		assert.Equal(t, newBalance, result.PettyCashBalance)
		assert.False(t, result.Active)
		branchRepo.AssertExpectations(t)
	})

	t.Run("partial update with only name", func(t *testing.T) {
		branchID := uuid.New()
		newName := "New Name Only"

		existingBranch := &domain.Branch{
			ID:               branchID,
			Name:             "Original Name",
			Address:          "Original Address",
			PettyCashBalance: 1000.00,
			Active:           true,
		}

		req := ports.UpdateBranchRequest{
			Name: &newName,
		}

		branchRepo.On("FindByID", ctx, branchID).Return(existingBranch, nil).Once()
		branchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Branch")).Return(nil).Once()

		result, err := service.Update(ctx, branchID, req)

		require.NoError(t, err)
		assert.Equal(t, newName, result.Name)
		assert.Equal(t, "Original Address", result.Address) // Unchanged
		assert.Equal(t, 1000.00, result.PettyCashBalance)   // Unchanged
		branchRepo.AssertExpectations(t)
	})

	t.Run("fails when branch not found", func(t *testing.T) {
		branchID := uuid.New()
		newName := "New Name"
		req := ports.UpdateBranchRequest{Name: &newName}

		branchRepo.On("FindByID", ctx, branchID).Return(nil, errors.New("branch not found")).Once()

		result, err := service.Update(ctx, branchID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		branchRepo.AssertExpectations(t)
	})

	t.Run("fails with empty name", func(t *testing.T) {
		branchID := uuid.New()
		emptyName := ""

		existingBranch := &domain.Branch{
			ID:   branchID,
			Name: "Original Name",
		}

		req := ports.UpdateBranchRequest{Name: &emptyName}

		branchRepo.On("FindByID", ctx, branchID).Return(existingBranch, nil).Once()

		result, err := service.Update(ctx, branchID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name is required")
		branchRepo.AssertExpectations(t)
	})
}

func TestBranchService_Deactivate(t *testing.T) {
	branchRepo := new(MockBranchRepository)
	service := NewBranchService(branchRepo)
	ctx := context.Background()

	t.Run("successful deactivation", func(t *testing.T) {
		branchID := uuid.New()
		existingBranch := &domain.Branch{
			ID:     branchID,
			Name:   "Test Branch",
			Active: true,
		}

		branchRepo.On("FindByID", ctx, branchID).Return(existingBranch, nil).Once()
		branchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Branch")).Return(nil).Once()

		err := service.Deactivate(ctx, branchID)

		require.NoError(t, err)
		branchRepo.AssertExpectations(t)
	})

	t.Run("fails when branch not found", func(t *testing.T) {
		branchID := uuid.New()

		branchRepo.On("FindByID", ctx, branchID).Return(nil, errors.New("branch not found")).Once()

		err := service.Deactivate(ctx, branchID)

		assert.Error(t, err)
		branchRepo.AssertExpectations(t)
	})
}
