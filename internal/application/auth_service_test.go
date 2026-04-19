package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/infrastructure/jwt"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserRepository is a mock implementation of ports.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, filter ports.UserFilter) ([]domain.User, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockTokenRepository is a mock implementation of ports.TokenRepository
type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, tokenHash, expiresAt)
	return args.Error(0)
}

func (m *MockTokenRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*ports.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.RefreshToken), args.Error(1)
}

func (m *MockTokenRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteUserTokens(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestAuthService_Login(t *testing.T) {
	userRepo := new(MockUserRepository)
	tokenRepo := new(MockTokenRepository)

	service := NewAuthService(userRepo, tokenRepo, "access-secret-32-bytes-long-key", "refresh-secret-32-bytes-long")

	ctx := context.Background()

	t.Run("successful login with valid credentials", func(t *testing.T) {
		password := "password123"
		passwordHash, _ := jwt.HashPassword(password)
		userID := uuid.New()

		user := &domain.User{
			ID:           userID,
			Username:     "testuser",
			PasswordHash: passwordHash,
			Role:         domain.UserRoleAdmin,
			Active:       true,
		}

		userRepo.On("FindByUsername", ctx, "testuser").Return(user, nil).Once()
		tokenRepo.On("SaveRefreshToken", ctx, userID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil).Once()

		result, err := service.Login(ctx, "testuser", password)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.User.ID)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
		assert.Empty(t, result.User.PasswordHash) // Should not expose password hash
		userRepo.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})

	t.Run("login fails with invalid password", func(t *testing.T) {
		passwordHash, _ := jwt.HashPassword("correctpassword")
		user := &domain.User{
			ID:           uuid.New(),
			Username:     "testuser",
			PasswordHash: passwordHash,
			Role:         domain.UserRoleAdmin,
			Active:       true,
		}

		userRepo.On("FindByUsername", ctx, "testuser").Return(user, nil).Once()

		result, err := service.Login(ctx, "testuser", "wrongpassword")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid credentials")
		userRepo.AssertExpectations(t)
	})

	t.Run("login fails with non-existent user", func(t *testing.T) {
		userRepo.On("FindByUsername", ctx, "nonexistent").Return(nil, errors.New("user not found")).Once()

		result, err := service.Login(ctx, "nonexistent", "password")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid credentials")
		userRepo.AssertExpectations(t)
	})

	t.Run("login fails with inactive user", func(t *testing.T) {
		password := "password123"
		passwordHash, _ := jwt.HashPassword(password)
		user := &domain.User{
			ID:           uuid.New(),
			Username:     "inactiveuser",
			PasswordHash: passwordHash,
			Role:         domain.UserRoleAdmin,
			Active:       false,
		}

		userRepo.On("FindByUsername", ctx, "inactiveuser").Return(user, nil).Once()

		result, err := service.Login(ctx, "inactiveuser", password)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "account is deactivated")
		userRepo.AssertExpectations(t)
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	userRepo := new(MockUserRepository)
	tokenRepo := new(MockTokenRepository)

	service := NewAuthService(userRepo, tokenRepo, "access-secret-32-bytes-long-key", "refresh-secret-32-bytes-long")

	ctx := context.Background()

	t.Run("successful token refresh", func(t *testing.T) {
		userID := uuid.New()
		user := &domain.User{
			ID:       userID,
			Username: "testuser",
			Role:     domain.UserRoleCajero,
			Active:   true,
		}

		// Generate a valid refresh token
		refreshToken, _ := jwt.GenerateRefreshToken("refresh-secret-32-bytes-long")
		tokenHash := jwt.HashRefreshToken(refreshToken)

		tokenRepo.On("FindRefreshToken", ctx, tokenHash).Return(&ports.RefreshToken{
			ID:        uuid.New(),
			UserID:    userID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}, nil).Once()

		userRepo.On("FindByID", ctx, userID).Return(user, nil).Once()

		// Should delete old token
		tokenRepo.On("DeleteRefreshToken", ctx, tokenHash).Return(nil).Once()

		// Should save new token
		tokenRepo.On("SaveRefreshToken", ctx, userID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil).Once()

		result, err := service.RefreshToken(ctx, refreshToken)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
		tokenRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("refresh fails with invalid token", func(t *testing.T) {
		result, err := service.RefreshToken(ctx, "invalid-token")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("refresh fails with revoked token", func(t *testing.T) {
		refreshToken, _ := jwt.GenerateRefreshToken("refresh-secret-32-bytes-long")
		tokenHash := jwt.HashRefreshToken(refreshToken)

		tokenRepo.On("FindRefreshToken", ctx, tokenHash).Return(nil, errors.New("token not found")).Once()

		result, err := service.RefreshToken(ctx, refreshToken)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "token has been revoked")
		tokenRepo.AssertExpectations(t)
	})
}

func TestAuthService_GetCurrentUser(t *testing.T) {
	userRepo := new(MockUserRepository)
	tokenRepo := new(MockTokenRepository)

	service := NewAuthService(userRepo, tokenRepo, "access-secret", "refresh-secret")

	ctx := context.Background()

	t.Run("returns user without password hash", func(t *testing.T) {
		userID := uuid.New()
		user := &domain.User{
			ID:           userID,
			Username:     "testuser",
			PasswordHash: "should-not-be-returned",
			Role:         domain.UserRoleAdmin,
			Active:       true,
		}

		userRepo.On("FindByID", ctx, userID).Return(user, nil).Once()

		result, err := service.GetCurrentUser(ctx, userID)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
		assert.Empty(t, result.PasswordHash)
		userRepo.AssertExpectations(t)
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		userID := uuid.New()

		userRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		result, err := service.GetCurrentUser(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		userRepo.AssertExpectations(t)
	})
}

func TestAuthService_Logout(t *testing.T) {
	userRepo := new(MockUserRepository)
	tokenRepo := new(MockTokenRepository)

	service := NewAuthService(userRepo, tokenRepo, "access-secret", "refresh-secret")

	ctx := context.Background()

	t.Run("successful logout deletes token", func(t *testing.T) {
		refreshToken := "some-refresh-token"
		tokenHash := jwt.HashRefreshToken(refreshToken)

		tokenRepo.On("DeleteRefreshToken", ctx, tokenHash).Return(nil).Once()

		err := service.Logout(ctx, refreshToken)

		require.NoError(t, err)
		tokenRepo.AssertExpectations(t)
	})

	t.Run("logout is idempotent - no error if token already deleted", func(t *testing.T) {
		refreshToken := "already-deleted-token"
		tokenHash := jwt.HashRefreshToken(refreshToken)

		// If token doesn't exist, DeleteRefreshToken should not return error
		tokenRepo.On("DeleteRefreshToken", ctx, tokenHash).Return(nil).Once()

		err := service.Logout(ctx, refreshToken)

		require.NoError(t, err)
		tokenRepo.AssertExpectations(t)
	})
}
