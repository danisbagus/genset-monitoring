package model

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken stores hashed refresh tokens for secure token rotation.
// The raw token is NEVER persisted; only its SHA-256 hash is stored.
type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"                       json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID"                              json:"-"`
	TokenHash string     `gorm:"not null;uniqueIndex"                           json:"-"`
	ExpiresAt time.Time  `gorm:"not null"                                       json:"expires_at"`
	Revoked   bool       `gorm:"not null;default:false"                         json:"revoked"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	UserAgent string     `gorm:"size:512"                                       json:"user_agent,omitempty"`
	IPAddress string     `gorm:"size:45"                                        json:"ip_address,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// TableName overrides the table name used by GORM.
func (RefreshToken) TableName() string { return "refresh_tokens" }

// IsExpired returns true if the token has passed its expiry time.
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().UTC().After(rt.ExpiresAt)
}

// IsValid returns true if the token is not revoked and not expired.
func (rt *RefreshToken) IsValid() bool {
	return !rt.Revoked && !rt.IsExpired()
}
