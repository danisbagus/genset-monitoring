import api from './axios'
import type { LoginResponse, RefreshResponse } from '@/types/auth'
import { API_ENDPOINTS } from './api.constant'

export const authApi = {
  /**
   * Login user with username and password
   */
  login: async (credentials: any): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>(API_ENDPOINTS.AUTH.LOGIN, credentials)
    return response.data
  },

  /**
   * Refresh access token using refresh token
   */
  refresh: async (refreshToken: string): Promise<RefreshResponse> => {
    const response = await api.post<RefreshResponse>(API_ENDPOINTS.AUTH.REFRESH, {
      refresh_token: refreshToken
    })
    return response.data
  },

  /**
   * Logout user by blacklisting the refresh token
   */
  logout: async (refreshToken: string): Promise<void> => {
    await api.post(API_ENDPOINTS.AUTH.LOGOUT, {
      refresh_token: refreshToken
    })
  }
}
