package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// PurchaseBatch represents a purchase order/receipt from a supplier.
type PurchaseBatch struct {
	ID            uuid.UUID  `json:"id"`
	SupplierID    *uuid.UUID `json:"supplier_id,omitempty"`
	InvoiceNumber *string    `json:"invoice_number,omitempty"`
	PurchaseDate  time.Time  `json:"purchase_date"`
	Notes         *string    `json:"notes,omitempty"`
	TotalCost     float64    `json:"total_cost"`
	Received      bool       `json:"received"`
	ReceivedAt    *time.Time `json:"received_at,omitempty"`
	ReceivedBy    *uuid.UUID `json:"received_by,omitempty"`
	CreatedBy     uuid.UUID  `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
}

// PurchaseBatchDetail represents an item in a purchase batch.
type PurchaseBatchDetail struct {
	ID              uuid.UUID `json:"id"`
	PurchaseBatchID uuid.UUID `json:"purchase_batch_id"`
	ProductID       uuid.UUID `json:"product_id"`
	Quantity        int       `json:"quantity"`
	UnitCost        float64   `json:"unit_cost"`
}

// PurchaseBatchWithDetails includes the batch and its line items.
type PurchaseBatchWithDetails struct {
	Batch   PurchaseBatch
	Details []PurchaseBatchDetail
}

// IsValid validates the purchase batch entity.
func (p PurchaseBatch) IsValid() error {
	if p.CreatedBy == uuid.Nil {
		return errors.New("created_by is required")
	}
	return nil
}

// IsReceived returns true if the batch has been received.
func (p PurchaseBatch) IsReceived() bool {
	return p.Received
}

// CanReceive returns true if the batch can be marked as received.
func (p PurchaseBatch) CanReceive() bool {
	return !p.Received
}

// IsValid validates the purchase batch detail.
func (d PurchaseBatchDetail) IsValid() error {
	if d.PurchaseBatchID == uuid.Nil {
		return errors.New("purchase_batch_id is required")
	}
	if d.ProductID == uuid.Nil {
		return errors.New("product_id is required")
	}
	if d.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if d.UnitCost <= 0 {
		return errors.New("unit_cost must be positive")
	}
	return nil
}

// TotalCost calculates the total cost for this detail line.
func (d PurchaseBatchDetail) TotalCost() float64 {
	return float64(d.Quantity) * d.UnitCost
}
