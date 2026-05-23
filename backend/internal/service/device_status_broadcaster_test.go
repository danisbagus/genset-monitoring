package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/pubsub"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"github.com/danisbagus/genset-monitoring/backend/internal/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type mockDashboardRepoForBroadcaster struct {
	repository.DashboardRepository
	getDeviceStateByIDFunc func(ctx context.Context, deviceID uuid.UUID) (*model.DeviceState, error)
}

func (m *mockDashboardRepoForBroadcaster) GetDeviceStateByID(ctx context.Context, deviceID uuid.UUID) (*model.DeviceState, error) {
	return m.getDeviceStateByIDFunc(ctx, deviceID)
}

func TestDeviceStatusBroadcaster_StartAndBroadcast(t *testing.T) {
	broker := pubsub.NewBroker()
	defer broker.Close()

	wsHub := websocket.NewHub(zap.NewNop())
	log := zap.NewNop()

	deviceID := uuid.New()
	lastSeen := time.Now().UTC()

	var wg sync.WaitGroup
	wg.Add(4) // 4 events to test: heartbeat, device CRUD (device object), device CRUD (UUID), telemetry created

	mockRepo := &mockDashboardRepoForBroadcaster{
		getDeviceStateByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.DeviceState, error) {
			if id != deviceID {
				t.Errorf("expected device ID %v, got %v", deviceID, id)
			}
			wg.Done()
			return &model.DeviceState{
				DeviceID:           id.String(),
				DeviceName:         "Genset Tower B",
				DeviceOnline:       true,
				EngineRunning:      true,
				FuelLevel:          85.5,
				CoolantTemperature: 90.0,
				LastSeenAt:         &lastSeen,
			}, nil
		},
	}

	broadcaster := NewDeviceStatusBroadcaster(broker, mockRepo, wsHub, log)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	broadcaster.Start(ctx)

	// 1. Publish heartbeat event
	broker.Publish(pubsub.TopicHeartbeatReceived, HeartbeatInput{
		DeviceID: deviceID,
	})

	// 2. Publish device CRUD event with device object
	broker.Publish(pubsub.TopicDeviceCRUD, &model.Device{
		Base: model.Base{ID: deviceID},
	})

	// 3. Publish device CRUD event with device ID
	broker.Publish(pubsub.TopicDeviceCRUD, deviceID)

	// 4. Publish telemetry engine created event
	broker.Publish(pubsub.TopicTelemetryEngineCreated, &EngineTelemetryOutput{
		DeviceID: deviceID,
	})

	// Wait with a timeout to verify GetDeviceStateByID is called for all events
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for broadcaster to process all events")
	}
}
