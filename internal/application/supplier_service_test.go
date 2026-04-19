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

// MockSupplierRepository is a mock implementation of ports.SupplierRepository
type MockSupplierRepository struct {
	mock.Mock
}

func (m *MockSupplierRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Supplier), args.Error(1)
}

func (m *MockSupplierRepository) List(ctx context.Context, filter ports.SupplierFilter) ([]domain.Supplier, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]domain.Supplier), args.Error(1)
}

func (m *MockSupplierRepository) Create(ctx context.Context, supplier *domain.Supplier) error {
	args := m.Called(ctx, supplier)
	return args.Error(0)
}

func (m *MockSupplierRepository) Update(ctx context.Context, supplier *domain.Supplier) error {
	args := m.Called(ctx, supplier)
	return args.Error(0)
}

func TestSupplierService_Create(t *testing.T) {
	supplierRepo := new(MockSupplierRepository)
	service := NewSupplierService(supplierRepo)
	ctx := context.Background()

	t.Run("successful supplier creation", func(t *testing.T) {
		contact := "John Doe"
		address := "123 Supply St"
		req := ports.CreateSupplierRequest{
			Name:    "Acme Supplies",
			Contact: &contact,
			Address: &address,
		}

		supplierRepo.On("Create", ctx, mock.AnythingOfType("*domain.Supplier")).Return(nil).Once()

		supplier, err := service.Create(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, supplier)
		assert.Equal(t, req.Name, supplier.Name)
		assert.Equal(t, &contact, supplier.Contact)
		assert.Equal(t, &address, supplier.Address)
		assert.True(t, supplier.Active)
		supplierRepo.AssertExpectations(t)
	})

	t.Run("successful creation without optional fields", func(t *testing.T) {
		req := ports.CreateSupplierRequest{
			Name: "Simple Supplier",
		}

		supplierRepo.On("Create", ctx, mock.AnythingOfType("*domain.Supplier")).Return(nil).Once()

		supplier, err := service.Create(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, supplier)
		assert.Equal(t, req.Name, supplier.Name)
		assert.Nil(t, supplier.Contact)
		assert.Nil(t, supplier.Address)
		supplierRepo.AssertExpectations(t)
	})

	t.Run("fails with empty name", func(t *testing.T) {
		req := ports.CreateSupplierRequest{
			Name: "",
		}

		supplier, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, supplier)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("fails when repository returns error", func(t *testing.T) {
		req := ports.CreateSupplierRequest{
			Name: "Test Supplier",
		}

		supplierRepo.On("Create", ctx, mock.AnythingOfType("*domain.Supplier")).Return(errors.New("db error")).Once()

		supplier, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, supplier)
		supplierRepo.AssertExpectations(t)
	})
}

func TestSupplierService_GetByID(t *testing.T) {
	supplierRepo := new(MockSupplierRepository)
	service := NewSupplierService(supplierRepo)
	ctx := context.Background()

	t.Run("returns supplier by id", func(t *testing.T) {
		supplierID := uuid.New()
		supplier := &domain.Supplier{
			ID:     supplierID,
			Name:   "Test Supplier",
			Active: true,
		}

		supplierRepo.On("FindByID", ctx, supplierID).Return(supplier, nil).Once()

		result, err := service.GetByID(ctx, supplierID)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, supplierID, result.ID)
		assert.Equal(t, "Test Supplier", result.Name)
		supplierRepo.AssertExpectations(t)
	})

	t.Run("returns error for non-existent supplier", func(t *testing.T) {
		supplierID := uuid.New()

		supplierRepo.On("FindByID", ctx, supplierID).Return(nil, errors.New("supplier not found")).Once()

		result, err := service.GetByID(ctx, supplierID)

		assert.Error(t, err)
		assert.Nil(t, result)
		supplierRepo.AssertExpectations(t)
	})
}

func TestSupplierService_List(t *testing.T) {
	supplierRepo := new(MockSupplierRepository)
	service := NewSupplierService(supplierRepo)
	ctx := context.Background()

	t.Run("returns list of suppliers", func(t *testing.T) {
		active := true
		filter := ports.SupplierFilter{
			Active: &active,
			Limit:  10,
			Offset: 0,
		}

		suppliers := []domain.Supplier{
			{ID: uuid.New(), Name: "Supplier 1", Active: true},
			{ID: uuid.New(), Name: "Supplier 2", Active: true},
		}

		supplierRepo.On("List", ctx, filter).Return(suppliers, nil).Once()

		result, err := service.List(ctx, filter)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Supplier 1", result[0].Name)
		supplierRepo.AssertExpectations(t)
	})

	t.Run("returns empty list when no suppliers", func(t *testing.T) {
		filter := ports.SupplierFilter{}

		supplierRepo.On("List", ctx, filter).Return([]domain.Supplier{}, nil).Once()

		result, err := service.List(ctx, filter)

		require.NoError(t, err)
		assert.Empty(t, result)
		supplierRepo.AssertExpectations(t)
	})
}

func TestSupplierService_Update(t *testing.T) {
	supplierRepo := new(MockSupplierRepository)
	service := NewSupplierService(supplierRepo)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		supplierID := uuid.New()
		newName := "Updated Name"
		newContact := "Updated Contact"
		newAddress := "Updated Address"
		newActive := false

		existingSupplier := &domain.Supplier{
			ID:      supplierID,
			Name:    "Original Name",
			Contact: strPtr("Original Contact"),
			Address: strPtr("Original Address"),
			Active:  true,
		}

		req := ports.UpdateSupplierRequest{
			Name:    &newName,
			Contact: &newContact,
			Address: &newAddress,
			Active:  &newActive,
		}

		supplierRepo.On("FindByID", ctx, supplierID).Return(existingSupplier, nil).Once()
		supplierRepo.On("Update", ctx, mock.AnythingOfType("*domain.Supplier")).Return(nil).Once()

		result, err := service.Update(ctx, supplierID, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, newName, result.Name)
		assert.Equal(t, &newContact, result.Contact)
		assert.Equal(t, &newAddress, result.Address)
		assert.False(t, result.Active)
		supplierRepo.AssertExpectations(t)
	})

	t.Run("partial update with only name", func(t *testing.T) {
		supplierID := uuid.New()
		newName := "New Name Only"

		existingSupplier := &domain.Supplier{
			ID:      supplierID,
			Name:    "Original Name",
			Contact: strPtr("Original Contact"),
			Active:  true,
		}

		req := ports.UpdateSupplierRequest{
			Name: &newName,
		}

		supplierRepo.On("FindByID", ctx, supplierID).Return(existingSupplier, nil).Once()
		supplierRepo.On("Update", ctx, mock.AnythingOfType("*domain.Supplier")).Return(nil).Once()

		result, err := service.Update(ctx, supplierID, req)

		require.NoError(t, err)
		assert.Equal(t, newName, result.Name)
		assert.Equal(t, "Original Contact", *result.Contact) // Unchanged
		supplierRepo.AssertExpectations(t)
	})

	t.Run("can clear optional fields with empty string", func(t *testing.T) {
		supplierID := uuid.New()
		emptyContact := ""

		existingSupplier := &domain.Supplier{
			ID:      supplierID,
			Name:    "Original Name",
			Contact: strPtr("Original Contact"),
			Active:  true,
		}

		req := ports.UpdateSupplierRequest{
			Contact: &emptyContact,
		}

		supplierRepo.On("FindByID", ctx, supplierID).Return(existingSupplier, nil).Once()
		supplierRepo.On("Update", ctx, mock.AnythingOfType("*domain.Supplier")).Return(nil).Once()

		result, err := service.Update(ctx, supplierID, req)

		require.NoError(t, err)
		// When pointer is set to empty string, we update it
		assert.Equal(t, &emptyContact, result.Contact)
		supplierRepo.AssertExpectations(t)
	})

	t.Run("fails when supplier not found", func(t *testing.T) {
		supplierID := uuid.New()
		newName := "New Name"
		req := ports.UpdateSupplierRequest{Name: &newName}

		supplierRepo.On("FindByID", ctx, supplierID).Return(nil, errors.New("supplier not found")).Once()

		result, err := service.Update(ctx, supplierID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		supplierRepo.AssertExpectations(t)
	})

	t.Run("fails with empty name", func(t *testing.T) {
		supplierID := uuid.New()
		emptyName := ""

		existingSupplier := &domain.Supplier{
			ID:   supplierID,
			Name: "Original Name",
		}

		req := ports.UpdateSupplierRequest{Name: &emptyName}

		supplierRepo.On("FindByID", ctx, supplierID).Return(existingSupplier, nil).Once()

		result, err := service.Update(ctx, supplierID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name is required")
		supplierRepo.AssertExpectations(t)
	})
}

func TestSupplierService_Deactivate(t *testing.T) {
	supplierRepo := new(MockSupplierRepository)
	service := NewSupplierService(supplierRepo)
	ctx := context.Background()

	t.Run("successful deactivation", func(t *testing.T) {
		supplierID := uuid.New()
		existingSupplier := &domain.Supplier{
			ID:     supplierID,
			Name:   "Test Supplier",
			Active: true,
		}

		supplierRepo.On("FindByID", ctx, supplierID).Return(existingSupplier, nil).Once()
		supplierRepo.On("Update", ctx, mock.AnythingOfType("*domain.Supplier")).Return(nil).Once()

		err := service.Deactivate(ctx, supplierID)

		require.NoError(t, err)
		supplierRepo.AssertExpectations(t)
	})

	t.Run("fails when supplier not found", func(t *testing.T) {
		supplierID := uuid.New()

		supplierRepo.On("FindByID", ctx, supplierID).Return(nil, errors.New("supplier not found")).Once()

		err := service.Deactivate(ctx, supplierID)

		assert.Error(t, err)
		supplierRepo.AssertExpectations(t)
	})
}

// Helper function
func strPtr(s string) *string {
	return &s
}
