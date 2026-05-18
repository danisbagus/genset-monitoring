import api from './axios'
import { API_ENDPOINTS } from './api.constant'
import type {
  DashboardSummary,
  DashboardDevicesResponse,
  DashboardAlertsResponse
} from '@/types/dashboard'

/**
 * Dashboard API Service
 * Handles all dashboard-related API requests
 */
export const dashboardApi = {
  /**
   * Get dashboard summary statistics
   */
  getSummary: async (): Promise<DashboardSummary> => {
    const response = await api.get(API_ENDPOINTS.DASHBOARD.SUMMARY)
    return response.data.data
  },

  /**
   * Get device status list with pagination
   */
  getDeviceStates: async (page = 1, limit = 10): Promise<DashboardDevicesResponse> => {
    const response = await api.get(API_ENDPOINTS.DASHBOARD.DEVICE_STATES, {
      params: { page, limit }
    })
    return response.data.data
  },

  /**
   * Get recent alerts list with pagination
   */
  getRecentAlerts: async (page = 1, limit = 10): Promise<DashboardAlertsResponse> => {
    const response = await api.get(API_ENDPOINTS.DASHBOARD.RECENT_ALERTS, {
      params: { page, limit }
    })
    return response.data.data
  }
}

export default dashboardApi
