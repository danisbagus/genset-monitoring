package service

import (
	"context"
	"fmt"
	"time"

	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"go.uber.org/zap"
)

// ── DTOs ─────────────────────────────────────────────────────────

// DashboardSummaryOutput DTO struct for Dashboard Summary
type DashboardSummaryOutput struct {
	TotalDevices   int64 `json:"total_devices"`
	OnlineDevices  int64 `json:"online_devices"`
	OfflineDevices int64 `json:"offline_devices"`
	RunningEngines int64 `json:"running_engines"`
	CriticalAlerts int64 `json:"critical_alerts"`
	WarningAlerts  int64 `json:"warning_alerts"`
}

// GetDeviceStatesQuery carries validated query parameters for device states list
type GetDeviceStatesQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// GetRecentAlertsQuery carries validated query parameters for recent alerts list
type GetRecentAlertsQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// DashboardDeviceStateOutput represents a single device state in the dashboard list
type DashboardDeviceStateOutput struct {
	DeviceID           string    `json:"device_id"`
	DeviceName         string    `json:"device_name"`
	DeviceOnline       bool      `json:"device_online"`
	EngineRunning      bool      `json:"engine_running"`
	FuelLevel          float64   `json:"fuel_level"`
	CoolantTemperature float64   `json:"coolant_temperature"`
	LastSeenAt         *time.Time `json:"last_seen_at"`
}

// DashboardDeviceStatesOutput wraps the device states list with pagination metadata
type DashboardDeviceStatesOutput struct {
	Devices    []DashboardDeviceStateOutput `json:"devices"`
	Pagination PaginationMeta               `json:"pagination"`
}

// RecentAlertOutput represents a single alert in the dashboard list
type RecentAlertOutput struct {
	AlertID      string    `json:"alert_id"`
	DeviceID     string    `json:"device_id"`
	DeviceName   string    `json:"device_name"`
	Severity     string    `json:"severity"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
	Acknowledged bool      `json:"acknowledged"`
}

// RecentAlertsOutput wraps the alerts list with pagination metadata
type RecentAlertsOutput struct {
	Alerts     []RecentAlertOutput `json:"alerts"`
	Pagination PaginationMeta      `json:"pagination"`
}

// ── Interface ─────────────────────────────────────────────────────

// DashboardService define interface for dashboard service
type DashboardService interface {
	GetSummary(ctx context.Context) (*DashboardSummaryOutput, error)
	GetDeviceStates(ctx context.Context, query GetDeviceStatesQuery) (*DashboardDeviceStatesOutput, error)
	GetRecentAlerts(ctx context.Context, query GetRecentAlertsQuery) (*RecentAlertsOutput, error)
}

// ── Implementation ────────────────────────────────────────────────

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
	log           *zap.Logger
}

// NewDashboardService create new instance of DashboardService
func NewDashboardService(dashboardRepo repository.DashboardRepository, log *zap.Logger) DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
		log:           log,
	}
}

// GetSummary get dashboard summary
func (s *dashboardService) GetSummary(ctx context.Context) (*DashboardSummaryOutput, error) {
	summary, err := s.dashboardRepo.GetSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("dashboardService.GetSummary: %w", err)
	}

	return &DashboardSummaryOutput{
		TotalDevices:   summary.TotalDevices,
		OnlineDevices:  summary.OnlineDevices,
		OfflineDevices: summary.OfflineDevices,
		RunningEngines: summary.RunningEngines,
		CriticalAlerts: summary.CriticalAlerts,
		WarningAlerts:  summary.WarningAlerts,
	}, nil
}

// GetDeviceStates get dashboard device states
func (s *dashboardService) GetDeviceStates(ctx context.Context, query GetDeviceStatesQuery) (*DashboardDeviceStatesOutput, error) {
	// Default value validation
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 5
	}

	result, err := s.dashboardRepo.GetDeviceStates(ctx, repository.DeviceStateFilter{
		Page:  query.Page,
		Limit: query.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("dashboardService.GetDeviceStates: %w", err)
	}

	devices := make([]DashboardDeviceStateOutput, 0, len(result.Devices))
	for _, d := range result.Devices {
		devices = append(devices, DashboardDeviceStateOutput{
			DeviceID:           d.DeviceID,
			DeviceName:         d.DeviceName,
			DeviceOnline:       d.DeviceOnline,
			EngineRunning:      d.EngineRunning,
			FuelLevel:          d.FuelLevel,
			CoolantTemperature: d.CoolantTemperature,
			LastSeenAt:         d.LastSeenAt,
		})
	}

	return &DashboardDeviceStatesOutput{
		Devices: devices,
		Pagination: PaginationMeta{
			Page:  query.Page,
			Limit: query.Limit,
			Total: result.Total,
		},
	}, nil
}

// GetRecentAlerts get dashboard recent alerts
func (s *dashboardService) GetRecentAlerts(ctx context.Context, query GetRecentAlertsQuery) (*RecentAlertsOutput, error) {
	// Default value validation
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 5
	}

	result, err := s.dashboardRepo.GetRecentAlerts(ctx, repository.RecentAlertFilter{
		Page:  query.Page,
		Limit: query.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("dashboardService.GetRecentAlerts: %w", err)
	}

	alerts := make([]RecentAlertOutput, 0, len(result.Alerts))
	for _, a := range result.Alerts {
		alerts = append(alerts, RecentAlertOutput{
			AlertID:      a.AlertID,
			DeviceID:     a.DeviceID,
			DeviceName:   a.DeviceName,
			Severity:     a.Severity,
			Message:      a.Message,
			CreatedAt:    a.CreatedAt,
			Acknowledged: a.Acknowledged,
		})
	}

	return &RecentAlertsOutput{
		Alerts: alerts,
		Pagination: PaginationMeta{
			Page:  query.Page,
			Limit: query.Limit,
			Total: result.Total,
		},
	}, nil
}
