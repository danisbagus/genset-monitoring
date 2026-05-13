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
	pkgvalidator "github.com/danisbagus/genset-monitoring/backend/pkg/validator"
)

type mockTelemetryService struct {
	service.TelemetryService
	createEngineFunc      func(ctx context.Context, input service.CreateEngineTelemetryInput) (*service.EngineTelemetryOutput, error)
	getLatestEngineFunc   func(ctx context.Context, deviceID uuid.UUID) (*service.EngineTelemetryOutput, error)
}

func (m *mockTelemetryService) CreateEngine(ctx context.Context, input service.CreateEngineTelemetryInput) (*service.EngineTelemetryOutput, error) {
	return m.createEngineFunc(ctx, input)
}

func (m *mockTelemetryService) GetLatestEngine(ctx context.Context, deviceID uuid.UUID) (*service.EngineTelemetryOutput, error) {
	return m.getLatestEngineFunc(ctx, deviceID)
}

func TestTelemetryHandler_CreateEngine_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTelemetryHandler(nil, pkgvalidator.New(), zap.NewNop())
	
	r := gin.Default()
	r.POST("/api/v1/devices/:deviceID/engine", h.CreateEngine)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/devices/invalid-uuid/engine", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestTelemetryHandler_CreateEngine_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTelemetryHandler(nil, pkgvalidator.New(), zap.NewNop())
	
	r := gin.Default()
	r.POST("/api/v1/devices/:deviceID/engine", h.CreateEngine)

	deviceID := uuid.New().String()
	body := []byte(`{ "speed": "invalid" }`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/devices/"+deviceID+"/engine", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestTelemetryHandler_CreateEngine_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	deviceID := uuid.New()
	
	mockSvc := &mockTelemetryService{
		createEngineFunc: func(ctx context.Context, input service.CreateEngineTelemetryInput) (*service.EngineTelemetryOutput, error) {
			return &service.EngineTelemetryOutput{DeviceID: input.DeviceID}, nil
		},
	}
	
	h := NewTelemetryHandler(mockSvc, pkgvalidator.New(), zap.NewNop())
	
	r := gin.Default()
	r.POST("/api/v1/devices/:deviceID/engine", h.CreateEngine)

	body, _ := json.Marshal(CreateEngineTelemetryRequest{
		Speed: int32Ptr(1500),
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/devices/"+deviceID.String()+"/engine", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func int32Ptr(i int32) *int32 { return &i }
