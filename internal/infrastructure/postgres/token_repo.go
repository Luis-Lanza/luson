package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// tokenRepository implements ports.TokenRepository using PostgreSQL.
type tokenRepository struct {
	db *pgxpool.Pool
}

// NewTokenRepository creates a new PostgreSQL token repository.
func NewTokenRepository(db *pgxpool.Pool) ports.TokenRepository {
	return &tokenRepository{db: db}
}

// SaveRefreshToken stores a hashed refresh token for a user.
func (r *tokenRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query, userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

// FindRefreshToken retrieves a refresh token by its hash.
func (r *tokenRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*ports.RefreshToken, error) {
	var token ports.RefreshToken

	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	err := pgxscan.Get(ctx, r.db, &token, query, tokenHash)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	return &token, nil
}

// DeleteRefreshToken removes a refresh token by its hash.
func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE token_hash = $1
	`

	result, err := r.db.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	// If no rows were affected, the token didn't exist
	// We don't return an error here for idempotency
	if result.RowsAffected() == 0 {
		return nil
	}

	return nil
}

// DeleteUserTokens removes all refresh tokens for a user (logout from all devices).
func (r *tokenRepository) DeleteUserTokens(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user tokens: %w", err)
	}

	return nil
}
