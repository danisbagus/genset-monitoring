package service

import (
	"context"
	"fmt"

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

// ── Interface ─────────────────────────────────────────────────────

// DashboardService define interface for dashboard service
type DashboardService interface {
	GetSummary(ctx context.Context) (*DashboardSummaryOutput, error)
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
