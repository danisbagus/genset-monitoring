package repository

import (
	"context"
	"fmt"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"gorm.io/gorm"
)

// DashboardRepository mendefinisikan kontrak untuk operasi data dashboard.
type DashboardRepository interface {
	GetSummary(ctx context.Context) (*model.DashboardSummary, error)
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
