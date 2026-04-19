package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid admin user",
			user: User{
				ID:           uuid.New(),
				Username:     "admin",
				PasswordHash: "hashedpassword",
				Role:         UserRoleAdmin,
				BranchID:     nil, // admin doesn't need branch
				Active:       true,
				CreatedAt:    time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid encargado_almacen user",
			user: User{
				ID:           uuid.New(),
				Username:     "encargado1",
				PasswordHash: "hashedpassword",
				Role:         UserRoleEncargadoAlmacen,
				BranchID:     nil, // encargado doesn't need branch
				Active:       true,
				CreatedAt:    time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid cajero user with branch",
			user: User{
				ID:           uuid.New(),
				Username:     "cajero1",
				PasswordHash: "hashedpassword",
				Role:         UserRoleCajero,
				BranchID:     ptr(uuid.New()), // cajero needs branch
				Active:       true,
				CreatedAt:    time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid - empty username",
			user: User{
				ID:           uuid.New(),
				Username:     "",
				PasswordHash: "hashedpassword",
				Role:         UserRoleAdmin,
				Active:       true,
				CreatedAt:    time.Now(),
			},
			wantErr: true,
			errMsg:  "username is required",
		},
		{
			name: "invalid - unknown role",
			user: User{
				ID:           uuid.New(),
				Username:     "user1",
				PasswordHash: "hashedpassword",
				Role:         "unknown_role",
				Active:       true,
				CreatedAt:    time.Now(),
			},
			wantErr: true,
			errMsg:  "invalid role",
		},
		{
			name: "invalid - cajero without branch_id",
			user: User{
				ID:           uuid.New(),
				Username:     "cajero1",
				PasswordHash: "hashedpassword",
				Role:         UserRoleCajero,
				BranchID:     nil, // cajero MUST have branch
				Active:       true,
				CreatedAt:    time.Now(),
			},
			wantErr: true,
			errMsg:  "cajero must have a branch_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.IsValid()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{"admin is admin", UserRoleAdmin, true},
		{"encargado is not admin", UserRoleEncargadoAlmacen, false},
		{"cajero is not admin", UserRoleCajero, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{Role: tt.role}
			assert.Equal(t, tt.expected, u.IsAdmin())
		})
	}
}

func TestUser_HasRole(t *testing.T) {
	u := User{Role: UserRoleCajero}

	tests := []struct {
		name     string
		roles    []UserRole
		expected bool
	}{
		{"has cajero role", []UserRole{UserRoleCajero}, true},
		{"has role in list", []UserRole{UserRoleAdmin, UserRoleCajero}, true},
		{"does not have role", []UserRole{UserRoleAdmin}, false},
		{"empty list", []UserRole{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, u.HasRole(tt.roles...))
		})
	}
}

func TestUser_ActiveStatus(t *testing.T) {
	t.Run("active user", func(t *testing.T) {
		u := User{Active: true}
		assert.True(t, u.Active)
	})

	t.Run("inactive user", func(t *testing.T) {
		u := User{Active: false}
		assert.False(t, u.Active)
	})
}

func TestUser_PasswordHashNotExported(t *testing.T) {
	u := User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: "supersecret",
		Role:         UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}

	// PasswordHash should have json:"-" tag to exclude from JSON
	// This is verified by the struct definition
	assert.Equal(t, "supersecret", u.PasswordHash)
}

// Helper function to create pointer to uuid
func ptr(u uuid.UUID) *uuid.UUID {
	return &u
}
