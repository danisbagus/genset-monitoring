package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
)

// ── DTOs ─────────────────────────────────────────────────────────

// MonitoringDeviceFilter carries validated query params for the monitoring list.
type MonitoringDeviceFilter struct {
	Search        string
	EngineRunning *bool
	Status        string  // device_status (active|inactive|maintenance)
	Online        *bool   // true = last_seen_at >= now()-5min, false = otherwise
	Limit         int
	Offset        int
	SortBy        string // column alias used in ORDER BY
	SortOrder     string // asc | desc
}

// MonitoringDeviceListResult is the raw result returned from GetMonitoringDevices.
type MonitoringDeviceListResult struct {
	Devices []model.MonitoringDevice
	Total   int64
}

// MonitoringRepository defines the contract for monitoring data operations.
type MonitoringRepository interface {
	// GetMonitoringDevices returns the paginated monitoring list.
	GetMonitoringDevices(ctx context.Context, filter MonitoringDeviceFilter) (*MonitoringDeviceListResult, error)

	// GetMonitoringDeviceByID returns the full realtime detail for a device.
	GetMonitoringDeviceByID(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringDevice, error)

	// GetMonitoringDeviceInfo returns basic device info (section A).
	GetMonitoringDeviceInfo(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringDeviceInfo, error)

	// GetMonitoringLatestState returns realtime state (section B).
	GetMonitoringLatestState(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringLatestState, error)

	// GetMonitoringConnectivity returns connectivity status (section C).
	GetMonitoringConnectivity(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringConnectivity, error)
}

type monitoringRepository struct {
	db *gorm.DB
}

// NewMonitoringRepository constructs a MonitoringRepository.
func NewMonitoringRepository(db *gorm.DB) MonitoringRepository {
	return &monitoringRepository{db: db}
}

// allowedMonitoringSortColumns maps user-facing sort_by values to safe SQL expressions.
var allowedMonitoringSortColumns = map[string]string{
	"updated_at":           "dls.updated_at",
	"last_seen_at":         "dls.last_seen_at",
	"telemetry_recorded_at": "dls.telemetry_recorded_at",
	"name":                 "d.name",
}

// GetMonitoringDevices returns a paginated list of monitoring devices.
// Single optimised query – no N+1.
func (r *monitoringRepository) GetMonitoringDevices(ctx context.Context, filter MonitoringDeviceFilter) (*MonitoringDeviceListResult, error) {
	// ── Build base WHERE conditions ──────────────────────────────
	conditions := []string{"d.deleted_at IS NULL"}
	args := []interface{}{}

	if filter.Search != "" {
		pattern := "%" + filter.Search + "%"
		conditions = append(conditions, "(d.name ILIKE ? OR d.device_code ILIKE ? OR d.serial_number ILIKE ?)")
		args = append(args, pattern, pattern, pattern)
	}

	if filter.Status != "" {
		conditions = append(conditions, "d.status = ?")
		args = append(args, filter.Status)
	}

	if filter.EngineRunning != nil {
		conditions = append(conditions, "COALESCE(dls.engine_running, false) = ?")
		args = append(args, *filter.EngineRunning)
	}

	if filter.Online != nil {
		if *filter.Online {
			conditions = append(conditions, "dls.last_seen_at >= (NOW() - INTERVAL '5 minutes')")
		} else {
			conditions = append(conditions, "(dls.last_seen_at IS NULL OR dls.last_seen_at < (NOW() - INTERVAL '5 minutes'))")
		}
	}

	whereClause := strings.Join(conditions, " AND ")

	// ── Count ────────────────────────────────────────────────────
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM devices d
		LEFT JOIN device_latest_state dls ON d.id = dls.device_id
		LEFT JOIN device_statuses ds     ON d.id = ds.device_id
		WHERE %s
	`, whereClause)

	var total int64
	if err := r.db.WithContext(ctx).Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("monitoringRepository.GetMonitoringDevices count: %w", err)
	}

	// ── ORDER BY ─────────────────────────────────────────────────
	sortCol := "dls.updated_at" // safe default
	if col, ok := allowedMonitoringSortColumns[filter.SortBy]; ok {
		sortCol = col
	}
	sortOrder := "DESC"
	if strings.ToLower(filter.SortOrder) == "asc" {
		sortOrder = "ASC"
	}
	orderClause := fmt.Sprintf("%s %s NULLS LAST", sortCol, sortOrder)

	// ── Fetch ────────────────────────────────────────────────────
	fetchQuery := fmt.Sprintf(`
		SELECT
			d.id                              AS device_id,
			d.device_code                     AS device_code,
			d.name                            AS device_name,
			d.serial_number                   AS serial_number,
			d.status                          AS device_status,
			COALESCE(dls.engine_running, false) AS engine_running,
			dls.speed,
			dls.coolant_temperature,
			dls.oil_pressure,
			dls.fuel_level,
			dls.batt_volt,
			dls.frequency,
			dls.total_va,
			dls.pf_avg,
			dls.telemetry_recorded_at,
			dls.last_seen_at,
			dls.last_online_at,
			dls.updated_at,
			ds.gsm_signal,
			ds.gps_connected,
			ds.server_connected,
			ds.can_connected,
			ds.rs485_connected,
			ds.sd_card_ok
		FROM devices d
		LEFT JOIN device_latest_state dls ON d.id = dls.device_id
		LEFT JOIN device_statuses ds      ON d.id = ds.device_id
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, orderClause)

	fetchArgs := append(args, filter.Limit, filter.Offset)

	var devices []model.MonitoringDevice
	if err := r.db.WithContext(ctx).Raw(fetchQuery, fetchArgs...).Scan(&devices).Error; err != nil {
		return nil, fmt.Errorf("monitoringRepository.GetMonitoringDevices fetch: %w", err)
	}

	return &MonitoringDeviceListResult{
		Devices: devices,
		Total:   total,
	}, nil
}

// GetMonitoringDeviceByID fetches the monitoring summary row for a single device
// (used internally; detail endpoint uses section-specific queries for clarity).
func (r *monitoringRepository) GetMonitoringDeviceByID(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringDevice, error) {
	query := `
		SELECT
			d.id                              AS device_id,
			d.device_code                     AS device_code,
			d.name                            AS device_name,
			d.serial_number                   AS serial_number,
			d.status                          AS device_status,
			COALESCE(dls.engine_running, false) AS engine_running,
			dls.speed,
			dls.coolant_temperature,
			dls.oil_pressure,
			dls.fuel_level,
			dls.batt_volt,
			dls.frequency,
			dls.total_va,
			dls.pf_avg,
			dls.telemetry_recorded_at,
			dls.last_seen_at,
			dls.last_online_at,
			dls.updated_at,
			ds.gsm_signal,
			ds.gps_connected,
			ds.server_connected,
			ds.can_connected,
			ds.rs485_connected,
			ds.sd_card_ok
		FROM devices d
		LEFT JOIN device_latest_state dls ON d.id = dls.device_id
		LEFT JOIN device_statuses ds      ON d.id = ds.device_id
		WHERE d.deleted_at IS NULL AND d.id = ?
	`

	var device model.MonitoringDevice
	if err := r.db.WithContext(ctx).Raw(query, deviceID).Scan(&device).Error; err != nil {
		return nil, fmt.Errorf("monitoringRepository.GetMonitoringDeviceByID: %w", err)
	}

	if device.DeviceID == "" {
		return nil, ErrMonitoringDeviceNotFound
	}

	return &device, nil
}

// GetMonitoringDeviceInfo fetches section A (devices table only).
func (r *monitoringRepository) GetMonitoringDeviceInfo(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringDeviceInfo, error) {
	query := `
		SELECT
			d.id::text         AS id,
			d.device_code      AS device_code,
			d.name             AS name,
			d.serial_number    AS serial_number,
			d.engine_id        AS engine_id,
			d.gsm_number       AS gsm_number,
			d.firmware_version AS firmware_version,
			d.status           AS status,
			d.metadata         AS metadata
		FROM devices d
		WHERE d.deleted_at IS NULL AND d.id = ?
	`

	var info struct {
		model.MonitoringDeviceInfo
		MetadataRaw []byte `gorm:"column:metadata"`
	}

	if err := r.db.WithContext(ctx).Raw(query, deviceID).Scan(&info).Error; err != nil {
		return nil, fmt.Errorf("monitoringRepository.GetMonitoringDeviceInfo: %w", err)
	}

	if info.ID == "" {
		return nil, ErrMonitoringDeviceNotFound
	}

	return &info.MonitoringDeviceInfo, nil
}

// GetMonitoringLatestState fetches section B (device_latest_state table).
func (r *monitoringRepository) GetMonitoringLatestState(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringLatestState, error) {
	query := `
		SELECT
			dls.engine_running,
			dls.speed,
			dls.coolant_temperature,
			dls.oil_pressure,
			dls.fuel_level,
			dls.batt_volt,
			dls.frequency,
			dls.total_va,
			dls.pf_avg,
			dls.telemetry_recorded_at,
			dls.last_seen_at,
			dls.last_online_at,
			dls.updated_at
		FROM device_latest_state dls
		WHERE dls.device_id = ?
	`

	var state model.MonitoringLatestState
	err := r.db.WithContext(ctx).Raw(query, deviceID).Scan(&state).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // Intentional: state may not exist yet
		}
		return nil, fmt.Errorf("monitoringRepository.GetMonitoringLatestState: %w", err)
	}

	return &state, nil
}

// GetMonitoringConnectivity fetches section C (device_statuses table).
func (r *monitoringRepository) GetMonitoringConnectivity(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringConnectivity, error) {
	query := `
		SELECT
			ds.gsm_signal,
			ds.gps_connected,
			ds.server_connected,
			ds.can_connected,
			ds.rs485_connected,
			ds.sd_card_ok,
			ds.last_seen
		FROM device_statuses ds
		WHERE ds.device_id = ?
	`

	var connectivity model.MonitoringConnectivity
	err := r.db.WithContext(ctx).Raw(query, deviceID).Scan(&connectivity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // Intentional: status may not exist yet
		}
		return nil, fmt.Errorf("monitoringRepository.GetMonitoringConnectivity: %w", err)
	}

	return &connectivity, nil
}

// ── Sentinel errors ───────────────────────────────────────────────

var ErrMonitoringDeviceNotFound = errors.New("monitoring device not found")
