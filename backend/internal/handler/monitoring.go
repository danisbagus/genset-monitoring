package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	"github.com/danisbagus/genset-monitoring/backend/pkg/response"
)

// MonitoringHandler handles HTTP endpoints for the realtime monitoring module.
type MonitoringHandler struct {
	monitoringSvc service.MonitoringService
	log           *zap.Logger
}

// NewMonitoringHandler constructs a MonitoringHandler.
func NewMonitoringHandler(monitoringSvc service.MonitoringService, log *zap.Logger) *MonitoringHandler {
	return &MonitoringHandler{
		monitoringSvc: monitoringSvc,
		log:           log,
	}
}

// ListDevices returns a paginated monitoring dashboard device list.
//
//	@Summary		Monitoring Device List
//	@Description	Returns a paginated list of devices with realtime monitoring summary data
//	@Tags			monitoring
//	@Produce		json
//	@Security		BearerAuth
//	@Param			search			query		string	false	"Search by device name, device code, or serial number"
//	@Param			engine_running	query		bool	false	"Filter by engine running status: true | false"
//	@Param			status			query		string	false	"Filter by device status: active | inactive | maintenance"
//	@Param			online			query		bool	false	"Filter by online status (last_seen_at >= now()-5min): true | false"
//	@Param			limit			query		int		false	"Items per page (default: 20, max: 100)"
//	@Param			offset			query		int		false	"Offset (default: 0)"
//	@Param			sort_by			query		string	false	"Sort column: updated_at | last_seen_at | telemetry_recorded_at | name"
//	@Param			sort_order		query		string	false	"Sort direction: asc | desc (default: desc)"
//	@Success		200				{object}	response.Response{data=service.MonitoringDeviceListOutput}	"Monitoring device list"
//	@Failure		400				{object}	response.ErrorResponse											"Bad request"
//	@Failure		401				{object}	response.ErrorResponse											"Unauthorized"
//	@Failure		500				{object}	response.ErrorResponse											"Internal server error"
//	@Router			/api/v1/monitoring/devices [get]
func (h *MonitoringHandler) ListDevices(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	query := service.MonitoringDeviceListQuery{
		Search:        c.Query("search"),
		EngineRunning: c.Query("engine_running"),
		Status:        c.Query("status"),
		Online:        c.Query("online"),
		Limit:         limit,
		Offset:        offset,
		SortBy:        c.Query("sort_by"),
		SortOrder:     c.Query("sort_order"),
	}

	out, err := h.monitoringSvc.ListDevices(c.Request.Context(), query)
	if err != nil {
		h.log.Error("failed to list monitoring devices", zap.Error(err))
		response.InternalServerError(c, "failed to list monitoring devices")
		return
	}

	response.OK(c, "monitoring devices retrieved successfully", out)
}

// GetDeviceDetail returns the full realtime monitoring detail for a single device.
//
//	@Summary		Monitoring Device Detail
//	@Description	Returns full realtime monitoring detail for a device (device info, latest state, connectivity, engine & electrical telemetry)
//	@Tags			monitoring
//	@Produce		json
//	@Security		BearerAuth
//	@Param			deviceID	path		string	true	"Device UUID"
//	@Success		200			{object}	response.Response{data=service.MonitoringDeviceDetailOutput}	"Monitoring device detail"
//	@Failure		400			{object}	response.ErrorResponse											"Invalid UUID"
//	@Failure		401			{object}	response.ErrorResponse											"Unauthorized"
//	@Failure		404			{object}	response.ErrorResponse											"Device not found"
//	@Failure		500			{object}	response.ErrorResponse											"Internal server error"
//	@Router			/api/v1/monitoring/devices/{deviceID} [get]
func (h *MonitoringHandler) GetDeviceDetail(c *gin.Context) {
	raw := c.Param("deviceID")
	deviceID, err := uuid.Parse(raw)
	if err != nil {
		response.BadRequest(c, "invalid deviceID, must be a valid UUID", nil)
		return
	}

	out, err := h.monitoringSvc.GetDeviceDetail(c.Request.Context(), deviceID)
	if err != nil {
		if errors.Is(err, service.ErrMonitoringDeviceNotFound) {
			response.NotFound(c, "device not found")
			return
		}
		h.log.Error("failed to get monitoring device detail",
			zap.String("deviceID", deviceID.String()),
			zap.Error(err),
		)
		response.InternalServerError(c, "failed to get monitoring device detail")
		return
	}

	response.OK(c, "monitoring device detail retrieved successfully", out)
}
