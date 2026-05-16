package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AlertSeverity defines alert severity levels.
type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "info"
	SeverityWarning  AlertSeverity = "warning"
	SeverityCritical AlertSeverity = "critical"
)

// AlertType defines alert types.
type AlertType string

const (
	TypeEngine       AlertType = "engine"
	TypeElectrical   AlertType = "electrical"
	TypeConnectivity AlertType = "connectivity"
)

// AlertStatus defines alert statuses.
type AlertStatus string

const (
	StatusActive       AlertStatus = "active"
	StatusAcknowledged AlertStatus = "acknowledged"
	StatusResolved     AlertStatus = "resolved"
)

// Alert represents a genset alert event.
type Alert struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt      time.Time     `json:"created_at"`
	DeviceID       uuid.UUID     `gorm:"type:uuid;not null;index" json:"device_id"`
	Device         *Device       `gorm:"foreignKey:DeviceID"       json:"device,omitempty"`
	Type           AlertType     `gorm:"not null"                  json:"type"`
	Severity       AlertSeverity `gorm:"not null"                  json:"severity"`
	Title          string        `gorm:"not null"                  json:"title"`
	Message        string        `gorm:"not null"                  json:"message"`
	MetricName     *string       `json:"metric_name,omitempty"`
	MetricValue    *float64      `json:"metric_value,omitempty"`
	ThresholdValue *float64      `json:"threshold_value,omitempty"`
	Status         AlertStatus   `gorm:"not null;default:'active'" json:"status"`
	AcknowledgedAt *time.Time    `json:"acknowledged_at,omitempty"`
	AcknowledgedBy *uuid.UUID    `gorm:"type:uuid"                 json:"acknowledged_by,omitempty"`
	ResolvedAt     *time.Time    `json:"resolved_at,omitempty"`
}

// BeforeCreate sets a UUID if not already set.
func (a *Alert) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
