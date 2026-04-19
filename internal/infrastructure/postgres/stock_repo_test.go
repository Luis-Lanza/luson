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

func TestStockRepository_Create(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "stock")
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewStockRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_stock",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	branchRepo := NewBranchRepository(db)
	branch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Test Branch Stock",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, branch)

	productRepo := NewProductRepository(db)
	product := &domain.Product{
		ID:           uuid.New(),
		Name:         "Stock Test Product",
		ProductType:  domain.ProductTypeAccesorio,
		MinSalePrice: 50.00,
		Active:       true,
		CreatedAt:    time.Now(),
		CreatedBy:    user.ID,
	}
	productRepo.Create(ctx, product)

	t.Run("creates stock entry", func(t *testing.T) {
		stock := &domain.Stock{
			ID:            uuid.New(),
			ProductID:     product.ID,
			ProductType:   domain.ProductTypeAccesorio,
			LocationType:  "branch",
			LocationID:    branch.ID,
			Quantity:      100,
			MinStockAlert: 10,
			UpdatedAt:     time.Now(),
		}

		err := repo.Create(ctx, stock)
		require.NoError(t, err)

		// Verify
		found, err := repo.FindByID(ctx, stock.ID)
		require.NoError(t, err)
		assert.Equal(t, stock.ProductID, found.ProductID)
		assert.Equal(t, stock.LocationID, found.LocationID)
		assert.Equal(t, 100, found.Quantity)
	})

	t.Run("fails with duplicate product and location", func(t *testing.T) {
		stock := &domain.Stock{
			ID:            uuid.New(),
			ProductID:     product.ID,
			ProductType:   domain.ProductTypeAccesorio,
			LocationType:  "branch",
			LocationID:    branch.ID,
			Quantity:      50,
			MinStockAlert: 5,
			UpdatedAt:     time.Now(),
		}

		err := repo.Create(ctx, stock)
		assert.Error(t, err) // Should fail due to unique constraint
	})
}

func TestStockRepository_FindByID(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "stock")
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewStockRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_find_stock",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	branchRepo := NewBranchRepository(db)
	branch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Test Branch Find Stock",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, branch)

	productRepo := NewProductRepository(db)
	product := &domain.Product{
		ID:           uuid.New(),
		Name:         "Find Stock Product",
		ProductType:  domain.ProductTypeAccesorio,
		MinSalePrice: 50.00,
		Active:       true,
		CreatedAt:    time.Now(),
		CreatedBy:    user.ID,
	}
	productRepo.Create(ctx, product)

	t.Run("finds existing stock", func(t *testing.T) {
		stock := &domain.Stock{
			ID:            uuid.New(),
			ProductID:     product.ID,
			ProductType:   domain.ProductTypeAccesorio,
			LocationType:  "branch",
			LocationID:    branch.ID,
			Quantity:      75,
			MinStockAlert: 15,
			UpdatedAt:     time.Now(),
		}

		err := repo.Create(ctx, stock)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, stock.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, stock.ID, found.ID)
		assert.Equal(t, 75, found.Quantity)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestStockRepository_FindByProductAndLocation(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "stock")
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewStockRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_pl_stock",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	branchRepo := NewBranchRepository(db)
	branch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Test Branch PL Stock",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, branch)

	productRepo := NewProductRepository(db)
	product := &domain.Product{
		ID:           uuid.New(),
		Name:         "PL Stock Product",
		ProductType:  domain.ProductTypeAccesorio,
		MinSalePrice: 50.00,
		Active:       true,
		CreatedAt:    time.Now(),
		CreatedBy:    user.ID,
	}
	productRepo.Create(ctx, product)

	t.Run("finds stock by product and location", func(t *testing.T) {
		stock := &domain.Stock{
			ID:            uuid.New(),
			ProductID:     product.ID,
			ProductType:   domain.ProductTypeAccesorio,
			LocationType:  "branch",
			LocationID:    branch.ID,
			Quantity:      200,
			MinStockAlert: 20,
			UpdatedAt:     time.Now(),
		}

		err := repo.Create(ctx, stock)
		require.NoError(t, err)

		found, err := repo.FindByProductAndLocation(ctx, product.ID, "branch", branch.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, stock.ID, found.ID)
		assert.Equal(t, 200, found.Quantity)
	})

	t.Run("returns nil for non-existent combination", func(t *testing.T) {
		found, err := repo.FindByProductAndLocation(ctx, uuid.New(), "branch", branch.ID)
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestStockRepository_Update(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "stock")
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewStockRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_update_stock",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	branchRepo := NewBranchRepository(db)
	branch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Test Branch Update Stock",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, branch)

	productRepo := NewProductRepository(db)
	product := &domain.Product{
		ID:           uuid.New(),
		Name:         "Update Stock Product",
		ProductType:  domain.ProductTypeAccesorio,
		MinSalePrice: 50.00,
		Active:       true,
		CreatedAt:    time.Now(),
		CreatedBy:    user.ID,
	}
	productRepo.Create(ctx, product)

	t.Run("updates stock quantity", func(t *testing.T) {
		stock := &domain.Stock{
			ID:            uuid.New(),
			ProductID:     product.ID,
			ProductType:   domain.ProductTypeAccesorio,
			LocationType:  "branch",
			LocationID:    branch.ID,
			Quantity:      100,
			MinStockAlert: 10,
			UpdatedAt:     time.Now(),
		}

		err := repo.Create(ctx, stock)
		require.NoError(t, err)

		// Update
		stock.Quantity = 50
		stock.MinStockAlert = 5

		err = repo.Update(ctx, stock)
		require.NoError(t, err)

		// Verify
		found, err := repo.FindByID(ctx, stock.ID)
		require.NoError(t, err)
		assert.Equal(t, 50, found.Quantity)
		assert.Equal(t, 5, found.MinStockAlert)
	})
}

func TestStockRepository_Delete(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "stock")
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewStockRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_delete_stock",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	branchRepo := NewBranchRepository(db)
	branch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Test Branch Delete Stock",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, branch)

	productRepo := NewProductRepository(db)
	product := &domain.Product{
		ID:           uuid.New(),
		Name:         "Delete Stock Product",
		ProductType:  domain.ProductTypeAccesorio,
		MinSalePrice: 50.00,
		Active:       true,
		CreatedAt:    time.Now(),
		CreatedBy:    user.ID,
	}
	productRepo.Create(ctx, product)

	t.Run("deletes stock entry", func(t *testing.T) {
		stock := &domain.Stock{
			ID:            uuid.New(),
			ProductID:     product.ID,
			ProductType:   domain.ProductTypeAccesorio,
			LocationType:  "branch",
			LocationID:    branch.ID,
			Quantity:      30,
			MinStockAlert: 3,
			UpdatedAt:     time.Now(),
		}

		err := repo.Create(ctx, stock)
		require.NoError(t, err)

		err = repo.Delete(ctx, stock.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.FindByID(ctx, stock.ID)
		assert.Error(t, err)
	})

	t.Run("returns error for non-existent stock", func(t *testing.T) {
		err := repo.Delete(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestStockRepository_ListByLocation(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "stock")
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewStockRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_list_stock",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	branchRepo := NewBranchRepository(db)
	branch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Test Branch List Stock",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, branch)

	productRepo := NewProductRepository(db)

	// Create products and stock
	products := []domain.Product{
		{ID: uuid.New(), Name: "Stock Product 1", ProductType: domain.ProductTypeAccesorio, MinSalePrice: 50.00, Active: true, CreatedAt: time.Now(), CreatedBy: user.ID},
		{ID: uuid.New(), Name: "Stock Product 2", ProductType: domain.ProductTypeAccesorio, MinSalePrice: 60.00, Active: true, CreatedAt: time.Now(), CreatedBy: user.ID},
	}

	for i := range products {
		productRepo.Create(ctx, &products[i])
	}

	stockItems := []domain.Stock{
		{ID: uuid.New(), ProductID: products[0].ID, ProductType: domain.ProductTypeAccesorio, LocationType: "branch", LocationID: branch.ID, Quantity: 100, MinStockAlert: 10, UpdatedAt: time.Now()},
		{ID: uuid.New(), ProductID: products[1].ID, ProductType: domain.ProductTypeAccesorio, LocationType: "branch", LocationID: branch.ID, Quantity: 5, MinStockAlert: 10, UpdatedAt: time.Now()}, // Low stock
	}

	for i := range stockItems {
		repo.Create(ctx, &stockItems[i])
	}

	t.Run("lists all stock for location", func(t *testing.T) {
		result, err := repo.ListByLocation(ctx, "branch", branch.ID, ports.StockFilter{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 2)
	})

	t.Run("filters low stock only", func(t *testing.T) {
		lowStock := true
		result, err := repo.ListByLocation(ctx, "branch", branch.ID, ports.StockFilter{LowStockOnly: &lowStock})
		require.NoError(t, err)
		for _, s := range result {
			assert.True(t, s.IsLowStock())
		}
	})
}

func TestStockRepository_ListByProduct(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "stock")
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewStockRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_listp_stock",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	branchRepo := NewBranchRepository(db)
	branch1 := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Test Branch 1",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branch2 := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Test Branch 2",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, branch1)
	branchRepo.Create(ctx, branch2)

	productRepo := NewProductRepository(db)
	product := &domain.Product{
		ID:           uuid.New(),
		Name:         "Multi Location Product",
		ProductType:  domain.ProductTypeAccesorio,
		MinSalePrice: 50.00,
		Active:       true,
		CreatedAt:    time.Now(),
		CreatedBy:    user.ID,
	}
	productRepo.Create(ctx, product)

	// Create stock at multiple locations
	stockItems := []domain.Stock{
		{ID: uuid.New(), ProductID: product.ID, ProductType: domain.ProductTypeAccesorio, LocationType: "branch", LocationID: branch1.ID, Quantity: 50, MinStockAlert: 5, UpdatedAt: time.Now()},
		{ID: uuid.New(), ProductID: product.ID, ProductType: domain.ProductTypeAccesorio, LocationType: "branch", LocationID: branch2.ID, Quantity: 30, MinStockAlert: 5, UpdatedAt: time.Now()},
	}

	for i := range stockItems {
		repo.Create(ctx, &stockItems[i])
	}

	t.Run("lists stock across all locations for product", func(t *testing.T) {
		result, err := repo.ListByProduct(ctx, product.ID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 2)
	})
}
