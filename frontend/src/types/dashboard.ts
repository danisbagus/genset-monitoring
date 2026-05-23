export interface DashboardSummary {
  critical_alerts: number
  offline_devices: number
  online_devices: number
  running_engines: number
  total_devices: number
  warning_alerts: number
}

export interface DeviceStatus {
  coolant_temperature: number
  device_id: string
  device_name: string
  device_online: boolean
  engine_running: boolean
  fuel_level: number
  last_seen_at: string
  isHighlighted?: boolean
}

export interface Alert {
  acknowledged: boolean
  alert_id: string
  created_at: string
  device_id: string
  device_name: string
  message: string
  severity: 'critical' | 'warning' | 'info' | string
}

export interface Pagination {
  limit: number
  page: number
  total: number
}

export interface DashboardDevicesResponse {
  devices: DeviceStatus[]
  pagination: Pagination
}

export interface DashboardAlertsResponse {
  alerts: Alert[]
  pagination: Pagination
}
