package crypto

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const maxBcryptPasswordBytes = 72

var (
	// ErrEmptyPassword indicates an empty password was supplied for hashing.
	ErrEmptyPassword = errors.New("crypto: password must not be empty")
	// ErrPasswordTooLong indicates a password exceeds bcrypt's 72-byte input limit.
	ErrPasswordTooLong = errors.New("crypto: password exceeds 72 bytes")
)

// HashPassword hashes a non-empty password using bcrypt with bcrypt.DefaultCost.
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}
	if len(password) > maxBcryptPasswordBytes {
		return "", fmt.Errorf("%w: %w", ErrPasswordTooLong, bcrypt.ErrPasswordTooLong)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("crypto: hash password: %w", err)
	}
	return string(hashed), nil
}

// VerifyPassword reports whether password matches the bcrypt hash.
func VerifyPassword(hash, password string) bool {
	if hash == "" || password == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
