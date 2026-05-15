package service

import (
	"context"
	"testing"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"go.uber.org/zap"
)

type mockDashboardRepo struct {
	repository.DashboardRepository
	getSummaryFunc func(ctx context.Context) (*model.DashboardSummary, error)
}

func (m *mockDashboardRepo) GetSummary(ctx context.Context) (*model.DashboardSummary, error) {
	return m.getSummaryFunc(ctx)
}

func TestDashboardService_GetSummary_Success(t *testing.T) {
	mockRepo := &mockDashboardRepo{
		getSummaryFunc: func(ctx context.Context) (*model.DashboardSummary, error) {
			return &model.DashboardSummary{
				TotalDevices:   20,
				OnlineDevices:  18,
				OfflineDevices: 2,
				RunningEngines: 12,
				CriticalAlerts: 1,
				WarningAlerts:  4,
			}, nil
		},
	}

	svc := NewDashboardService(mockRepo, zap.NewNop())
	output, err := svc.GetSummary(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.TotalDevices != 20 {
		t.Errorf("expected total_devices 20, got %d", output.TotalDevices)
	}
	if output.OnlineDevices != 18 {
		t.Errorf("expected online_devices 18, got %d", output.OnlineDevices)
	}
	if output.OfflineDevices != 2 {
		t.Errorf("expected offline_devices 2, got %d", output.OfflineDevices)
	}
	if output.RunningEngines != 12 {
		t.Errorf("expected running_engines 12, got %d", output.RunningEngines)
	}
	if output.CriticalAlerts != 1 {
		t.Errorf("expected critical_alerts 1, got %d", output.CriticalAlerts)
	}
	if output.WarningAlerts != 4 {
		t.Errorf("expected warning_alerts 4, got %d", output.WarningAlerts)
	}
}
