package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	"github.com/danisbagus/genset-monitoring/backend/pkg/response"
	"github.com/danisbagus/genset-monitoring/backend/pkg/validator"
)

// ── Request DTOs ─────────────────────────────────────────────────

type CreateEngineTelemetryRequest struct {
	Speed                     *int32   `json:"speed"`
	FuelRate                  *float32 `json:"fuel_rate"`
	RatedPower                *float32 `json:"rated_power"`
	RatedSpeed                *int32   `json:"rated_speed"`
	OilFilterOutPressure      *float32 `json:"oil_filter_out_pressure"`
	DesiredOperatingSpeed     *int32   `json:"desired_operating_speed"`
	OilPressure               *float32 `json:"oil_pressure"`
	CoolantTemperature        *float32 `json:"coolant_temperature"`
	TripFuel                  *float32 `json:"trip_fuel"`
	TotalFuel                 *float32 `json:"total_fuel"`
	AvgFuelRate               *float32 `json:"avg_fuel_rate"`
	IntakeManifoldPressure    *float32 `json:"intake_manifold_pressure"`
	IntakeManifoldTemperature *float32 `json:"intake_manifold_temperature"`
	KeyswitchBattPotential   *float32 `json:"keyswitch_batt_potential"`
	BattVolt                  *float32 `json:"batt_volt"`
	EcuTemperature            *float32 `json:"ecu_temperature"`
	RunTime                   *int64   `json:"run_time"`
	FuelLevelTop              *float32 `json:"fuel_level_top"`
	FuelLevelBottom           *float32 `json:"fuel_level_bottom"`
	FuelLevelPressure1        *float32 `json:"fuel_level_pressure_1"`
	FuelLevelPressure2        *float32 `json:"fuel_level_pressure_2"`
	TurboPressure             *float32 `json:"turbo_pressure"`
}

type CreateElectricalTelemetryRequest struct {
	ChargeAltVolt *float32 `json:"charge_alt_volt"`
	Frequency     *float32 `json:"frequency"`
	L1NVolt       *float32 `json:"l1_n_volt"`
	L2NVolt       *float32 `json:"l2_n_volt"`
	L3NVolt       *float32 `json:"l3_n_volt"`
	L1L2Volt      *float32 `json:"l1_l2_volt"`
	L2L3Volt      *float32 `json:"l2_l3_volt"`
	L3L1Volt      *float32 `json:"l3_l1_volt"`
	L1Curr        *float32 `json:"l1_curr"`
	L2Curr        *float32 `json:"l2_curr"`
	L3Curr        *float32 `json:"l3_curr"`
	EarthCurr     *float32 `json:"earth_curr"`
	L1Va          *float32 `json:"l1_va"`
	L2Va          *float32 `json:"l2_va"`
	L3Va          *float32 `json:"l3_va"`
	TotalVa       *float32 `json:"total_va"`
	L1Var         *float32 `json:"l1_var"`
	L2Var         *float32 `json:"l2_var"`
	L3Var         *float32 `json:"l3_var"`
	TotalVar      *float32 `json:"total_var"`
	PfL1          *float32 `json:"pf_l1"`
	PfL2          *float32 `json:"pf_l2"`
	PfL3          *float32 `json:"pf_l3"`
	PfAvg         *float32 `json:"pf_avg"`
	PercentFv     *float32 `json:"percent_fv"`
	PercentFp     *float32 `json:"percent_fp"`
}

// ── Handler ───────────────────────────────────────────────────────

type TelemetryHandler struct {
	telemetrySvc service.TelemetryService
	validator    *validator.Validator
	log          *zap.Logger
}

func NewTelemetryHandler(telemetrySvc service.TelemetryService, v *validator.Validator, log *zap.Logger) *TelemetryHandler {
	return &TelemetryHandler{
		telemetrySvc: telemetrySvc,
		validator:    v,
		log:          log,
	}
}

// ── Endpoints ─────────────────────────────────────────────────────

// CreateEngine registers new engine telemetry.
//
//	@Summary		Create Engine Telemetry
//	@Description	Registers new engine telemetry for a specific device
//	@Tags			telemetry
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			deviceID	path		string							true	"Device UUID"
//	@Param			body		body		CreateEngineTelemetryRequest	true	"Engine telemetry payload"
//	@Success		201			{object}	response.Response{data=service.EngineTelemetryOutput}	"Telemetry created"
//	@Failure		400			{object}	response.ErrorResponse									"Invalid UUID or request body"
//	@Failure		401			{object}	response.ErrorResponse									"Unauthorized"
//	@Failure		404			{object}	response.ErrorResponse									"Device not found"
//	@Failure		500			{object}	response.ErrorResponse									"Internal server error"
//	@Router			/api/v1/devices/{deviceID}/engine [post]
func (h *TelemetryHandler) CreateEngine(c *gin.Context) {
	deviceID, ok := parseDeviceID(c)
	if !ok {
		return
	}

	var req CreateEngineTelemetryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	out, err := h.telemetrySvc.CreateEngine(c.Request.Context(), service.CreateEngineTelemetryInput{
		DeviceID:                  deviceID,
		Speed:                     req.Speed,
		FuelRate:                  req.FuelRate,
		RatedPower:                req.RatedPower,
		RatedSpeed:                req.RatedSpeed,
		OilFilterOutPressure:      req.OilFilterOutPressure,
		DesiredOperatingSpeed:     req.DesiredOperatingSpeed,
		OilPressure:               req.OilPressure,
		CoolantTemperature:        req.CoolantTemperature,
		TripFuel:                  req.TripFuel,
		TotalFuel:                 req.TotalFuel,
		AvgFuelRate:               req.AvgFuelRate,
		IntakeManifoldPressure:    req.IntakeManifoldPressure,
		IntakeManifoldTemperature: req.IntakeManifoldTemperature,
		KeyswitchBattPotential:   req.KeyswitchBattPotential,
		BattVolt:                  req.BattVolt,
		EcuTemperature:            req.EcuTemperature,
		RunTime:                  req.RunTime,
		FuelLevelTop:              req.FuelLevelTop,
		FuelLevelBottom:           req.FuelLevelBottom,
		FuelLevelPressure1:       req.FuelLevelPressure1,
		FuelLevelPressure2:       req.FuelLevelPressure2,
		TurboPressure:             req.TurboPressure,
	})

	if err != nil {
		h.handleTelemetryError(c, err)
		return
	}

	response.Created(c, "engine telemetry created successfully", out)
}

// GetLatestEngine returns the most recent engine telemetry.
//
//	@Summary		Get Latest Engine Telemetry
//	@Description	Returns the most recent engine telemetry for a specific device
//	@Tags			telemetry
//	@Produce		json
//	@Security		BearerAuth
//	@Param			deviceID	path		string	true	"Device UUID"
//	@Success		200			{object}	response.Response{data=service.EngineTelemetryOutput}	"Latest telemetry"
//	@Failure		400			{object}	response.ErrorResponse									"Invalid UUID"
//	@Failure		401			{object}	response.ErrorResponse									"Unauthorized"
//	@Failure		404			{object}	response.ErrorResponse									"Telemetry not found"
//	@Failure		500			{object}	response.ErrorResponse									"Internal server error"
//	@Router			/api/v1/devices/{deviceID}/engine/latest [get]
func (h *TelemetryHandler) GetLatestEngine(c *gin.Context) {
	deviceID, ok := parseDeviceID(c)
	if !ok {
		return
	}

	out, err := h.telemetrySvc.GetLatestEngine(c.Request.Context(), deviceID)
	if err != nil {
		h.handleTelemetryError(c, err)
		return
	}

	response.OK(c, "latest engine telemetry retrieved successfully", out)
}

// CreateElectrical registers new electrical telemetry.
//
//	@Summary		Create Electrical Telemetry
//	@Description	Registers new electrical telemetry for a specific device
//	@Tags			telemetry
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			deviceID	path		string								true	"Device UUID"
//	@Param			body		body		CreateElectricalTelemetryRequest	true	"Electrical telemetry payload"
//	@Success		201			{object}	response.Response{data=service.ElectricalTelemetryOutput}	"Telemetry created"
//	@Failure		400			{object}	response.ErrorResponse										"Invalid UUID or request body"
//	@Failure		401			{object}	response.ErrorResponse										"Unauthorized"
//	@Failure		404			{object}	response.ErrorResponse										"Device not found"
//	@Failure		500			{object}	response.ErrorResponse										"Internal server error"
//	@Router			/api/v1/devices/{deviceID}/electrical [post]
func (h *TelemetryHandler) CreateElectrical(c *gin.Context) {
	deviceID, ok := parseDeviceID(c)
	if !ok {
		return
	}

	var req CreateElectricalTelemetryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	out, err := h.telemetrySvc.CreateElectrical(c.Request.Context(), service.CreateElectricalTelemetryInput{
		DeviceID:      deviceID,
		ChargeAltVolt: req.ChargeAltVolt,
		Frequency:     req.Frequency,
		L1NVolt:       req.L1NVolt,
		L2NVolt:       req.L2NVolt,
		L3NVolt:       req.L3NVolt,
		L1L2Volt:      req.L1L2Volt,
		L2L3Volt:      req.L2L3Volt,
		L3L1Volt:      req.L3L1Volt,
		L1Curr:        req.L1Curr,
		L2Curr:        req.L2Curr,
		L3Curr:        req.L3Curr,
		EarthCurr:     req.EarthCurr,
		L1Va:          req.L1Va,
		L2Va:          req.L2Va,
		L3Va:          req.L3Va,
		TotalVa:       req.TotalVa,
		L1Var:         req.L1Var,
		L2Var:         req.L2Var,
		L3Var:         req.L3Var,
		TotalVar:      req.TotalVar,
		PfL1:          req.PfL1,
		PfL2:          req.PfL2,
		PfL3:          req.PfL3,
		PfAvg:         req.PfAvg,
		PercentFv:     req.PercentFv,
		PercentFp:     req.PercentFp,
	})

	if err != nil {
		h.handleTelemetryError(c, err)
		return
	}

	response.Created(c, "electrical telemetry created successfully", out)
}

// GetLatestElectrical returns the most recent electrical telemetry.
//
//	@Summary		Get Latest Electrical Telemetry
//	@Description	Returns the most recent electrical telemetry for a specific device
//	@Tags			telemetry
//	@Produce		json
//	@Security		BearerAuth
//	@Param			deviceID	path		string	true	"Device UUID"
//	@Success		200			{object}	response.Response{data=service.ElectricalTelemetryOutput}	"Latest telemetry"
//	@Failure		400			{object}	response.ErrorResponse										"Invalid UUID"
//	@Failure		401			{object}	response.ErrorResponse										"Unauthorized"
//	@Failure		404			{object}	response.ErrorResponse										"Telemetry not found"
//	@Failure		500			{object}	response.ErrorResponse										"Internal server error"
//	@Router			/api/v1/devices/{deviceID}/electrical/latest [get]
func (h *TelemetryHandler) GetLatestElectrical(c *gin.Context) {
	deviceID, ok := parseDeviceID(c)
	if !ok {
		return
	}

	out, err := h.telemetrySvc.GetLatestElectrical(c.Request.Context(), deviceID)
	if err != nil {
		h.handleTelemetryError(c, err)
		return
	}

	response.OK(c, "latest electrical telemetry retrieved successfully", out)
}

// ── Helpers ───────────────────────────────────────────────────────

func (h *TelemetryHandler) handleTelemetryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrDeviceNotFound):
		response.NotFound(c, "device not found")
	case errors.Is(err, service.ErrTelemetryNotFound):
		response.NotFound(c, "telemetry data not found")
	default:
		h.log.Error("unexpected telemetry error", zap.Error(err))
		response.InternalServerError(c, "an unexpected error occurred")
	}
}
