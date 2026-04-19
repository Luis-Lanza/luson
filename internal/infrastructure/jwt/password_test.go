package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Run("creates bcrypt hash with cost 12", func(t *testing.T) {
		hash, err := HashPassword("mypassword123")
		require.NoError(t, err)
		assert.NotEmpty(t, hash)

		// Verify it's a valid bcrypt hash
		cost, err := bcrypt.Cost([]byte(hash))
		require.NoError(t, err)
		assert.Equal(t, bcrypt.DefaultCost, cost)
	})

	t.Run("returns different hash each time (unique salt)", func(t *testing.T) {
		hash1, err := HashPassword("samepassword")
		require.NoError(t, err)

		hash2, err := HashPassword("samepassword")
		require.NoError(t, err)

		// Hashes should be different due to unique salts
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("handles empty password", func(t *testing.T) {
		hash, err := HashPassword("")
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("handles long passwords", func(t *testing.T) {
		longPassword := "this-is-a-very-long-password-with-many-characters-1234567890-!@#$%^&*()"
		hash, err := HashPassword(longPassword)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)

		// Verify it can be checked
		match := CheckPassword(longPassword, hash)
		assert.True(t, match)
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("returns true for correct password", func(t *testing.T) {
		password := "correctpassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		match := CheckPassword(password, hash)
		assert.True(t, match)
	})

	t.Run("returns false for incorrect password", func(t *testing.T) {
		password := "correctpassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		match := CheckPassword("wrongpassword", hash)
		assert.False(t, match)
	})

	t.Run("returns false for empty password check", func(t *testing.T) {
		password := "somepassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		match := CheckPassword("", hash)
		assert.False(t, match)
	})

	t.Run("returns false for invalid hash", func(t *testing.T) {
		match := CheckPassword("password", "not-a-valid-hash")
		assert.False(t, match)
	})

	t.Run("consistent behavior across multiple checks", func(t *testing.T) {
		password := "testpassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		// Check multiple times, should always return true
		for i := 0; i < 5; i++ {
			assert.True(t, CheckPassword(password, hash))
		}
	})
}
