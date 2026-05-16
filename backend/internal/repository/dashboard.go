package repository

import (
	"context"
	"fmt"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"gorm.io/gorm"
)

// DeviceStateFilter carries parameters for device states queries.
type DeviceStateFilter struct {
	Page  int
	Limit int
}

// DeviceStateResult is the raw result returned from GetDeviceStates.
type DeviceStateResult struct {
	Devices []model.DeviceState
	Total   int64
}

// RecentAlertFilter carries parameters for recent alerts queries.
type RecentAlertFilter struct {
	Page  int
	Limit int
}

// RecentAlertResult is the raw result returned from GetRecentAlerts.
type RecentAlertResult struct {
	Alerts []model.RecentAlert
	Total  int64
}

// DashboardRepository mendefinisikan kontrak untuk operasi data dashboard.
type DashboardRepository interface {
	GetSummary(ctx context.Context) (*model.DashboardSummary, error)
	GetDeviceStates(ctx context.Context, filter DeviceStateFilter) (*DeviceStateResult, error)
	GetRecentAlerts(ctx context.Context, filter RecentAlertFilter) (*RecentAlertResult, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository create new instance
func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// GetSummary get dashboard summary
func (r *dashboardRepository) GetSummary(ctx context.Context) (*model.DashboardSummary, error) {
	var summary model.DashboardSummary

	// Single query using subquery for efficiency
	query := `
		SELECT
			(SELECT COUNT(*) FROM devices WHERE deleted_at IS NULL) as total_devices,
			COALESCE((SELECT COUNT(*) FROM device_latest_state dls JOIN devices d ON dls.device_id = d.id WHERE d.deleted_at IS NULL AND dls.last_online_at >= (NOW() - INTERVAL '5 minutes')), 0) as online_devices,
			COALESCE((SELECT COUNT(*) FROM device_latest_state dls JOIN devices d ON dls.device_id = d.id WHERE d.deleted_at IS NULL AND dls.last_online_at < (NOW() - INTERVAL '5 minutes')), 0) as offline_devices,
			COALESCE((SELECT COUNT(*) FROM device_latest_state dls JOIN devices d ON dls.device_id = d.id WHERE d.deleted_at IS NULL AND dls.engine_running = true AND dls.last_online_at >= (NOW() - INTERVAL '5 minutes')), 0) as running_engines,
			COALESCE((SELECT COUNT(*) FROM alerts a JOIN devices d ON a.device_id = d.id WHERE d.deleted_at IS NULL AND a.status = 'active' AND a.severity = 'critical'), 0) as critical_alerts,
			COALESCE((SELECT COUNT(*) FROM alerts a JOIN devices d ON a.device_id = d.id WHERE d.deleted_at IS NULL AND a.status = 'active' AND a.severity = 'warning'), 0) as warning_alerts
	`

	if err := r.db.WithContext(ctx).Raw(query).Scan(&summary).Error; err != nil {
		return nil, fmt.Errorf("dashboardRepository.GetSummary: %w", err)
	}

	return &summary, nil
}

// GetDeviceStates returns a paginated list of device states for the dashboard.
func (r *dashboardRepository) GetDeviceStates(ctx context.Context, filter DeviceStateFilter) (*DeviceStateResult, error) {
	var total int64
	if err := r.db.WithContext(ctx).Table("devices").Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, fmt.Errorf("dashboardRepository.GetDeviceStates count: %w", err)
	}

	offset := (filter.Page - 1) * filter.Limit
	var devices []model.DeviceState

	query := `
		SELECT 
			d.id AS device_id,
			d.name AS device_name,
			COALESCE(dls.last_online_at >= (NOW() - INTERVAL '5 minutes'), false) AS device_online,
			COALESCE(dls.engine_running = true AND dls.last_online_at >= (NOW() - INTERVAL '5 minutes'), false) AS engine_running,
			COALESCE(dls.fuel_level, 0) AS fuel_level,
			COALESCE(dls.coolant_temperature, 0) AS coolant_temperature,
			dls.last_seen_at
		FROM devices d
		LEFT JOIN device_latest_state dls ON d.id = dls.device_id
		WHERE d.deleted_at IS NULL
		ORDER BY dls.last_seen_at DESC NULLS LAST
		LIMIT ? OFFSET ?
	`

	if err := r.db.WithContext(ctx).Raw(query, filter.Limit, offset).Scan(&devices).Error; err != nil {
		return nil, fmt.Errorf("dashboardRepository.GetDeviceStates fetch: %w", err)
	}

	return &DeviceStateResult{
		Devices: devices,
		Total:   total,
	}, nil
}

// GetRecentAlerts returns a paginated list of recent alerts for the dashboard.
func (r *dashboardRepository) GetRecentAlerts(ctx context.Context, filter RecentAlertFilter) (*RecentAlertResult, error) {
	var total int64
	// Only count alerts for devices that are not deleted
	countQuery := `
		SELECT COUNT(*) 
		FROM alerts a
		JOIN devices d ON a.device_id = d.id
		WHERE d.deleted_at IS NULL
	`
	if err := r.db.WithContext(ctx).Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("dashboardRepository.GetRecentAlerts count: %w", err)
	}

	offset := (filter.Page - 1) * filter.Limit
	var alerts []model.RecentAlert

	query := `
		SELECT 
			a.id AS alert_id,
			a.device_id,
			d.name AS device_name,
			a.severity,
			a.message,
			a.created_at,
			(a.status != 'active') AS acknowledged
		FROM alerts a
		JOIN devices d ON a.device_id = d.id
		WHERE d.deleted_at IS NULL
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?
	`

	if err := r.db.WithContext(ctx).Raw(query, filter.Limit, offset).Scan(&alerts).Error; err != nil {
		return nil, fmt.Errorf("dashboardRepository.GetRecentAlerts fetch: %w", err)
	}

	return &RecentAlertResult{
		Alerts: alerts,
		Total:  total,
	}, nil
}
