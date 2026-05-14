import axios from 'axios'
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/auth.store'
import { isTokenExpired } from '@/utils/token'
import router from '@/router'

import { API_CONFIG } from '@/config/api.config'
import { API_ENDPOINTS } from './api.constant'

// Create axios instance
const api: AxiosInstance = axios.create({
  baseURL: API_CONFIG.BASE_URL,
  timeout: API_CONFIG.TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
})

// Queue for failed requests during token refresh
let isRefreshing = false
let failedQueue: any[] = []

/**
 * Clears the failed requests queue, rejecting all pending promises.
 */
export const clearFailedQueue = (error: any = new Error('Session expired')) => {
  failedQueue.forEach((prom) => prom.reject(error))
  failedQueue = []
  isRefreshing = false
}

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(token)
    }
  })
  failedQueue = []
}

// Request Interceptor
api.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const authStore = useAuthStore()
    
    // Skip expiration check for login/refresh/register endpoints
    const isAuthRequest = [
      API_ENDPOINTS.AUTH.LOGIN,
      API_ENDPOINTS.AUTH.REFRESH,
      API_ENDPOINTS.AUTH.REGISTER
    ].some(endpoint => config.url?.includes(endpoint))

    if (authStore.accessToken && !isAuthRequest) {
      // Check if token is expired or about to expire
      if (isTokenExpired(authStore.expiresAt)) {
        if (!isRefreshing) {
          try {
            const newToken = await authStore.refreshSession()
            if (config.headers) {
              config.headers.Authorization = `Bearer ${newToken}`
            }
          } catch (error) {
            // Refresh failed, logout is handled in store
            return Promise.reject(error)
          }
        } else {
          // If already refreshing, wait for it
          return new Promise((resolve, reject) => {
            failedQueue.push({ resolve, reject })
          })
            .then((token) => {
              if (config.headers) {
                config.headers.Authorization = `Bearer ${token}`
              }
              return config
            })
            .catch((err) => Promise.reject(err))
        }
      } else {
        if (config.headers) {
          config.headers.Authorization = `Bearer ${authStore.accessToken}`
        }
      }
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response Interceptor
api.interceptors.response.use(
  (response) => {
    return response
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }
    const authStore = useAuthStore()

    // Handle Network Error / Timeout
    if (!error.response) {
      console.error('Network Error:', error.message)
      return Promise.reject(error)
    }

    // Handle 401 Unauthorized (Token Expired)
    if (error.response.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // If already refreshing, add to queue
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        })
          .then((token) => {
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`
            }
            return api(originalRequest)
          })
          .catch((err) => {
            return Promise.reject(err)
          })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        const newToken = await authStore.refreshSession()
        processQueue(null, newToken)
        
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`
        }
        return api(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError, null)
        authStore.logout(true) // Force logout on refresh failure
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    // Handle other errors (403, 404, 500)
    return Promise.reject(error)
  }
)

export default api
