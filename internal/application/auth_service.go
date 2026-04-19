package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/infrastructure/jwt"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/google/uuid"
)

// authService implements the ports.AuthService interface.
type authService struct {
	userRepo   ports.UserRepository
	tokenRepo  ports.TokenRepository
	jwtSecret  string
	jwtRefresh string
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(
	userRepo ports.UserRepository,
	tokenRepo ports.TokenRepository,
	jwtSecret string,
	jwtRefresh string,
) ports.AuthService {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtSecret:  jwtSecret,
		jwtRefresh: jwtRefresh,
	}
}

// Login authenticates a user and returns tokens.
func (s *authService) Login(ctx context.Context, username, password string) (*ports.LoginResult, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !jwt.CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.Active {
		return nil, errors.New("account is deactivated")
	}

	// Generate tokens
	var branchIDStr *string
	if user.BranchID != nil {
		id := user.BranchID.String()
		branchIDStr = &id
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, string(user.Role), branchIDStr, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(s.jwtRefresh)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash and store refresh token
	tokenHash := jwt.HashRefreshToken(refreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days
	if err := s.tokenRepo.SaveRefreshToken(ctx, user.ID, tokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// Return user without password hash
	return &ports.LoginResult{
		User: domain.User{
			ID:        user.ID,
			Username:  user.Username,
			Role:      user.Role,
			BranchID:  user.BranchID,
			Active:    user.Active,
			CreatedAt: user.CreatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken validates a refresh token and issues new tokens.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*ports.TokenPair, error) {
	// Validate the refresh token
	_, err := jwt.ValidateRefreshToken(refreshToken, s.jwtRefresh)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Check if token exists in database (not revoked)
	tokenHash := jwt.HashRefreshToken(refreshToken)
	storedToken, err := s.tokenRepo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("token has been revoked")
	}

	// Check if token is expired
	if storedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate new tokens
	var branchIDStr *string
	if user.BranchID != nil {
		id := user.BranchID.String()
		branchIDStr = &id
	}

	newAccessToken, err := jwt.GenerateAccessToken(user.ID, string(user.Role), branchIDStr, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := jwt.GenerateRefreshToken(s.jwtRefresh)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Delete old token (token rotation)
	if err := s.tokenRepo.DeleteRefreshToken(ctx, tokenHash); err != nil {
		// Log error but don't fail - user still gets new tokens
		// In production, this should be logged
		_ = err
	}

	// Store new refresh token
	newTokenHash := jwt.HashRefreshToken(newRefreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.tokenRepo.SaveRefreshToken(ctx, user.ID, newTokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &ports.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// GetCurrentUser retrieves the current user by ID.
func (s *authService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Return user without password hash
	return &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		BranchID:  user.BranchID,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Logout revokes a refresh token.
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := jwt.HashRefreshToken(refreshToken)

	// Delete the token - if it doesn't exist, it's already revoked (idempotent)
	if err := s.tokenRepo.DeleteRefreshToken(ctx, tokenHash); err != nil {
		// If token doesn't exist, we consider logout successful
		// This makes logout idempotent
		return nil
	}

	return nil
}
