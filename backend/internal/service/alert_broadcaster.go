package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/pubsub"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"github.com/danisbagus/genset-monitoring/backend/internal/websocket"
)

// AlertBroadcastPayload is the WebSocket event payload for dashboard.alert.received.
type AlertBroadcastPayload struct {
	AlertID      string    `json:"alert_id"`
	DeviceID     string    `json:"device_id"`
	DeviceName   string    `json:"device_name"`
	Severity     string    `json:"severity"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
	Acknowledged bool      `json:"acknowledged"`
}

// AlertBroadcaster acts as a bridge between the Pub/Sub broker,
// the Device repository, and the WebSocket Hub to push real-time
// alert notifications to all connected dashboard clients.
type AlertBroadcaster struct {
	broker     pubsub.Broker
	deviceRepo repository.DeviceRepository
	hub        *websocket.Hub
	log        *zap.Logger
}

// NewAlertBroadcaster creates a new instance of AlertBroadcaster.
func NewAlertBroadcaster(
	broker pubsub.Broker,
	deviceRepo repository.DeviceRepository,
	hub *websocket.Hub,
	log *zap.Logger,
) *AlertBroadcaster {
	return &AlertBroadcaster{
		broker:     broker,
		deviceRepo: deviceRepo,
		hub:        hub,
		log:        log,
	}
}

// Start launches the background listener in a context-aware goroutine.
// Unlike summary/device-status broadcasters, alert events are not debounced —
// every alert is a discrete, user-visible incident that must be delivered immediately.
func (ab *AlertBroadcaster) Start(ctx context.Context) {
	alertCh := ab.broker.Subscribe(pubsub.TopicAlertCreated)

	ab.log.Info("Alert broadcaster initialized and listening to topics")

	go func() {
		for {
			select {
			case <-ctx.Done():
				ab.log.Info("Alert broadcaster stopping due to context cancellation")
				return

			case val, ok := <-alertCh:
				if !ok {
					return
				}
				ab.handleEvent(ctx, val)
			}
		}
	}()
}

// handleEvent validates the incoming pubsub payload and dispatches the broadcast
// in a separate goroutine so the subscriber loop is never blocked.
func (ab *AlertBroadcaster) handleEvent(ctx context.Context, val interface{}) {
	alert, ok := val.(*model.Alert)
	if !ok || alert == nil {
		ab.log.Warn("Alert broadcaster: unexpected payload type",
			zap.String("type", typeOf(val)))
		return
	}

	// Spin off a goroutine so the subscriber channel is not blocked
	// during the device-name lookup.
	go ab.broadcastAlert(ctx, alert)
}

// broadcastAlert resolves the device name and emits the dashboard.alert.received
// WebSocket event to all connected clients.
func (ab *AlertBroadcaster) broadcastAlert(ctx context.Context, alert *model.Alert) {
	ab.log.Debug("Alert broadcast triggered",
		zap.String("alert_id", alert.ID.String()),
		zap.String("device_id", alert.DeviceID.String()))

	// Resolve device name with a bounded timeout so a slow DB does not
	// stall the broadcaster goroutine indefinitely.
	fetchCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	device, err := ab.deviceRepo.FindByID(fetchCtx, alert.DeviceID)
	if err != nil {
		ab.log.Error("Alert broadcaster: failed to fetch device for alert",
			zap.Error(err),
			zap.String("alert_id", alert.ID.String()),
			zap.String("device_id", alert.DeviceID.String()))
		return
	}

	payload := &AlertBroadcastPayload{
		AlertID:      alert.ID.String(),
		DeviceID:     alert.DeviceID.String(),
		DeviceName:   device.Name,
		Severity:     string(alert.Severity),
		Message:      alert.Message,
		CreatedAt:    alert.CreatedAt,
		Acknowledged: alert.Status != model.StatusActive,
	}

	msg, err := websocket.NewMessage("dashboard.alert.created", payload)
	if err != nil {
		ab.log.Error("Alert broadcaster: failed to marshal alert payload",
			zap.Error(err),
			zap.String("alert_id", alert.ID.String()))
		return
	}

	// Non-blocking broadcast to all connected WebSocket clients.
	ab.hub.Broadcast(msg)
	ab.log.Info("Alert broadcasted successfully via websocket",
		zap.String("alert_id", payload.AlertID),
		zap.String("device_id", payload.DeviceID),
		zap.String("device_name", payload.DeviceName),
		zap.String("severity", payload.Severity),
		zap.Bool("acknowledged", payload.Acknowledged))
}

// typeOf returns the runtime type string of any value (used for warning logs).
func typeOf(v interface{}) string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%T", v)
}
