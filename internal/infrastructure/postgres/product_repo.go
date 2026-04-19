package postgres

import (
	"context"
	"fmt"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// productRepository implements ports.ProductRepository using PostgreSQL.
type productRepository struct {
	db *pgxpool.Pool
}

// NewProductRepository creates a new PostgreSQL product repository.
func NewProductRepository(db *pgxpool.Pool) ports.ProductRepository {
	return &productRepository{db: db}
}

// FindByID finds a product by its ID.
func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product

	query := `
		SELECT id, name, description, product_type, brand, model, voltage, amperage,
		       battery_type, polarity, acid_liters, vehicle_type, min_sale_price,
		       effective_date, previous_price, active, created_at, created_by
		FROM products
		WHERE id = $1
	`

	err := pgxscan.Get(ctx, r.db, &product, query, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find product by id: %w", err)
	}

	return &product, nil
}

// FindByName finds a product by its exact name.
func (r *productRepository) FindByName(ctx context.Context, name string) (*domain.Product, error) {
	var product domain.Product

	query := `
		SELECT id, name, description, product_type, brand, model, voltage, amperage,
		       battery_type, polarity, acid_liters, vehicle_type, min_sale_price,
		       effective_date, previous_price, active, created_at, created_by
		FROM products
		WHERE name = $1
	`

	err := pgxscan.Get(ctx, r.db, &product, query, name)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product by name: %w", err)
	}

	return &product, nil
}

// Create inserts a new product into the database.
func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (id, name, description, product_type, brand, model, voltage, amperage,
		                     battery_type, polarity, acid_liters, vehicle_type, min_sale_price,
		                     effective_date, previous_price, active, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := r.db.Exec(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.ProductType,
		product.Brand,
		product.Model,
		product.Voltage,
		product.Amperage,
		product.BatteryType,
		product.Polarity,
		product.AcidLiters,
		product.VehicleType,
		product.MinSalePrice,
		product.EffectiveDate,
		product.PreviousPrice,
		product.Active,
		product.CreatedAt,
		product.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// Update modifies an existing product.
func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET name = $2, description = $3, product_type = $4, brand = $5, model = $6,
		    voltage = $7, amperage = $8, battery_type = $9, polarity = $10,
		    acid_liters = $11, vehicle_type = $12, min_sale_price = $13,
		    effective_date = $14, previous_price = $15, active = $16
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.ProductType,
		product.Brand,
		product.Model,
		product.Voltage,
		product.Amperage,
		product.BatteryType,
		product.Polarity,
		product.AcidLiters,
		product.VehicleType,
		product.MinSalePrice,
		product.EffectiveDate,
		product.PreviousPrice,
		product.Active,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// List retrieves products with optional filtering and pagination.
func (r *productRepository) List(ctx context.Context, filter ports.ProductFilter) ([]domain.Product, error) {
	query := `
		SELECT id, name, description, product_type, brand, model, voltage, amperage,
		       battery_type, polarity, acid_liters, vehicle_type, min_sale_price,
		       effective_date, previous_price, active, created_at, created_by
		FROM products
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filter.ProductType != nil {
		query += fmt.Sprintf(" AND product_type = $%d", argIdx)
		args = append(args, *filter.ProductType)
		argIdx++
	}

	if filter.Brand != nil {
		query += fmt.Sprintf(" AND brand = $%d", argIdx)
		args = append(args, *filter.Brand)
		argIdx++
	}

	if filter.VehicleType != nil {
		query += fmt.Sprintf(" AND vehicle_type = $%d", argIdx)
		args = append(args, *filter.VehicleType)
		argIdx++
	}

	if filter.BatteryType != nil {
		query += fmt.Sprintf(" AND battery_type = $%d", argIdx)
		args = append(args, *filter.BatteryType)
		argIdx++
	}

	if filter.Active != nil {
		query += fmt.Sprintf(" AND active = $%d", argIdx)
		args = append(args, *filter.Active)
		argIdx++
	}

	if filter.Search != nil && *filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d OR brand ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+*filter.Search+"%")
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
	}

	var products []domain.Product
	err := pgxscan.Select(ctx, r.db, &products, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return products, nil
}
