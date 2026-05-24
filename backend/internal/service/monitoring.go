package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
)

// ── DTOs ─────────────────────────────────────────────────────────

// MonitoringDeviceListQuery carries validated query params from the HTTP layer.
type MonitoringDeviceListQuery struct {
	Search        string `form:"search"`
	EngineRunning string `form:"engine_running"` // "true" | "false" | ""
	Status        string `form:"status"`
	Online        string `form:"online"` // "true" | "false" | ""
	Limit         int    `form:"limit"`
	Offset        int    `form:"offset"`
	SortBy        string `form:"sort_by"`
	SortOrder     string `form:"sort_order"`
}

// MonitoringDeviceItem is a single device row in the list response.
type MonitoringDeviceItem struct {
	DeviceID            string     `json:"device_id"`
	DeviceCode          string     `json:"device_code"`
	DeviceName          string     `json:"device_name"`
	SerialNumber        string     `json:"serial_number"`
	EngineRunning       bool       `json:"engine_running"`
	Speed               *int32     `json:"speed"`
	CoolantTemperature  *float32   `json:"coolant_temperature"`
	OilPressure         *float32   `json:"oil_pressure"`
	FuelLevel           *float32   `json:"fuel_level"`
	BattVolt            *float32   `json:"batt_volt"`
	Frequency           *float32   `json:"frequency"`
	TotalVa             *float32   `json:"total_va"`
	PfAvg               *float32   `json:"pf_avg"`
	TelemetryRecordedAt *time.Time `json:"telemetry_recorded_at"`
	LastSeenAt          *time.Time `json:"last_seen_at"`
	LastOnlineAt        *time.Time `json:"last_online_at"`
	GSMSignal           *int16     `json:"gsm_signal"`
	GPSConnected        *bool      `json:"gps_connected"`
	ServerConnected     *bool      `json:"server_connected"`
	CANConnected        *bool      `json:"can_connected"`
	RS485Connected      *bool      `json:"rs485_connected"`
	SDCardOK            *bool      `json:"sd_card_ok"`
	DeviceStatus        string     `json:"device_status"`
}

// MonitoringDeviceListOutput wraps the list with pagination.
type MonitoringDeviceListOutput struct {
	Devices    []MonitoringDeviceItem `json:"devices"`
	Pagination PaginationMeta         `json:"pagination"`
}

// ── Detail DTOs ───────────────────────────────────────────────────

// MonitoringDeviceInfoOutput is section A of the detail response.
type MonitoringDeviceInfoOutput struct {
	ID              string      `json:"id"`
	DeviceCode      string      `json:"device_code"`
	Name            string      `json:"name"`
	SerialNumber    string      `json:"serial_number"`
	EngineID        string      `json:"engine_id"`
	GSMNumber       string      `json:"gsm_number"`
	FirmwareVersion string      `json:"firmware_version"`
	Status          string      `json:"status"`
	Metadata        interface{} `json:"metadata,omitempty"`
}

// MonitoringLatestStateOutput is section B of the detail response.
type MonitoringLatestStateOutput struct {
	EngineRunning       bool       `json:"engine_running"`
	Speed               *int32     `json:"speed"`
	CoolantTemperature  *float32   `json:"coolant_temperature"`
	OilPressure         *float32   `json:"oil_pressure"`
	FuelLevel           *float32   `json:"fuel_level"`
	BattVolt            *float32   `json:"batt_volt"`
	Frequency           *float32   `json:"frequency"`
	TotalVa             *float32   `json:"total_va"`
	PfAvg               *float32   `json:"pf_avg"`
	TelemetryRecordedAt *time.Time `json:"telemetry_recorded_at"`
	LastSeenAt          *time.Time `json:"last_seen_at"`
	LastOnlineAt        *time.Time `json:"last_online_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
}

// MonitoringConnectivityOutput is section C of the detail response.
type MonitoringConnectivityOutput struct {
	GSMSignal       int16     `json:"gsm_signal"`
	GPSConnected    bool      `json:"gps_connected"`
	ServerConnected bool      `json:"server_connected"`
	CANConnected    bool      `json:"can_connected"`
	RS485Connected  bool      `json:"rs485_connected"`
	SDCardOK        bool      `json:"sd_card_ok"`
	LastSeen        time.Time `json:"last_seen"`
}

// MonitoringDeviceDetailOutput is the complete response for GET /monitoring/devices/:device_id.
type MonitoringDeviceDetailOutput struct {
	DeviceInfo          MonitoringDeviceInfoOutput    `json:"device_info"`
	LatestState         *MonitoringLatestStateOutput  `json:"latest_state"`
	Connectivity        *MonitoringConnectivityOutput `json:"connectivity"`
	EngineTelemetry     *EngineTelemetryOutput        `json:"engine_telemetry"`
	ElectricalTelemetry *ElectricalTelemetryOutput    `json:"electrical_telemetry"`
}

// ── Allowed sort columns allow-list ──────────────────────────────

var monitoringAllowedSortBy = map[string]bool{
	"updated_at":            true,
	"last_seen_at":          true,
	"telemetry_recorded_at": true,
	"name":                  true,
}

// ── Interface ─────────────────────────────────────────────────────

// MonitoringService defines the contract for monitoring endpoints.
type MonitoringService interface {
	ListDevices(ctx context.Context, query MonitoringDeviceListQuery) (*MonitoringDeviceListOutput, error)
	GetDeviceDetail(ctx context.Context, deviceID uuid.UUID) (*MonitoringDeviceDetailOutput, error)
}

// ── Implementation ────────────────────────────────────────────────

type monitoringService struct {
	monitoringRepo repository.MonitoringRepository
	telemetryRepo  repository.TelemetryRepository
	log            *zap.Logger
}

// NewMonitoringService constructs a MonitoringService.
func NewMonitoringService(
	monitoringRepo repository.MonitoringRepository,
	telemetryRepo repository.TelemetryRepository,
	log *zap.Logger,
) MonitoringService {
	return &monitoringService{
		monitoringRepo: monitoringRepo,
		telemetryRepo:  telemetryRepo,
		log:            log,
	}
}

// ListDevices returns a paginated, filtered monitoring device list.
func (s *monitoringService) ListDevices(ctx context.Context, query MonitoringDeviceListQuery) (*MonitoringDeviceListOutput, error) {
	// ── Defaults & clamping ──────────────────────────────────────
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	// ── Validate sort_by allow-list ──────────────────────────────
	sortBy := ""
	if monitoringAllowedSortBy[query.SortBy] {
		sortBy = query.SortBy
	}

	// ── Validate sort_order ──────────────────────────────────────
	sortOrder := "desc"
	if query.SortOrder == "asc" {
		sortOrder = "asc"
	}

	// ── Parse boolean filters ────────────────────────────────────
	var engineRunning *bool
	switch query.EngineRunning {
	case "true":
		v := true
		engineRunning = &v
	case "false":
		v := false
		engineRunning = &v
	}

	var online *bool
	switch query.Online {
	case "true":
		v := true
		online = &v
	case "false":
		v := false
		online = &v
	}

	result, err := s.monitoringRepo.GetMonitoringDevices(ctx, repository.MonitoringDeviceFilter{
		Search:        query.Search,
		EngineRunning: engineRunning,
		Status:        query.Status,
		Online:        online,
		Limit:         query.Limit,
		Offset:        query.Offset,
		SortBy:        sortBy,
		SortOrder:     sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("monitoringService.ListDevices: %w", err)
	}

	items := make([]MonitoringDeviceItem, 0, len(result.Devices))
	for _, d := range result.Devices {
		items = append(items, monitoringDeviceToItem(d))
	}

	return &MonitoringDeviceListOutput{
		Devices: items,
		Pagination: PaginationMeta{
			Limit: query.Limit,
			Total: result.Total,
		},
	}, nil
}

// GetDeviceDetail returns the full realtime detail for a single device.
// Sections D & E use the existing TelemetryRepository to avoid code duplication.
func (s *monitoringService) GetDeviceDetail(ctx context.Context, deviceID uuid.UUID) (*MonitoringDeviceDetailOutput, error) {
	// ── Section A: Device Info ───────────────────────────────────
	info, err := s.monitoringRepo.GetMonitoringDeviceInfo(ctx, deviceID)
	if err != nil {
		if errors.Is(err, repository.ErrMonitoringDeviceNotFound) {
			return nil, ErrMonitoringDeviceNotFound
		}
		return nil, fmt.Errorf("monitoringService.GetDeviceDetail info: %w", err)
	}

	out := &MonitoringDeviceDetailOutput{
		DeviceInfo: monitoringDeviceInfoToOutput(info),
	}

	// ── Section B: Latest Device State ──────────────────────────
	latestState, err := s.monitoringRepo.GetMonitoringLatestState(ctx, deviceID)
	if err != nil {
		s.log.Warn("failed to get monitoring latest state", zap.String("device_id", deviceID.String()), zap.Error(err))
	} else if latestState != nil {
		out.LatestState = monitoringLatestStateToOutput(latestState)
	}

	// ── Section C: Connectivity ──────────────────────────────────
	connectivity, err := s.monitoringRepo.GetMonitoringConnectivity(ctx, deviceID)
	if err != nil {
		s.log.Warn("failed to get monitoring connectivity", zap.String("device_id", deviceID.String()), zap.Error(err))
	} else if connectivity != nil {
		out.Connectivity = monitoringConnectivityToOutput(connectivity)
	}

	// ── Section D: Latest Engine Telemetry ──────────────────────
	engineTelemetry, err := s.telemetryRepo.GetLatestEngine(ctx, deviceID)
	if err != nil {
		if !errors.Is(err, repository.ErrTelemetryNotFound) {
			s.log.Warn("failed to get latest engine telemetry", zap.String("device_id", deviceID.String()), zap.Error(err))
		}
	} else {
		out.EngineTelemetry = engineToOutput(engineTelemetry)
	}

	// ── Section E: Latest Electrical Telemetry ──────────────────
	electricalTelemetry, err := s.telemetryRepo.GetLatestElectrical(ctx, deviceID)
	if err != nil {
		if !errors.Is(err, repository.ErrTelemetryNotFound) {
			s.log.Warn("failed to get latest electrical telemetry", zap.String("device_id", deviceID.String()), zap.Error(err))
		}
	} else {
		out.ElectricalTelemetry = electricalToOutput(electricalTelemetry)
	}

	return out, nil
}

// ── Mappers ───────────────────────────────────────────────────────

func monitoringDeviceToItem(d model.MonitoringDevice) MonitoringDeviceItem {
	return MonitoringDeviceItem{
		DeviceID:            d.DeviceID,
		DeviceCode:          d.DeviceCode,
		DeviceName:          d.DeviceName,
		SerialNumber:        d.SerialNumber,
		EngineRunning:       d.EngineRunning,
		Speed:               d.Speed,
		CoolantTemperature:  d.CoolantTemperature,
		OilPressure:         d.OilPressure,
		FuelLevel:           d.FuelLevel,
		BattVolt:            d.BattVolt,
		Frequency:           d.Frequency,
		TotalVa:             d.TotalVa,
		PfAvg:               d.PfAvg,
		TelemetryRecordedAt: d.TelemetryRecordedAt,
		LastSeenAt:          d.LastSeenAt,
		LastOnlineAt:        d.LastOnlineAt,
		GSMSignal:           d.GSMSignal,
		GPSConnected:        d.GPSConnected,
		ServerConnected:     d.ServerConnected,
		CANConnected:        d.CANConnected,
		RS485Connected:      d.RS485Connected,
		SDCardOK:            d.SDCardOK,
		DeviceStatus:        d.DeviceStatus,
	}
}

func monitoringDeviceInfoToOutput(info *model.MonitoringDeviceInfo) MonitoringDeviceInfoOutput {
	out := MonitoringDeviceInfoOutput{
		ID:              info.ID,
		DeviceCode:      info.DeviceCode,
		Name:            info.Name,
		SerialNumber:    info.SerialNumber,
		EngineID:        info.EngineID,
		GSMNumber:       info.GSMNumber,
		FirmwareVersion: info.FirmwareVersion,
		Status:          info.Status,
	}

	if info.Metadata != nil {
		if raw, ok := info.Metadata.([]byte); ok && len(raw) > 0 {
			var m interface{}
			if err := json.Unmarshal(raw, &m); err == nil {
				out.Metadata = m
			}
		} else {
			out.Metadata = info.Metadata
		}
	}

	return out
}

func monitoringLatestStateToOutput(s *model.MonitoringLatestState) *MonitoringLatestStateOutput {
	return &MonitoringLatestStateOutput{
		EngineRunning:       s.EngineRunning,
		Speed:               s.Speed,
		CoolantTemperature:  s.CoolantTemperature,
		OilPressure:         s.OilPressure,
		FuelLevel:           s.FuelLevel,
		BattVolt:            s.BattVolt,
		Frequency:           s.Frequency,
		TotalVa:             s.TotalVa,
		PfAvg:               s.PfAvg,
		TelemetryRecordedAt: s.TelemetryRecordedAt,
		LastSeenAt:          s.LastSeenAt,
		LastOnlineAt:        s.LastOnlineAt,
		UpdatedAt:           s.UpdatedAt,
	}
}

func monitoringConnectivityToOutput(c *model.MonitoringConnectivity) *MonitoringConnectivityOutput {
	return &MonitoringConnectivityOutput{
		GSMSignal:       c.GSMSignal,
		GPSConnected:    c.GPSConnected,
		ServerConnected: c.ServerConnected,
		CANConnected:    c.CANConnected,
		RS485Connected:  c.RS485Connected,
		SDCardOK:        c.SDCardOK,
		LastSeen:        c.LastSeen,
	}
}

// ── Sentinel errors ───────────────────────────────────────────────

var ErrMonitoringDeviceNotFound = errors.New("monitoring device not found")
