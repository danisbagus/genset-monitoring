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
	// UpsertLatestState updates the realtime snapshot of a device.
	UpsertLatestState(ctx context.Context, state *model.DeviceLatestState) error
	// FindByDeviceID retrieves the current status for a specific device.
	FindByDeviceID(ctx context.Context, deviceID uuid.UUID) (*model.DeviceStatus, error)
	// FindLatestStateByDeviceID retrieves the latest state snapshot for a specific device.
	FindLatestStateByDeviceID(ctx context.Context, deviceID uuid.UUID) (*model.DeviceLatestState, error)
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
			"last_seen", "gsm_signal", "gps_connected",
			"server_connected", "can_connected", "rs485_connected",
			"sd_card_ok", "updated_at",
		}),
	}).Create(status).Error

	if err != nil {
		return fmt.Errorf("deviceStatusRepository.Upsert: %w", err)
	}
	return nil
}

// UpsertLatestState implements the complex partial update logic for device_latest_state.
func (r *deviceStatusRepository) UpsertLatestState(ctx context.Context, state *model.DeviceLatestState) error {
	query := `
		INSERT INTO device_latest_state (
			device_id, 
			speed, coolant_temperature, oil_pressure, fuel_level,
			frequency, total_va, pf_avg, 
			engine_running,
			telemetry_recorded_at, last_online_at, last_seen_at, updated_at
		) VALUES (
			?, 
			?, ?, ?, ?,
			?, ?, ?,
			?, 
			?, ?, NOW(), NOW()
		)
		ON CONFLICT (device_id) DO UPDATE SET
			speed = COALESCE(EXCLUDED.speed, device_latest_state.speed),
			coolant_temperature = COALESCE(EXCLUDED.coolant_temperature, device_latest_state.coolant_temperature),
			oil_pressure = COALESCE(EXCLUDED.oil_pressure, device_latest_state.oil_pressure),
			fuel_level = COALESCE(EXCLUDED.fuel_level, device_latest_state.fuel_level),
			frequency = COALESCE(EXCLUDED.frequency, device_latest_state.frequency),
			total_va = COALESCE(EXCLUDED.total_va, device_latest_state.total_va),
			pf_avg = COALESCE(EXCLUDED.pf_avg, device_latest_state.pf_avg),
			engine_running = COALESCE(
				CASE WHEN EXCLUDED.speed IS NOT NULL THEN EXCLUDED.speed >= 300 ELSE NULL END, 
				device_latest_state.engine_running
			),
			telemetry_recorded_at = COALESCE(EXCLUDED.telemetry_recorded_at, device_latest_state.telemetry_recorded_at),
			last_online_at = COALESCE(EXCLUDED.last_online_at, device_latest_state.last_online_at),
			last_seen_at = NOW(),
			updated_at = NOW()
	`

	// Note: engine_running in VALUES is only used for initial INSERT.
	// We calculate it here for the initial insert case.
	var initialEngineRunning bool
	if state.Speed != nil {
		initialEngineRunning = *state.Speed >= 300
	}

	err := r.db.WithContext(ctx).Exec(query,
		state.DeviceID,
		state.Speed, state.CoolantTemperature, state.OilPressure, state.FuelLevel,
		state.Frequency, state.TotalVa, state.PfAvg,
		initialEngineRunning,
		state.TelemetryRecordedAt,
		state.LastOnlineAt,
	).Error

	if err != nil {
		return fmt.Errorf("deviceStatusRepository.UpsertLatestState: %w", err)
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

// FindLatestStateByDeviceID returns the latest state record for a device or ErrDeviceStatusNotFound.
func (r *deviceStatusRepository) FindLatestStateByDeviceID(ctx context.Context, deviceID uuid.UUID) (*model.DeviceLatestState, error) {
	var state model.DeviceLatestState
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		First(&state).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceStatusNotFound
		}
		return nil, fmt.Errorf("deviceStatusRepository.FindLatestStateByDeviceID: %w", err)
	}
	return &state, nil
}

// ── Sentinel errors ───────────────────────────────────────────────

var ErrDeviceStatusNotFound = errors.New("device status not found")
