package service

import (
	"context"
	"time"

	"github.com/danisbagus/genset-monitoring/backend/internal/pubsub"
	"github.com/danisbagus/genset-monitoring/backend/internal/websocket"
	"go.uber.org/zap"
)

// SummaryBroadcaster acts as a bridge between the Pub/Sub broker,
// the Dashboard service, and the WebSocket Hub to update clients in real-time.
type SummaryBroadcaster struct {
	broker       pubsub.Broker
	dashboardSvc DashboardService
	hub          *websocket.Hub
	log          *zap.Logger
}

// NewSummaryBroadcaster creates a new instance of SummaryBroadcaster.
func NewSummaryBroadcaster(
	broker pubsub.Broker,
	dashboardSvc DashboardService,
	hub *websocket.Hub,
	log *zap.Logger,
) *SummaryBroadcaster {
	return &SummaryBroadcaster{
		broker:       broker,
		dashboardSvc: dashboardSvc,
		hub:          hub,
		log:          log,
	}
}

// Start launches the background listener in a context-aware goroutine.
// It uses a robust debouncer to ensure rapid updates (e.g. dozens of heartbeats)
// are aggregated and do not thrash the database.
func (sb *SummaryBroadcaster) Start(ctx context.Context) {
	deviceCh := sb.broker.Subscribe(pubsub.TopicDeviceCRUD)
	alertCh := sb.broker.Subscribe(pubsub.TopicAlertCreated)
	heartbeatCh := sb.broker.Subscribe(pubsub.TopicHeartbeatReceived)
	telemetryCh := sb.broker.Subscribe(pubsub.TopicTelemetryEngineCreated)

	sb.log.Info("Dashboard summary broadcaster initialized and listening to topics")

	go func() {
		var delayTimer *time.Timer
		var updateTriggered bool
		var triggerSource string
		var timerCh <-chan time.Time

		// Ensure the timer is cleanly disposed when exiting
		defer func() {
			if delayTimer != nil {
				delayTimer.Stop()
			}
		}()

		// Aggregation window for debouncing updates
		const debounceInterval = 1 * time.Second

		for {
			select {
			case <-ctx.Done():
				sb.log.Info("Dashboard summary broadcaster stopping due to context cancellation")
				return

			case <-deviceCh:
				triggerSource = "device_crud"
				if timerCh == nil {
					delayTimer = time.NewTimer(debounceInterval)
					timerCh = delayTimer.C
				} else {
					updateTriggered = true
				}

			case <-alertCh:
				triggerSource = "alert_created"
				if timerCh == nil {
					delayTimer = time.NewTimer(debounceInterval)
					timerCh = delayTimer.C
				} else {
					updateTriggered = true
				}

			case <-heartbeatCh:
				triggerSource = "heartbeat_received"
				if timerCh == nil {
					delayTimer = time.NewTimer(debounceInterval)
					timerCh = delayTimer.C
				} else {
					updateTriggered = true
				}

			case <-telemetryCh:
				triggerSource = "telemetry_created"
				if timerCh == nil {
					delayTimer = time.NewTimer(debounceInterval)
					timerCh = delayTimer.C
				} else {
					updateTriggered = true
				}

			case <-timerCh:
				// Debounce delay finished, trigger the real summary broadcast
				sb.broadcastSummary(ctx, triggerSource)

				// Reset current timer
				if delayTimer != nil {
					delayTimer.Stop()
				}
				timerCh = nil

				// If other events queued up during the delay, schedule another update
				if updateTriggered {
					updateTriggered = false
					delayTimer = time.NewTimer(debounceInterval)
					timerCh = delayTimer.C
				}
			}
		}
	}()
}

// broadcastSummary fetches the dashboard summary using the existing service
// and broadcasts it as a 'dashboard.summary.updated' websocket event.
func (sb *SummaryBroadcaster) broadcastSummary(ctx context.Context, triggerSource string) {
	sb.log.Debug("Dashboard summary broadcast triggered", zap.String("source", triggerSource))

	// Use a 5-second timeout for the DB fetch to keep operations snappy
	fetchCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	summary, err := sb.dashboardSvc.GetSummary(fetchCtx)
	if err != nil {
		sb.log.Error("Failed to fetch dashboard summary for websocket broadcast",
			zap.Error(err),
			zap.String("source", triggerSource))
		return
	}

	msg, err := websocket.NewMessage("dashboard.summary.updated", summary)
	if err != nil {
		sb.log.Error("Failed to marshal dashboard summary for websocket broadcast",
			zap.Error(err),
			zap.String("source", triggerSource))
		return
	}

	// Non-blocking broadcast
	sb.hub.Broadcast(msg)
	sb.log.Info("Dashboard summary broadcasted successfully",
		zap.String("source", triggerSource),
		zap.Int64("total_devices", summary.TotalDevices),
		zap.Int64("online_devices", summary.OnlineDevices),
		zap.Int64("offline_devices", summary.OfflineDevices),
		zap.Int64("running_engines", summary.RunningEngines),
		zap.Int64("critical_alerts", summary.CriticalAlerts),
		zap.Int64("warning_alerts", summary.WarningAlerts),
	)
}
