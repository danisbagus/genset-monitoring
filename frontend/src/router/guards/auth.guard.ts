import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'

export const setupAuthGuard = (
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) => {
  const authStore = useAuthStore()
  const isAuthenticated = authStore.isAuthenticated

  // Check if route requires auth
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)
  const isGuestOnly = to.matched.some(record => record.meta.guestOnly === true)

  if (requiresAuth && !isAuthenticated) {
    // Save the path they were trying to access
    next({ 
      name: 'login', 
      query: { redirect: to.fullPath } 
    })
  } else if (isGuestOnly && isAuthenticated) {
    // Redirect authenticated users to dashboard if they try to access login
    next({ name: 'dashboard' })
  } else {
    next()
  }
}
