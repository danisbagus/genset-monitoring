import api from './axios'
import { API_ENDPOINTS } from './api.constant'
import type {
  MonitoringDeviceFilters,
  MonitoringDeviceListOutput,
  MonitoringDeviceDetailOutput
} from '@/types/monitoring'

/**
 * Monitoring API Service
 * Handles all monitoring-related API requests
 */
export const monitoringApi = {
  /**
   * Get paginated device list with realtime monitoring summary data.
   * Uses offset-based pagination (not page-based).
   */
  getDevices: async (filters: Partial<MonitoringDeviceFilters> = {}): Promise<MonitoringDeviceListOutput> => {
    const params: Record<string, unknown> = {}

    if (filters.search) params.search = filters.search
    if (filters.online !== '' && filters.online !== undefined) params.online = filters.online
    if (filters.engine_running !== '' && filters.engine_running !== undefined) params.engine_running = filters.engine_running
    if (filters.status) params.status = filters.status
    if (filters.sort_by) params.sort_by = filters.sort_by
    if (filters.sort_order) params.sort_order = filters.sort_order
    params.limit = filters.limit ?? 20
    params.offset = filters.offset ?? 0

    const response = await api.get(API_ENDPOINTS.MONITORING.DEVICES, { params })
    return response.data.data
  },

  /**
   * Get full realtime monitoring detail for a single device.
   */
  getDeviceDetail: async (deviceId: string): Promise<MonitoringDeviceDetailOutput> => {
    const response = await api.get(API_ENDPOINTS.MONITORING.DEVICE_DETAIL(deviceId))
    return response.data.data
  }
}

export default monitoringApi
