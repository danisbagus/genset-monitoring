import type { ApiResponse } from './api'

export interface User {
  id: string
  username: string
  email: string
  role: 'admin' | 'user'
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
  expires_in: number // expiry timestamp in ms
}

export interface AuthData extends AuthTokens {
  user: User
}

export interface LoginResponse extends ApiResponse<AuthData> {}

export interface RefreshResponse extends ApiResponse<AuthTokens> {}

export interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  expiresAt: number | null
  isAuthenticated: boolean
  isLoading: boolean
}
