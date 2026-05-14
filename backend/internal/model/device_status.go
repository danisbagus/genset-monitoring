package model

import (
	"time"

	"github.com/google/uuid"
)

// DeviceStatus represents the real-time health and connectivity state of a device.
type DeviceStatus struct {
	DeviceID        uuid.UUID `gorm:"primaryKey;type:uuid" json:"device_id"`
	GSMSignal       int16     `gorm:"type:smallint" json:"gsm_signal"`
	GPSConnected    bool      `gorm:"not null;default:false" json:"gps_connected"`
	ServerConnected bool      `gorm:"not null;default:false" json:"server_connected"`
	CANConnected    bool      `gorm:"not null;default:false" json:"can_connected"`
	RS485Connected  bool      `gorm:"not null;default:false" json:"rs485_connected"`
	SDCardOK        bool      `gorm:"not null;default:false" json:"sd_card_ok"`
	LastSeen        time.Time `gorm:"not null;default:now()" json:"last_seen"`
	CreatedAt       time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"not null;default:now()" json:"updated_at"`

	// Associations
	Device *Device `gorm:"foreignKey:DeviceID" json:"-"`
}

// TableName overrides the GORM table name convention.
func (DeviceStatus) TableName() string {
	return "device_statuses"
}

// DeviceLatestState represents a realtime snapshot of the latest device telemetry.
type DeviceLatestState struct {
	DeviceID            uuid.UUID  `gorm:"primaryKey;type:uuid" json:"device_id"`
	EngineRunning       bool       `gorm:"not null;default:false" json:"engine_running"`
	Speed               *int32     `json:"speed"`
	CoolantTemperature  *float32   `json:"coolant_temperature"`
	OilPressure         *float32   `json:"oil_pressure"`
	FuelLevel           *float32   `json:"fuel_level"`
	Frequency           *float32   `json:"frequency"`
	TotalVa             *float32   `json:"total_va"`
	PfAvg               *float32   `json:"pf_avg"`
	TelemetryRecordedAt *time.Time `json:"telemetry_recorded_at"`
	LastOnlineAt        *time.Time `json:"last_online_at"`
	LastSeenAt          time.Time  `gorm:"not null;default:now()" json:"last_seen_at"`
	UpdatedAt           time.Time  `gorm:"not null;default:now()" json:"updated_at"`

	// Associations
	Device *Device `gorm:"foreignKey:DeviceID" json:"-"`
}

// TableName overrides the GORM table name convention.
func (DeviceLatestState) TableName() string {
	return "device_latest_state"
}
