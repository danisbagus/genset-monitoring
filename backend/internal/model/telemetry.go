package model

import "github.com/google/uuid"

// Telemetry represents a single telemetry reading from a genset device.
type Telemetry struct {
	Base
	DeviceID    uuid.UUID `gorm:"type:uuid;not null;index"    json:"device_id"`
	Device      *Device   `gorm:"foreignKey:DeviceID"          json:"device,omitempty"`
	Voltage     float64   `json:"voltage"`     // Volts
	Current     float64   `json:"current"`     // Amperes
	Frequency   float64   `json:"frequency"`   // Hz
	Power       float64   `json:"power"`       // kW
	Temperature float64   `json:"temperature"` // Celsius
	FuelLevel   float64   `json:"fuel_level"`  // Percentage
	RPM         int       `json:"rpm"`
	Status      string    `gorm:"default:'unknown'" json:"status"`
}
