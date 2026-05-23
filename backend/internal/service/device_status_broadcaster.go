package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/pubsub"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"github.com/danisbagus/genset-monitoring/backend/internal/websocket"
)

// DeviceStatusBroadcastPayload is the public representation of a device's realtime dashboard status.
type DeviceStatusBroadcastPayload struct {
	DeviceID           string     `json:"device_id"`
	DeviceName         string     `json:"device_name"`
	DeviceOnline       bool       `json:"device_online"`
	EngineRunning      bool       `json:"engine_running"`
	FuelLevel          float64    `json:"fuel_level"`
	CoolantTemperature float64    `json:"coolant_temperature"`
	LastSeenAt         *time.Time `json:"last_seen_at"`
}

// DeviceStatusBroadcaster acts as a bridge between the Pub/Sub broker,
// the Dashboard repository, and the WebSocket Hub to update clients on device status changes in real-time.
type DeviceStatusBroadcaster struct {
	broker        pubsub.Broker
	dashboardRepo repository.DashboardRepository
	hub           *websocket.Hub
	log           *zap.Logger
}

// NewDeviceStatusBroadcaster creates a new instance of DeviceStatusBroadcaster.
func NewDeviceStatusBroadcaster(
	broker pubsub.Broker,
	dashboardRepo repository.DashboardRepository,
	hub *websocket.Hub,
	log *zap.Logger,
) *DeviceStatusBroadcaster {
	return &DeviceStatusBroadcaster{
		broker:        broker,
		dashboardRepo: dashboardRepo,
		hub:           hub,
		log:           log,
	}
}

// Start launches the background listener in a context-aware goroutine.
func (dsb *DeviceStatusBroadcaster) Start(ctx context.Context) {
	deviceCh := dsb.broker.Subscribe(pubsub.TopicDeviceCRUD)
	heartbeatCh := dsb.broker.Subscribe(pubsub.TopicHeartbeatReceived)
	telemetryCh := dsb.broker.Subscribe(pubsub.TopicTelemetryEngineCreated)

	dsb.log.Info("Device status broadcaster initialized and listening to topics")

	go func() {
		for {
			select {
			case <-ctx.Done():
				dsb.log.Info("Device status broadcaster stopping due to context cancellation")
				return

			case val, ok := <-deviceCh:
				if !ok {
					return
				}
				dsb.handleEvent(ctx, val, "device_crud")

			case val, ok := <-heartbeatCh:
				if !ok {
					return
				}
				dsb.handleEvent(ctx, val, "heartbeat_received")

			case val, ok := <-telemetryCh:
				if !ok {
					return
				}
				dsb.handleEvent(ctx, val, "telemetry_created")
			}
		}
	}()
}

// handleEvent extracts the device ID from the event and triggers the broadcast.
func (dsb *DeviceStatusBroadcaster) handleEvent(ctx context.Context, val interface{}, triggerSource string) {
	if val == nil {
		return
	}

	var deviceID uuid.UUID

	switch payload := val.(type) {
	case *model.Device:
		if payload != nil {
			deviceID = payload.ID
		}
	case uuid.UUID:
		deviceID = payload
	case HeartbeatInput:
		deviceID = payload.DeviceID
	case *EngineTelemetryOutput:
		if payload != nil {
			deviceID = payload.DeviceID
		}
	default:
		dsb.log.Warn("Device status broadcaster: unknown payload type received",
			zap.String("source", triggerSource),
			zap.String("type", fmt.Sprintf("%T", val)))
		return
	}

	if deviceID == uuid.Nil {
		return
	}

	// Trigger broadcast
	go dsb.broadcastDeviceStatus(ctx, deviceID, triggerSource)
}

// broadcastDeviceStatus fetches the latest device state and broadcasts it
// as a 'dashboard.devices.updated' websocket event.
func (dsb *DeviceStatusBroadcaster) broadcastDeviceStatus(ctx context.Context, deviceID uuid.UUID, triggerSource string) {
	dsb.log.Debug("Device status broadcast triggered",
		zap.String("device_id", deviceID.String()),
		zap.String("source", triggerSource))

	// Use a 5-second timeout for the DB fetch to keep operations snappy
	fetchCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	state, err := dsb.dashboardRepo.GetDeviceStateByID(fetchCtx, deviceID)
	if err != nil {
		dsb.log.Error("Failed to fetch device state for websocket broadcast",
			zap.Error(err),
			zap.String("device_id", deviceID.String()),
			zap.String("source", triggerSource))
		return
	}

	// If device was soft deleted, or does not exist, state might be empty
	if state.DeviceID == "" {
		dsb.log.Debug("Device state not found, likely deleted",
			zap.String("device_id", deviceID.String()),
			zap.String("source", triggerSource))
		return
	}

	// Map to clean payload DTO
	payload := &DeviceStatusBroadcastPayload{
		DeviceID:           state.DeviceID,
		DeviceName:         state.DeviceName,
		DeviceOnline:       state.DeviceOnline,
		EngineRunning:      state.EngineRunning,
		FuelLevel:          state.FuelLevel,
		CoolantTemperature: state.CoolantTemperature,
		LastSeenAt:         state.LastSeenAt,
	}

	msg, err := websocket.NewMessage("dashboard.devices.updated", payload)
	if err != nil {
		dsb.log.Error("Failed to marshal device status for websocket broadcast",
			zap.Error(err),
			zap.String("device_id", deviceID.String()),
			zap.String("source", triggerSource))
		return
	}

	// Non-blocking broadcast
	dsb.hub.Broadcast(msg)
	dsb.log.Info("Device status broadcasted successfully via websocket",
		zap.String("device_id", deviceID.String()),
		zap.String("source", triggerSource),
		zap.String("device_name", payload.DeviceName),
		zap.Bool("device_online", payload.DeviceOnline),
		zap.Bool("engine_running", payload.EngineRunning))
}
