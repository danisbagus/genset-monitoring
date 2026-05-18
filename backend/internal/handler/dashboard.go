package handler

import (
	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	"github.com/danisbagus/genset-monitoring/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DashboardHandler menangani endpoint HTTP terkait dashboard.
type DashboardHandler struct {
	dashboardSvc service.DashboardService
	log          *zap.Logger
}

// NewDashboardHandler membuat instance baru dari DashboardHandler.
func NewDashboardHandler(dashboardSvc service.DashboardService, log *zap.Logger) *DashboardHandler {
	return &DashboardHandler{
		dashboardSvc: dashboardSvc,
		log:          log,
	}
}

// GetSummary get dashboard summary
//
//	@Summary		Dashboard Summary
//	@Description	Gets a summary of statistical data for the main dashboard (total devices, connection status, engines running, and active alerts)
//	@Tags			dashboard
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200		{object}	response.Response{data=service.DashboardSummaryOutput}	"Dashboard Summary berhasil diambil"
//	@Failure		401		{object}	response.ErrorResponse								"Unauthorized"
//	@Failure		500		{object}	response.ErrorResponse								"Internal server error"
//	@Router			/api/v1/dashboard/summary [get]
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	out, err := h.dashboardSvc.GetSummary(c.Request.Context())
	if err != nil {
		h.log.Error("failed to get dashboard summary", zap.Error(err))
		response.InternalServerError(c, "failed to get dashboard summary")
		return
	}

	response.OK(c, "dashboard summary retrieved successfully", out)
}

// GetDeviceStates get dashboard device states
//
//	@Summary		Dashboard Device States
//	@Description	Gets a paginated list of device states for the main dashboard
//	@Tags			dashboard
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int	false	"Page number (default: 1)"
//	@Param			limit	query		int	false	"Limit number (default: 5)"
//	@Success		200		{object}	response.Response{data=service.DashboardDeviceStatesOutput}	"Dashboard Device States berhasil diambil"
//	@Failure		401		{object}	response.ErrorResponse									"Unauthorized"
//	@Failure		500		{object}	response.ErrorResponse									"Internal server error"
//	@Router			/api/v1/dashboard/device-states [get]
func (h *DashboardHandler) GetDeviceStates(c *gin.Context) {
	var query service.GetDeviceStatesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.log.Warn("invalid query params", zap.Error(err))
	}

	out, err := h.dashboardSvc.GetDeviceStates(c.Request.Context(), query)
	if err != nil {
		h.log.Error("failed to get dashboard device states", zap.Error(err))
		response.InternalServerError(c, "failed to get dashboard device states")
		return
	}

	response.OK(c, "dashboard device states retrieved successfully", out)
}

// GetRecentAlerts get dashboard recent alerts
//
//	@Summary		Dashboard Recent Alerts
//	@Description	Gets a paginated list of recent alerts for the main dashboard
//	@Tags			dashboard
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int	false	"Page number (default: 1)"
//	@Param			limit	query		int	false	"Limit number (default: 5)"
//	@Success		200		{object}	service.RecentAlertsOutput	"Dashboard Recent Alerts berhasil diambil"
//	@Failure		401		{object}	response.ErrorResponse		"Unauthorized"
//	@Failure		500		{object}	response.ErrorResponse		"Internal server error"
//	@Router			/api/v1/dashboard/recent-alerts [get]
func (h *DashboardHandler) GetRecentAlerts(c *gin.Context) {
	var query service.GetRecentAlertsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.log.Warn("invalid query params", zap.Error(err))
	}

	out, err := h.dashboardSvc.GetRecentAlerts(c.Request.Context(), query)
	if err != nil {
		h.log.Error("failed to get dashboard recent alerts", zap.Error(err))
		response.InternalServerError(c, "failed to get dashboard recent alerts")
		return
	}

	response.OK(c, "dashboard recent alerts retrieved successfully", out)
}
