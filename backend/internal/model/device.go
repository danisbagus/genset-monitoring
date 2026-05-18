package model

// DeviceLifecycle represents the lifecycle status of a device.
type DeviceLifecycle string

const (
	DeviceLifecycleActive      DeviceLifecycle = "active"
	DeviceLifecycleInactive    DeviceLifecycle = "inactive"
	DeviceLifecycleMaintenance DeviceLifecycle = "maintenance"
)

// Device represents a physical genset device registered in the system.
//
// Security note: HashedDeviceToken stores the SHA-256 hash of the opaque
// device token. The raw token is returned only once at creation and is never
// persisted. This mirrors the refresh-token pattern used in the auth module.
type Device struct {
	Base

	// Identity
	DeviceCode   string `gorm:"column:device_code;uniqueIndex;not null" json:"device_code"`
	Name         string `gorm:"column:name;not null"                   json:"name"`
	SerialNumber string `gorm:"column:serial_number;uniqueIndex;not null" json:"serial_number"`
	EngineID     string `gorm:"column:engine_id"                        json:"engine_id"`
	GSMNumber    string `gorm:"column:gsm_number"                       json:"gsm_number"`

	// Firmware
	FirmwareVersion string `gorm:"column:firmware_version" json:"firmware_version"`

	// Lifecycle
	Status DeviceLifecycle `gorm:"column:status;type:device_status;default:active" json:"status"`

	// Auth token (hash stored, raw returned once at creation only)
	HashedDeviceToken string `gorm:"column:hashed_device_token" json:"-"`

	// Extensible metadata (stored as JSONB)
	Metadata []byte `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
}

// TableName overrides the GORM table name convention.
func (Device) TableName() string { return "devices" }
