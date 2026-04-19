package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// testDB holds a connection pool for integration tests
var testDB *pgxpool.Pool

// getTestDB returns the test database connection pool
func getTestDB(t *testing.T) *pgxpool.Pool {
	if testDB != nil {
		return testDB
	}

	// Use environment variables or defaults for test database
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/battery_pos_test?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Skipf("Skipping test: could not connect to test database: %v", err)
		return nil
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("Skipping test: could not ping test database: %v", err)
		return nil
	}

	testDB = pool
	return testDB
}

// cleanupTable truncates a table for clean test state
func cleanupTable(t *testing.T, tableName string) {
	db := getTestDB(t)
	if db == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use CASCADE to handle foreign key constraints
	_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName))
	if err != nil {
		t.Logf("Warning: could not cleanup table %s: %v", tableName, err)
	}
}
