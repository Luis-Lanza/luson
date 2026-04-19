package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_FindByUsername(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("finds existing user", func(t *testing.T) {
		// Create a user first
		user := &domain.User{
			ID:           uuid.New(),
			Username:     "testuser",
			PasswordHash: "hashedpassword",
			Role:         domain.UserRoleAdmin,
			Active:       true,
			CreatedAt:    time.Now(),
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Find by username
		found, err := repo.FindByUsername(ctx, "testuser")
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "testuser", found.Username)
		assert.Equal(t, domain.UserRoleAdmin, found.Role)
	})

	t.Run("returns nil for non-existent user", func(t *testing.T) {
		found, err := repo.FindByUsername(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("finds user by ID", func(t *testing.T) {
		user := &domain.User{
			ID:           uuid.New(),
			Username:     "testuser",
			PasswordHash: "hashedpassword",
			Role:         domain.UserRoleCajero,
			Active:       true,
			CreatedAt:    time.Now(),
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, user.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, "testuser", found.Username)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestUserRepository_Create(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("creates user with all fields", func(t *testing.T) {
		user := &domain.User{
			ID:           uuid.New(),
			Username:     "newuser",
			PasswordHash: "hashedpassword123",
			Role:         domain.UserRoleEncargadoAlmacen,
			Active:       true,
			CreatedAt:    time.Now(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Verify user was created
		found, err := repo.FindByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "newuser", found.Username)
		assert.Equal(t, "hashedpassword123", found.PasswordHash)
	})

	t.Run("fails with duplicate username", func(t *testing.T) {
		user1 := &domain.User{
			ID:           uuid.New(),
			Username:     "duplicateuser",
			PasswordHash: "hash1",
			Role:         domain.UserRoleAdmin,
			Active:       true,
			CreatedAt:    time.Now(),
		}
		err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := &domain.User{
			ID:           uuid.New(),
			Username:     "duplicateuser",
			PasswordHash: "hash2",
			Role:         domain.UserRoleCajero,
			Active:       true,
			CreatedAt:    time.Now(),
		}
		err = repo.Create(ctx, user2)
		assert.Error(t, err) // Should fail due to unique constraint
	})
}

func TestUserRepository_Update(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("updates user fields", func(t *testing.T) {
		user := &domain.User{
			ID:           uuid.New(),
			Username:     "updateuser",
			PasswordHash: "oldhash",
			Role:         domain.UserRoleAdmin,
			Active:       true,
			CreatedAt:    time.Now(),
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Update fields
		user.PasswordHash = "newhash"
		user.Active = false

		err = repo.Update(ctx, user)
		require.NoError(t, err)

		// Verify update
		found, err := repo.FindByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "newhash", found.PasswordHash)
		assert.False(t, found.Active)
	})
}

func TestUserRepository_List(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "users")

	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test users
	users := []domain.User{
		{ID: uuid.New(), Username: "admin1", PasswordHash: "hash", Role: domain.UserRoleAdmin, Active: true, CreatedAt: time.Now()},
		{ID: uuid.New(), Username: "cajero1", PasswordHash: "hash", Role: domain.UserRoleCajero, Active: true, CreatedAt: time.Now()},
		{ID: uuid.New(), Username: "inactive1", PasswordHash: "hash", Role: domain.UserRoleAdmin, Active: false, CreatedAt: time.Now()},
	}

	for i := range users {
		err := repo.Create(ctx, &users[i])
		require.NoError(t, err)
	}

	t.Run("lists all users", func(t *testing.T) {
		result, err := repo.List(ctx, ports.UserFilter{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 3)
	})

	t.Run("filters by role", func(t *testing.T) {
		adminRole := string(domain.UserRoleAdmin)
		result, err := repo.List(ctx, ports.UserFilter{Role: &adminRole})
		require.NoError(t, err)
		for _, u := range result {
			assert.Equal(t, domain.UserRoleAdmin, u.Role)
		}
	})

	t.Run("filters by active status", func(t *testing.T) {
		active := true
		result, err := repo.List(ctx, ports.UserFilter{Active: &active})
		require.NoError(t, err)
		for _, u := range result {
			assert.True(t, u.Active)
		}
	})

	t.Run("respects pagination", func(t *testing.T) {
		result, err := repo.List(ctx, ports.UserFilter{Limit: 2, Offset: 0})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(result), 2)
	})
}
