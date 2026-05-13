package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
)

type mockTelemetryRepo struct {
	repository.TelemetryRepository
	createEngineFunc      func(ctx context.Context, telemetry *model.EngineTelemetry) error
	getLatestEngineFunc   func(ctx context.Context, deviceID uuid.UUID) (*model.EngineTelemetry, error)
	createElectricalFunc  func(ctx context.Context, telemetry *model.ElectricalTelemetry) error
	getLatestElectricalFunc func(ctx context.Context, deviceID uuid.UUID) (*model.ElectricalTelemetry, error)
}

func (m *mockTelemetryRepo) CreateEngine(ctx context.Context, telemetry *model.EngineTelemetry) error {
	return m.createEngineFunc(ctx, telemetry)
}

func (m *mockTelemetryRepo) GetLatestEngine(ctx context.Context, deviceID uuid.UUID) (*model.EngineTelemetry, error) {
	return m.getLatestEngineFunc(ctx, deviceID)
}

func (m *mockTelemetryRepo) CreateElectrical(ctx context.Context, telemetry *model.ElectricalTelemetry) error {
	return m.createElectricalFunc(ctx, telemetry)
}

func (m *mockTelemetryRepo) GetLatestElectrical(ctx context.Context, deviceID uuid.UUID) (*model.ElectricalTelemetry, error) {
	return m.getLatestElectricalFunc(ctx, deviceID)
}

type mockDeviceRepo struct {
	repository.DeviceRepository
	findByIDFunc func(ctx context.Context, id uuid.UUID) (*model.Device, error)
}

func (m *mockDeviceRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	return m.findByIDFunc(ctx, id)
}

func TestTelemetryService_CreateEngine_Success(t *testing.T) {
	deviceID := uuid.New()
	mockT := &mockTelemetryRepo{
		createEngineFunc: func(ctx context.Context, telemetry *model.EngineTelemetry) error {
			return nil
		},
	}
	mockD := &mockDeviceRepo{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Device, error) {
			return &model.Device{Base: model.Base{ID: deviceID}}, nil
		},
	}
	svc := NewTelemetryService(mockT, mockD, zap.NewNop())

	input := CreateEngineTelemetryInput{
		DeviceID: deviceID,
	}

	output, err := svc.CreateEngine(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.DeviceID != deviceID {
		t.Errorf("expected device ID %v, got %v", deviceID, output.DeviceID)
	}
}

func TestTelemetryService_GetLatestEngine_Success(t *testing.T) {
	deviceID := uuid.New()
	mockT := &mockTelemetryRepo{
		getLatestEngineFunc: func(ctx context.Context, id uuid.UUID) (*model.EngineTelemetry, error) {
			return &model.EngineTelemetry{DeviceID: id}, nil
		},
	}
	svc := NewTelemetryService(mockT, nil, zap.NewNop())

	output, err := svc.GetLatestEngine(context.Background(), deviceID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.DeviceID != deviceID {
		t.Errorf("expected device ID %v, got %v", deviceID, output.DeviceID)
	}
}

func TestTelemetryService_GetLatestEngine_NotFound(t *testing.T) {
	deviceID := uuid.New()
	mockT := &mockTelemetryRepo{
		getLatestEngineFunc: func(ctx context.Context, id uuid.UUID) (*model.EngineTelemetry, error) {
			return nil, repository.ErrTelemetryNotFound
		},
	}
	svc := NewTelemetryService(mockT, nil, zap.NewNop())

	_, err := svc.GetLatestEngine(context.Background(), deviceID)
	if err == nil || err != ErrTelemetryNotFound {
		t.Fatalf("expected ErrTelemetryNotFound, got %v", err)
	}
}
