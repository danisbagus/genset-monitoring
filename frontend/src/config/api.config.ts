/**
 * Centralized API configuration
 * Uses Vite's import.meta.env for environment-specific values
 */

export const API_CONFIG = {
  // Base URL for REST API requests
  // Defaults to /api which works with Vite proxy in dev and Nginx in prod
  BASE_URL: import.meta.env.VITE_API_BASE_URL || '/api',

  // Base URL for WebSocket connections
  // Defaults to /ws
  WS_URL: import.meta.env.VITE_WS_BASE_URL || '/ws',

  // Timeout for API requests in milliseconds
  TIMEOUT: 15000,
}

export default API_CONFIG
