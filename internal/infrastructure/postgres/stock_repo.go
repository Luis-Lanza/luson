package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// stockRepository implements ports.StockRepository using PostgreSQL.
type stockRepository struct {
	db *pgxpool.Pool
}

// NewStockRepository creates a new PostgreSQL stock repository.
func NewStockRepository(db *pgxpool.Pool) ports.StockRepository {
	return &stockRepository{db: db}
}

// FindByID finds a stock entry by its ID.
func (r *stockRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
	var stock domain.Stock

	query := `
		SELECT id, product_id, product_type, location_type, location_id, quantity, min_stock_alert, updated_at
		FROM stock
		WHERE id = $1
	`

	err := pgxscan.Get(ctx, r.db, &stock, query, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, fmt.Errorf("stock not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find stock by id: %w", err)
	}

	return &stock, nil
}

// FindByProductAndLocation finds stock for a specific product at a specific location.
func (r *stockRepository) FindByProductAndLocation(ctx context.Context, productID uuid.UUID, locationType string, locationID uuid.UUID) (*domain.Stock, error) {
	var stock domain.Stock

	query := `
		SELECT id, product_id, product_type, location_type, location_id, quantity, min_stock_alert, updated_at
		FROM stock
		WHERE product_id = $1 AND location_type = $2 AND location_id = $3
	`

	err := pgxscan.Get(ctx, r.db, &stock, query, productID, locationType, locationID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find stock by product and location: %w", err)
	}

	return &stock, nil
}

// Create inserts a new stock entry into the database.
func (r *stockRepository) Create(ctx context.Context, stock *domain.Stock) error {
	query := `
		INSERT INTO stock (id, product_id, product_type, location_type, location_id, quantity, min_stock_alert, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		stock.ID,
		stock.ProductID,
		stock.ProductType,
		stock.LocationType,
		stock.LocationID,
		stock.Quantity,
		stock.MinStockAlert,
		stock.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create stock: %w", err)
	}

	return nil
}

// Update modifies an existing stock entry.
func (r *stockRepository) Update(ctx context.Context, stock *domain.Stock) error {
	query := `
		UPDATE stock
		SET quantity = $2, min_stock_alert = $3, updated_at = $4
		WHERE id = $1
	`

	stock.UpdatedAt = time.Now()

	result, err := r.db.Exec(ctx, query,
		stock.ID,
		stock.Quantity,
		stock.MinStockAlert,
		stock.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("stock not found")
	}

	return nil
}

// Delete removes a stock entry from the database.
func (r *stockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM stock WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("stock not found")
	}

	return nil
}

// ListByLocation retrieves stock entries for a specific location.
func (r *stockRepository) ListByLocation(ctx context.Context, locationType string, locationID uuid.UUID, filter ports.StockFilter) ([]domain.Stock, error) {
	query := `
		SELECT id, product_id, product_type, location_type, location_id, quantity, min_stock_alert, updated_at
		FROM stock
		WHERE location_type = $1 AND location_id = $2
	`
	args := []interface{}{locationType, locationID}
	argIdx := 3

	if filter.LowStockOnly != nil && *filter.LowStockOnly {
		query += fmt.Sprintf(" AND quantity <= min_stock_alert")
	}

	query += " ORDER BY updated_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
	}

	var stock []domain.Stock
	err := pgxscan.Select(ctx, r.db, &stock, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock by location: %w", err)
	}

	return stock, nil
}

// ListByProduct retrieves stock entries for a specific product across all locations.
func (r *stockRepository) ListByProduct(ctx context.Context, productID uuid.UUID) ([]domain.Stock, error) {
	query := `
		SELECT id, product_id, product_type, location_type, location_id, quantity, min_stock_alert, updated_at
		FROM stock
		WHERE product_id = $1
		ORDER BY updated_at DESC
	`

	var stock []domain.Stock
	err := pgxscan.Select(ctx, r.db, &stock, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock by product: %w", err)
	}

	return stock, nil
}
