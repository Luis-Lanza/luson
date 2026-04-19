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

// supplierRepository implements ports.SupplierRepository using PostgreSQL.
type supplierRepository struct {
	db *pgxpool.Pool
}

// NewSupplierRepository creates a new PostgreSQL supplier repository.
func NewSupplierRepository(db *pgxpool.Pool) ports.SupplierRepository {
	return &supplierRepository{db: db}
}

// FindByID finds a supplier by its ID.
func (r *supplierRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	var supplier domain.Supplier

	query := `
		SELECT id, name, contact, address, active, created_at
		FROM suppliers
		WHERE id = $1
	`

	err := pgxscan.Get(ctx, r.db, &supplier, query, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, fmt.Errorf("supplier not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find supplier by id: %w", err)
	}

	return &supplier, nil
}

// Create inserts a new supplier into the database.
func (r *supplierRepository) Create(ctx context.Context, supplier *domain.Supplier) error {
	query := `
		INSERT INTO suppliers (id, name, contact, address, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		supplier.ID,
		supplier.Name,
		supplier.Contact,
		supplier.Address,
		supplier.Active,
		supplier.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create supplier: %w", err)
	}

	return nil
}

// Update modifies an existing supplier.
func (r *supplierRepository) Update(ctx context.Context, supplier *domain.Supplier) error {
	query := `
		UPDATE suppliers
		SET name = $2, contact = $3, address = $4, active = $5
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		supplier.ID,
		supplier.Name,
		supplier.Contact,
		supplier.Address,
		supplier.Active,
	)
	if err != nil {
		return fmt.Errorf("failed to update supplier: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("supplier not found")
	}

	return nil
}

// List retrieves suppliers with optional filtering and pagination.
func (r *supplierRepository) List(ctx context.Context, filter ports.SupplierFilter) ([]domain.Supplier, error) {
	query := `
		SELECT id, name, contact, address, active, created_at
		FROM suppliers
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filter.Active != nil {
		query += fmt.Sprintf(" AND active = $%d", argIdx)
		args = append(args, *filter.Active)
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

	var suppliers []domain.Supplier
	err := pgxscan.Select(ctx, r.db, &suppliers, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list suppliers: %w", err)
	}

	return suppliers, nil
}
