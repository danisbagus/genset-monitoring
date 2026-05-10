package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
)

// DeviceStatusRepository defines the contract for device status persistence.
type DeviceStatusRepository interface {
	// Upsert updates an existing status or creates a new one in a single query.
	Upsert(ctx context.Context, status *model.DeviceStatus) error
	// FindByDeviceID retrieves the current status for a specific device.
	FindByDeviceID(ctx context.Context, deviceID uuid.UUID) (*model.DeviceStatus, error)
}

type deviceStatusRepository struct {
	db *gorm.DB
}

// NewDeviceStatusRepository constructs a DeviceStatusRepository.
func NewDeviceStatusRepository(db *gorm.DB) DeviceStatusRepository {
	return &deviceStatusRepository{db: db}
}

// Upsert uses PostgreSQL's ON CONFLICT clause for high-performance atomic updates.
func (r *deviceStatusRepository) Upsert(ctx context.Context, status *model.DeviceStatus) error {
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"is_online", "last_seen", "gsm_signal", "gps_connected",
			"server_connected", "can_connected", "rs485_connected",
			"sd_card_ok", "updated_at",
		}),
	}).Create(status).Error

	if err != nil {
		return fmt.Errorf("deviceStatusRepository.Upsert: %w", err)
	}
	return nil
}

// FindByDeviceID returns the status record for a device or ErrDeviceStatusNotFound.
func (r *deviceStatusRepository) FindByDeviceID(ctx context.Context, deviceID uuid.UUID) (*model.DeviceStatus, error) {
	var status model.DeviceStatus
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		First(&status).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceStatusNotFound
		}
		return nil, fmt.Errorf("deviceStatusRepository.FindByDeviceID: %w", err)
	}
	return &status, nil
}

// ── Sentinel errors ───────────────────────────────────────────────

var ErrDeviceStatusNotFound = errors.New("device status not found")
