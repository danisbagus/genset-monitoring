package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	"github.com/danisbagus/genset-monitoring/backend/pkg/response"
	"github.com/danisbagus/genset-monitoring/backend/pkg/validator"
)

// ── Request DTOs ─────────────────────────────────────────────────

// HeartbeatRequest is the payload from the physical device.
type HeartbeatRequest struct {
	GSMSignal       int16     `json:"gsm_signal"       validate:"required,min=-120,max=0"`
	GPSConnected    bool      `json:"gps_connected"    validate:"omitempty"`
	ServerConnected bool      `json:"server_connected" validate:"omitempty"`
	CANConnected    bool      `json:"can_connected"    validate:"omitempty"`
	RS485Connected  bool      `json:"rs485_connected"  validate:"omitempty"`
	SDCardOK        bool      `json:"sd_card_ok"       validate:"omitempty"`
	Timestamp       time.Time `json:"timestamp"        validate:"omitempty"`
}

// ── Handler ───────────────────────────────────────────────────────

type DeviceStatusHandler struct {
	statusSvc service.DeviceStatusService
	validator *validator.Validator
	log       *zap.Logger
}

func NewDeviceStatusHandler(statusSvc service.DeviceStatusService, v *validator.Validator, log *zap.Logger) *DeviceStatusHandler {
	return &DeviceStatusHandler{
		statusSvc: statusSvc,
		validator: v,
		log:       log,
	}
}

// ── Endpoints ─────────────────────────────────────────────────────

// GetStatus returns the current health and connectivity of a device.
//
//	@Summary		Get Device Status
//	@Description	Returns real-time hardware health and connectivity status
//	@Tags			status
//	@Produce		json
//	@Security		BearerAuth
//	@Param			deviceID	path		string	true	"Device UUID"
//	@Success		200			{object}	response.Response{data=service.DeviceStatusOutput}	"Status retrieved"
//	@Failure		404			{object}	response.ErrorResponse								"Status not found"
//	@Failure		500			{object}	response.ErrorResponse								"Internal error"
//	@Router			/api/v1/devices/{deviceID}/status [get]
func (h *DeviceStatusHandler) GetStatus(c *gin.Context) {
	deviceID, err := uuid.Parse(c.Param("deviceID"))
	if err != nil {
		response.BadRequest(c, "invalid device ID", nil)
		return
	}

	out, err := h.statusSvc.GetStatus(c.Request.Context(), deviceID)
	if err != nil {
		h.handleStatusError(c, err)
		return
	}

	response.OK(c, "status retrieved successfully", out)
}

// Heartbeat records a new health check from the physical device.
//
//	@Summary		Device Heartbeat
//	@Description	Receives health check from device and updates online status
//	@Tags			status
//	@Accept			json
//	@Produce		json
//	@Param			deviceID	path		string				true	"Device UUID"
//	@Param			body		body		HeartbeatRequest	true	"Heartbeat payload"
//	@Success		200			{object}	response.Response	"Heartbeat recorded"
//	@Failure		400			{object}	response.ErrorResponse	"Invalid payload"
//	@Failure		404			{object}	response.ErrorResponse	"Device not found"
//	@Failure		500			{object}	response.ErrorResponse	"Internal error"
//	@Router			/api/v1/devices/{deviceID}/heartbeat [post]
func (h *DeviceStatusHandler) Heartbeat(c *gin.Context) {
	deviceID, err := uuid.Parse(c.Param("deviceID"))
	if err != nil {
		response.BadRequest(c, "invalid device ID", nil)
		return
	}

	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	err = h.statusSvc.Heartbeat(c.Request.Context(), service.HeartbeatInput{
		DeviceID:        deviceID,
		GSMSignal:       req.GSMSignal,
		GPSConnected:    req.GPSConnected,
		ServerConnected: req.ServerConnected,
		CANConnected:    req.CANConnected,
		RS485Connected:  req.RS485Connected,
		SDCardOK:        req.SDCardOK,
		Timestamp:       req.Timestamp,
	})

	if err != nil {
		h.handleStatusError(c, err)
		return
	}

	response.OK(c, "heartbeat recorded", nil)
}

// ── Error mapping ─────────────────────────────────────────────────

func (h *DeviceStatusHandler) handleStatusError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrStatusNotFound),
		errors.Is(err, service.ErrDeviceNotFound):
		response.NotFound(c, err.Error())
	default:
		h.log.Error("unexpected status error", zap.Error(err))
		response.InternalServerError(c, "an unexpected error occurred")
	}
}
