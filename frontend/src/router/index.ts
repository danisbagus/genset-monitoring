import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { setupAuthGuard } from './guards/auth.guard'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/pages/auth/LoginPage.vue'),
    meta: { 
      title: 'Login',
      requiresAuth: false,
      guestOnly: true,
      layout: 'auth'
    }
  },
  {
    path: '/',
    name: 'dashboard',
    component: () => import('@/pages/dashboard/DashboardPage.vue'),
    meta: { 
      title: 'Dashboard',
      requiresAuth: true
    }
  },
  {
    path: '/devices/:id',
    name: 'device-detail',
    component: () => import('@/pages/device/DeviceDetailPage.vue'),
    props: true,
    meta: { 
      title: 'Device Details',
      requiresAuth: true
    }
  },

  // ── Monitoring Module ────────────────────────────────────────────────────────
  {
    path: '/monitoring/devices',
    name: 'monitoring-devices',
    component: () => import('@/pages/monitoring/MonitoringDevicesPage.vue'),
    meta: {
      title: 'Device Monitoring',
      requiresAuth: true
    }
  },
  {
    path: '/monitoring/devices/:deviceId',
    name: 'monitoring-device-detail',
    component: () => import('@/pages/monitoring/MonitoringDeviceDetailPage.vue'),
    props: true,
    meta: {
      title: 'Device Detail',
      requiresAuth: true
    }
  },
  {
    path: '/monitoring/alerts',
    name: 'monitoring-alerts',
    component: () => import('@/pages/monitoring/MonitoringAlertsPage.vue'),
    meta: {
      title: 'Alert Monitoring',
      requiresAuth: true
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  // Set page title
  document.title = `${to.meta.title} | Genset Monitoring` || 'Genset Monitoring'
  
  // Execute auth guard
  setupAuthGuard(to, from, next)
})

export default router
