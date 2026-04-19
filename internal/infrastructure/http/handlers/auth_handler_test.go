package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Luis-Lanza/luson/internal/domain"
	"github.com/Luis-Lanza/luson/internal/infrastructure/http/dto"
	"github.com/Luis-Lanza/luson/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of ports.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, username, password string) (*ports.LoginResult, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.LoginResult), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*ports.TokenPair, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.TokenPair), args.Error(1)
}

func (m *MockAuthService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func setupAuthHandlerTest() (*gin.Engine, *MockAuthService, *AuthHandler) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	router := gin.New()
	return router, mockService, handler
}

func TestAuthHandler_Login(t *testing.T) {
	router, mockService, handler := setupAuthHandlerTest()

	router.POST("/api/auth/login", handler.Login)

	t.Run("successful login", func(t *testing.T) {
		userID := uuid.New()
		mockService.On("Login", mock.Anything, "testuser", "password123").Return(&ports.LoginResult{
			User: domain.User{
				ID:       userID,
				Username: "testuser",
				Role:     domain.UserRoleAdmin,
				Active:   true,
			},
			AccessToken:  "access-token-123",
			RefreshToken: "refresh-token-456",
		}, nil).Once()

		reqBody := dto.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		data := response.Data.(map[string]interface{})
		assert.Equal(t, "access-token-123", data["access_token"])
		assert.Equal(t, "refresh-token-456", data["refresh_token"])

		user := data["user"].(map[string]interface{})
		assert.Equal(t, "testuser", user["username"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockService.On("Login", mock.Anything, "wronguser", "wrongpass").Return(nil, errors.New("invalid credentials")).Once()

		reqBody := dto.LoginRequest{
			Username: "wronguser",
			Password: "wrongpass",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("validation error - missing fields", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "testuser",
			// missing password
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Refresh(t *testing.T) {
	router, mockService, handler := setupAuthHandlerTest()

	router.POST("/api/auth/refresh", handler.Refresh)

	t.Run("successful token refresh", func(t *testing.T) {
		mockService.On("RefreshToken", mock.Anything, "valid-refresh-token").Return(&ports.TokenPair{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
		}, nil).Once()

		reqBody := dto.RefreshRequest{
			RefreshToken: "valid-refresh-token",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		mockService.On("RefreshToken", mock.Anything, "invalid-token").Return(nil, errors.New("invalid token")).Once()

		reqBody := dto.RefreshRequest{
			RefreshToken: "invalid-token",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestAuthHandler_Me(t *testing.T) {
	router, mockService, handler := setupAuthHandlerTest()

	router.GET("/api/auth/me", func(c *gin.Context) {
		// Simulate auth middleware setting user_id
		c.Set("user_id", uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"))
		handler.Me(c)
	})

	t.Run("returns current user", func(t *testing.T) {
		userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		mockService.On("GetCurrentUser", mock.Anything, userID).Return(&domain.User{
			ID:       userID,
			Username: "testuser",
			Role:     domain.UserRoleAdmin,
			Active:   true,
		}, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/auth/me", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		mockService.AssertExpectations(t)
	})

	t.Run("returns 401 when user_id not in context", func(t *testing.T) {
		routerNoAuth := gin.New()
		routerNoAuth.GET("/api/auth/me", handler.Me)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/auth/me", nil)

		routerNoAuth.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	router, mockService, handler := setupAuthHandlerTest()

	router.POST("/api/auth/logout", handler.Logout)

	t.Run("successful logout", func(t *testing.T) {
		mockService.On("Logout", mock.Anything, "refresh-token-to-revoke").Return(nil).Once()

		reqBody := dto.RefreshRequest{
			RefreshToken: "refresh-token-to-revoke",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/logout", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		mockService.AssertExpectations(t)
	})
}
