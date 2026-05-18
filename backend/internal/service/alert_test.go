package service

import (
	"context"
	"errors"
	"testing"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/pubsub"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type mockAlertRepoImpl struct {
	repository.AlertRepository
	createFunc func(ctx context.Context, alert *model.Alert) error
}

func (m *mockAlertRepoImpl) Create(ctx context.Context, alert *model.Alert) error {
	return m.createFunc(ctx, alert)
}

type mockAlertDeviceRepoImpl struct {
	repository.DeviceRepository
	findByIDFunc func(ctx context.Context, id uuid.UUID) (*model.Device, error)
}

func (m *mockAlertDeviceRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	return m.findByIDFunc(ctx, id)
}

func TestAlertService_Create_Success(t *testing.T) {
	deviceID := uuid.New()
	mockAlertRepo := &mockAlertRepoImpl{
		createFunc: func(ctx context.Context, alert *model.Alert) error {
			alert.ID = uuid.New()
			return nil
		},
	}
	mockDeviceRepo := &mockAlertDeviceRepoImpl{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Device, error) {
			return &model.Device{Base: model.Base{ID: deviceID}}, nil
		},
	}

	svc := NewAlertService(mockAlertRepo, mockDeviceRepo, pubsub.NewBroker(), zap.NewNop())
	input := CreateAlertInput{
		DeviceID: deviceID,
		Type:     model.TypeEngine,
		Severity: model.SeverityCritical,
		Title:    "High Temp",
		Message:  "Critical engine temperature",
	}

	out, err := svc.Create(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if out.AlertID == uuid.Nil {
		t.Errorf("expected alert_id, got nil")
	}
}

func TestAlertService_Create_DeviceNotFound(t *testing.T) {
	mockDeviceRepo := &mockAlertDeviceRepoImpl{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Device, error) {
			return nil, repository.ErrDeviceNotFound
		},
	}

	svc := NewAlertService(nil, mockDeviceRepo, pubsub.NewBroker(), zap.NewNop())
	input := CreateAlertInput{
		DeviceID: uuid.New(),
	}

	_, err := svc.Create(context.Background(), input)

	if !errors.Is(err, ErrDeviceNotFound) {
		t.Errorf("expected ErrDeviceNotFound, got %v", err)
	}
}
