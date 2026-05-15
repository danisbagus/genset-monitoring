package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	"github.com/danisbagus/genset-monitoring/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type mockDashboardService struct {
	service.DashboardService
	getSummaryFunc func(ctx context.Context) (*service.DashboardSummaryOutput, error)
}

func (m *mockDashboardService) GetSummary(ctx context.Context) (*service.DashboardSummaryOutput, error) {
	return m.getSummaryFunc(ctx)
}

func TestDashboardHandler_GetSummary_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockDashboardService{
		getSummaryFunc: func(ctx context.Context) (*service.DashboardSummaryOutput, error) {
			return &service.DashboardSummaryOutput{
				TotalDevices:   20,
				OnlineDevices:  18,
				OfflineDevices: 2,
				RunningEngines: 12,
				CriticalAlerts: 1,
				WarningAlerts:  4,
			}, nil
		},
	}

	h := NewDashboardHandler(mockSvc, zap.NewNop())

	r := gin.Default()
	r.GET("/api/v1/dashboard/summary", h.GetSummary)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var res response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data, ok := res.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("data response is not map[string]interface{}")
	}

	if int64(data["total_devices"].(float64)) != 20 {
		t.Errorf("expected total_devices 20, got %v", data["total_devices"])
	}
}
