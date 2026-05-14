export interface ApiResponse<T> {
  data: T
  message: string
  success: boolean
}

export interface ApiError {
  message: string
  success: boolean
  errors?: Record<string, string[]>
}
