// Package hashutil provides secure password hashing via bcrypt
// and opaque token hashing via SHA-256.
package hashutil

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost is the bcrypt work factor. 12 is OWASP-recommended minimum for 2024.
	// Increase over time as hardware improves; each +1 doubles the work.
	BcryptCost = 12

	// RefreshTokenBytes is the number of cryptographically random bytes in
	// the opaque refresh token before hex-encoding (= 64 hex chars).
	RefreshTokenBytes = 32
)

// HashPassword hashes a plain-text password using bcrypt.
func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("hashutil: failed to hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPassword returns nil if plain matches the bcrypt hash, else an error.
func CheckPassword(plain, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)); err != nil {
		return ErrPasswordMismatch
	}
	return nil
}

// GenerateRefreshToken generates a cryptographically secure opaque token
// (hex-encoded random bytes) and its SHA-256 hash for DB storage.
// Returns (rawToken, tokenHash, error).
// Only rawToken is sent to the client; only tokenHash is stored in the DB.
func GenerateRefreshToken() (rawToken, tokenHash string, err error) {
	b := make([]byte, RefreshTokenBytes)
	if _, err = rand.Read(b); err != nil {
		return "", "", fmt.Errorf("hashutil: failed to generate random bytes: %w", err)
	}
	rawToken = hex.EncodeToString(b)
	tokenHash = HashSHA256(rawToken)
	return rawToken, tokenHash, nil
}

// HashSHA256 returns the hex-encoded SHA-256 hash of the input string.
// Used to store/lookup refresh tokens without exposing raw values.
func HashSHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// ErrPasswordMismatch is returned when bcrypt comparison fails.
var ErrPasswordMismatch = fmt.Errorf("hashutil: password does not match")
