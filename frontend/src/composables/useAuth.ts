import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth.store'
import { useRouter } from 'vue-router'

export const useAuth = () => {
  const authStore = useAuthStore()
  const router = useRouter()
  
  const { 
    user, 
    accessToken, 
    isAuthenticated, 
    isLoading,
    userRole 
  } = storeToRefs(authStore)

  const handleLogin = async (credentials: any) => {
    try {
      const response = await authStore.login(credentials)
      if (response.success) {
        // Redirect will be handled by the component or here
        return { success: true, message: response.message }
      }
      return { success: false, message: response.message }
    } catch (error: any) {
      const message = error.response?.data?.message || error.message || 'Login failed'
      return { success: false, message }
    }
  }

  const handleLogout = async () => {
    await authStore.logout()
  }

  return {
    user,
    accessToken,
    isAuthenticated,
    isLoading,
    userRole,
    login: handleLogin,
    logout: handleLogout
  }
}
