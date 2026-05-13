import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/pages/auth/LoginPage.vue'),
    meta: { title: 'Login' }
  },
  {
    path: '/',
    name: 'dashboard',
    component: () => import('@/pages/dashboard/DashboardPage.vue'),
    meta: { title: 'Dashboard' }
  },
  {
    path: '/devices/:id',
    name: 'device-detail',
    component: () => import('@/pages/device/DeviceDetailPage.vue'),
    props: true,
    meta: { title: 'Device Details' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  document.title = `${to.meta.title} | Genset Monitoring` || 'Genset Monitoring'
  next()
})

export default router
