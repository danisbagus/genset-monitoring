package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
)

var (
	ErrTelemetryNotFound = errors.New("telemetry not found")
)

// TelemetryRepository defines the contract for telemetry data persistence.
type TelemetryRepository interface {
	CreateEngine(ctx context.Context, telemetry *model.EngineTelemetry) error
	GetLatestEngine(ctx context.Context, deviceID uuid.UUID) (*model.EngineTelemetry, error)
	CreateElectrical(ctx context.Context, telemetry *model.ElectricalTelemetry) error
	GetLatestElectrical(ctx context.Context, deviceID uuid.UUID) (*model.ElectricalTelemetry, error)
}

type telemetryRepository struct {
	db *gorm.DB
}

// NewTelemetryRepository constructs a TelemetryRepository.
func NewTelemetryRepository(db *gorm.DB) TelemetryRepository {
	return &telemetryRepository{db: db}
}

// CreateEngine saves a new engine telemetry record.
func (r *telemetryRepository) CreateEngine(ctx context.Context, telemetry *model.EngineTelemetry) error {
	if err := r.db.WithContext(ctx).Create(telemetry).Error; err != nil {
		return fmt.Errorf("telemetryRepository.CreateEngine: %w", err)
	}
	return nil
}

// GetLatestEngine retrieves the most recent engine telemetry for a device.
func (r *telemetryRepository) GetLatestEngine(ctx context.Context, deviceID uuid.UUID) (*model.EngineTelemetry, error) {
	var telemetry model.EngineTelemetry
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("created_at DESC").
		First(&telemetry).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTelemetryNotFound
		}
		return nil, fmt.Errorf("telemetryRepository.GetLatestEngine: %w", err)
	}

	return &telemetry, nil
}

// CreateElectrical saves a new electrical telemetry record.
func (r *telemetryRepository) CreateElectrical(ctx context.Context, telemetry *model.ElectricalTelemetry) error {
	if err := r.db.WithContext(ctx).Create(telemetry).Error; err != nil {
		return fmt.Errorf("telemetryRepository.CreateElectrical: %w", err)
	}
	return nil
}

// GetLatestElectrical retrieves the most recent electrical telemetry for a device.
func (r *telemetryRepository) GetLatestElectrical(ctx context.Context, deviceID uuid.UUID) (*model.ElectricalTelemetry, error) {
	var telemetry model.ElectricalTelemetry
	err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("created_at DESC").
		First(&telemetry).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTelemetryNotFound
		}
		return nil, fmt.Errorf("telemetryRepository.GetLatestElectrical: %w", err)
	}

	return &telemetry, nil
}
