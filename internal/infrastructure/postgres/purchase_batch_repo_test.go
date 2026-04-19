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

func TestPurchaseBatchRepository_Create(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "purchase_batch_details")
	defer cleanupTable(t, "purchase_batches")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "suppliers")

	repo := NewPurchaseBatchRepository(db)
	ctx := context.Background()

	// Create user
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_batch",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	// Create supplier
	supplierRepo := NewSupplierRepository(db)
	supplier := &domain.Supplier{
		ID:        uuid.New(),
		Name:      "Test Supplier",
		Active:    true,
		CreatedAt: time.Now(),
	}
	supplierRepo.Create(ctx, supplier)

	t.Run("creates purchase batch", func(t *testing.T) {
		invoiceNumber := "INV-001"
		notes := "Test purchase"

		batch := &domain.PurchaseBatch{
			ID:            uuid.New(),
			SupplierID:    &supplier.ID,
			InvoiceNumber: &invoiceNumber,
			PurchaseDate:  time.Now(),
			Notes:         &notes,
			TotalCost:     1000.00,
			Received:      false,
			CreatedBy:     user.ID,
			CreatedAt:     time.Now(),
		}

		err := repo.Create(ctx, batch)
		require.NoError(t, err)

		// Verify
		found, err := repo.FindByID(ctx, batch.ID)
		require.NoError(t, err)
		assert.Equal(t, batch.ID, found.ID)
		assert.Equal(t, "INV-001", *found.InvoiceNumber)
		assert.Equal(t, 1000.00, found.TotalCost)
		assert.False(t, found.Received)
	})

	t.Run("creates purchase batch without supplier", func(t *testing.T) {
		batch := &domain.PurchaseBatch{
			ID:            uuid.New(),
			SupplierID:    nil,
			InvoiceNumber: strPtr("INV-002"),
			PurchaseDate:  time.Now(),
			TotalCost:     500.00,
			Received:      false,
			CreatedBy:     user.ID,
			CreatedAt:     time.Now(),
		}

		err := repo.Create(ctx, batch)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, batch.ID)
		require.NoError(t, err)
		assert.Nil(t, found.SupplierID)
	})
}

func TestPurchaseBatchRepository_CreateWithDetails(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "purchase_batch_details")
	defer cleanupTable(t, "purchase_batches")
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "suppliers")
	defer cleanupTable(t, "branches")

	repo := NewPurchaseBatchRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_batch_details",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	supplierRepo := NewSupplierRepository(db)
	supplier := &domain.Supplier{
		ID:        uuid.New(),
		Name:      "Test Supplier Details",
		Active:    true,
		CreatedAt: time.Now(),
	}
	supplierRepo.Create(ctx, supplier)

	productRepo := NewProductRepository(db)
	product1 := &domain.Product{
		ID:           uuid.New(),
		Name:         "Batch Product 1",
		ProductType:  domain.ProductTypeAccesorio,
		MinSalePrice: 100.00,
		Active:       true,
		CreatedAt:    time.Now(),
		CreatedBy:    user.ID,
	}
	product2 := &domain.Product{
		ID:           uuid.New(),
		Name:         "Batch Product 2",
		ProductType:  domain.ProductTypeAccesorio,
		MinSalePrice: 150.00,
		Active:       true,
		CreatedAt:    time.Now(),
		CreatedBy:    user.ID,
	}
	productRepo.Create(ctx, product1)
	productRepo.Create(ctx, product2)

	t.Run("creates purchase batch with details", func(t *testing.T) {
		batch := &domain.PurchaseBatch{
			ID:            uuid.New(),
			SupplierID:    &supplier.ID,
			InvoiceNumber: strPtr("INV-DETAILS-001"),
			PurchaseDate:  time.Now(),
			TotalCost:     850.00,
			Received:      false,
			CreatedBy:     user.ID,
			CreatedAt:     time.Now(),
		}

		details := []domain.PurchaseBatchDetail{
			{
				ID:              uuid.New(),
				PurchaseBatchID: batch.ID,
				ProductID:       product1.ID,
				Quantity:        5,
				UnitCost:        80.00,
			},
			{
				ID:              uuid.New(),
				PurchaseBatchID: batch.ID,
				ProductID:       product2.ID,
				Quantity:        3,
				UnitCost:        150.00,
			},
		}

		// Use concrete type to access method
		concreteRepo := repo.(*purchaseBatchRepository)
		err := concreteRepo.CreateWithDetails(ctx, batch, details)
		require.NoError(t, err)

		// Verify batch
		foundBatch, err := repo.FindByID(ctx, batch.ID)
		require.NoError(t, err)
		assert.Equal(t, batch.ID, foundBatch.ID)

		// Verify details
		foundDetails, err := concreteRepo.findDetailsByBatchID(ctx, batch.ID)
		require.NoError(t, err)
		assert.Len(t, foundDetails, 2)
	})
}

func TestPurchaseBatchRepository_FindByID(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "purchase_batch_details")
	defer cleanupTable(t, "purchase_batches")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "suppliers")

	repo := NewPurchaseBatchRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_find_batch",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestPurchaseBatchRepository_MarkAsReceived(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "purchase_batch_details")
	defer cleanupTable(t, "purchase_batches")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "suppliers")

	repo := NewPurchaseBatchRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_receive",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	supplierRepo := NewSupplierRepository(db)
	supplier := &domain.Supplier{
		ID:        uuid.New(),
		Name:      "Test Supplier Receive",
		Active:    true,
		CreatedAt: time.Now(),
	}
	supplierRepo.Create(ctx, supplier)

	t.Run("marks batch as received", func(t *testing.T) {
		batch := &domain.PurchaseBatch{
			ID:            uuid.New(),
			SupplierID:    &supplier.ID,
			InvoiceNumber: strPtr("INV-RECV-001"),
			PurchaseDate:  time.Now(),
			TotalCost:     500.00,
			Received:      false,
			CreatedBy:     user.ID,
			CreatedAt:     time.Now(),
		}

		err := repo.Create(ctx, batch)
		require.NoError(t, err)

		err = repo.MarkAsReceived(ctx, batch.ID, user.ID)
		require.NoError(t, err)

		// Verify
		found, err := repo.FindByID(ctx, batch.ID)
		require.NoError(t, err)
		assert.True(t, found.Received)
		assert.NotNil(t, found.ReceivedAt)
		assert.NotNil(t, found.ReceivedBy)
		assert.Equal(t, user.ID, *found.ReceivedBy)
	})

	t.Run("returns error for non-existent batch", func(t *testing.T) {
		err := repo.MarkAsReceived(ctx, uuid.New(), user.ID)
		assert.Error(t, err)
	})
}

func TestPurchaseBatchRepository_List(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "purchase_batch_details")
	defer cleanupTable(t, "purchase_batches")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "suppliers")

	repo := NewPurchaseBatchRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_list_batch",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	supplierRepo := NewSupplierRepository(db)
	supplier := &domain.Supplier{
		ID:        uuid.New(),
		Name:      "Test Supplier List",
		Active:    true,
		CreatedAt: time.Now(),
	}
	supplierRepo.Create(ctx, supplier)

	// Create test batches
	batches := []domain.PurchaseBatch{
		{ID: uuid.New(), SupplierID: &supplier.ID, InvoiceNumber: strPtr("INV-LIST-001"), PurchaseDate: time.Now(), TotalCost: 100.00, Received: false, CreatedBy: user.ID, CreatedAt: time.Now()},
		{ID: uuid.New(), SupplierID: &supplier.ID, InvoiceNumber: strPtr("INV-LIST-002"), PurchaseDate: time.Now(), TotalCost: 200.00, Received: true, CreatedBy: user.ID, CreatedAt: time.Now()},
	}

	for i := range batches {
		repo.Create(ctx, &batches[i])
	}

	t.Run("lists all batches", func(t *testing.T) {
		result, err := repo.List(ctx, ports.PurchaseBatchFilter{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 2)
	})

	t.Run("filters by received status", func(t *testing.T) {
		received := true
		result, err := repo.List(ctx, ports.PurchaseBatchFilter{Received: &received})
		require.NoError(t, err)
		for _, b := range result {
			assert.True(t, b.Received)
		}
	})

	t.Run("filters by supplier", func(t *testing.T) {
		result, err := repo.List(ctx, ports.PurchaseBatchFilter{SupplierID: &supplier.ID})
		require.NoError(t, err)
		for _, b := range result {
			assert.Equal(t, supplier.ID, *b.SupplierID)
		}
	})

	t.Run("respects pagination", func(t *testing.T) {
		result, err := repo.List(ctx, ports.PurchaseBatchFilter{Limit: 1, Offset: 0})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(result), 1)
	})
}
