package repository

import (
	"context"
	"fmt"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"gorm.io/gorm"
)

// AlertRepository defines the persistence contract for alerts.
type AlertRepository interface {
	Create(ctx context.Context, alert *model.Alert) error
}

type alertRepository struct {
	db *gorm.DB
}

// NewAlertRepository constructs an AlertRepository backed by GORM.
func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &alertRepository{db: db}
}

// Create inserts a new alert record.
func (r *alertRepository) Create(ctx context.Context, alert *model.Alert) error {
	if err := r.db.WithContext(ctx).Create(alert).Error; err != nil {
		return fmt.Errorf("alertRepository.Create: %w", err)
	}
	return nil
}
