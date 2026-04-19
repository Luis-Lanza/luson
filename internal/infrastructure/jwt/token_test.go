package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAccessToken(t *testing.T) {
	secret := "test-secret-key-that-is-32-bytes-long-for-hs256"
	userID := uuid.New()
	role := "admin"
	branchID := uuid.New().String()

	t.Run("generates valid access token", func(t *testing.T) {
		token, err := GenerateAccessToken(userID, role, &branchID, secret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("token contains correct claims", func(t *testing.T) {
		token, err := GenerateAccessToken(userID, role, &branchID, secret)
		require.NoError(t, err)

		claims, err := ValidateAccessToken(token, secret)
		require.NoError(t, err)

		assert.Equal(t, userID.String(), claims.UserID)
		assert.Equal(t, role, claims.Role)
		assert.Equal(t, branchID, *claims.BranchID)
	})

	t.Run("token expires in 15 minutes", func(t *testing.T) {
		token, err := GenerateAccessToken(userID, role, &branchID, secret)
		require.NoError(t, err)

		claims, err := ValidateAccessToken(token, secret)
		require.NoError(t, err)

		// Check expiration is approximately 15 minutes from now
		expectedExpiry := time.Now().Add(15 * time.Minute)
		diff := claims.ExpiresAt.Time.Sub(expectedExpiry)
		assert.Less(t, diff.Abs(), 5*time.Second) // Within 5 seconds
	})

	t.Run("nil branch_id is allowed", func(t *testing.T) {
		token, err := GenerateAccessToken(userID, "admin", nil, secret)
		require.NoError(t, err)

		claims, err := ValidateAccessToken(token, secret)
		require.NoError(t, err)

		assert.Nil(t, claims.BranchID)
	})
}

func TestGenerateRefreshToken(t *testing.T) {
	secret := "test-refresh-secret-that-is-32-bytes-long-for-hs256"

	t.Run("generates valid refresh token", func(t *testing.T) {
		token, err := GenerateRefreshToken(secret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generates unique tokens each time", func(t *testing.T) {
		token1, err := GenerateRefreshToken(secret)
		require.NoError(t, err)

		token2, err := GenerateRefreshToken(secret)
		require.NoError(t, err)

		assert.NotEqual(t, token1, token2)
	})

	t.Run("refresh token expires in 7 days", func(t *testing.T) {
		token, err := GenerateRefreshToken(secret)
		require.NoError(t, err)

		claims, err := ValidateRefreshToken(token, secret)
		require.NoError(t, err)

		// Check expiration is approximately 7 days from now
		expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
		diff := claims.ExpiresAt.Time.Sub(expectedExpiry)
		assert.Less(t, diff.Abs(), 5*time.Second) // Within 5 seconds
	})
}

func TestValidateAccessToken(t *testing.T) {
	secret := "test-secret-key-that-is-32-bytes-long-for-hs256"
	wrongSecret := "wrong-secret-key-that-is-32-bytes-long!!"
	userID := uuid.New()
	role := "cajero"
	branchID := uuid.New().String()

	t.Run("validates correct token", func(t *testing.T) {
		token, err := GenerateAccessToken(userID, role, &branchID, secret)
		require.NoError(t, err)

		claims, err := ValidateAccessToken(token, secret)
		require.NoError(t, err)
		assert.NotNil(t, claims)
	})

	t.Run("fails with invalid signature", func(t *testing.T) {
		token, err := GenerateAccessToken(userID, role, &branchID, secret)
		require.NoError(t, err)

		_, err = ValidateAccessToken(token, wrongSecret)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature")
	})

	t.Run("fails with expired token", func(t *testing.T) {
		// Create an expired token manually
		now := time.Now()
		expiredClaims := AccessTokenClaims{
			UserID:   userID.String(),
			Role:     role,
			BranchID: &branchID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		tokenString, err := token.SignedString([]byte(secret))
		require.NoError(t, err)

		_, err = ValidateAccessToken(tokenString, secret)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expired")
	})

	t.Run("fails with malformed token", func(t *testing.T) {
		_, err := ValidateAccessToken("not.a.valid.token", secret)
		assert.Error(t, err)
	})
}

func TestHashRefreshToken(t *testing.T) {
	t.Run("generates consistent hash", func(t *testing.T) {
		token := "my-refresh-token-123"
		hash1 := HashRefreshToken(token)
		hash2 := HashRefreshToken(token)

		assert.Equal(t, hash1, hash2)
		assert.Equal(t, 64, len(hash1)) // SHA-256 hex is 64 chars
	})

	t.Run("different tokens produce different hashes", func(t *testing.T) {
		hash1 := HashRefreshToken("token-1")
		hash2 := HashRefreshToken("token-2")

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("hash is hexadecimal", func(t *testing.T) {
		hash := HashRefreshToken("any-token")
		assert.Regexp(t, "^[a-f0-9]{64}$", hash)
	})
}

func TestTokenClaims(t *testing.T) {
	t.Run("access token claims structure", func(t *testing.T) {
		branchID := uuid.New().String()
		claims := AccessTokenClaims{
			UserID:   uuid.New().String(),
			Role:     "admin",
			BranchID: &branchID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		assert.NotEmpty(t, claims.UserID)
		assert.Equal(t, "admin", claims.Role)
		assert.NotNil(t, claims.BranchID)
	})
}
