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

type EngineTelemetryOutput struct {
	DeviceID                  uuid.UUID `json:"device_id"`
	Speed                     *int32    `json:"speed"`
	FuelRate                  *float32  `json:"fuel_rate"`
	RatedPower                *float32  `json:"rated_power"`
	RatedSpeed                *int32    `json:"rated_speed"`
	OilFilterOutPressure      *float32  `json:"oil_filter_out_pressure"`
	DesiredOperatingSpeed     *int32    `json:"desired_operating_speed"`
	OilPressure               *float32  `json:"oil_pressure"`
	CoolantTemperature        *float32  `json:"coolant_temperature"`
	TripFuel                  *float32  `json:"trip_fuel"`
	TotalFuel                 *float32  `json:"total_fuel"`
	AvgFuelRate               *float32  `json:"avg_fuel_rate"`
	IntakeManifoldPressure    *float32  `json:"intake_manifold_pressure"`
	IntakeManifoldTemperature *float32  `json:"intake_manifold_temperature"`
	KeyswitchBattPotential    *float32  `json:"keyswitch_batt_potential"`
	BattVolt                  *float32  `json:"batt_volt"`
	EcuTemperature            *float32  `json:"ecu_temperature"`
	RunTime                   *int64    `json:"run_time"`
	FuelLevelTop              *float32  `json:"fuel_level_top"`
	FuelLevelBottom           *float32  `json:"fuel_level_bottom"`
	FuelLevelPressure1        *float32  `json:"fuel_level_pressure_1"`
	FuelLevelPressure2        *float32  `json:"fuel_level_pressure_2"`
	TurboPressure             *float32  `json:"turbo_pressure"`
	CreatedAt                 time.Time `json:"created_at"`
}

type ElectricalTelemetryOutput struct {
	DeviceID      uuid.UUID `json:"device_id"`
	ChargeAltVolt *float32  `json:"charge_alt_volt"`
	Frequency     *float32  `json:"frequency"`
	L1NVolt       *float32  `json:"l1_n_volt"`
	L2NVolt       *float32  `json:"l2_n_volt"`
	L3NVolt       *float32  `json:"l3_n_volt"`
	L1L2Volt      *float32  `json:"l1_l2_volt"`
	L2L3Volt      *float32  `json:"l2_l3_volt"`
	L3L1Volt      *float32  `json:"l3_l1_volt"`
	L1Curr        *float32  `json:"l1_curr"`
	L2Curr        *float32  `json:"l2_curr"`
	L3Curr        *float32  `json:"l3_curr"`
	EarthCurr     *float32  `json:"earth_curr"`
	L1Va          *float32  `json:"l1_va"`
	L2Va          *float32  `json:"l2_va"`
	L3Va          *float32  `json:"l3_va"`
	TotalVa       *float32  `json:"total_va"`
	L1Var         *float32  `json:"l1_var"`
	L2Var         *float32  `json:"l2_var"`
	L3Var         *float32  `json:"l3_var"`
	TotalVar      *float32  `json:"total_var"`
	PfL1          *float32  `json:"pf_l1"`
	PfL2          *float32  `json:"pf_l2"`
	PfL3          *float32  `json:"pf_l3"`
	PfAvg         *float32  `json:"pf_avg"`
	PercentFv     *float32  `json:"percent_fv"`
	PercentFp     *float32  `json:"percent_fp"`
	CreatedAt     time.Time `json:"created_at"`
}

// ── Input structs ─────────────────────────────────────────────────

type CreateEngineTelemetryInput struct {
	DeviceID                  uuid.UUID
	Speed                     *int32
	FuelRate                  *float32
	RatedPower                *float32
	RatedSpeed                *int32
	OilFilterOutPressure      *float32
	DesiredOperatingSpeed     *int32
	OilPressure               *float32
	CoolantTemperature        *float32
	TripFuel                  *float32
	TotalFuel                 *float32
	AvgFuelRate               *float32
	IntakeManifoldPressure    *float32
	IntakeManifoldTemperature *float32
	KeyswitchBattPotential    *float32
	BattVolt                  *float32
	EcuTemperature            *float32
	RunTime                   *int64
	FuelLevelTop              *float32
	FuelLevelBottom           *float32
	FuelLevelPressure1        *float32
	FuelLevelPressure2        *float32
	TurboPressure             *float32
}

type CreateElectricalTelemetryInput struct {
	DeviceID      uuid.UUID
	ChargeAltVolt *float32
	Frequency     *float32
	L1NVolt       *float32
	L2NVolt       *float32
	L3NVolt       *float32
	L1L2Volt      *float32
	L2L3Volt      *float32
	L3L1Volt      *float32
	L1Curr        *float32
	L2Curr        *float32
	L3Curr        *float32
	EarthCurr     *float32
	L1Va          *float32
	L2Va          *float32
	L3Va          *float32
	TotalVa       *float32
	L1Var         *float32
	L2Var         *float32
	L3Var         *float32
	TotalVar      *float32
	PfL1          *float32
	PfL2          *float32
	PfL3          *float32
	PfAvg         *float32
	PercentFv     *float32
	PercentFp     *float32
}

// ── Interface ─────────────────────────────────────────────────────

type TelemetryService interface {
	CreateEngine(ctx context.Context, input CreateEngineTelemetryInput) (*EngineTelemetryOutput, error)
	GetLatestEngine(ctx context.Context, deviceID uuid.UUID) (*EngineTelemetryOutput, error)
	CreateElectrical(ctx context.Context, input CreateElectricalTelemetryInput) (*ElectricalTelemetryOutput, error)
	GetLatestElectrical(ctx context.Context, deviceID uuid.UUID) (*ElectricalTelemetryOutput, error)
}

// ── Implementation ────────────────────────────────────────────────

type telemetryService struct {
	telemetryRepo repository.TelemetryRepository
	deviceRepo    repository.DeviceRepository
	wsHub         *websocket.Hub
	log           *zap.Logger
}

func NewTelemetryService(
	telemetryRepo repository.TelemetryRepository,
	deviceRepo repository.DeviceRepository,
	wsHub *websocket.Hub,
	log *zap.Logger,
) TelemetryService {
	return &telemetryService{
		telemetryRepo: telemetryRepo,
		deviceRepo:    deviceRepo,
		wsHub:         wsHub,
		log:           log,
	}
}

func (s *telemetryService) CreateEngine(ctx context.Context, input CreateEngineTelemetryInput) (*EngineTelemetryOutput, error) {
	// Check if device exists
	if _, err := s.deviceRepo.FindByID(ctx, input.DeviceID); err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("telemetryService.CreateEngine: %w", err)
	}

	telemetry := &model.EngineTelemetry{
		DeviceID:                  input.DeviceID,
		Speed:                     input.Speed,
		FuelRate:                  input.FuelRate,
		RatedPower:                input.RatedPower,
		RatedSpeed:                input.RatedSpeed,
		OilFilterOutPressure:      input.OilFilterOutPressure,
		DesiredOperatingSpeed:     input.DesiredOperatingSpeed,
		OilPressure:               input.OilPressure,
		CoolantTemperature:        input.CoolantTemperature,
		TripFuel:                  input.TripFuel,
		TotalFuel:                 input.TotalFuel,
		AvgFuelRate:               input.AvgFuelRate,
		IntakeManifoldPressure:    input.IntakeManifoldPressure,
		IntakeManifoldTemperature: input.IntakeManifoldTemperature,
		KeyswitchBattPotential:   input.KeyswitchBattPotential,
		BattVolt:                  input.BattVolt,
		EcuTemperature:            input.EcuTemperature,
		RunTime:                  input.RunTime,
		FuelLevelTop:              input.FuelLevelTop,
		FuelLevelBottom:           input.FuelLevelBottom,
		FuelLevelPressure1:       input.FuelLevelPressure1,
		FuelLevelPressure2:       input.FuelLevelPressure2,
		TurboPressure:             input.TurboPressure,
	}

	if err := s.telemetryRepo.CreateEngine(ctx, telemetry); err != nil {
		return nil, fmt.Errorf("telemetryService.CreateEngine: %w", err)
	}

	output := engineToOutput(telemetry)

	// Broadcast via websocket
	go func() {
		msg, err := websocket.NewMessage("engine.telemetry.updated", output)
		if err != nil {
			s.log.Warn("Failed to marshal engine telemetry for websocket",
				zap.Error(err),
				zap.String("device_id", input.DeviceID.String()))
			return
		}

		s.wsHub.Broadcast(msg)
		s.log.Info("Engine telemetry broadcasted via websocket",
			zap.String("device_id", input.DeviceID.String()))
	}()

	return output, nil
}

func (s *telemetryService) GetLatestEngine(ctx context.Context, deviceID uuid.UUID) (*EngineTelemetryOutput, error) {
	telemetry, err := s.telemetryRepo.GetLatestEngine(ctx, deviceID)
	if err != nil {
		if errors.Is(err, repository.ErrTelemetryNotFound) {
			return nil, ErrTelemetryNotFound
		}
		return nil, fmt.Errorf("telemetryService.GetLatestEngine: %w", err)
	}

	return engineToOutput(telemetry), nil
}

func (s *telemetryService) CreateElectrical(ctx context.Context, input CreateElectricalTelemetryInput) (*ElectricalTelemetryOutput, error) {
	// Check if device exists
	if _, err := s.deviceRepo.FindByID(ctx, input.DeviceID); err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("telemetryService.CreateElectrical: %w", err)
	}

	telemetry := &model.ElectricalTelemetry{
		DeviceID:      input.DeviceID,
		ChargeAltVolt: input.ChargeAltVolt,
		Frequency:     input.Frequency,
		L1NVolt:       input.L1NVolt,
		L2NVolt:       input.L2NVolt,
		L3NVolt:       input.L3NVolt,
		L1L2Volt:      input.L1L2Volt,
		L2L3Volt:      input.L2L3Volt,
		L3L1Volt:      input.L3L1Volt,
		L1Curr:        input.L1Curr,
		L2Curr:        input.L2Curr,
		L3Curr:        input.L3Curr,
		EarthCurr:     input.EarthCurr,
		L1Va:          input.L1Va,
		L2Va:          input.L2Va,
		L3Va:          input.L3Va,
		TotalVa:       input.TotalVa,
		L1Var:         input.L1Var,
		L2Var:         input.L2Var,
		L3Var:         input.L3Var,
		TotalVar:      input.TotalVar,
		PfL1:          input.PfL1,
		PfL2:          input.PfL2,
		PfL3:          input.PfL3,
		PfAvg:         input.PfAvg,
		PercentFv:     input.PercentFv,
		PercentFp:     input.PercentFp,
	}

	if err := s.telemetryRepo.CreateElectrical(ctx, telemetry); err != nil {
		return nil, fmt.Errorf("telemetryService.CreateElectrical: %w", err)
	}

	output := electricalToOutput(telemetry)

	// Broadcast via websocket
	go func() {
		msg, err := websocket.NewMessage("electrical.telemetry.updated", output)
		if err != nil {
			s.log.Warn("Failed to marshal electrical telemetry for websocket",
				zap.Error(err),
				zap.String("device_id", input.DeviceID.String()))
			return
		}

		s.wsHub.Broadcast(msg)
		s.log.Info("Electrical telemetry broadcasted via websocket",
			zap.String("device_id", input.DeviceID.String()))
	}()

	return output, nil
}

func (s *telemetryService) GetLatestElectrical(ctx context.Context, deviceID uuid.UUID) (*ElectricalTelemetryOutput, error) {
	telemetry, err := s.telemetryRepo.GetLatestElectrical(ctx, deviceID)
	if err != nil {
		if errors.Is(err, repository.ErrTelemetryNotFound) {
			return nil, ErrTelemetryNotFound
		}
		return nil, fmt.Errorf("telemetryService.GetLatestElectrical: %w", err)
	}

	return electricalToOutput(telemetry), nil
}

// ── Helpers ───────────────────────────────────────────────────────

func engineToOutput(m *model.EngineTelemetry) *EngineTelemetryOutput {
	return &EngineTelemetryOutput{
		DeviceID:                  m.DeviceID,
		Speed:                     m.Speed,
		FuelRate:                  m.FuelRate,
		RatedPower:                m.RatedPower,
		RatedSpeed:                m.RatedSpeed,
		OilFilterOutPressure:      m.OilFilterOutPressure,
		DesiredOperatingSpeed:     m.DesiredOperatingSpeed,
		OilPressure:               m.OilPressure,
		CoolantTemperature:        m.CoolantTemperature,
		TripFuel:                  m.TripFuel,
		TotalFuel:                 m.TotalFuel,
		AvgFuelRate:               m.AvgFuelRate,
		IntakeManifoldPressure:    m.IntakeManifoldPressure,
		IntakeManifoldTemperature: m.IntakeManifoldTemperature,
		KeyswitchBattPotential:   m.KeyswitchBattPotential,
		BattVolt:                  m.BattVolt,
		EcuTemperature:            m.EcuTemperature,
		RunTime:                  m.RunTime,
		FuelLevelTop:              m.FuelLevelTop,
		FuelLevelBottom:           m.FuelLevelBottom,
		FuelLevelPressure1:       m.FuelLevelPressure1,
		FuelLevelPressure2:       m.FuelLevelPressure2,
		TurboPressure:             m.TurboPressure,
		CreatedAt:                 m.CreatedAt,
	}
}

func electricalToOutput(m *model.ElectricalTelemetry) *ElectricalTelemetryOutput {
	return &ElectricalTelemetryOutput{
		DeviceID:      m.DeviceID,
		ChargeAltVolt: m.ChargeAltVolt,
		Frequency:     m.Frequency,
		L1NVolt:       m.L1NVolt,
		L2NVolt:       m.L2NVolt,
		L3NVolt:       m.L3NVolt,
		L1L2Volt:      m.L1L2Volt,
		L2L3Volt:      m.L2L3Volt,
		L3L1Volt:      m.L3L1Volt,
		L1Curr:        m.L1Curr,
		L2Curr:        m.L2Curr,
		L3Curr:        m.L3Curr,
		EarthCurr:     m.EarthCurr,
		L1Va:          m.L1Va,
		L2Va:          m.L2Va,
		L3Va:          m.L3Va,
		TotalVa:       m.TotalVa,
		L1Var:         m.L1Var,
		L2Var:         m.L2Var,
		L3Var:         m.L3Var,
		TotalVar:      m.TotalVar,
		PfL1:          m.PfL1,
		PfL2:          m.PfL2,
		PfL3:          m.PfL3,
		PfAvg:         m.PfAvg,
		PercentFv:     m.PercentFv,
		PercentFp:     m.PercentFp,
		CreatedAt:     m.CreatedAt,
	}
}

// ── Sentinel errors ───────────────────────────────────────────────

var (
	ErrTelemetryNotFound = errors.New("telemetry not found")
)
