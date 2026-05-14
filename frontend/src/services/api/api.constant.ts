/**
 * API Endpoint Constants
 * These endpoints are relative to the axios baseURL (usually /api)
 */

export const API_ENDPOINTS = {
  AUTH: {
    // Note: v1 prefix is included here as part of the service path
    LOGIN: '/v1/auth/login',
    REGISTER: '/v1/auth/register',
    REFRESH: '/v1/auth/refresh',
    LOGOUT: '/v1/auth/logout',
    ME: '/v1/auth/me',
  },
  DEVICES: {
    LIST: '/v1/devices',
    CREATE: '/v1/devices',
    DETAIL: (id: string) => `/v1/devices/${id}`,
    UPDATE: (id: string) => `/v1/devices/${id}`,
    DELETE: (id: string) => `/v1/devices/${id}`,
    STATUS: (id: string) => `/v1/devices/${id}/status`,
    HEARTBEAT: (id: string) => `/v1/devices/${id}/heartbeat`,
  },
  TELEMETRY: {
    ENGINE_CREATE: (id: string) => `/v1/devices/${id}/engine`,
    ENGINE_LATEST: (id: string) => `/v1/devices/${id}/engine/latest`,
    ELECTRICAL_CREATE: (id: string) => `/v1/devices/${id}/electrical`,
    ELECTRICAL_LATEST: (id: string) => `/v1/devices/${id}/electrical/latest`,
  },
} as const
