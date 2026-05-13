import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'

const config: AxiosRequestConfig = {
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
}

const axiosInstance: AxiosInstance = axios.create(config)

// Request Interceptor
axiosInstance.interceptors.request.use(
  (config) => {
    // You can add auth token here in the future
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response Interceptor
axiosInstance.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    // Handle global errors here
    return Promise.reject(error)
  }
)

export default axiosInstance
