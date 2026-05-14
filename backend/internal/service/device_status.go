package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"github.com/danisbagus/genset-monitoring/backend/internal/websocket"
)

// ── DTOs ─────────────────────────────────────────────────────────

// DeviceStatusOutput is the public representation of a device's status.
type DeviceStatusOutput struct {
	DeviceID        uuid.UUID `json:"device_id"`
	LastSeen        time.Time `json:"last_seen"`
	GSMSignal       int16     `json:"gsm_signal"`
	GPSConnected    bool      `json:"gps_connected"`
	ServerConnected bool      `json:"server_connected"`
	CANConnected    bool      `json:"can_connected"`
	RS485Connected  bool      `json:"rs485_connected"`
	SDCardOK        bool      `json:"sd_card_ok"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// HeartbeatInput carries the data payload from the device.
type HeartbeatInput struct {
	DeviceID        uuid.UUID
	GSMSignal       int16
	GPSConnected    bool
	ServerConnected bool
	CANConnected    bool
	RS485Connected  bool
	SDCardOK        bool
}

// ── Interface ─────────────────────────────────────────────────────

type DeviceStatusService interface {
	GetStatus(ctx context.Context, deviceID uuid.UUID) (*DeviceStatusOutput, error)
	Heartbeat(ctx context.Context, input HeartbeatInput) error
}

// ── Implementation ────────────────────────────────────────────────

type deviceStatusService struct {
	statusRepo repository.DeviceStatusRepository
	deviceRepo repository.DeviceRepository
	wsHub      *websocket.Hub
	log        *zap.Logger
}

func NewDeviceStatusService(
	statusRepo repository.DeviceStatusRepository,
	deviceRepo repository.DeviceRepository,
	wsHub *websocket.Hub,
	log *zap.Logger,
) DeviceStatusService {
	return &deviceStatusService{
		statusRepo: statusRepo,
		deviceRepo: deviceRepo,
		wsHub:      wsHub,
		log:        log,
	}
}

func (s *deviceStatusService) GetStatus(ctx context.Context, deviceID uuid.UUID) (*DeviceStatusOutput, error) {
	status, err := s.statusRepo.FindByDeviceID(ctx, deviceID)
	if err != nil {
		if errors.Is(err, repository.ErrDeviceStatusNotFound) {
			return nil, ErrStatusNotFound
		}
		return nil, fmt.Errorf("deviceStatusService.GetStatus: %w", err)
	}

	return &DeviceStatusOutput{
		DeviceID:        status.DeviceID,
		LastSeen:        status.LastSeen,
		GSMSignal:       status.GSMSignal,
		GPSConnected:    status.GPSConnected,
		ServerConnected: status.ServerConnected,
		CANConnected:    status.CANConnected,
		RS485Connected:  status.RS485Connected,
		SDCardOK:        status.SDCardOK,
		UpdatedAt:       status.UpdatedAt,
	}, nil
}

func (s *deviceStatusService) Heartbeat(ctx context.Context, input HeartbeatInput) error {
	// 1. Verify device exists (optional but recommended for security)
	if _, err := s.deviceRepo.FindByID(ctx, input.DeviceID); err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return ErrDeviceNotFound
		}
		return fmt.Errorf("deviceStatusService.Heartbeat: device validation: %w", err)
	}

	// 2. Prepare status update
	now := time.Now().UTC()
	lastSeen := now

	status := &model.DeviceStatus{
		DeviceID:        input.DeviceID,
		LastSeen:        lastSeen,
		GSMSignal:       input.GSMSignal,
		GPSConnected:    input.GPSConnected,
		ServerConnected: input.ServerConnected,
		CANConnected:    input.CANConnected,
		RS485Connected:  input.RS485Connected,
		SDCardOK:        input.SDCardOK,
		UpdatedAt:       now,
	}

	// 3. Upsert status record
	if err := s.statusRepo.Upsert(ctx, status); err != nil {
		return fmt.Errorf("deviceStatusService.Heartbeat: upsert: %w", err)
	}

	// 4. Update latest state (realtime snapshot)
	err := s.statusRepo.UpsertLatestState(ctx, &model.DeviceLatestState{
		DeviceID:     input.DeviceID,
		LastOnlineAt: &now,
	})
	if err != nil {
		s.log.Error("Failed to update device latest state from heartbeat",
			zap.Error(err),
			zap.String("device_id", input.DeviceID.String()))
	}

	// Broadcast via websocket
	go func() {
		output := &DeviceStatusOutput{
			DeviceID:        status.DeviceID,
			LastSeen:        status.LastSeen,
			GSMSignal:       status.GSMSignal,
			GPSConnected:    status.GPSConnected,
			ServerConnected: status.ServerConnected,
			CANConnected:    status.CANConnected,
			RS485Connected:  status.RS485Connected,
			SDCardOK:        status.SDCardOK,
			UpdatedAt:       status.UpdatedAt,
		}

		msg, err := websocket.NewMessage("device.status.updated", output)
		if err != nil {
			s.log.Warn("Failed to marshal device status for websocket",
				zap.Error(err),
				zap.String("device_id", input.DeviceID.String()))
			return
		}

		s.wsHub.Broadcast(msg)
		s.log.Info("Device status broadcasted via websocket",
			zap.String("device_id", input.DeviceID.String()))
	}()

	s.log.Debug("heartbeat processed", zap.String("device_id", input.DeviceID.String()))
	return nil
}

// ── Sentinel errors ───────────────────────────────────────────────

var (
	ErrStatusNotFound = errors.New("device status record not found")
)
