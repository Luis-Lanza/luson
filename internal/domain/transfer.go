package domain

import (
	"time"

	"github.com/google/uuid"
)

// TransferStatus represents the status of a stock transfer.
type TransferStatus string

const (
	TransferStatusPendiente TransferStatus = "pendiente"
	TransferStatusAprobada  TransferStatus = "aprobada"
	TransferStatusRechazada TransferStatus = "rechazada"
	TransferStatusEnviada   TransferStatus = "enviada"
	TransferStatusRecibida  TransferStatus = "recibida"
	TransferStatusCancelada TransferStatus = "cancelada"
)

// Transfer represents a stock transfer between locations.
type Transfer struct {
	ID              uuid.UUID      `json:"id"`
	OriginType      string         `json:"origin_type"`
	OriginID        uuid.UUID      `json:"origin_id"`
	DestinationType string         `json:"destination_type"`
	DestinationID   uuid.UUID      `json:"destination_id"`
	Status          TransferStatus `json:"status"`
	RequestedBy     uuid.UUID      `json:"requested_by"`
	ApprovedBy      *uuid.UUID     `json:"approved_by,omitempty"`
	RejectedBy      *uuid.UUID     `json:"rejected_by,omitempty"`
	SentBy          *uuid.UUID     `json:"sent_by,omitempty"`
	ReceivedBy      *uuid.UUID     `json:"received_by,omitempty"`
	Notes           *string        `json:"notes,omitempty"`
	RejectionReason *string        `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// TransferDetail represents an item in a transfer.
type TransferDetail struct {
	ID         uuid.UUID `json:"id"`
	TransferID uuid.UUID `json:"transfer_id"`
	ProductID  uuid.UUID `json:"product_id"`
	Quantity   int       `json:"quantity"`
}

// TransferWithDetails includes the transfer and its line items.
type TransferWithDetails struct {
	Transfer Transfer
	Details  []TransferDetail
}

// CanTransition checks if a status transition is valid.
func (t Transfer) CanTransition(to TransferStatus) bool {
	return CanTransferTransition(t.Status, to)
}

// IsTerminal returns true if the transfer is in a terminal state.
func (t Transfer) IsTerminal() bool {
	switch t.Status {
	case TransferStatusRecibida, TransferStatusRechazada, TransferStatusCancelada:
		return true
	default:
		return false
	}
}

// CanTransferTransition determines if a transfer can move from one status to another.
func CanTransferTransition(from, to TransferStatus) bool {
	switch from {
	case TransferStatusPendiente:
		return to == TransferStatusAprobada || to == TransferStatusRechazada || to == TransferStatusCancelada
	case TransferStatusAprobada:
		return to == TransferStatusEnviada || to == TransferStatusCancelada
	case TransferStatusEnviada:
		return to == TransferStatusRecibida || to == TransferStatusCancelada
	case TransferStatusRecibida, TransferStatusRechazada, TransferStatusCancelada:
		return false // Terminal states
	default:
		return false
	}
}

// IsValidTransferStatus checks if a string is a valid transfer status.
func IsValidTransferStatus(status string) bool {
	switch TransferStatus(status) {
	case TransferStatusPendiente, TransferStatusAprobada, TransferStatusRechazada,
		TransferStatusEnviada, TransferStatusRecibida, TransferStatusCancelada:
		return true
	default:
		return false
	}
}
