package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	"github.com/danisbagus/genset-monitoring/backend/pkg/validator"
)

type mockAlertService struct {
	service.AlertService
	createFunc func(ctx context.Context, input service.CreateAlertInput) (*service.AlertCreatedOutput, error)
}

func (m *mockAlertService) Create(ctx context.Context, input service.CreateAlertInput) (*service.AlertCreatedOutput, error) {
	return m.createFunc(ctx, input)
}

func TestAlertHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceID := uuid.New().String()
	alertID := uuid.New()

	mockSvc := &mockAlertService{
		createFunc: func(ctx context.Context, input service.CreateAlertInput) (*service.AlertCreatedOutput, error) {
			return &service.AlertCreatedOutput{AlertID: alertID}, nil
		},
	}

	v := validator.New()
	h := NewAlertHandler(mockSvc, v, zap.NewNop())

	r := gin.Default()
	r.POST("/api/v1/alerts", h.Create)

	reqBody := CreateAlertRequest{
		DeviceID: deviceID,
		Type:     "engine",
		Severity: "critical",
		Title:    "High Temp",
		Message:  "Critical engine temperature",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/alerts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)

	data := res["data"].(map[string]interface{})
	if data["alert_id"] != alertID.String() {
		t.Errorf("expected alert_id %s, got %s", alertID.String(), data["alert_id"])
	}
}

func TestAlertHandler_Create_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	v := validator.New()
	h := NewAlertHandler(nil, v, zap.NewNop())

	r := gin.Default()
	r.POST("/api/v1/alerts", h.Create)

	reqBody := CreateAlertRequest{
		DeviceID: "invalid-uuid",
		Type:     "invalid-type",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/alerts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
