package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Supplier represents a product supplier.
type Supplier struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Contact   *string   `json:"contact,omitempty"`
	Address   *string   `json:"address,omitempty"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

// IsValid validates the supplier entity according to business rules.
func (s Supplier) IsValid() error {
	if s.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
