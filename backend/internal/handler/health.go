package handler

import (
	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	"github.com/danisbagus/genset-monitoring/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	svc service.HealthService
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(svc service.HealthService) *HealthHandler {
	return &HealthHandler{svc: svc}
}

// Check performs health checks on all infrastructure components.
//
//	@Summary		Health Check
//	@Description	Returns the health status of the service and its infrastructure dependencies
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	response.Response{data=service.HealthStatus}	"Service healthy"
//	@Failure		503	{object}	response.Response{data=service.HealthStatus}	"Service degraded"
//	@Router			/api/v1/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	ctx := c.Request.Context()

	status, err := h.svc.Check(ctx)
	if err != nil {
		response.InternalServerError(c, "health check failed")
		return
	}

	if !status.Healthy() {
		c.JSON(503, response.Response{
			Success: false,
			Message: "service degraded",
			Data:    status,
		})
		return
	}

	response.OK(c, "service healthy", status)
}
