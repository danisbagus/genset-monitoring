package model

import "time"

// DashboardSummary represent data summary for dashboard
type DashboardSummary struct {
	TotalDevices   int64 `json:"total_devices" gorm:"column:total_devices"`
	OnlineDevices  int64 `json:"online_devices" gorm:"column:online_devices"`
	OfflineDevices int64 `json:"offline_devices" gorm:"column:offline_devices"`
	RunningEngines int64 `json:"running_engines" gorm:"column:running_engines"`
	CriticalAlerts int64 `json:"critical_alerts" gorm:"column:critical_alerts"`
	WarningAlerts  int64 `json:"warning_alerts" gorm:"column:warning_alerts"`
}

// DeviceState represent state of a single device for dashboard
type DeviceState struct {
	DeviceID           string     `gorm:"column:device_id"`
	DeviceName         string     `gorm:"column:device_name"`
	DeviceOnline       bool       `gorm:"column:device_online"`
	EngineRunning      bool       `gorm:"column:engine_running"`
	FuelLevel          float64    `gorm:"column:fuel_level"`
	CoolantTemperature float64    `gorm:"column:coolant_temperature"`
	LastSeenAt         *time.Time `gorm:"column:last_seen_at"`
}

// RecentAlert represent a single alert for dashboard
type RecentAlert struct {
	AlertID      string    `gorm:"column:alert_id"`
	DeviceID     string    `gorm:"column:device_id"`
	DeviceName   string    `gorm:"column:device_name"`
	Severity     string    `gorm:"column:severity"`
	Message      string    `gorm:"column:message"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	Acknowledged bool      `gorm:"column:acknowledged"`
}
