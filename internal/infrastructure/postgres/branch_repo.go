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

// branchRepository implements ports.BranchRepository using PostgreSQL.
type branchRepository struct {
	db *pgxpool.Pool
}

// NewBranchRepository creates a new PostgreSQL branch repository.
func NewBranchRepository(db *pgxpool.Pool) ports.BranchRepository {
	return &branchRepository{db: db}
}

// FindByID finds a branch by its ID.
func (r *branchRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error) {
	var branch domain.Branch

	query := `
		SELECT id, name, address, petty_cash_balance, active, created_at
		FROM branches
		WHERE id = $1
	`

	err := pgxscan.Get(ctx, r.db, &branch, query, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, fmt.Errorf("branch not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find branch by id: %w", err)
	}

	return &branch, nil
}

// Create inserts a new branch into the database.
func (r *branchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	query := `
		INSERT INTO branches (id, name, address, petty_cash_balance, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		branch.ID,
		branch.Name,
		branch.Address,
		branch.PettyCashBalance,
		branch.Active,
		branch.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

// Update modifies an existing branch.
func (r *branchRepository) Update(ctx context.Context, branch *domain.Branch) error {
	query := `
		UPDATE branches
		SET name = $2, address = $3, petty_cash_balance = $4, active = $5
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		branch.ID,
		branch.Name,
		branch.Address,
		branch.PettyCashBalance,
		branch.Active,
	)
	if err != nil {
		return fmt.Errorf("failed to update branch: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("branch not found")
	}

	return nil
}

// List retrieves branches with optional filtering and pagination.
func (r *branchRepository) List(ctx context.Context, filter ports.BranchFilter) ([]domain.Branch, error) {
	query := `
		SELECT id, name, address, petty_cash_balance, active, created_at
		FROM branches
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

	var branches []domain.Branch
	err := pgxscan.Select(ctx, r.db, &branches, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	return branches, nil
}
