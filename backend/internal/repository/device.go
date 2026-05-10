package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
)

// ── DTOs used only in the repository layer ────────────────────────

// DeviceFilter carries the validated parameters for list queries.
type DeviceFilter struct {
	Search   string // matches device_code OR name (ILIKE)
	IsOnline *bool  // nil = no filter, true/false = filter by online state
	Status   string // empty = no filter; values: active | inactive | maintenance
	Page     int    // 1-based
	Limit    int    // rows per page (max enforced at service layer)
	SortBy   string // column name (allow-listed in service layer)
	SortDir  string // "asc" | "desc"
}

// DeviceListResult is the raw result returned from FindAll.
type DeviceListResult struct {
	Devices []*model.Device
	Total   int64
}

// ── Interface ─────────────────────────────────────────────────────

// DeviceRepository defines the persistence contract for devices.
type DeviceRepository interface {
	// Create inserts a new device inside an optional caller-supplied
	// transaction. Pass nil to let the repository create its own session.
	Create(ctx context.Context, device *model.Device) error

	// FindAll returns a paginated, filtered, sorted list of devices.
	FindAll(ctx context.Context, filter DeviceFilter) (*DeviceListResult, error)

	// FindByID returns the device with the given UUID.
	// Returns ErrDeviceNotFound when absent or soft-deleted.
	FindByID(ctx context.Context, id uuid.UUID) (*model.Device, error)

	// Update applies a partial map of column→value updates for the
	// device with the given UUID. Uses GORM's Model().Updates() which
	// automatically sets updated_at.
	Update(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error

	// SoftDelete sets deleted_at to NOW() for the given device.
	SoftDelete(ctx context.Context, id uuid.UUID) error

	// ExistsByDeviceCode checks uniqueness among non-deleted rows.
	ExistsByDeviceCode(ctx context.Context, code string, excludeID *uuid.UUID) (bool, error)

	// ExistsBySerialNumber checks uniqueness among non-deleted rows.
	ExistsBySerialNumber(ctx context.Context, serial string, excludeID *uuid.UUID) (bool, error)
}

// ── Implementation ────────────────────────────────────────────────

type deviceRepository struct {
	db *gorm.DB
}

// NewDeviceRepository constructs a DeviceRepository backed by GORM.
func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

// Create inserts a new device record.
func (r *deviceRepository) Create(ctx context.Context, device *model.Device) error {
	if err := r.db.WithContext(ctx).Create(device).Error; err != nil {
		return fmt.Errorf("deviceRepository.Create: %w", err)
	}
	return nil
}

// FindAll returns a paginated list with total count.
//
// Query optimisations:
//   - Uses trigram GIN indexes via ILIKE for search.
//   - Partial indexes on is_online and status (deleted_at IS NULL predicate).
//   - COUNT uses a separate query so the main fetch does not carry an
//     expensive COUNT(*) OVER().
func (r *deviceRepository) FindAll(ctx context.Context, filter DeviceFilter) (*DeviceListResult, error) {
	base := r.db.WithContext(ctx).
		Model(&model.Device{}).
		Where("deleted_at IS NULL")

	// ── Search ────────────────────────────────────────────────────
	if filter.Search != "" {
		pattern := "%" + filter.Search + "%"
		base = base.Where("device_code ILIKE ? OR name ILIKE ? OR engine_id ILIKE ? OR gsm_number ILIKE ?", pattern, pattern, pattern, pattern)
	}

	// ── Online filter ─────────────────────────────────────────────
	if filter.IsOnline != nil {
		base = base.Where("is_online = ?", *filter.IsOnline)
	}

	// ── Status filter ─────────────────────────────────────────────
	if filter.Status != "" {
		base = base.Where("status = ?", filter.Status)
	}

	// ── Total count (re-uses same WHERE clauses) ───────────────────
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("deviceRepository.FindAll count: %w", err)
	}

	// ── Sorting ───────────────────────────────────────────────────
	sortCol := "created_at" // safe default
	if filter.SortBy != "" {
		sortCol = filter.SortBy
	}
	sortDir := "desc"
	if filter.SortDir == "asc" {
		sortDir = "asc"
	}
	orderClause := fmt.Sprintf("%s %s", sortCol, sortDir)

	// ── Pagination ────────────────────────────────────────────────
	offset := (filter.Page - 1) * filter.Limit

	var devices []*model.Device
	if err := base.
		Order(orderClause).
		Limit(filter.Limit).
		Offset(offset).
		Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("deviceRepository.FindAll fetch: %w", err)
	}

	return &DeviceListResult{Devices: devices, Total: total}, nil
}

// FindByID returns a single device by UUID.
func (r *deviceRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	var device model.Device
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("deviceRepository.FindByID: %w", err)
	}
	return &device, nil
}

// Update applies a selective column update via a plain map to avoid
// zero-value gotchas from struct-based GORM updates.
func (r *deviceRepository) Update(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&model.Device{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(fields)
	if result.Error != nil {
		return fmt.Errorf("deviceRepository.Update: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

// SoftDelete sets deleted_at without removing the row.
func (r *deviceRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&model.Device{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()"))
	if result.Error != nil {
		return fmt.Errorf("deviceRepository.SoftDelete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

// ExistsByDeviceCode returns true when an active device with the given code exists.
// Pass a non-nil excludeID to skip the record being updated (for PATCH uniqueness checks).
func (r *deviceRepository) ExistsByDeviceCode(ctx context.Context, code string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).
		Model(&model.Device{}).
		Where("device_code = ? AND deleted_at IS NULL", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("deviceRepository.ExistsByDeviceCode: %w", err)
	}
	return count > 0, nil
}

// ExistsBySerialNumber returns true when an active device with the given serial exists.
func (r *deviceRepository) ExistsBySerialNumber(ctx context.Context, serial string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).
		Model(&model.Device{}).
		Where("serial_number = ? AND deleted_at IS NULL", serial)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("deviceRepository.ExistsBySerialNumber: %w", err)
	}
	return count > 0, nil
}

// ── Sentinel errors ───────────────────────────────────────────────

// ErrDeviceNotFound is returned when a device lookup yields no result.
var ErrDeviceNotFound = errors.New("device not found")
