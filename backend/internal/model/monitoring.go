package model

import "time"

// MonitoringDevice represents a single device row in the monitoring dashboard list.
type MonitoringDevice struct {
	DeviceID            string     `gorm:"column:device_id"`
	DeviceCode          string     `gorm:"column:device_code"`
	DeviceName          string     `gorm:"column:device_name"`
	SerialNumber        string     `gorm:"column:serial_number"`
	EngineRunning       bool       `gorm:"column:engine_running"`
	Speed               *int32     `gorm:"column:speed"`
	CoolantTemperature  *float32   `gorm:"column:coolant_temperature"`
	OilPressure         *float32   `gorm:"column:oil_pressure"`
	FuelLevel           *float32   `gorm:"column:fuel_level"`
	BattVolt            *float32   `gorm:"column:batt_volt"`
	Frequency           *float32   `gorm:"column:frequency"`
	TotalVa             *float32   `gorm:"column:total_va"`
	PfAvg               *float32   `gorm:"column:pf_avg"`
	TelemetryRecordedAt *time.Time `gorm:"column:telemetry_recorded_at"`
	LastSeenAt          *time.Time `gorm:"column:last_seen_at"`
	LastOnlineAt        *time.Time `gorm:"column:last_online_at"`
	GSMSignal           *int16     `gorm:"column:gsm_signal"`
	GPSConnected        *bool      `gorm:"column:gps_connected"`
	ServerConnected     *bool      `gorm:"column:server_connected"`
	CANConnected        *bool      `gorm:"column:can_connected"`
	RS485Connected      *bool      `gorm:"column:rs485_connected"`
	SDCardOK            *bool      `gorm:"column:sd_card_ok"`
	DeviceStatus        string     `gorm:"column:device_status"`
	UpdatedAt           *time.Time `gorm:"column:updated_at"`
}

// MonitoringDeviceDetail holds the full realtime detail for a single device.
type MonitoringDeviceDetail struct {
	// A. Device Information
	DeviceInfo MonitoringDeviceInfo

	// B. Latest Device State
	LatestState *MonitoringLatestState

	// C. Device Connectivity Status
	Connectivity *MonitoringConnectivity

	// D. Latest Engine Telemetry
	LatestEngineTelemetry *EngineTelemetry

	// E. Latest Electrical Telemetry
	LatestElectricalTelemetry *ElectricalTelemetry
}

// MonitoringDeviceInfo is the basic device information section (from devices table).
type MonitoringDeviceInfo struct {
	ID              string      `gorm:"column:id"`
	DeviceCode      string      `gorm:"column:device_code"`
	Name            string      `gorm:"column:name"`
	SerialNumber    string      `gorm:"column:serial_number"`
	EngineID        string      `gorm:"column:engine_id"`
	GSMNumber       string      `gorm:"column:gsm_number"`
	FirmwareVersion string      `gorm:"column:firmware_version"`
	Status          string      `gorm:"column:status"`
	Metadata        interface{} `gorm:"-"`
}

// MonitoringLatestState holds realtime snapshot from device_latest_state.
type MonitoringLatestState struct {
	EngineRunning       bool       `gorm:"column:engine_running"`
	Speed               *int32     `gorm:"column:speed"`
	CoolantTemperature  *float32   `gorm:"column:coolant_temperature"`
	OilPressure         *float32   `gorm:"column:oil_pressure"`
	FuelLevel           *float32   `gorm:"column:fuel_level"`
	BattVolt            *float32   `gorm:"column:batt_volt"`
	Frequency           *float32   `gorm:"column:frequency"`
	TotalVa             *float32   `gorm:"column:total_va"`
	PfAvg               *float32   `gorm:"column:pf_avg"`
	TelemetryRecordedAt *time.Time `gorm:"column:telemetry_recorded_at"`
	LastSeenAt          *time.Time `gorm:"column:last_seen_at"`
	LastOnlineAt        *time.Time `gorm:"column:last_online_at"`
	UpdatedAt           *time.Time `gorm:"column:updated_at"`
}

// MonitoringConnectivity holds device connectivity/health status from device_statuses.
type MonitoringConnectivity struct {
	GSMSignal       int16     `gorm:"column:gsm_signal"`
	GPSConnected    bool      `gorm:"column:gps_connected"`
	ServerConnected bool      `gorm:"column:server_connected"`
	CANConnected    bool      `gorm:"column:can_connected"`
	RS485Connected  bool      `gorm:"column:rs485_connected"`
	SDCardOK        bool      `gorm:"column:sd_card_ok"`
	LastSeen        time.Time `gorm:"column:last_seen"`
}
