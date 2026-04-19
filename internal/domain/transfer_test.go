package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCanTransferTransition(t *testing.T) {
	tests := []struct {
		name     string
		from     TransferStatus
		to       TransferStatus
		expected bool
	}{
		// From Pendiente
		{"pendiente -> aprobada", TransferStatusPendiente, TransferStatusAprobada, true},
		{"pendiente -> rechazada", TransferStatusPendiente, TransferStatusRechazada, true},
		{"pendiente -> cancelada", TransferStatusPendiente, TransferStatusCancelada, true},
		{"pendiente -> enviada", TransferStatusPendiente, TransferStatusEnviada, false},
		{"pendiente -> recibida", TransferStatusPendiente, TransferStatusRecibida, false},
		{"pendiente -> pendiente", TransferStatusPendiente, TransferStatusPendiente, false},

		// From Aprobada
		{"aprobada -> enviada", TransferStatusAprobada, TransferStatusEnviada, true},
		{"aprobada -> cancelada", TransferStatusAprobada, TransferStatusCancelada, true},
		{"aprobada -> recibida", TransferStatusAprobada, TransferStatusRecibida, false},
		{"aprobada -> pendiente", TransferStatusAprobada, TransferStatusPendiente, false},
		{"aprobada -> rechazada", TransferStatusAprobada, TransferStatusRechazada, false},
		{"aprobada -> aprobada", TransferStatusAprobada, TransferStatusAprobada, false},

		// From Enviada
		{"enviada -> recibida", TransferStatusEnviada, TransferStatusRecibida, true},
		{"enviada -> cancelada", TransferStatusEnviada, TransferStatusCancelada, true},
		{"enviada -> pendiente", TransferStatusEnviada, TransferStatusPendiente, false},
		{"enviada -> aprobada", TransferStatusEnviada, TransferStatusAprobada, false},
		{"enviada -> rechazada", TransferStatusEnviada, TransferStatusRechazada, false},
		{"enviada -> enviada", TransferStatusEnviada, TransferStatusEnviada, false},

		// Terminal states
		{"recibida -> pendiente", TransferStatusRecibida, TransferStatusPendiente, false},
		{"recibida -> aprobada", TransferStatusRecibida, TransferStatusAprobada, false},
		{"recibida -> recibida", TransferStatusRecibida, TransferStatusRecibida, false},
		{"rechazada -> pendiente", TransferStatusRechazada, TransferStatusPendiente, false},
		{"rechazada -> aprobada", TransferStatusRechazada, TransferStatusAprobada, false},
		{"rechazada -> rechazada", TransferStatusRechazada, TransferStatusRechazada, false},
		{"cancelada -> pendiente", TransferStatusCancelada, TransferStatusPendiente, false},
		{"cancelada -> aprobada", TransferStatusCancelada, TransferStatusAprobada, false},
		{"cancelada -> cancelada", TransferStatusCancelada, TransferStatusCancelada, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanTransferTransition(tt.from, tt.to)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTransfer_CanTransition(t *testing.T) {
	transfer := Transfer{
		ID:     uuid.New(),
		Status: TransferStatusPendiente,
	}

	assert.True(t, transfer.CanTransition(TransferStatusAprobada))
	assert.True(t, transfer.CanTransition(TransferStatusRechazada))
	assert.False(t, transfer.CanTransition(TransferStatusRecibida))
}

func TestTransfer_IsTerminal(t *testing.T) {
	tests := []struct {
		name     string
		status   TransferStatus
		expected bool
	}{
		{"recibida is terminal", TransferStatusRecibida, true},
		{"rechazada is terminal", TransferStatusRechazada, true},
		{"cancelada is terminal", TransferStatusCancelada, true},
		{"pendiente is not terminal", TransferStatusPendiente, false},
		{"aprobada is not terminal", TransferStatusAprobada, false},
		{"enviada is not terminal", TransferStatusEnviada, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfer := Transfer{Status: tt.status}
			assert.Equal(t, tt.expected, transfer.IsTerminal())
		})
	}
}

func TestIsValidTransferStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"pendiente is valid", "pendiente", true},
		{"aprobada is valid", "aprobada", true},
		{"rechazada is valid", "rechazada", true},
		{"enviada is valid", "enviada", true},
		{"recibida is valid", "recibida", true},
		{"cancelada is valid", "cancelada", true},
		{"unknown is invalid", "unknown", false},
		{"empty is invalid", "", false},
		{"APROBADA is invalid (case sensitive)", "APROBADA", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsValidTransferStatus(tt.status))
		})
	}
}

func TestTransferStatus_Constants(t *testing.T) {
	assert.Equal(t, TransferStatus("pendiente"), TransferStatusPendiente)
	assert.Equal(t, TransferStatus("aprobada"), TransferStatusAprobada)
	assert.Equal(t, TransferStatus("rechazada"), TransferStatusRechazada)
	assert.Equal(t, TransferStatus("enviada"), TransferStatusEnviada)
	assert.Equal(t, TransferStatus("recibida"), TransferStatusRecibida)
	assert.Equal(t, TransferStatus("cancelada"), TransferStatusCancelada)
}

func TestTransferWithDetails(t *testing.T) {
	transfer := Transfer{
		ID:              uuid.New(),
		OriginType:      "branch",
		OriginID:        uuid.New(),
		DestinationType: "branch",
		DestinationID:   uuid.New(),
		Status:          TransferStatusPendiente,
		RequestedBy:     uuid.New(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	details := []TransferDetail{
		{
			ID:         uuid.New(),
			TransferID: transfer.ID,
			ProductID:  uuid.New(),
			Quantity:   5,
		},
		{
			ID:         uuid.New(),
			TransferID: transfer.ID,
			ProductID:  uuid.New(),
			Quantity:   3,
		},
	}

	transferWithDetails := TransferWithDetails{
		Transfer: transfer,
		Details:  details,
	}

	assert.Equal(t, transfer.ID, transferWithDetails.Transfer.ID)
	assert.Len(t, transferWithDetails.Details, 2)
	assert.Equal(t, 5, transferWithDetails.Details[0].Quantity)
}
