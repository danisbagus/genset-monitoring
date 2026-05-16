package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	"github.com/danisbagus/genset-monitoring/backend/pkg/response"
	"github.com/danisbagus/genset-monitoring/backend/pkg/validator"
)

// ── Request DTOs ─────────────────────────────────────────────────

// CreateAlertRequest is the expected JSON body for POST /api/v1/alerts.
type CreateAlertRequest struct {
	DeviceID       string   `json:"device_id"       validate:"required,uuid"`
	Type           string   `json:"type"            validate:"required,oneof=engine electrical connectivity"`
	Severity       string   `json:"severity"        validate:"required,oneof=critical warning info"`
	Title          string   `json:"title"           validate:"required"`
	Message        string   `json:"message"         validate:"required"`
	MetricName     *string  `json:"metric_name"     validate:"omitempty"`
	MetricValue    *float64 `json:"metric_value"    validate:"omitempty"`
	ThresholdValue *float64 `json:"threshold_value" validate:"omitempty"`
}

// ── Handler ───────────────────────────────────────────────────────

// AlertHandler handles all alert-management HTTP endpoints.
type AlertHandler struct {
	alertSvc  service.AlertService
	validator *validator.Validator
	log       *zap.Logger
}

// NewAlertHandler creates a new AlertHandler.
func NewAlertHandler(alertSvc service.AlertService, v *validator.Validator, log *zap.Logger) *AlertHandler {
	return &AlertHandler{
		alertSvc:  alertSvc,
		validator: v,
		log:       log,
	}
}

// Create inserts a new alert record.
//
//	@Summary		Create Alert
//	@Description	Inserts a new alert record into the alerts table.
//	@Tags			alerts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		CreateAlertRequest										true	"Create alert payload"
//	@Success		201		{object}	response.Response{data=service.AlertCreatedOutput}		"Alert created"
//	@Failure		400		{object}	response.ErrorResponse									"Invalid request body or validation error"
//	@Failure		401		{object}	response.ErrorResponse									"Unauthorized"
//	@Failure		404		{object}	response.ErrorResponse									"Device not found"
//	@Failure		500		{object}	response.ErrorResponse									"Internal server error"
//	@Router			/api/v1/alerts [post]
func (h *AlertHandler) Create(c *gin.Context) {
	var req CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		response.BadRequest(c, "validation failed", errs)
		return
	}

	deviceID, _ := uuid.Parse(req.DeviceID)

	out, err := h.alertSvc.Create(c.Request.Context(), service.CreateAlertInput{
		DeviceID:       deviceID,
		Type:           model.AlertType(req.Type),
		Severity:       model.AlertSeverity(req.Severity),
		Title:          req.Title,
		Message:        req.Message,
		MetricName:     req.MetricName,
		MetricValue:    req.MetricValue,
		ThresholdValue: req.ThresholdValue,
	})
	if err != nil {
		h.handleAlertError(c, err)
		return
	}

	response.Created(c, "alert created successfully", out)
}

// handleAlertError maps service sentinel errors to HTTP status codes.
func (h *AlertHandler) handleAlertError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrDeviceNotFound):
		response.NotFound(c, err.Error())
	default:
		h.log.Error("unexpected alert error", zap.Error(err))
		response.InternalServerError(c, "an unexpected error occurred")
	}
}
