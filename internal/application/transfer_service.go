package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/google/uuid"
)

// transferService implements the ports.TransferService interface.
type transferService struct {
	transferRepo ports.TransferRepository
	productRepo  ports.ProductRepository
	stockRepo    ports.StockRepository
}

// NewTransferService creates a new instance of TransferService.
func NewTransferService(transferRepo ports.TransferRepository, productRepo ports.ProductRepository, stockRepo ports.StockRepository) ports.TransferService {
	return &transferService{
		transferRepo: transferRepo,
		productRepo:  productRepo,
		stockRepo:    stockRepo,
	}
}

// Create creates a new transfer with items.
func (s *transferService) Create(ctx context.Context, req ports.CreateTransferRequest) (*domain.Transfer, error) {
	// Validate input
	if req.OriginID == uuid.Nil {
		return nil, domain.ErrTransferInvalidOrigin
	}
	if req.DestinationID == uuid.Nil {
		return nil, domain.ErrTransferInvalidDestination
	}
	if req.OriginType == "" || req.DestinationType == "" {
		return nil, domain.ErrTransferInvalidLocationType
	}
	if req.OriginID == req.DestinationID && req.OriginType == req.DestinationType {
		return nil, domain.ErrTransferSameLocation
	}
	if req.RequestedBy == uuid.Nil {
		return nil, errors.New("requested_by is required")
	}
	if len(req.Items) == 0 {
		return nil, domain.ErrTransferNoItems
	}

	// Validate all products exist and have sufficient stock at origin
	for i, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("item %d: %w", i, domain.ErrTransferInvalidQuantity)
		}

		_, err := s.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("item %d: product not found: %w", i, err)
		}

		// Check stock at origin
		stock, err := s.stockRepo.FindByProductAndLocation(ctx, item.ProductID, req.OriginType, req.OriginID)
		if err != nil {
			return nil, fmt.Errorf("item %d: stock not found at origin: %w", i, err)
		}
		if stock == nil {
			return nil, fmt.Errorf("item %d: no stock at origin location", i)
		}
		if !stock.CanFulfill(item.Quantity) {
			return nil, fmt.Errorf("item %d: %w", i, domain.ErrStockInsufficientQuantity)
		}
	}

	// Create transfer entity
	now := time.Now()
	transfer := &domain.Transfer{
		ID:              uuid.New(),
		OriginType:      req.OriginType,
		OriginID:        req.OriginID,
		DestinationType: req.DestinationType,
		DestinationID:   req.DestinationID,
		Status:          domain.TransferStatusPendiente,
		RequestedBy:     req.RequestedBy,
		Notes:           req.Notes,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Create details
	details := make([]domain.TransferDetail, len(req.Items))
	for i, item := range req.Items {
		details[i] = domain.TransferDetail{
			ID:         uuid.New(),
			TransferID: transfer.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
		}
	}

	// Save to repository
	if err := s.transferRepo.Create(ctx, transfer, details); err != nil {
		return nil, fmt.Errorf("failed to create transfer: %w", err)
	}

	return transfer, nil
}

// GetByID retrieves a transfer by ID with its details.
func (s *transferService) GetByID(ctx context.Context, id uuid.UUID) (*domain.TransferWithDetails, error) {
	transferWithDetails, err := s.transferRepo.FindWithDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("transfer not found: %w", err)
	}

	return transferWithDetails, nil
}

// List retrieves transfers with filtering.
func (s *transferService) List(ctx context.Context, filter ports.TransferFilter) ([]domain.Transfer, error) {
	transfers, err := s.transferRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list transfers: %w", err)
	}

	return transfers, nil
}

// Approve approves a pending transfer.
func (s *transferService) Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error {
	transfer, err := s.transferRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("transfer not found: %w", err)
	}

	if !transfer.CanTransition(domain.TransferStatusAprobada) {
		return domain.ErrTransferInvalidTransition
	}

	if err := s.transferRepo.UpdateStatus(ctx, id, domain.TransferStatusAprobada, &approvedBy); err != nil {
		return fmt.Errorf("failed to approve transfer: %w", err)
	}

	return nil
}

// Reject rejects a pending transfer.
func (s *transferService) Reject(ctx context.Context, id uuid.UUID, rejectedBy uuid.UUID, reason string) error {
	transfer, err := s.transferRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("transfer not found: %w", err)
	}

	if !transfer.CanTransition(domain.TransferStatusRechazada) {
		return domain.ErrTransferInvalidTransition
	}

	// TODO: Store rejection reason (would need to update the transfer entity)
	_ = reason

	if err := s.transferRepo.UpdateStatus(ctx, id, domain.TransferStatusRechazada, &rejectedBy); err != nil {
		return fmt.Errorf("failed to reject transfer: %w", err)
	}

	return nil
}

// MarkAsSent marks an approved transfer as sent.
func (s *transferService) MarkAsSent(ctx context.Context, id uuid.UUID, sentBy uuid.UUID) error {
	transfer, err := s.transferRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("transfer not found: %w", err)
	}

	if !transfer.CanTransition(domain.TransferStatusEnviada) {
		return domain.ErrTransferInvalidTransition
	}

	if err := s.transferRepo.UpdateStatus(ctx, id, domain.TransferStatusEnviada, &sentBy); err != nil {
		return fmt.Errorf("failed to mark transfer as sent: %w", err)
	}

	return nil
}

// MarkAsReceived marks a sent transfer as received.
func (s *transferService) MarkAsReceived(ctx context.Context, id uuid.UUID, receivedBy uuid.UUID) error {
	transfer, err := s.transferRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("transfer not found: %w", err)
	}

	if !transfer.CanTransition(domain.TransferStatusRecibida) {
		return domain.ErrTransferInvalidTransition
	}

	if err := s.transferRepo.UpdateStatus(ctx, id, domain.TransferStatusRecibida, &receivedBy); err != nil {
		return fmt.Errorf("failed to mark transfer as received: %w", err)
	}

	// TODO: Update stock quantities (decrease at origin, increase at destination)

	return nil
}

// Cancel cancels a transfer.
func (s *transferService) Cancel(ctx context.Context, id uuid.UUID) error {
	transfer, err := s.transferRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("transfer not found: %w", err)
	}

	if !transfer.CanTransition(domain.TransferStatusCancelada) {
		return domain.ErrTransferInvalidTransition
	}

	if err := s.transferRepo.UpdateStatus(ctx, id, domain.TransferStatusCancelada, nil); err != nil {
		return fmt.Errorf("failed to cancel transfer: %w", err)
	}

	return nil
}
