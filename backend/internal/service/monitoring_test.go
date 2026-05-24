package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
)

// ── Mock implementations ──────────────────────────────────────────

type mockMonitoringRepo struct {
	repository.MonitoringRepository

	getMonitoringDevicesFunc      func(ctx context.Context, filter repository.MonitoringDeviceFilter) (*repository.MonitoringDeviceListResult, error)
	getMonitoringDeviceByIDFunc   func(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringDevice, error)
	getMonitoringDeviceInfoFunc   func(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringDeviceInfo, error)
	getMonitoringLatestStateFunc  func(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringLatestState, error)
	getMonitoringConnectivityFunc func(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringConnectivity, error)
}

func (m *mockMonitoringRepo) GetMonitoringDevices(ctx context.Context, filter repository.MonitoringDeviceFilter) (*repository.MonitoringDeviceListResult, error) {
	return m.getMonitoringDevicesFunc(ctx, filter)
}

func (m *mockMonitoringRepo) GetMonitoringDeviceByID(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringDevice, error) {
	return m.getMonitoringDeviceByIDFunc(ctx, deviceID)
}

func (m *mockMonitoringRepo) GetMonitoringDeviceInfo(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringDeviceInfo, error) {
	return m.getMonitoringDeviceInfoFunc(ctx, deviceID)
}

func (m *mockMonitoringRepo) GetMonitoringLatestState(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringLatestState, error) {
	return m.getMonitoringLatestStateFunc(ctx, deviceID)
}

func (m *mockMonitoringRepo) GetMonitoringConnectivity(ctx context.Context, deviceID uuid.UUID) (*model.MonitoringConnectivity, error) {
	return m.getMonitoringConnectivityFunc(ctx, deviceID)
}

type mockTelemetryRepoForMonitoring struct {
	repository.TelemetryRepository

	getLatestEngineFunc      func(ctx context.Context, deviceID uuid.UUID) (*model.EngineTelemetry, error)
	getLatestElectricalFunc  func(ctx context.Context, deviceID uuid.UUID) (*model.ElectricalTelemetry, error)
}

func (m *mockTelemetryRepoForMonitoring) GetLatestEngine(ctx context.Context, deviceID uuid.UUID) (*model.EngineTelemetry, error) {
	return m.getLatestEngineFunc(ctx, deviceID)
}

func (m *mockTelemetryRepoForMonitoring) GetLatestElectrical(ctx context.Context, deviceID uuid.UUID) (*model.ElectricalTelemetry, error) {
	return m.getLatestElectricalFunc(ctx, deviceID)
}

// ── Tests: ListDevices ────────────────────────────────────────────

func TestMonitoringService_ListDevices_Success(t *testing.T) {
	now := time.Now()

	mockRepo := &mockMonitoringRepo{
		getMonitoringDevicesFunc: func(ctx context.Context, filter repository.MonitoringDeviceFilter) (*repository.MonitoringDeviceListResult, error) {
			return &repository.MonitoringDeviceListResult{
				Devices: []model.MonitoringDevice{
					{
						DeviceID:      "device-1",
						DeviceCode:    "GS-001",
						DeviceName:    "Genset Alpha",
						SerialNumber:  "SN-001",
						EngineRunning: true,
						DeviceStatus:  "active",
						LastSeenAt:    &now,
					},
				},
				Total: 1,
			}, nil
		},
	}

	svc := NewMonitoringService(mockRepo, &mockTelemetryRepoForMonitoring{}, zap.NewNop())
	out, err := svc.ListDevices(context.Background(), MonitoringDeviceListQuery{
		Limit:  20,
		Offset: 0,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(out.Devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(out.Devices))
	}
	if out.Devices[0].DeviceCode != "GS-001" {
		t.Errorf("expected device_code GS-001, got %s", out.Devices[0].DeviceCode)
	}
	if out.Devices[0].EngineRunning != true {
		t.Errorf("expected engine_running true, got %v", out.Devices[0].EngineRunning)
	}
	if out.Pagination.Total != 1 {
		t.Errorf("expected total 1, got %d", out.Pagination.Total)
	}
}

func TestMonitoringService_ListDevices_DefaultLimitClamped(t *testing.T) {
	var capturedFilter repository.MonitoringDeviceFilter

	mockRepo := &mockMonitoringRepo{
		getMonitoringDevicesFunc: func(ctx context.Context, filter repository.MonitoringDeviceFilter) (*repository.MonitoringDeviceListResult, error) {
			capturedFilter = filter
			return &repository.MonitoringDeviceListResult{Devices: []model.MonitoringDevice{}, Total: 0}, nil
		},
	}

	svc := NewMonitoringService(mockRepo, &mockTelemetryRepoForMonitoring{}, zap.NewNop())

	// Limit 0 should default to 20
	_, _ = svc.ListDevices(context.Background(), MonitoringDeviceListQuery{Limit: 0})
	if capturedFilter.Limit != 20 {
		t.Errorf("expected default limit 20, got %d", capturedFilter.Limit)
	}

	// Limit exceeding 100 should clamp to 20
	_, _ = svc.ListDevices(context.Background(), MonitoringDeviceListQuery{Limit: 999})
	if capturedFilter.Limit != 20 {
		t.Errorf("expected clamped limit 20, got %d", capturedFilter.Limit)
	}
}

func TestMonitoringService_ListDevices_BooleanFilters(t *testing.T) {
	var capturedFilter repository.MonitoringDeviceFilter

	mockRepo := &mockMonitoringRepo{
		getMonitoringDevicesFunc: func(ctx context.Context, filter repository.MonitoringDeviceFilter) (*repository.MonitoringDeviceListResult, error) {
			capturedFilter = filter
			return &repository.MonitoringDeviceListResult{Devices: []model.MonitoringDevice{}, Total: 0}, nil
		},
	}

	svc := NewMonitoringService(mockRepo, &mockTelemetryRepoForMonitoring{}, zap.NewNop())

	// engine_running=true & online=false
	_, _ = svc.ListDevices(context.Background(), MonitoringDeviceListQuery{
		Limit:         10,
		EngineRunning: "true",
		Online:        "false",
	})

	if capturedFilter.EngineRunning == nil || *capturedFilter.EngineRunning != true {
		t.Errorf("expected engine_running filter true, got %v", capturedFilter.EngineRunning)
	}
	if capturedFilter.Online == nil || *capturedFilter.Online != false {
		t.Errorf("expected online filter false, got %v", capturedFilter.Online)
	}
}

func TestMonitoringService_ListDevices_SortByAllowList(t *testing.T) {
	var capturedFilter repository.MonitoringDeviceFilter

	mockRepo := &mockMonitoringRepo{
		getMonitoringDevicesFunc: func(ctx context.Context, filter repository.MonitoringDeviceFilter) (*repository.MonitoringDeviceListResult, error) {
			capturedFilter = filter
			return &repository.MonitoringDeviceListResult{Devices: []model.MonitoringDevice{}, Total: 0}, nil
		},
	}

	svc := NewMonitoringService(mockRepo, &mockTelemetryRepoForMonitoring{}, zap.NewNop())

	// Disallowed sort_by should be rejected (empty passed to repo)
	_, _ = svc.ListDevices(context.Background(), MonitoringDeviceListQuery{
		Limit:  10,
		SortBy: "hacked_column; DROP TABLE devices;--",
	})
	if capturedFilter.SortBy != "" {
		t.Errorf("expected empty sortBy for disallowed column, got %q", capturedFilter.SortBy)
	}

	// Allowed sort_by should pass through
	_, _ = svc.ListDevices(context.Background(), MonitoringDeviceListQuery{
		Limit:  10,
		SortBy: "name",
	})
	if capturedFilter.SortBy != "name" {
		t.Errorf("expected sortBy 'name', got %q", capturedFilter.SortBy)
	}
}

// ── Tests: GetDeviceDetail ────────────────────────────────────────

func TestMonitoringService_GetDeviceDetail_Success(t *testing.T) {
	deviceID := uuid.New()
	now := time.Now()

	speed := int32(1500)
	engMock := &mockTelemetryRepoForMonitoring{
		getLatestEngineFunc: func(ctx context.Context, id uuid.UUID) (*model.EngineTelemetry, error) {
			return &model.EngineTelemetry{
				DeviceID:  id,
				Speed:     &speed,
				CreatedAt: now,
			}, nil
		},
		getLatestElectricalFunc: func(ctx context.Context, id uuid.UUID) (*model.ElectricalTelemetry, error) {
			return nil, repository.ErrTelemetryNotFound
		},
	}

	mockRepo := &mockMonitoringRepo{
		getMonitoringDeviceInfoFunc: func(ctx context.Context, id uuid.UUID) (*model.MonitoringDeviceInfo, error) {
			return &model.MonitoringDeviceInfo{
				ID:         id.String(),
				DeviceCode: "GS-001",
				Name:       "Genset Alpha",
				Status:     "active",
			}, nil
		},
		getMonitoringLatestStateFunc: func(ctx context.Context, id uuid.UUID) (*model.MonitoringLatestState, error) {
			return &model.MonitoringLatestState{
				EngineRunning: true,
				Speed:         &speed,
				LastSeenAt:    &now,
			}, nil
		},
		getMonitoringConnectivityFunc: func(ctx context.Context, id uuid.UUID) (*model.MonitoringConnectivity, error) {
			return &model.MonitoringConnectivity{
				GSMSignal:       3,
				ServerConnected: true,
				LastSeen:        now,
			}, nil
		},
	}

	svc := NewMonitoringService(mockRepo, engMock, zap.NewNop())
	out, err := svc.GetDeviceDetail(context.Background(), deviceID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.DeviceInfo.DeviceCode != "GS-001" {
		t.Errorf("expected device_code GS-001, got %s", out.DeviceInfo.DeviceCode)
	}
	if out.LatestState == nil {
		t.Fatal("expected latest_state to be non-nil")
	}
	if !out.LatestState.EngineRunning {
		t.Errorf("expected engine_running true")
	}
	if out.Connectivity == nil {
		t.Fatal("expected connectivity to be non-nil")
	}
	if out.Connectivity.GSMSignal != 3 {
		t.Errorf("expected gsm_signal 3, got %d", out.Connectivity.GSMSignal)
	}
	if out.EngineTelemetry == nil {
		t.Fatal("expected engine telemetry to be non-nil")
	}
	// Electrical telemetry should be nil (not found)
	if out.ElectricalTelemetry != nil {
		t.Errorf("expected electrical telemetry nil, got non-nil")
	}
}

func TestMonitoringService_GetDeviceDetail_DeviceNotFound(t *testing.T) {
	mockRepo := &mockMonitoringRepo{
		getMonitoringDeviceInfoFunc: func(ctx context.Context, id uuid.UUID) (*model.MonitoringDeviceInfo, error) {
			return nil, repository.ErrMonitoringDeviceNotFound
		},
	}

	svc := NewMonitoringService(mockRepo, &mockTelemetryRepoForMonitoring{}, zap.NewNop())
	_, err := svc.GetDeviceDetail(context.Background(), uuid.New())

	if !errors.Is(err, ErrMonitoringDeviceNotFound) {
		t.Errorf("expected ErrMonitoringDeviceNotFound, got %v", err)
	}
}

func TestMonitoringService_GetDeviceDetail_PartialData_NilState(t *testing.T) {
	deviceID := uuid.New()

	mockRepo := &mockMonitoringRepo{
		getMonitoringDeviceInfoFunc: func(ctx context.Context, id uuid.UUID) (*model.MonitoringDeviceInfo, error) {
			return &model.MonitoringDeviceInfo{
				ID:     id.String(),
				Status: "active",
			}, nil
		},
		getMonitoringLatestStateFunc: func(ctx context.Context, id uuid.UUID) (*model.MonitoringLatestState, error) {
			return nil, nil // device has no state yet
		},
		getMonitoringConnectivityFunc: func(ctx context.Context, id uuid.UUID) (*model.MonitoringConnectivity, error) {
			return nil, nil // device has no connectivity record yet
		},
	}

	engMock := &mockTelemetryRepoForMonitoring{
		getLatestEngineFunc: func(ctx context.Context, id uuid.UUID) (*model.EngineTelemetry, error) {
			return nil, repository.ErrTelemetryNotFound
		},
		getLatestElectricalFunc: func(ctx context.Context, id uuid.UUID) (*model.ElectricalTelemetry, error) {
			return nil, repository.ErrTelemetryNotFound
		},
	}

	svc := NewMonitoringService(mockRepo, engMock, zap.NewNop())
	out, err := svc.GetDeviceDetail(context.Background(), deviceID)

	if err != nil {
		t.Fatalf("expected no error for partial data, got %v", err)
	}
	if out.LatestState != nil {
		t.Errorf("expected nil LatestState for new device, got non-nil")
	}
	if out.Connectivity != nil {
		t.Errorf("expected nil Connectivity for new device, got non-nil")
	}
	if out.EngineTelemetry != nil {
		t.Errorf("expected nil EngineTelemetry, got non-nil")
	}
}
