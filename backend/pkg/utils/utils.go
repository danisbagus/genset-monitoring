package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/google/uuid"
)

// NewUUID generates a new v4 UUID string.
func NewUUID() string {
	return uuid.New().String()
}

// ParseUUID parses and validates a UUID string.
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// NowUTC returns the current time in UTC.
func NowUTC() time.Time {
	return time.Now().UTC()
}

// GenerateRandomHex generates a cryptographically random hex string of n bytes.
func GenerateRandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// SlugFromString converts a string to a URL-friendly slug.
func SlugFromString(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
