package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
)

// RefreshTokenRepository defines the contract for refresh token persistence.
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	FindByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	RevokeByHash(ctx context.Context, tokenHash string) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) (int64, error)
}

// refreshTokenRepository implements RefreshTokenRepository using GORM.
type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository.
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create inserts a new refresh token record.
func (r *refreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("refreshTokenRepository.Create: %w", err)
	}
	return nil
}

// FindByHash retrieves a refresh token by its SHA-256 hash.
// Returns ErrRefreshTokenNotFound if the hash is unknown.
func (r *refreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("refreshTokenRepository.FindByHash: %w", err)
	}
	return &token, nil
}

// RevokeByHash marks a single refresh token as revoked and records the revocation time.
func (r *refreshTokenRepository) RevokeByHash(ctx context.Context, tokenHash string) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("token_hash = ? AND revoked = FALSE", tokenHash).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": now,
		})

	if result.Error != nil {
		return fmt.Errorf("refreshTokenRepository.RevokeByHash: %w", result.Error)
	}
	return nil
}

// RevokeAllByUserID revokes every active refresh token for a user (logout-all / security wipe).
func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked = FALSE", userID).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": now,
		})

	if result.Error != nil {
		return fmt.Errorf("refreshTokenRepository.RevokeAllByUserID: %w", result.Error)
	}
	return nil
}

// DeleteExpired hard-deletes expired token records (for periodic cleanup jobs).
func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now().UTC()).
		Delete(&model.RefreshToken{})

	if result.Error != nil {
		return 0, fmt.Errorf("refreshTokenRepository.DeleteExpired: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// ── Sentinel errors ───────────────────────────────────────────────

// ErrRefreshTokenNotFound is returned when a token hash lookup yields no result.
var ErrRefreshTokenNotFound = errors.New("refresh token not found")
