package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPurchaseBatch_IsValid(t *testing.T) {
	supplierID := uuid.New()
	invoiceNumber := "INV-001"
	notes := "Notas de compra"

	tests := []struct {
		name    string
		batch   PurchaseBatch
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid batch",
			batch: PurchaseBatch{
				ID:            uuid.New(),
				SupplierID:    &supplierID,
				InvoiceNumber: &invoiceNumber,
				PurchaseDate:  time.Now(),
				Notes:         &notes,
				TotalCost:     1000.00,
				Received:      false,
				CreatedBy:     uuid.New(),
				CreatedAt:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid batch without supplier",
			batch: PurchaseBatch{
				ID:            uuid.New(),
				SupplierID:    nil,
				InvoiceNumber: &invoiceNumber,
				PurchaseDate:  time.Now(),
				TotalCost:     500.00,
				Received:      false,
				CreatedBy:     uuid.New(),
				CreatedAt:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid - missing created_by",
			batch: PurchaseBatch{
				ID:           uuid.New(),
				PurchaseDate: time.Now(),
				TotalCost:    100.00,
				Received:     false,
				CreatedBy:    uuid.Nil,
				CreatedAt:    time.Now(),
			},
			wantErr: true,
			errMsg:  "created_by is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.batch.IsValid()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPurchaseBatch_IsReceived(t *testing.T) {
	tests := []struct {
		name     string
		received bool
		expected bool
	}{
		{"received", true, true},
		{"not received", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batch := PurchaseBatch{Received: tt.received}
			assert.Equal(t, tt.expected, batch.IsReceived())
		})
	}
}

func TestPurchaseBatch_CanReceive(t *testing.T) {
	tests := []struct {
		name     string
		received bool
		expected bool
	}{
		{"can receive when not received", false, true},
		{"cannot receive when already received", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batch := PurchaseBatch{Received: tt.received}
			assert.Equal(t, tt.expected, batch.CanReceive())
		})
	}
}

func TestPurchaseBatchDetail_IsValid(t *testing.T) {
	batchID := uuid.New()
	productID := uuid.New()

	tests := []struct {
		name    string
		detail  PurchaseBatchDetail
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid detail",
			detail: PurchaseBatchDetail{
				ID:              uuid.New(),
				PurchaseBatchID: batchID,
				ProductID:       productID,
				Quantity:        10,
				UnitCost:        150.00,
			},
			wantErr: false,
		},
		{
			name: "invalid - nil purchase_batch_id",
			detail: PurchaseBatchDetail{
				ID:              uuid.New(),
				PurchaseBatchID: uuid.Nil,
				ProductID:       productID,
				Quantity:        10,
				UnitCost:        150.00,
			},
			wantErr: true,
			errMsg:  "purchase_batch_id is required",
		},
		{
			name: "invalid - nil product_id",
			detail: PurchaseBatchDetail{
				ID:              uuid.New(),
				PurchaseBatchID: batchID,
				ProductID:       uuid.Nil,
				Quantity:        10,
				UnitCost:        150.00,
			},
			wantErr: true,
			errMsg:  "product_id is required",
		},
		{
			name: "invalid - zero quantity",
			detail: PurchaseBatchDetail{
				ID:              uuid.New(),
				PurchaseBatchID: batchID,
				ProductID:       productID,
				Quantity:        0,
				UnitCost:        150.00,
			},
			wantErr: true,
			errMsg:  "quantity must be positive",
		},
		{
			name: "invalid - negative quantity",
			detail: PurchaseBatchDetail{
				ID:              uuid.New(),
				PurchaseBatchID: batchID,
				ProductID:       productID,
				Quantity:        -5,
				UnitCost:        150.00,
			},
			wantErr: true,
			errMsg:  "quantity must be positive",
		},
		{
			name: "invalid - zero unit_cost",
			detail: PurchaseBatchDetail{
				ID:              uuid.New(),
				PurchaseBatchID: batchID,
				ProductID:       productID,
				Quantity:        10,
				UnitCost:        0,
			},
			wantErr: true,
			errMsg:  "unit_cost must be positive",
		},
		{
			name: "invalid - negative unit_cost",
			detail: PurchaseBatchDetail{
				ID:              uuid.New(),
				PurchaseBatchID: batchID,
				ProductID:       productID,
				Quantity:        10,
				UnitCost:        -50.00,
			},
			wantErr: true,
			errMsg:  "unit_cost must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.detail.IsValid()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPurchaseBatchDetail_TotalCost(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		unitCost float64
		expected float64
	}{
		{"normal calculation", 10, 150.00, 1500.00},
		{"single item", 1, 200.00, 200.00},
		{"zero quantity", 0, 150.00, 0},
		{"decimal cost", 3, 99.99, 299.97},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail := PurchaseBatchDetail{
				Quantity: tt.quantity,
				UnitCost: tt.unitCost,
			}
			assert.InDelta(t, tt.expected, detail.TotalCost(), 0.01)
		})
	}
}

func TestPurchaseBatchWithDetails(t *testing.T) {
	batch := PurchaseBatch{
		ID:        uuid.New(),
		TotalCost: 1000.00,
		CreatedBy: uuid.New(),
	}

	details := []PurchaseBatchDetail{
		{
			ID:              uuid.New(),
			PurchaseBatchID: batch.ID,
			ProductID:       uuid.New(),
			Quantity:        5,
			UnitCost:        100.00,
		},
		{
			ID:              uuid.New(),
			PurchaseBatchID: batch.ID,
			ProductID:       uuid.New(),
			Quantity:        3,
			UnitCost:        150.00,
		},
	}

	batchWithDetails := PurchaseBatchWithDetails{
		Batch:   batch,
		Details: details,
	}

	assert.Equal(t, batch.ID, batchWithDetails.Batch.ID)
	assert.Len(t, batchWithDetails.Details, 2)
	assert.Equal(t, 5, batchWithDetails.Details[0].Quantity)
}
