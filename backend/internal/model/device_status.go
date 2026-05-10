package model

import (
	"time"

	"github.com/google/uuid"
)

// DeviceStatus represents the real-time health and connectivity state of a device.
type DeviceStatus struct {
	DeviceID        uuid.UUID `gorm:"primaryKey;type:uuid" json:"device_id"`
	IsOnline        bool      `gorm:"not null;default:false" json:"is_online"`
	LastSeen        time.Time `gorm:"not null;default:now()" json:"last_seen"`
	GSMSignal       int16     `gorm:"type:smallint" json:"gsm_signal"`
	GPSConnected    bool      `gorm:"not null;default:false" json:"gps_connected"`
	ServerConnected bool      `gorm:"not null;default:false" json:"server_connected"`
	CANConnected    bool      `gorm:"not null;default:false" json:"can_connected"`
	RS485Connected  bool      `gorm:"not null;default:false" json:"rs485_connected"`
	SDCardOK        bool      `gorm:"not null;default:false" json:"sd_card_ok"`
	CreatedAt       time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"not null;default:now()" json:"updated_at"`

	// Associations
	Device *Device `gorm:"foreignKey:DeviceID" json:"-"`
}

// TableName overrides the GORM table name convention.
func (DeviceStatus) TableName() string {
	return "device_statuses"
}
