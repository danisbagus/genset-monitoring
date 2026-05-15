package model

// DashboardSummary represent data summary for dashboard
type DashboardSummary struct {
	TotalDevices   int64 `json:"total_devices" gorm:"column:total_devices"`
	OnlineDevices  int64 `json:"online_devices" gorm:"column:online_devices"`
	OfflineDevices int64 `json:"offline_devices" gorm:"column:offline_devices"`
	RunningEngines int64 `json:"running_engines" gorm:"column:running_engines"`
	CriticalAlerts int64 `json:"critical_alerts" gorm:"column:critical_alerts"`
	WarningAlerts  int64 `json:"warning_alerts" gorm:"column:warning_alerts"`
}
