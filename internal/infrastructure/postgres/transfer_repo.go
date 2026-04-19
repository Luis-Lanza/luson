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

// transferRepository implements ports.TransferRepository using PostgreSQL.
type transferRepository struct {
	db *pgxpool.Pool
}

// NewTransferRepository creates a new PostgreSQL transfer repository.
func NewTransferRepository(db *pgxpool.Pool) ports.TransferRepository {
	return &transferRepository{db: db}
}

// FindByID finds a transfer by its ID.
func (r *transferRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Transfer, error) {
	var transfer domain.Transfer

	query := `
		SELECT id, origin_type, origin_id, destination_type, destination_id, status,
		       requested_by, approved_by, rejected_by, sent_by, received_by,
		       notes, rejection_reason, created_at, updated_at
		FROM transfers
		WHERE id = $1
	`

	err := pgxscan.Get(ctx, r.db, &transfer, query, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, fmt.Errorf("transfer not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find transfer by id: %w", err)
	}

	return &transfer, nil
}

// FindWithDetails retrieves a transfer with all its line items.
func (r *transferRepository) FindWithDetails(ctx context.Context, id uuid.UUID) (*domain.TransferWithDetails, error) {
	// Get the transfer
	transfer, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get the details
	details, err := r.findDetailsByTransferID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.TransferWithDetails{
		Transfer: *transfer,
		Details:  details,
	}, nil
}

// findDetailsByTransferID retrieves all details for a transfer.
func (r *transferRepository) findDetailsByTransferID(ctx context.Context, transferID uuid.UUID) ([]domain.TransferDetail, error) {
	query := `
		SELECT id, transfer_id, product_id, quantity
		FROM transfer_details
		WHERE transfer_id = $1
		ORDER BY id
	`

	var details []domain.TransferDetail
	err := pgxscan.Select(ctx, r.db, &details, query, transferID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transfer details: %w", err)
	}

	return details, nil
}

// Create inserts a new transfer with its details into the database.
func (r *transferRepository) Create(ctx context.Context, transfer *domain.Transfer, details []domain.TransferDetail) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert the transfer
	transferQuery := `
		INSERT INTO transfers (id, origin_type, origin_id, destination_type, destination_id, status,
		                      requested_by, approved_by, rejected_by, sent_by, received_by,
		                      notes, rejection_reason, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err = tx.Exec(ctx, transferQuery,
		transfer.ID,
		transfer.OriginType,
		transfer.OriginID,
		transfer.DestinationType,
		transfer.DestinationID,
		transfer.Status,
		transfer.RequestedBy,
		transfer.ApprovedBy,
		transfer.RejectedBy,
		transfer.SentBy,
		transfer.ReceivedBy,
		transfer.Notes,
		transfer.RejectionReason,
		transfer.CreatedAt,
		transfer.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create transfer: %w", err)
	}

	// Insert details
	detailQuery := `
		INSERT INTO transfer_details (id, transfer_id, product_id, quantity)
		VALUES ($1, $2, $3, $4)
	`

	for _, detail := range details {
		_, err = tx.Exec(ctx, detailQuery,
			detail.ID,
			transfer.ID,
			detail.ProductID,
			detail.Quantity,
		)
		if err != nil {
			return fmt.Errorf("failed to create transfer detail: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// UpdateStatus updates the status of a transfer and the corresponding user field.
func (r *transferRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransferStatus, userID *uuid.UUID) error {
	now := time.Now()

	var query string
	var args []interface{}

	switch status {
	case domain.TransferStatusAprobada:
		query = `UPDATE transfers SET status = $2, approved_by = $3, updated_at = $4 WHERE id = $1`
		args = []interface{}{id, status, *userID, now}
	case domain.TransferStatusRechazada:
		query = `UPDATE transfers SET status = $2, rejected_by = $3, updated_at = $4 WHERE id = $1`
		args = []interface{}{id, status, *userID, now}
	case domain.TransferStatusEnviada:
		query = `UPDATE transfers SET status = $2, sent_by = $3, updated_at = $4 WHERE id = $1`
		args = []interface{}{id, status, *userID, now}
	case domain.TransferStatusRecibida:
		query = `UPDATE transfers SET status = $2, received_by = $3, updated_at = $4 WHERE id = $1`
		args = []interface{}{id, status, *userID, now}
	case domain.TransferStatusCancelada:
		query = `UPDATE transfers SET status = $2, updated_at = $3 WHERE id = $1`
		args = []interface{}{id, status, now}
	default:
		query = `UPDATE transfers SET status = $2, updated_at = $3 WHERE id = $1`
		args = []interface{}{id, status, now}
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update transfer status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("transfer not found")
	}

	return nil
}

// List retrieves transfers with optional filtering and pagination.
func (r *transferRepository) List(ctx context.Context, filter ports.TransferFilter) ([]domain.Transfer, error) {
	query := `
		SELECT id, origin_type, origin_id, destination_type, destination_id, status,
		       requested_by, approved_by, rejected_by, sent_by, received_by,
		       notes, rejection_reason, created_at, updated_at
		FROM transfers
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filter.OriginType != nil {
		query += fmt.Sprintf(" AND origin_type = $%d", argIdx)
		args = append(args, *filter.OriginType)
		argIdx++
	}

	if filter.OriginID != nil {
		query += fmt.Sprintf(" AND origin_id = $%d", argIdx)
		args = append(args, *filter.OriginID)
		argIdx++
	}

	if filter.DestinationType != nil {
		query += fmt.Sprintf(" AND destination_type = $%d", argIdx)
		args = append(args, *filter.DestinationType)
		argIdx++
	}

	if filter.DestinationID != nil {
		query += fmt.Sprintf(" AND destination_id = $%d", argIdx)
		args = append(args, *filter.DestinationID)
		argIdx++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}

	if filter.RequestedBy != nil {
		query += fmt.Sprintf(" AND requested_by = $%d", argIdx)
		args = append(args, *filter.RequestedBy)
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

	var transfers []domain.Transfer
	err := pgxscan.Select(ctx, r.db, &transfers, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list transfers: %w", err)
	}

	return transfers, nil
}
