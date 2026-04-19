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

func TestTransferRepository_Create(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "transfer_details")
	defer cleanupTable(t, "transfers")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewTransferRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_transfer",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	branchRepo := NewBranchRepository(db)
	originBranch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Origin Branch",
		Active:    true,
		CreatedAt: time.Now(),
	}
	destBranch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Destination Branch",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, originBranch)
	branchRepo.Create(ctx, destBranch)

	productRepo := NewProductRepository(db)
	product := &domain.Product{
		ID:           uuid.New(),
		Name:         "Transfer Product",
		ProductType:  domain.ProductTypeAccesorio,
		MinSalePrice: 100.00,
		Active:       true,
		CreatedAt:    time.Now(),
		CreatedBy:    user.ID,
	}
	productRepo.Create(ctx, product)

	t.Run("creates transfer with details", func(t *testing.T) {
		notes := "Urgent transfer"
		transfer := &domain.Transfer{
			ID:              uuid.New(),
			OriginType:      "branch",
			OriginID:        originBranch.ID,
			DestinationType: "branch",
			DestinationID:   destBranch.ID,
			Status:          domain.TransferStatusPendiente,
			RequestedBy:     user.ID,
			Notes:           &notes,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		details := []domain.TransferDetail{
			{
				ID:         uuid.New(),
				TransferID: transfer.ID,
				ProductID:  product.ID,
				Quantity:   10,
			},
		}

		err := repo.Create(ctx, transfer, details)
		require.NoError(t, err)

		// Verify
		found, err := repo.FindByID(ctx, transfer.ID)
		require.NoError(t, err)
		assert.Equal(t, transfer.ID, found.ID)
		assert.Equal(t, domain.TransferStatusPendiente, found.Status)
		assert.Equal(t, originBranch.ID, found.OriginID)
		assert.Equal(t, destBranch.ID, found.DestinationID)

		// Verify details
		foundWithDetails, err := repo.FindWithDetails(ctx, transfer.ID)
		require.NoError(t, err)
		assert.Len(t, foundWithDetails.Details, 1)
		assert.Equal(t, 10, foundWithDetails.Details[0].Quantity)
	})
}

func TestTransferRepository_FindByID(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "transfer_details")
	defer cleanupTable(t, "transfers")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewTransferRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_find_transfer",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	branchRepo := NewBranchRepository(db)
	originBranch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Origin Branch Find",
		Active:    true,
		CreatedAt: time.Now(),
	}
	destBranch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Dest Branch Find",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, originBranch)
	branchRepo.Create(ctx, destBranch)

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestTransferRepository_UpdateStatus(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "transfer_details")
	defer cleanupTable(t, "transfers")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewTransferRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	requester := &domain.User{
		ID:           uuid.New(),
		Username:     "requester",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	approver := &domain.User{
		ID:           uuid.New(),
		Username:     "approver",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, requester)
	userRepo.Create(ctx, approver)

	branchRepo := NewBranchRepository(db)
	originBranch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Origin Branch Status",
		Active:    true,
		CreatedAt: time.Now(),
	}
	destBranch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Dest Branch Status",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, originBranch)
	branchRepo.Create(ctx, destBranch)

	t.Run("updates to aprobada", func(t *testing.T) {
		transfer := &domain.Transfer{
			ID:              uuid.New(),
			OriginType:      "branch",
			OriginID:        originBranch.ID,
			DestinationType: "branch",
			DestinationID:   destBranch.ID,
			Status:          domain.TransferStatusPendiente,
			RequestedBy:     requester.ID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		err := repo.Create(ctx, transfer, []domain.TransferDetail{})
		require.NoError(t, err)

		err = repo.UpdateStatus(ctx, transfer.ID, domain.TransferStatusAprobada, &approver.ID)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, transfer.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.TransferStatusAprobada, found.Status)
		assert.NotNil(t, found.ApprovedBy)
		assert.Equal(t, approver.ID, *found.ApprovedBy)
	})

	t.Run("updates to enviada", func(t *testing.T) {
		sender := &domain.User{
			ID:           uuid.New(),
			Username:     "sender",
			PasswordHash: "hash",
			Role:         domain.UserRoleAdmin,
			Active:       true,
			CreatedAt:    time.Now(),
		}
		userRepo.Create(ctx, sender)

		transfer := &domain.Transfer{
			ID:              uuid.New(),
			OriginType:      "branch",
			OriginID:        originBranch.ID,
			DestinationType: "branch",
			DestinationID:   destBranch.ID,
			Status:          domain.TransferStatusAprobada,
			RequestedBy:     requester.ID,
			ApprovedBy:      &approver.ID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		err := repo.Create(ctx, transfer, []domain.TransferDetail{})
		require.NoError(t, err)

		err = repo.UpdateStatus(ctx, transfer.ID, domain.TransferStatusEnviada, &sender.ID)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, transfer.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.TransferStatusEnviada, found.Status)
		assert.NotNil(t, found.SentBy)
	})

	t.Run("updates to recibida", func(t *testing.T) {
		receiver := &domain.User{
			ID:           uuid.New(),
			Username:     "receiver",
			PasswordHash: "hash",
			Role:         domain.UserRoleAdmin,
			Active:       true,
			CreatedAt:    time.Now(),
		}
		userRepo.Create(ctx, receiver)

		senderID := uuid.New()
		transfer := &domain.Transfer{
			ID:              uuid.New(),
			OriginType:      "branch",
			OriginID:        originBranch.ID,
			DestinationType: "branch",
			DestinationID:   destBranch.ID,
			Status:          domain.TransferStatusEnviada,
			RequestedBy:     requester.ID,
			SentBy:          &senderID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		err := repo.Create(ctx, transfer, []domain.TransferDetail{})
		require.NoError(t, err)

		err = repo.UpdateStatus(ctx, transfer.ID, domain.TransferStatusRecibida, &receiver.ID)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, transfer.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.TransferStatusRecibida, found.Status)
		assert.NotNil(t, found.ReceivedBy)
	})

	t.Run("updates to cancelada", func(t *testing.T) {
		transfer := &domain.Transfer{
			ID:              uuid.New(),
			OriginType:      "branch",
			OriginID:        originBranch.ID,
			DestinationType: "branch",
			DestinationID:   destBranch.ID,
			Status:          domain.TransferStatusPendiente,
			RequestedBy:     requester.ID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		err := repo.Create(ctx, transfer, []domain.TransferDetail{})
		require.NoError(t, err)

		err = repo.UpdateStatus(ctx, transfer.ID, domain.TransferStatusCancelada, nil)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, transfer.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.TransferStatusCancelada, found.Status)
	})

	t.Run("returns error for non-existent transfer", func(t *testing.T) {
		err := repo.UpdateStatus(ctx, uuid.New(), domain.TransferStatusAprobada, &approver.ID)
		assert.Error(t, err)
	})
}

func TestTransferRepository_List(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer cleanupTable(t, "transfer_details")
	defer cleanupTable(t, "transfers")
	defer cleanupTable(t, "users")
	defer cleanupTable(t, "branches")

	repo := NewTransferRepository(db)
	ctx := context.Background()

	// Create dependencies
	userRepo := NewUserRepository(db)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser_list_transfer",
		PasswordHash: "hash",
		Role:         domain.UserRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
	}
	userRepo.Create(ctx, user)

	branchRepo := NewBranchRepository(db)
	originBranch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Origin Branch List",
		Active:    true,
		CreatedAt: time.Now(),
	}
	destBranch := &domain.Branch{
		ID:        uuid.New(),
		Name:      "Dest Branch List",
		Active:    true,
		CreatedAt: time.Now(),
	}
	branchRepo.Create(ctx, originBranch)
	branchRepo.Create(ctx, destBranch)

	// Create test transfers
	transfers := []domain.Transfer{
		{ID: uuid.New(), OriginType: "branch", OriginID: originBranch.ID, DestinationType: "branch", DestinationID: destBranch.ID, Status: domain.TransferStatusPendiente, RequestedBy: user.ID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), OriginType: "branch", OriginID: originBranch.ID, DestinationType: "branch", DestinationID: destBranch.ID, Status: domain.TransferStatusAprobada, RequestedBy: user.ID, ApprovedBy: &user.ID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for i := range transfers {
		repo.Create(ctx, &transfers[i], []domain.TransferDetail{})
	}

	t.Run("lists all transfers", func(t *testing.T) {
		result, err := repo.List(ctx, ports.TransferFilter{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 2)
	})

	t.Run("filters by status", func(t *testing.T) {
		status := string(domain.TransferStatusPendiente)
		result, err := repo.List(ctx, ports.TransferFilter{Status: &status})
		require.NoError(t, err)
		for _, tr := range result {
			assert.Equal(t, domain.TransferStatusPendiente, tr.Status)
		}
	})

	t.Run("filters by origin", func(t *testing.T) {
		originType := "branch"
		result, err := repo.List(ctx, ports.TransferFilter{OriginType: &originType, OriginID: &originBranch.ID})
		require.NoError(t, err)
		for _, tr := range result {
			assert.Equal(t, originBranch.ID, tr.OriginID)
		}
	})

	t.Run("filters by requested_by", func(t *testing.T) {
		result, err := repo.List(ctx, ports.TransferFilter{RequestedBy: &user.ID})
		require.NoError(t, err)
		for _, tr := range result {
			assert.Equal(t, user.ID, tr.RequestedBy)
		}
	})

	t.Run("respects pagination", func(t *testing.T) {
		result, err := repo.List(ctx, ports.TransferFilter{Limit: 1, Offset: 0})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(result), 1)
	})
}
