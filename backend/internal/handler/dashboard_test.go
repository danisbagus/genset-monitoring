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
	getSummaryFunc      func(ctx context.Context) (*service.DashboardSummaryOutput, error)
	getDeviceStatesFunc func(ctx context.Context, query service.GetDeviceStatesQuery) (*service.DashboardDeviceStatesOutput, error)
	getRecentAlertsFunc func(ctx context.Context, query service.GetRecentAlertsQuery) (*service.RecentAlertsOutput, error)
}

func (m *mockDashboardService) GetSummary(ctx context.Context) (*service.DashboardSummaryOutput, error) {
	return m.getSummaryFunc(ctx)
}

func (m *mockDashboardService) GetDeviceStates(ctx context.Context, query service.GetDeviceStatesQuery) (*service.DashboardDeviceStatesOutput, error) {
	return m.getDeviceStatesFunc(ctx, query)
}

func (m *mockDashboardService) GetRecentAlerts(ctx context.Context, query service.GetRecentAlertsQuery) (*service.RecentAlertsOutput, error) {
	return m.getRecentAlertsFunc(ctx, query)
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

func TestDashboardHandler_GetDeviceStates_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockDashboardService{
		getDeviceStatesFunc: func(ctx context.Context, query service.GetDeviceStatesQuery) (*service.DashboardDeviceStatesOutput, error) {
			return &service.DashboardDeviceStatesOutput{
				Devices: []service.DashboardDeviceStateOutput{
					{
						DeviceID:     "device-1",
						DeviceName:   "Genset 1",
						DeviceOnline: true,
					},
				},
				Pagination: service.PaginationMeta{
					Page:  1,
					Limit: 5,
					Total: 1,
				},
			}, nil
		},
	}

	h := NewDashboardHandler(mockSvc, zap.NewNop())

	r := gin.Default()
	r.GET("/api/v1/dashboard/device-states", h.GetDeviceStates)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/device-states?page=1&limit=5", nil)
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

	devices := data["devices"].([]interface{})
	if len(devices) != 1 {
		t.Errorf("expected 1 device, got %d", len(devices))
	}

	pagination := data["pagination"].(map[string]interface{})
	if int64(pagination["total"].(float64)) != 1 {
		t.Errorf("expected total 1, got %v", pagination["total"])
	}
}

func TestDashboardHandler_GetRecentAlerts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockDashboardService{
		getRecentAlertsFunc: func(ctx context.Context, query service.GetRecentAlertsQuery) (*service.RecentAlertsOutput, error) {
			return &service.RecentAlertsOutput{
				Alerts: []service.RecentAlertOutput{
					{
						AlertID:    "alert-1",
						DeviceName: "Genset 1",
					},
				},
				Pagination: service.PaginationMeta{
					Page:  1,
					Limit: 5,
					Total: 1,
				},
			}, nil
		},
	}

	h := NewDashboardHandler(mockSvc, zap.NewNop())

	r := gin.Default()
	r.GET("/api/v1/dashboard/recent-alerts", h.GetRecentAlerts)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/recent-alerts?page=1&limit=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var res service.RecentAlertsOutput
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(res.Alerts) != 1 {
		t.Errorf("expected 1 alert, got %d", len(res.Alerts))
	}

	if res.Alerts[0].DeviceName != "Genset 1" {
		t.Errorf("expected device name Genset 1, got %s", res.Alerts[0].DeviceName)
	}

	if res.Pagination.Total != 1 {
		t.Errorf("expected total 1, got %d", res.Pagination.Total)
	}
}
