package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
)

// UserRepository defines the contract for user persistence operations.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// userRepository implements UserRepository using GORM.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user record.
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("userRepository.Create: %w", err)
	}
	return nil
}

// FindByID returns a user by its UUID primary key. Returns ErrUserNotFound if absent.
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByID: %w", err)
	}
	return &user, nil
}

// FindByUsername returns a user by username. Returns ErrUserNotFound if absent.
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("username = ? AND deleted_at IS NULL", username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByUsername: %w", err)
	}
	return &user, nil
}

// FindByEmail returns a user by email. Returns ErrUserNotFound if absent.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByEmail: %w", err)
	}
	return &user, nil
}

// ExistsByUsername returns true if a non-deleted user with the given username exists.
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("username = ? AND deleted_at IS NULL", username).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("userRepository.ExistsByUsername: %w", err)
	}
	return count > 0, nil
}

// ExistsByEmail returns true if a non-deleted user with the given email exists.
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("userRepository.ExistsByEmail: %w", err)
	}
	return count > 0, nil
}

// ── Sentinel errors ───────────────────────────────────────────────

// ErrUserNotFound is returned when a user lookup yields no result.
var ErrUserNotFound = errors.New("user not found")
