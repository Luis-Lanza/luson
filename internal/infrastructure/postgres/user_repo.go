package postgres

import (
	"context"
	"fmt"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// userRepository implements ports.UserRepository using PostgreSQL.
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new PostgreSQL user repository.
func NewUserRepository(db *pgxpool.Pool) ports.UserRepository {
	return &userRepository{db: db}
}

// FindByUsername finds a user by their username.
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, username, password_hash, role, branch_id, active, created_at
		FROM users
		WHERE username = $1
	`

	err := pgxscan.Get(ctx, r.db, &user, query, username)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return &user, nil
}

// FindByID finds a user by their ID.
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, username, password_hash, role, branch_id, active, created_at
		FROM users
		WHERE id = $1
	`

	err := pgxscan.Get(ctx, r.db, &user, query, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &user, nil
}

// Create inserts a new user into the database.
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, username, password_hash, role, branch_id, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.BranchID,
		user.Active,
		user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update modifies an existing user.
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET username = $2, password_hash = $3, role = $4, branch_id = $5, active = $6
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.BranchID,
		user.Active,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List retrieves users with optional filtering and pagination.
func (r *userRepository) List(ctx context.Context, filter ports.UserFilter) ([]domain.User, error) {
	query := `
		SELECT id, username, password_hash, role, branch_id, active, created_at
		FROM users
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filter.Role != nil {
		query += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, *filter.Role)
		argIdx++
	}

	if filter.BranchID != nil {
		query += fmt.Sprintf(" AND branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}

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

	var users []domain.User
	err := pgxscan.Select(ctx, r.db, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// NotFound returns true if the error indicates a record was not found.
func NotFound(err error) bool {
	return err == pgx.ErrNoRows
}
