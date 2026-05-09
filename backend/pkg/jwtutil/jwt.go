// Package jwtutil provides JWT access token generation and validation utilities.
// It wraps github.com/golang-jwt/jwt/v5 and enforces HS256 signing.
package jwtutil

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims holds the JWT payload embedded in every access token.
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Manager handles JWT signing and parsing.
type Manager struct {
	secret     []byte
	expiration time.Duration
}

// NewManager creates a new JWT Manager.
// secret must be a sufficiently long random value (≥32 bytes recommended).
func NewManager(secret string, expiration time.Duration) *Manager {
	return &Manager{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

// Generate creates a signed JWT access token for the given user attributes.
func (m *Manager) Generate(userID, username, role string) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiration)),
			ID:        uuid.New().String(), // jti — prevents token reuse after revocation (future)
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("jwtutil: failed to sign token: %w", err)
	}
	return signed, nil
}

// Parse validates and parses a JWT access token, returning its claims.
// Returns an error if the token is invalid, expired, or tampered.
func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwtutil: unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		// Distinguish expiry from other errors for cleaner API responses
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ExpirationSeconds returns the configured expiration duration in seconds.
func (m *Manager) ExpirationSeconds() int64 {
	return int64(m.expiration.Seconds())
}

// Sentinel errors for JWT operations.
var (
	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token has expired")
)
