package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user in the system.
type UserRole string

const (
	UserRoleAdmin            UserRole = "admin"
	UserRoleEncargadoAlmacen UserRole = "encargado_almacen"
	UserRoleCajero           UserRole = "cajero"
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"` // Never serialized to JSON
	Role         UserRole   `json:"role"`
	BranchID     *uuid.UUID `json:"branch_id,omitempty"`
	Active       bool       `json:"active"`
	CreatedAt    time.Time  `json:"created_at"`
}

// IsValid validates the user entity according to business rules.
func (u User) IsValid() error {
	if u.Username == "" {
		return errors.New("username is required")
	}

	switch u.Role {
	case UserRoleAdmin, UserRoleEncargadoAlmacen, UserRoleCajero:
		// valid roles
	default:
		return fmt.Errorf("invalid role: %s", u.Role)
	}

	if u.Role == UserRoleCajero && u.BranchID == nil {
		return errors.New("cajero must have a branch_id")
	}

	return nil
}

// IsAdmin returns true if the user has admin role.
func (u User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// HasRole returns true if the user has any of the given roles.
func (u User) HasRole(roles ...UserRole) bool {
	for _, r := range roles {
		if u.Role == r {
			return true
		}
	}
	return false
}
