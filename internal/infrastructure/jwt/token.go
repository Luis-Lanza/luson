package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AccessTokenClaims represents the claims in an access token.
type AccessTokenClaims struct {
	UserID   string  `json:"user_id"`
	Role     string  `json:"role"`
	BranchID *string `json:"branch_id,omitempty"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents the claims in a refresh token.
type RefreshTokenClaims struct {
	TokenID string `json:"token_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a new JWT access token for a user.
// The token expires in 15 minutes and uses HS256 signing.
func GenerateAccessToken(userID uuid.UUID, role string, branchID *string, secret string) (string, error) {
	now := time.Now()
	claims := AccessTokenClaims{
		UserID:   userID.String(),
		Role:     role,
		BranchID: branchID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken creates a new JWT refresh token.
// The token expires in 7 days and uses HS256 signing.
func GenerateRefreshToken(secret string) (string, error) {
	now := time.Now()
	claims := RefreshTokenClaims{
		TokenID: uuid.New().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateAccessToken validates an access token and returns its claims.
func ValidateAccessToken(tokenString string, secret string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// ValidateRefreshToken validates a refresh token and returns its claims.
func ValidateRefreshToken(tokenString string, secret string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// HashRefreshToken creates a SHA-256 hash of a refresh token.
// This is used to store the token hash in the database instead of the raw token.
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
