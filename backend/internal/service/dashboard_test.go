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
	getSummaryFunc      func(ctx context.Context) (*model.DashboardSummary, error)
	getDeviceStatesFunc func(ctx context.Context, filter repository.DeviceStateFilter) (*repository.DeviceStateResult, error)
	getRecentAlertsFunc func(ctx context.Context, filter repository.RecentAlertFilter) (*repository.RecentAlertResult, error)
}

func (m *mockDashboardRepo) GetSummary(ctx context.Context) (*model.DashboardSummary, error) {
	return m.getSummaryFunc(ctx)
}

func (m *mockDashboardRepo) GetDeviceStates(ctx context.Context, filter repository.DeviceStateFilter) (*repository.DeviceStateResult, error) {
	return m.getDeviceStatesFunc(ctx, filter)
}

func (m *mockDashboardRepo) GetRecentAlerts(ctx context.Context, filter repository.RecentAlertFilter) (*repository.RecentAlertResult, error) {
	return m.getRecentAlertsFunc(ctx, filter)
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

func TestDashboardService_GetDeviceStates_Success(t *testing.T) {
	mockRepo := &mockDashboardRepo{
		getDeviceStatesFunc: func(ctx context.Context, filter repository.DeviceStateFilter) (*repository.DeviceStateResult, error) {
			return &repository.DeviceStateResult{
				Devices: []model.DeviceState{
					{
						DeviceID:     "device-1",
						DeviceName:   "Genset 1",
						DeviceOnline: true,
					},
				},
				Total: 1,
			}, nil
		},
	}

	svc := NewDashboardService(mockRepo, zap.NewNop())
	output, err := svc.GetDeviceStates(context.Background(), GetDeviceStatesQuery{Page: 1, Limit: 5})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.Devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(output.Devices))
	}

	if output.Devices[0].DeviceName != "Genset 1" {
		t.Errorf("expected device name Genset 1, got %s", output.Devices[0].DeviceName)
	}

	if output.Pagination.Total != 1 {
		t.Errorf("expected total pagination 1, got %d", output.Pagination.Total)
	}
}

func TestDashboardService_GetRecentAlerts_Success(t *testing.T) {
	mockRepo := &mockDashboardRepo{
		getRecentAlertsFunc: func(ctx context.Context, filter repository.RecentAlertFilter) (*repository.RecentAlertResult, error) {
			return &repository.RecentAlertResult{
				Alerts: []model.RecentAlert{
					{
						AlertID:      "alert-1",
						DeviceID:     "device-1",
						DeviceName:   "Genset 1",
						Severity:     "critical",
						Message:      "High Temp",
						Acknowledged: false,
					},
				},
				Total: 1,
			}, nil
		},
	}

	svc := NewDashboardService(mockRepo, zap.NewNop())
	output, err := svc.GetRecentAlerts(context.Background(), GetRecentAlertsQuery{Page: 1, Limit: 5})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.Alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(output.Alerts))
	}

	if output.Alerts[0].DeviceName != "Genset 1" {
		t.Errorf("expected device name Genset 1, got %s", output.Alerts[0].DeviceName)
	}

	if output.Pagination.Total != 1 {
		t.Errorf("expected total pagination 1, got %d", output.Pagination.Total)
	}
}
