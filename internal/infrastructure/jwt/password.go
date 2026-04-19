package jwt

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher provides password hashing operations.
type PasswordHasher struct{}

// NewPasswordHasher creates a new password hasher.
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

// HashPassword creates a bcrypt hash of a password.
// Uses bcrypt.DefaultCost (currently 10) for hashing.
func (p *PasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a plaintext password with a bcrypt hash.
// Returns true if the password matches the hash.
func (p *PasswordHasher) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Backward compatibility - package-level functions
var defaultHasher = NewPasswordHasher()

// HashPassword creates a bcrypt hash of a password (package-level function).
func HashPassword(password string) (string, error) {
	return defaultHasher.HashPassword(password)
}

// CheckPassword compares a plaintext password with a bcrypt hash (package-level function).
func CheckPassword(password, hash string) bool {
	return defaultHasher.CheckPassword(password, hash)
}
