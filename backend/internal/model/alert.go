package model

import (
	"time"

	"github.com/google/uuid"
)

// AlertSeverity defines alert severity levels.
type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "info"
	SeverityWarning  AlertSeverity = "warning"
	SeverityCritical AlertSeverity = "critical"
)

// Alert represents a genset alert event.
type Alert struct {
	Base
	DeviceID    uuid.UUID     `gorm:"type:uuid;not null;index" json:"device_id"`
	Device      *Device       `gorm:"foreignKey:DeviceID"       json:"device,omitempty"`
	Severity    AlertSeverity `gorm:"not null"                  json:"severity"`
	Message     string        `gorm:"not null"                  json:"message"`
	ResolvedAt  *time.Time    `json:"resolved_at,omitempty"`
	IsResolved  bool          `gorm:"default:false"             json:"is_resolved"`
}
