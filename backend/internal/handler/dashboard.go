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
