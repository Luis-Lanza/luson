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

// strPtr returns a pointer to a string
func strPtr(s string) *string {
	return &s
}

func TestProductRepository_Create(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")

	repo := NewProductRepository(db)
	ctx := context.Background()

	// Create a user first (required for created_by)
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_product",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("creates battery product", func(t *testing.T) {
		brand := "Yuasa"
		model := "YTX9-BS"
		voltage := 12.0
		amperage := 9.0
		batteryType := domain.BatteryTypeSeca
		polarity := domain.PolarityDerecha
		vehicleType := domain.VehicleTypeAuto

		product := &domain.Product{
			ID:           uuid.New(),
			Name:         "Batería Yuasa YTX9-BS Test",
			Description:  strPtr("Batería de prueba"),
			ProductType:  domain.ProductTypeBateria,
			Brand:        &brand,
			Model:        &model,
			Voltage:      &voltage,
			Amperage:     &amperage,
			BatteryType:  &batteryType,
			Polarity:     &polarity,
			VehicleType:  &vehicleType,
			MinSalePrice: 150.00,
			Active:       true,
			CreatedAt:    time.Now(),
			CreatedBy:    user.ID,
		}

		err := repo.Create(ctx, product)
		require.NoError(t, err)

		// Verify
		found, err := repo.FindByID(ctx, product.ID)
		require.NoError(t, err)
		assert.Equal(t, product.Name, found.Name)
		assert.Equal(t, brand, *found.Brand)
		assert.Equal(t, domain.ProductTypeBateria, found.ProductType)
	})

	t.Run("creates accessory product", func(t *testing.T) {
		product := &domain.Product{
			ID:           uuid.New(),
			Name:         "Cargador de Batería Test",
			Description:  strPtr("Cargador universal"),
			ProductType:  domain.ProductTypeAccesorio,
			MinSalePrice: 45.00,
			Active:       true,
			CreatedAt:    time.Now(),
			CreatedBy:    user.ID,
		}

		err := repo.Create(ctx, product)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, product.ID)
		require.NoError(t, err)
		assert.Equal(t, product.Name, found.Name)
		assert.Equal(t, domain.ProductTypeAccesorio, found.ProductType)
	})

	t.Run("fails with duplicate name", func(t *testing.T) {
		product := &domain.Product{
			ID:           uuid.New(),
			Name:         "Duplicate Product",
			ProductType:  domain.ProductTypeAccesorio,
			MinSalePrice: 100.00,
			Active:       true,
			CreatedAt:    time.Now(),
			CreatedBy:    user.ID,
		}

		err := repo.Create(ctx, product)
		require.NoError(t, err)

		// Try to create another with same name
		product2 := &domain.Product{
			ID:           uuid.New(),
			Name:         "Duplicate Product",
			ProductType:  domain.ProductTypeAccesorio,
			MinSalePrice: 200.00,
			Active:       true,
			CreatedAt:    time.Now(),
			CreatedBy:    user.ID,
		}

		err = repo.Create(ctx, product2)
		assert.Error(t, err) // Should fail due to unique constraint
	})
}

func TestProductRepository_FindByID(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")

	repo := NewProductRepository(db)
	ctx := context.Background()

	// Create user
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_find",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	t.Run("finds existing product", func(t *testing.T) {
		product := &domain.Product{
			ID:           uuid.New(),
			Name:         "Find Test Product",
			ProductType:  domain.ProductTypeAccesorio,
			MinSalePrice: 50.00,
			Active:       true,
			CreatedAt:    time.Now(),
			CreatedBy:    user.ID,
		}

		err := repo.Create(ctx, product)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, product.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, product.Name, found.Name)
		assert.Equal(t, product.ID, found.ID)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestProductRepository_FindByName(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")

	repo := NewProductRepository(db)
	ctx := context.Background()

	// Create user
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_name",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	t.Run("finds by exact name", func(t *testing.T) {
		product := &domain.Product{
			ID:           uuid.New(),
			Name:         "Exact Name Product",
			ProductType:  domain.ProductTypeAccesorio,
			MinSalePrice: 75.00,
			Active:       true,
			CreatedAt:    time.Now(),
			CreatedBy:    user.ID,
		}

		err := repo.Create(ctx, product)
		require.NoError(t, err)

		found, err := repo.FindByName(ctx, "Exact Name Product")
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, product.Name, found.Name)
	})

	t.Run("returns nil for non-existent name", func(t *testing.T) {
		found, err := repo.FindByName(ctx, "Non Existent Product Name XYZ")
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestProductRepository_Update(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")

	repo := NewProductRepository(db)
	ctx := context.Background()

	// Create user
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_update",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	t.Run("updates product fields", func(t *testing.T) {
		product := &domain.Product{
			ID:           uuid.New(),
			Name:         "Update Test Product",
			ProductType:  domain.ProductTypeAccesorio,
			MinSalePrice: 100.00,
			Active:       true,
			CreatedAt:    time.Now(),
			CreatedBy:    user.ID,
		}

		err := repo.Create(ctx, product)
		require.NoError(t, err)

		// Update
		newPrice := 150.00
		product.MinSalePrice = newPrice
		product.Name = "Updated Name"

		err = repo.Update(ctx, product)
		require.NoError(t, err)

		// Verify
		found, err := repo.FindByID(ctx, product.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
		assert.Equal(t, newPrice, found.MinSalePrice)
	})

	t.Run("returns error for non-existent product", func(t *testing.T) {
		product := &domain.Product{
			ID:           uuid.New(),
			Name:         "Non Existent",
			ProductType:  domain.ProductTypeAccesorio,
			MinSalePrice: 50.00,
			Active:       true,
			CreatedAt:    time.Now(),
			CreatedBy:    user.ID,
		}

		err := repo.Update(ctx, product)
		assert.Error(t, err)
	})
}

func TestProductRepository_List(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "products")
	defer cleanupTable(t, "users")

	repo := NewProductRepository(db)
	ctx := context.Background()

	// Create user
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_list",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	// Create test products
	brand := "Yuasa"
	vehicleType := domain.VehicleTypeAuto

	products := []domain.Product{
		{ID: uuid.New(), Name: "Active Battery", ProductType: domain.ProductTypeBateria, Brand: &brand, VehicleType: &vehicleType, BatteryType: (*domain.BatteryType)(strPtr("seca")), Model: strPtr("Model1"), MinSalePrice: 100.00, Active: true, CreatedAt: time.Now(), CreatedBy: user.ID},
		{ID: uuid.New(), Name: "Inactive Battery", ProductType: domain.ProductTypeBateria, Brand: &brand, VehicleType: &vehicleType, BatteryType: (*domain.BatteryType)(strPtr("seca")), Model: strPtr("Model2"), MinSalePrice: 100.00, Active: false, CreatedAt: time.Now(), CreatedBy: user.ID},
		{ID: uuid.New(), Name: "Active Accessory", ProductType: domain.ProductTypeAccesorio, MinSalePrice: 50.00, Active: true, CreatedAt: time.Now(), CreatedBy: user.ID},
	}

	for i := range products {
		err := repo.Create(ctx, &products[i])
		require.NoError(t, err)
	}

	t.Run("lists all products", func(t *testing.T) {
		result, err := repo.List(ctx, ports.ProductFilter{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 3)
	})

	t.Run("filters by product type", func(t *testing.T) {
		prodType := string(domain.ProductTypeBateria)
		result, err := repo.List(ctx, ports.ProductFilter{ProductType: &prodType})
		require.NoError(t, err)
		for _, p := range result {
			assert.Equal(t, domain.ProductTypeBateria, p.ProductType)
		}
	})

	t.Run("filters by active status", func(t *testing.T) {
		active := true
		result, err := repo.List(ctx, ports.ProductFilter{Active: &active})
		require.NoError(t, err)
		for _, p := range result {
			assert.True(t, p.Active)
		}
	})

	t.Run("respects pagination", func(t *testing.T) {
		result, err := repo.List(ctx, ports.ProductFilter{Limit: 2, Offset: 0})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(result), 2)
	})

	t.Run("searches by name", func(t *testing.T) {
		search := "Active Battery"
		result, err := repo.List(ctx, ports.ProductFilter{Search: &search})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 1)
	})
}
