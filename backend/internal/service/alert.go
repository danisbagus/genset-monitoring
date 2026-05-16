package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ── DTOs ─────────────────────────────────────────────────────────

// CreateAlertInput is the payload required to create an alert.
type CreateAlertInput struct {
	DeviceID       uuid.UUID
	Type           model.AlertType
	Severity       model.AlertSeverity
	Title          string
	Message        string
	MetricName     *string
	MetricValue    *float64
	ThresholdValue *float64
}

// AlertCreatedOutput is returned after a successful alert creation.
type AlertCreatedOutput struct {
	AlertID uuid.UUID `json:"alert_id"`
}

// ── Interface ─────────────────────────────────────────────────────

// AlertService defines the contract for alert management operations.
type AlertService interface {
	Create(ctx context.Context, input CreateAlertInput) (*AlertCreatedOutput, error)
}

// ── Implementation ────────────────────────────────────────────────

type alertService struct {
	alertRepo  repository.AlertRepository
	deviceRepo repository.DeviceRepository
	log        *zap.Logger
}

// NewAlertService constructs an AlertService.
func NewAlertService(alertRepo repository.AlertRepository, deviceRepo repository.DeviceRepository, log *zap.Logger) AlertService {
	return &alertService{
		alertRepo:  alertRepo,
		deviceRepo: deviceRepo,
		log:        log,
	}
}

// Create inserts a new alert record.
func (s *alertService) Create(ctx context.Context, input CreateAlertInput) (*AlertCreatedOutput, error) {
	// Ensure device exists and is not deleted
	if _, err := s.deviceRepo.FindByID(ctx, input.DeviceID); err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("alertService.Create: check device: %w", err)
	}

	alert := &model.Alert{
		DeviceID:       input.DeviceID,
		Type:           input.Type,
		Severity:       input.Severity,
		Title:          input.Title,
		Message:        input.Message,
		MetricName:     input.MetricName,
		MetricValue:    input.MetricValue,
		ThresholdValue: input.ThresholdValue,
		Status:         model.StatusActive,
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		return nil, fmt.Errorf("alertService.Create: %w", err)
	}

	s.log.Info("alert created",
		zap.String("alert_id", alert.ID.String()),
		zap.String("device_id", alert.DeviceID.String()),
	)

	return &AlertCreatedOutput{
		AlertID: alert.ID,
	}, nil
}

// ── Sentinel errors ───────────────────────────────────────────────

var (
	ErrAlertNotFound = errors.New("alert not found")
)
