/**
 * Monitoring Module Types
 * Based on swagger.json MonitoringDeviceItem, MonitoringDeviceDetailOutput, etc.
 */

// ── List Filters ─────────────────────────────────────────────────────────────

export interface MonitoringDeviceFilters {
  search?: string
  online?: boolean | ''
  engine_running?: boolean | ''
  status?: 'active' | 'inactive' | 'maintenance' | ''
  sort_by?: 'updated_at' | 'last_seen_at' | 'telemetry_recorded_at' | 'name' | ''
  sort_order?: 'asc' | 'desc'
  limit: number
  offset: number
}

// ── List Response ─────────────────────────────────────────────────────────────

export interface MonitoringDeviceItem {
  device_id: string
  device_name: string
  device_code: string
  serial_number: string
  device_status: string         // active | inactive | maintenance
  engine_running: boolean
  last_seen_at: string          // ISO timestamp
  last_online_at: string        // ISO timestamp
  telemetry_recorded_at: string // ISO timestamp

  // Connectivity
  server_connected: boolean
  gps_connected: boolean
  can_connected: boolean
  rs485_connected: boolean
  sd_card_ok: boolean
  gsm_signal: number

  // Engine snapshot
  speed: number                 // RPM
  oil_pressure: number
  coolant_temperature: number
  fuel_level: number
  batt_volt: number

  // Electrical snapshot
  total_va: number
  frequency: number
  pf_avg: number
}

export interface MonitoringPagination {
  limit: number
  page: number
  total: number
}

export interface MonitoringDeviceListOutput {
  devices: MonitoringDeviceItem[]
  pagination: MonitoringPagination
}

// ── Detail Response ───────────────────────────────────────────────────────────

export interface MonitoringDeviceInfoOutput {
  id: string
  name: string
  device_code: string
  serial_number: string
  status: string
  engine_id: string
  firmware_version: string
  gsm_number: string
  metadata: unknown
}

export interface MonitoringLatestStateOutput {
  engine_running: boolean
  last_seen_at: string
  last_online_at: string
  telemetry_recorded_at: string
  updated_at: string

  // Engine snapshot
  speed: number
  oil_pressure: number
  coolant_temperature: number
  fuel_level: number
  batt_volt: number

  // Electrical snapshot
  total_va: number
  frequency: number
  pf_avg: number
}

export interface MonitoringConnectivityOutput {
  server_connected: boolean
  gps_connected: boolean
  can_connected: boolean
  rs485_connected: boolean
  sd_card_ok: boolean
  gsm_signal: number
  last_seen: string
}

export interface EngineTelemetryOutput {
  device_id: string
  created_at: string
  speed: number
  oil_pressure: number
  coolant_temperature: number
  fuel_level_top: number
  fuel_level_bottom: number
  fuel_level_pressure_1: number
  fuel_level_pressure_2: number
  fuel_rate: number
  avg_fuel_rate: number
  total_fuel: number
  trip_fuel: number
  batt_volt: number
  keyswitch_batt_potential: number
  oil_filter_out_pressure: number
  turbo_pressure: number
  intake_manifold_pressure: number
  intake_manifold_temperature: number
  ecu_temperature: number
  rated_power: number
  rated_speed: number
  desired_operating_speed: number
  run_time: number
}

export interface ElectricalTelemetryOutput {
  device_id: string
  created_at: string
  l1_n_volt: number
  l2_n_volt: number
  l3_n_volt: number
  l1_l2_volt: number
  l2_l3_volt: number
  l3_l1_volt: number
  l1_curr: number
  l2_curr: number
  l3_curr: number
  earth_curr: number
  frequency: number
  pf_l1: number
  pf_l2: number
  pf_l3: number
  pf_avg: number
  l1_va: number
  l2_va: number
  l3_va: number
  total_va: number
  l1_var: number
  l2_var: number
  l3_var: number
  total_var: number
  charge_alt_volt: number
  percent_fp: number
  percent_fv: number
}

export interface MonitoringDeviceDetailOutput {
  device_info: MonitoringDeviceInfoOutput
  latest_state: MonitoringLatestStateOutput
  connectivity: MonitoringConnectivityOutput
  engine_telemetry: EngineTelemetryOutput | null
  electrical_telemetry: ElectricalTelemetryOutput | null
}
