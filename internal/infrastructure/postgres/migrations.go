package postgres

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

// Migration represents a single database migration.
type Migration struct {
	Version string
	Name    string
	UpSQL   string
	DownSQL string
}

// RunMigrations executes all pending up migrations from the migrations directory.
func (db *DB) RunMigrations(migrationsDir string) error {
	// Create migrations tracking table if not exists
	if err := db.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied, err := db.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get all migration files
	migrations, err := db.loadMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Apply pending migrations
	for _, m := range migrations {
		if _, ok := applied[m.Version]; ok {
			continue // Already applied
		}

		if err := db.applyMigration(m); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", m.Version, err)
		}
	}

	return nil
}

func (db *DB) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`
	_, err := db.Pool.Exec(context.Background(), query)
	return err
}

func (db *DB) getAppliedMigrations() (map[string]bool, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

func (db *DB) loadMigrations(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrationsMap := make(map[string]*Migration)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		// Parse filename: 000001_init.up.sql or 000001_init.down.sql
		parts := strings.Split(name, ".")
		if len(parts) != 3 {
			continue
		}

		baseParts := strings.Split(parts[0], "_")
		if len(baseParts) < 2 {
			continue
		}

		version := baseParts[0]
		direction := parts[1] // up or down

		if migrationsMap[version] == nil {
			migrationsMap[version] = &Migration{Version: version}
		}

		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", name, err)
		}

		if direction == "up" {
			migrationsMap[version].UpSQL = string(content)
		} else if direction == "down" {
			migrationsMap[version].DownSQL = string(content)
		}
	}

	// Convert to slice and sort by version
	var migrations []Migration
	for _, m := range migrationsMap {
		if m.UpSQL != "" { // Only include if has up migration
			migrations = append(migrations, *m)
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		vi, _ := strconv.Atoi(migrations[i].Version)
		vj, _ := strconv.Atoi(migrations[j].Version)
		return vi < vj
	})

	return migrations, nil
}

func (db *DB) applyMigration(m Migration) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Execute migration
	if _, err := tx.Exec(ctx, m.UpSQL); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record migration
	if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", m.Version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RollbackMigration reverts the last applied migration.
func (db *DB) RollbackMigration(migrationsDir string) error {
	// Get last applied migration
	var version string
	err := db.Pool.QueryRow(context.Background(),
		"SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT 1").Scan(&version)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Load migration files
	migrations, err := db.loadMigrations(migrationsDir)
	if err != nil {
		return err
	}

	// Find the migration
	var targetMigration *Migration
	for i := range migrations {
		if migrations[i].Version == version {
			targetMigration = &migrations[i]
			break
		}
	}

	if targetMigration == nil || targetMigration.DownSQL == "" {
		return fmt.Errorf("down migration not found for version %s", version)
	}

	// Execute rollback
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, targetMigration.DownSQL); err != nil {
		return fmt.Errorf("failed to execute rollback: %w", err)
	}

	if _, err := tx.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	return tx.Commit(ctx)
}
