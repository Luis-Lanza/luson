package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Branch represents a store branch or warehouse location.
type Branch struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Address          string    `json:"address,omitempty"`
	PettyCashBalance float64   `json:"petty_cash_balance"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"created_at"`
}

// IsValid validates the branch entity according to business rules.
func (b Branch) IsValid() error {
	if b.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
