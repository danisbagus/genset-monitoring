import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/services/api/auth.api'
import type { User, AuthState, AuthData } from '@/types/auth'
import { getStorage, setStorage, removeStorage, STORAGE_KEYS } from '@/utils/storage'
import { calculateExpiry } from '@/utils/token'
import { clearFailedQueue } from '@/services/api/axios'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(getStorage<User>(STORAGE_KEYS.USER_DATA))
  const accessToken = ref<string | null>(getStorage<string>('access_token'))
  const refreshToken = ref<string | null>(getStorage<string>('refresh_token'))
  const expiresAt = ref<number | null>(getStorage<number>('expires_at'))
  const isLoading = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!accessToken.value && !!user.value)
  const userRole = computed(() => user.value?.role || null)

  // Actions
  const setAuthData = (data: AuthData) => {
    user.value = data.user
    accessToken.value = data.access_token
    refreshToken.value = data.refresh_token
    expiresAt.value = calculateExpiry(data.expires_in)

    // Persist
    setStorage(STORAGE_KEYS.USER_DATA, user.value)
    setStorage('access_token', accessToken.value)
    setStorage('refresh_token', refreshToken.value)
    setStorage('expires_at', expiresAt.value)
  }

  const login = async (credentials: any) => {
    isLoading.value = true
    try {
      const response = await authApi.login(credentials)
      if (response.success && response.data) {
        setAuthData(response.data)
      }
      return response
    } finally {
      isLoading.value = false
    }
  }

  const refreshSession = async () => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await authApi.refresh(refreshToken.value)
      if (response.success && response.data) {
        accessToken.value = response.data.access_token
        refreshToken.value = response.data.refresh_token
        expiresAt.value = calculateExpiry(response.data.expires_in)

        setStorage('access_token', accessToken.value)
        setStorage('refresh_token', refreshToken.value)
        setStorage('expires_at', expiresAt.value)
        
        return accessToken.value
      }
      throw new Error('Refresh failed')
    } catch (error) {
      logout(true) // Force logout on refresh error
      throw error
    }
  }

  const logout = async (force = false) => {
    const currentRefreshToken = refreshToken.value
    
    // Clear local state first for better UX
    const performLocalCleanup = () => {
      user.value = null
      accessToken.value = null
      refreshToken.value = null
      expiresAt.value = null

      removeStorage(STORAGE_KEYS.USER_DATA)
      removeStorage('access_token')
      removeStorage('refresh_token')
      removeStorage('expires_at')
      
      // Notify other tabs
      localStorage.setItem('auth_logout', Date.now().toString())
      
      clearFailedQueue()
    }

    if (!force && currentRefreshToken) {
      try {
        await authApi.logout(currentRefreshToken)
      } catch (error) {
        console.error('Logout API failed:', error)
      }
    }

    performLocalCleanup()
    router.push('/login')
  }

  /**
   * Listener for multi-tab logout sync
   */
  const setupStorageListener = () => {
    window.addEventListener('storage', (event) => {
      if (event.key === 'auth_logout') {
        logout(true)
      }
    })
  }

  return {
    user,
    accessToken,
    refreshToken,
    expiresAt,
    isLoading,
    isAuthenticated,
    userRole,
    login,
    logout,
    refreshSession,
    setupStorageListener
  }
})
