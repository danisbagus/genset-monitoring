package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	"github.com/danisbagus/genset-monitoring/backend/pkg/response"
	"github.com/danisbagus/genset-monitoring/backend/pkg/validator"
)

// ── Request DTOs ─────────────────────────────────────────────────

// CreateDeviceRequest is the expected JSON body for POST /api/v1/devices.
type CreateDeviceRequest struct {
	DeviceCode   string `json:"device_code"   validate:"required,min=1,max=100"`
	Name         string `json:"name"          validate:"required,min=1,max=100"`
	SerialNumber string `json:"serial_number" validate:"required,min=1,max=100"`
	EngineID     string `json:"engine_id"     validate:"omitempty,max=100"`
	GSMNumber    string `json:"gsm_number"    validate:"omitempty,max=50"`
}

// UpdateDeviceRequest is the expected JSON body for PATCH /api/v1/devices/:deviceID.
// All fields are optional; only provided (non-null) fields are applied.
type UpdateDeviceRequest struct {
	Name            *string `json:"name"             validate:"omitempty,min=1,max=100"`
	EngineID        *string `json:"engine_id"        validate:"omitempty,max=100"`
	GSMNumber       *string `json:"gsm_number"       validate:"omitempty,max=50"`
	FirmwareVersion *string `json:"firmware_version" validate:"omitempty,max=50"`
	Status          *string `json:"status"           validate:"omitempty,oneof=active inactive maintenance"`
}

// ── Handler ───────────────────────────────────────────────────────

// DeviceHandler handles all device-management HTTP endpoints.
type DeviceHandler struct {
	deviceSvc service.DeviceService
	validator *validator.Validator
	log       *zap.Logger
}

// NewDeviceHandler creates a new DeviceHandler.
func NewDeviceHandler(deviceSvc service.DeviceService, v *validator.Validator, log *zap.Logger) *DeviceHandler {
	return &DeviceHandler{
		deviceSvc: deviceSvc,
		validator: v,
		log:       log,
	}
}

// ── Endpoints ─────────────────────────────────────────────────────

// List returns a paginated, filtered list of devices.
//
//	@Summary		List Devices
//	@Description	Returns a paginated list of devices with optional search, status, and online/offline filter
//	@Tags			devices
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int		false	"Page number (default: 1)"
//	@Param			limit	query		int		false	"Items per page (default: 10, max: 100)"
//	@Param			search	query		string	false	"Search by device_code or name"
//	@Param			status	query		string	false	"Filter by status: active | inactive | maintenance"
//	@Param			online	query		string	false	"Filter by connectivity: true | false"
//	@Param			sort_by	query		string	false	"Sort column: created_at | name | device_code | last_seen | firmware_version"
//	@Param			sort_dir	query		string	false	"Sort direction: asc | desc (default: desc)"
//	@Success		200		{object}	response.Response{data=service.DeviceListOutput}	"Device list"
//	@Failure		401		{object}	response.ErrorResponse								"Unauthorized"
//	@Failure		500		{object}	response.ErrorResponse								"Internal server error"
//	@Router			/api/v1/devices [get]
func (h *DeviceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Parse optional boolean "online" query param
	var online *bool
	if raw := c.Query("online"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			response.BadRequest(c, "invalid value for 'online' parameter, use true or false", nil)
			return
		}
		online = &v
	}

	out, err := h.deviceSvc.List(c.Request.Context(), service.DeviceListInput{
		Search:  c.Query("search"),
		Status:  c.Query("status"),
		Online:  online,
		Page:    page,
		Limit:   limit,
		SortBy:  c.Query("sort_by"),
		SortDir: c.Query("sort_dir"),
	})
	if err != nil {
		h.handleDeviceError(c, err)
		return
	}

	response.OK(c, "devices retrieved successfully", out)
}

// GetByID returns the full detail of a single device.
//
//	@Summary		Get Device Detail
//	@Description	Returns full detail of a device by its UUID
//	@Tags			devices
//	@Produce		json
//	@Security		BearerAuth
//	@Param			deviceID	path		string	true	"Device UUID"
//	@Success		200			{object}	response.Response{data=service.DeviceDetail}	"Device detail"
//	@Failure		400			{object}	response.ErrorResponse							"Invalid UUID"
//	@Failure		401			{object}	response.ErrorResponse							"Unauthorized"
//	@Failure		404			{object}	response.ErrorResponse							"Device not found"
//	@Failure		500			{object}	response.ErrorResponse							"Internal server error"
//	@Router			/api/v1/devices/{deviceID} [get]
func (h *DeviceHandler) GetByID(c *gin.Context) {
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}

	out, err := h.deviceSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleDeviceError(c, err)
		return
	}

	response.OK(c, "device retrieved successfully", out)
}

// Create registers a new device.
//
//	@Summary		Create Device
//	@Description	Registers a new device; returns its UUID and a one-time device token
//	@Tags			devices
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		CreateDeviceRequest										true	"Create device payload"
//	@Success		201		{object}	response.Response{data=service.DeviceCreatedOutput}		"Device created"
//	@Failure		400		{object}	response.ErrorResponse									"Invalid request body"
//	@Failure		401		{object}	response.ErrorResponse									"Unauthorized"
//	@Failure		409		{object}	response.ErrorResponse									"device_code or serial_number already in use"
//	@Failure		422		{object}	response.ErrorResponse									"Validation error"
//	@Failure		500		{object}	response.ErrorResponse									"Internal server error"
//	@Router			/api/v1/devices [post]
func (h *DeviceHandler) Create(c *gin.Context) {
	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	out, err := h.deviceSvc.Create(c.Request.Context(), service.CreateDeviceInput{
		DeviceCode:   req.DeviceCode,
		Name:         req.Name,
		SerialNumber: req.SerialNumber,
		EngineID:     req.EngineID,
		GSMNumber:    req.GSMNumber,
	})
	if err != nil {
		h.handleDeviceError(c, err)
		return
	}

	response.Created(c, "device created successfully", out)
}

// Update partially updates a device's mutable fields.
//
//	@Summary		Update Device
//	@Description	Partially updates device fields (name, firmware_version, status)
//	@Tags			devices
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			deviceID	path		string				true	"Device UUID"
//	@Param			body		body		UpdateDeviceRequest	true	"Update device payload"
//	@Success		200			{object}	response.Response	"Device updated"
//	@Failure		400			{object}	response.ErrorResponse	"Invalid UUID or request body"
//	@Failure		401			{object}	response.ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	response.ErrorResponse	"Device not found"
//	@Failure		422			{object}	response.ErrorResponse	"Validation error"
//	@Failure		500			{object}	response.ErrorResponse	"Internal server error"
//	@Router			/api/v1/devices/{deviceID} [patch]
func (h *DeviceHandler) Update(c *gin.Context) {
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}

	var req UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	input := service.UpdateDeviceInput{
		Name:            req.Name,
		EngineID:        req.EngineID,
		GSMNumber:       req.GSMNumber,
		FirmwareVersion: req.FirmwareVersion,
	}
	if req.Status != nil {
		s := model.DeviceLifecycle(*req.Status)
		input.Status = &s
	}

	if err := h.deviceSvc.Update(c.Request.Context(), id, input); err != nil {
		h.handleDeviceError(c, err)
		return
	}

	response.OK(c, "device updated successfully", nil)
}

// Delete soft-deletes a device.
//
//	@Summary		Delete Device
//	@Description	Soft-deletes a device by UUID, preserving the audit trail
//	@Tags			devices
//	@Produce		json
//	@Security		BearerAuth
//	@Param			deviceID	path		string					true	"Device UUID"
//	@Success		200			{object}	response.Response		"Device deleted"
//	@Failure		400			{object}	response.ErrorResponse	"Invalid UUID"
//	@Failure		401			{object}	response.ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	response.ErrorResponse	"Device not found"
//	@Failure		500			{object}	response.ErrorResponse	"Internal server error"
//	@Router			/api/v1/devices/{deviceID} [delete]
func (h *DeviceHandler) Delete(c *gin.Context) {
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}

	if err := h.deviceSvc.Delete(c.Request.Context(), id); err != nil {
		h.handleDeviceError(c, err)
		return
	}

	response.OK(c, "device deleted successfully", nil)
}

// ── Helpers ───────────────────────────────────────────────────────

// parseDeviceID extracts and validates the :deviceID path parameter.
// Returns false and writes the error response when the UUID is invalid.
func parseDeviceID(c *gin.Context) (uuid.UUID, bool) {
	raw := c.Param("deviceID")
	id, err := uuid.Parse(raw)
	if err != nil {
		response.BadRequest(c, "invalid device ID, must be a valid UUID", nil)
		return uuid.Nil, false
	}
	return id, true
}

// handleDeviceError maps service sentinel errors to HTTP status codes.
func (h *DeviceHandler) handleDeviceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrDeviceNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, service.ErrDeviceCodeExists),
		errors.Is(err, service.ErrSerialNumberExists):
		c.JSON(http.StatusConflict, response.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	default:
		h.log.Error("unexpected device error", zap.Error(err))
		response.InternalServerError(c, "an unexpected error occurred")
	}
}
