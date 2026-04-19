package application

import (
	"context"
	"errors"
	"testing"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockPasswordHasher is a mock for password hashing
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) CheckPassword(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

func TestUserService_Create(t *testing.T) {
	userRepo := new(MockUserRepository)
	passwordHasher := new(MockPasswordHasher)

	service := NewUserService(userRepo, passwordHasher)
	ctx := context.Background()

	t.Run("successful user creation", func(t *testing.T) {
		branchID := uuid.New()
		req := ports.CreateUserRequest{
			Username: "newuser",
			Password: "password123",
			Role:     "cajero",
			BranchID: &branchID,
		}

		passwordHasher.On("HashPassword", req.Password).Return("hashedpassword", nil).Once()
		userRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		user, err := service.Create(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Username, user.Username)
		assert.Equal(t, domain.UserRole(req.Role), user.Role)
		assert.Equal(t, &branchID, user.BranchID)
		assert.True(t, user.Active)
		assert.Empty(t, user.PasswordHash) // Should not be exposed
		userRepo.AssertExpectations(t)
		passwordHasher.AssertExpectations(t)
	})

	t.Run("fails with empty username", func(t *testing.T) {
		req := ports.CreateUserRequest{
			Username: "",
			Password: "password123",
			Role:     "admin",
		}

		user, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "username is required")
	})

	t.Run("fails with empty password", func(t *testing.T) {
		req := ports.CreateUserRequest{
			Username: "testuser",
			Password: "",
			Role:     "admin",
		}

		user, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "password is required")
	})

	t.Run("fails with invalid role", func(t *testing.T) {
		req := ports.CreateUserRequest{
			Username: "testuser",
			Password: "password123",
			Role:     "invalid_role",
		}

		user, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid role")
	})

	t.Run("cajero role requires branch_id", func(t *testing.T) {
		req := ports.CreateUserRequest{
			Username: "testuser",
			Password: "password123",
			Role:     "cajero",
			BranchID: nil,
		}

		user, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "cajero must have a branch_id")
	})

	t.Run("fails when password hashing fails", func(t *testing.T) {
		branchID := uuid.New()
		req := ports.CreateUserRequest{
			Username: "newuser",
			Password: "password123",
			Role:     "admin",
			BranchID: &branchID,
		}

		passwordHasher.On("HashPassword", req.Password).Return("", errors.New("hash failed")).Once()

		user, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		passwordHasher.AssertExpectations(t)
	})

	t.Run("fails when repository returns error", func(t *testing.T) {
		branchID := uuid.New()
		req := ports.CreateUserRequest{
			Username: "newuser",
			Password: "password123",
			Role:     "admin",
			BranchID: &branchID,
		}

		passwordHasher.On("HashPassword", req.Password).Return("hashedpassword", nil).Once()
		userRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(errors.New("db error")).Once()

		user, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		userRepo.AssertExpectations(t)
		passwordHasher.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	userRepo := new(MockUserRepository)
	passwordHasher := new(MockPasswordHasher)

	service := NewUserService(userRepo, passwordHasher)
	ctx := context.Background()

	t.Run("returns user by id", func(t *testing.T) {
		userID := uuid.New()
		user := &domain.User{
			ID:       userID,
			Username: "testuser",
			Role:     domain.UserRoleAdmin,
			Active:   true,
		}

		userRepo.On("FindByID", ctx, userID).Return(user, nil).Once()

		result, err := service.GetByID(ctx, userID)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, "testuser", result.Username)
		userRepo.AssertExpectations(t)
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		userID := uuid.New()

		userRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		result, err := service.GetByID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_List(t *testing.T) {
	userRepo := new(MockUserRepository)
	passwordHasher := new(MockPasswordHasher)

	service := NewUserService(userRepo, passwordHasher)
	ctx := context.Background()

	t.Run("returns list of users", func(t *testing.T) {
		active := true
		filter := ports.UserFilter{
			Active: &active,
			Limit:  10,
			Offset: 0,
		}

		users := []domain.User{
			{ID: uuid.New(), Username: "user1", Role: domain.UserRoleAdmin, Active: true},
			{ID: uuid.New(), Username: "user2", Role: domain.UserRoleCajero, Active: true},
		}

		userRepo.On("List", ctx, filter).Return(users, nil).Once()

		result, err := service.List(ctx, filter)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "user1", result[0].Username)
		userRepo.AssertExpectations(t)
	})

	t.Run("returns empty list when no users", func(t *testing.T) {
		filter := ports.UserFilter{}

		userRepo.On("List", ctx, filter).Return([]domain.User{}, nil).Once()

		result, err := service.List(ctx, filter)

		require.NoError(t, err)
		assert.Empty(t, result)
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	userRepo := new(MockUserRepository)
	passwordHasher := new(MockPasswordHasher)

	service := NewUserService(userRepo, passwordHasher)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		userID := uuid.New()
		branchID := uuid.New()
		newRole := "encargado_almacen"
		newActive := false

		existingUser := &domain.User{
			ID:       userID,
			Username: "testuser",
			Role:     domain.UserRoleAdmin,
			Active:   true,
		}

		req := ports.UpdateUserRequest{
			Role:     &newRole,
			BranchID: &branchID,
			Active:   &newActive,
		}

		userRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
		userRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		result, err := service.Update(ctx, userID, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.UserRole(newRole), result.Role)
		assert.Equal(t, &branchID, result.BranchID)
		assert.False(t, result.Active)
		userRepo.AssertExpectations(t)
	})

	t.Run("fails when user not found", func(t *testing.T) {
		userID := uuid.New()
		newRole := "admin"
		req := ports.UpdateUserRequest{Role: &newRole}

		userRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		result, err := service.Update(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdatePassword(t *testing.T) {
	userRepo := new(MockUserRepository)
	passwordHasher := new(MockPasswordHasher)

	service := NewUserService(userRepo, passwordHasher)
	ctx := context.Background()

	t.Run("successful password update", func(t *testing.T) {
		userID := uuid.New()
		existingUser := &domain.User{
			ID:           userID,
			Username:     "testuser",
			PasswordHash: "oldhash",
			Role:         domain.UserRoleAdmin,
			Active:       true,
		}

		userRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
		passwordHasher.On("HashPassword", "newpassword123").Return("newhash", nil).Once()
		userRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		err := service.UpdatePassword(ctx, userID, "newpassword123")

		require.NoError(t, err)
		userRepo.AssertExpectations(t)
		passwordHasher.AssertExpectations(t)
	})

	t.Run("fails with empty password", func(t *testing.T) {
		userID := uuid.New()

		err := service.UpdatePassword(ctx, userID, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password is required")
	})

	t.Run("fails when user not found", func(t *testing.T) {
		userID := uuid.New()

		userRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		err := service.UpdatePassword(ctx, userID, "newpassword")

		assert.Error(t, err)
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_Deactivate(t *testing.T) {
	userRepo := new(MockUserRepository)
	passwordHasher := new(MockPasswordHasher)

	service := NewUserService(userRepo, passwordHasher)
	ctx := context.Background()

	t.Run("successful deactivation", func(t *testing.T) {
		userID := uuid.New()
		existingUser := &domain.User{
			ID:       userID,
			Username: "testuser",
			Role:     domain.UserRoleAdmin,
			Active:   true,
		}

		userRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
		userRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		err := service.Deactivate(ctx, userID)

		require.NoError(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("fails when user not found", func(t *testing.T) {
		userID := uuid.New()

		userRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		err := service.Deactivate(ctx, userID)

		assert.Error(t, err)
		userRepo.AssertExpectations(t)
	})
}
