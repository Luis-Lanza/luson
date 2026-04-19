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

// purchaseBatchRepository implements ports.PurchaseBatchRepository using PostgreSQL.
type purchaseBatchRepository struct {
	db *pgxpool.Pool
}

// NewPurchaseBatchRepository creates a new PostgreSQL purchase batch repository.
func NewPurchaseBatchRepository(db *pgxpool.Pool) ports.PurchaseBatchRepository {
	return &purchaseBatchRepository{db: db}
}

// FindByID finds a purchase batch by its ID.
func (r *purchaseBatchRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.PurchaseBatch, error) {
	var batch domain.PurchaseBatch

	query := `
		SELECT id, supplier_id, invoice_number, purchase_date, notes, total_cost,
		       received, received_at, received_by, created_by, created_at
		FROM purchase_batches
		WHERE id = $1
	`

	err := pgxscan.Get(ctx, r.db, &batch, query, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, fmt.Errorf("purchase batch not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find purchase batch by id: %w", err)
	}

	return &batch, nil
}

// findDetailsByBatchID retrieves all details for a purchase batch.
func (r *purchaseBatchRepository) findDetailsByBatchID(ctx context.Context, batchID uuid.UUID) ([]domain.PurchaseBatchDetail, error) {
	query := `
		SELECT id, purchase_batch_id, product_id, quantity, unit_cost
		FROM purchase_batch_details
		WHERE purchase_batch_id = $1
		ORDER BY id
	`

	var details []domain.PurchaseBatchDetail
	err := pgxscan.Select(ctx, r.db, &details, query, batchID)
	if err != nil {
		return nil, fmt.Errorf("failed to find purchase batch details: %w", err)
	}

	return details, nil
}

// Create inserts a new purchase batch with its details into the database.
func (r *purchaseBatchRepository) Create(ctx context.Context, batch *domain.PurchaseBatch) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert the batch
	batchQuery := `
		INSERT INTO purchase_batches (id, supplier_id, invoice_number, purchase_date, notes, total_cost,
		                              received, received_at, received_by, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.Exec(ctx, batchQuery,
		batch.ID,
		batch.SupplierID,
		batch.InvoiceNumber,
		batch.PurchaseDate,
		batch.Notes,
		batch.TotalCost,
		batch.Received,
		batch.ReceivedAt,
		batch.ReceivedBy,
		batch.CreatedBy,
		batch.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create purchase batch: %w", err)
	}

	return tx.Commit(ctx)
}

// CreateWithDetails inserts a new purchase batch with its line items.
func (r *purchaseBatchRepository) CreateWithDetails(ctx context.Context, batch *domain.PurchaseBatch, details []domain.PurchaseBatchDetail) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert the batch
	batchQuery := `
		INSERT INTO purchase_batches (id, supplier_id, invoice_number, purchase_date, notes, total_cost,
		                              received, received_at, received_by, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.Exec(ctx, batchQuery,
		batch.ID,
		batch.SupplierID,
		batch.InvoiceNumber,
		batch.PurchaseDate,
		batch.Notes,
		batch.TotalCost,
		batch.Received,
		batch.ReceivedAt,
		batch.ReceivedBy,
		batch.CreatedBy,
		batch.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create purchase batch: %w", err)
	}

	// Insert details
	detailQuery := `
		INSERT INTO purchase_batch_details (id, purchase_batch_id, product_id, quantity, unit_cost)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, detail := range details {
		_, err = tx.Exec(ctx, detailQuery,
			detail.ID,
			batch.ID,
			detail.ProductID,
			detail.Quantity,
			detail.UnitCost,
		)
		if err != nil {
			return fmt.Errorf("failed to create purchase batch detail: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// MarkAsReceived marks a purchase batch as received.
func (r *purchaseBatchRepository) MarkAsReceived(ctx context.Context, id uuid.UUID, receivedBy uuid.UUID) error {
	query := `
		UPDATE purchase_batches
		SET received = true, received_at = $2, received_by = $3
		WHERE id = $1
	`

	now := time.Now()
	result, err := r.db.Exec(ctx, query, id, now, receivedBy)
	if err != nil {
		return fmt.Errorf("failed to mark purchase batch as received: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("purchase batch not found")
	}

	return nil
}

// List retrieves purchase batches with optional filtering and pagination.
func (r *purchaseBatchRepository) List(ctx context.Context, filter ports.PurchaseBatchFilter) ([]domain.PurchaseBatch, error) {
	query := `
		SELECT id, supplier_id, invoice_number, purchase_date, notes, total_cost,
		       received, received_at, received_by, created_by, created_at
		FROM purchase_batches
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filter.SupplierID != nil {
		query += fmt.Sprintf(" AND supplier_id = $%d", argIdx)
		args = append(args, *filter.SupplierID)
		argIdx++
	}

	if filter.Received != nil {
		query += fmt.Sprintf(" AND received = $%d", argIdx)
		args = append(args, *filter.Received)
		argIdx++
	}

	if filter.FromDate != nil {
		query += fmt.Sprintf(" AND purchase_date >= $%d", argIdx)
		args = append(args, *filter.FromDate)
		argIdx++
	}

	if filter.ToDate != nil {
		query += fmt.Sprintf(" AND purchase_date <= $%d", argIdx)
		args = append(args, *filter.ToDate)
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

	var batches []domain.PurchaseBatch
	err := pgxscan.Select(ctx, r.db, &batches, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list purchase batches: %w", err)
	}

	return batches, nil
}
