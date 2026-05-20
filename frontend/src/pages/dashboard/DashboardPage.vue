<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import SummaryCards from '../../components/SummaryCards.vue'
import DeviceStatusTable from '../../components/DeviceStatusTable.vue'
import LatestAlertsCardList from '../../components/LatestAlertsCardList.vue'
import BaseErrorMessage from '../../components/BaseErrorMessage.vue'
import { dashboardApi } from '../../services/api/dashboard.api'
import type { DashboardSummary, DeviceStatus, Alert, Pagination } from '../../types/dashboard'
import { useWebsocket } from '../../composables/useWebsocket'
import { useRelativeTime } from '../../composables/useRelativeTime'

// WebSocket & Update Status
const { status: wsStatus, connect: connectWs, onMessage } = useWebsocket()
const updated_at = ref<Date>(new Date())
const { relativeText: lastUpdateText } = useRelativeTime(updated_at)

// Data State
const summaryData = ref<DashboardSummary>({
  critical_alerts: 0,
  offline_devices: 0,
  online_devices: 0,
  running_engines: 0,
  total_devices: 0,
  warning_alerts: 0
})

const deviceData = ref<DeviceStatus[]>([])
const alertData = ref<Alert[]>([])

const devicePagination = ref<Pagination>({
  limit: 10,
  page: 1,
  total: 0
})

const alertPagination = ref<Pagination>({
  limit: 10,
  page: 1,
  total: 0
})

const isLoading = ref(false)
const error = ref<string | null>(null)

const fetchData = async () => {
  // If only changing pages, we might not want to set global isLoading to true 
  // but let's keep it simple for now as the components handle their own loading props
  isLoading.value = true
  error.value = null
  
  try {
    const [summary, devicesResponse, alertsResponse] = await Promise.all([
      dashboardApi.getSummary(),
      dashboardApi.getDeviceStates(devicePagination.value.page, devicePagination.value.limit),
      dashboardApi.getRecentAlerts(alertPagination.value.page, alertPagination.value.limit)
    ])

    summaryData.value = summary
    console.log("devicesResponse", devicesResponse)
    console.log("alertsResponse", alertsResponse)
    deviceData.value = devicesResponse.devices
    devicePagination.value = devicesResponse.pagination
    
    alertData.value = alertsResponse.alerts
    alertPagination.value = alertsResponse.pagination
    
    // API /summary was successfully called; update updated_at
    updated_at.value = new Date()
  } catch (err: any) {
    console.error('Failed to fetch dashboard data:', err)
    error.value = 'Failed to load dashboard data. Please try again later.'
  } finally {
    isLoading.value = false
  }
}

// Listen to WebSocket messages and register the cleanup callback
const unsubscribeWs = onMessage((message: any) => {
  if (message && message.event === 'dashboard.summary.updated' && message.data) {
    summaryData.value = message.data
    
    // Update the last updated time using the timestamp from the WS message
    const wsTimestamp = message.timestamp || message.ts
    if (wsTimestamp) {
      updated_at.value = new Date(wsTimestamp)
    } else {
      updated_at.value = new Date()
    }
  }
})

const handleDevicePageChange = (page: number) => {
  devicePagination.value.page = page
  fetchData()
}

const handleAlertPageChange = (page: number) => {
  alertPagination.value.page = page
  fetchData()
}

const refreshData = () => {
  fetchData()
}

onMounted(() => {
  fetchData()
  connectWs()
})

// Clean up websocket listener when component is unmounted
onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs()
  }
})
</script>

<template>
  <div class="space-y-6 sm:space-y-8">
    <!-- Header Section -->
    <header class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-xl sm:text-2xl md:text-3xl font-bold text-slate-900 dark:text-white flex items-center gap-3">
          Dashboard Overview
          <span class="flex h-2.5 w-2.5 relative">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"></span>
          </span>
        </h1>
        <div class="flex flex-wrap items-center gap-x-4 gap-y-1 mt-1">
          <p class="text-sm text-slate-500 dark:text-slate-400">Real-time status of your genset fleet</p>
          <div class="hidden sm:flex items-center gap-2">
            <span class="h-1 w-1 rounded-full bg-slate-300 dark:bg-slate-700"></span>
            <p class="text-xs font-medium text-slate-400 dark:text-slate-500 italic">{{ lastUpdateText }}</p>
          </div>
        </div>
      </div>
      
      <div class="flex items-center gap-3">
        <div class="flex items-center space-x-2 bg-white dark:bg-slate-900 px-3 py-1.5 rounded-full border border-slate-200 dark:border-slate-800 shadow-sm transition-all duration-300">
          <div class="flex h-2 w-2 relative">
            <span 
              v-if="wsStatus === 'online'" 
              class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"
            ></span>
            <span 
              class="relative inline-flex rounded-full h-2 w-2"
              :class="{
                'bg-emerald-500': wsStatus === 'online',
                'bg-amber-500': wsStatus === 'reconnecting',
                'bg-rose-500': wsStatus === 'offline'
              }"
            ></span>
          </div>
          <span class="text-[10px] font-bold uppercase tracking-wider" :class="{
            'text-emerald-600 dark:text-emerald-400': wsStatus === 'online',
            'text-amber-600 dark:text-amber-400': wsStatus === 'reconnecting',
            'text-rose-600 dark:text-rose-400': wsStatus === 'offline'
          }">
            {{ wsStatus === 'online' ? 'Live' : wsStatus === 'reconnecting' ? 'Reconnecting' : 'Disconnected' }}
          </span>
        </div>

        <button 
          @click="refreshData"
          class="bg-white dark:bg-slate-900 p-2 rounded-lg border border-slate-200 dark:border-slate-800 shadow-sm text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors group"
          title="Refresh Data"
        >
          <svg 
            xmlns="http://www.w3.org/2000/svg" 
            width="18" 
            height="18" 
            viewBox="0 0 24 24" 
            fill="none" 
            stroke="currentColor" 
            stroke-width="2.5" 
            stroke-linecap="round" 
            stroke-linejoin="round"
            class="group-active:rotate-180 transition-transform duration-500"
          >
            <path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/>
            <path d="M21 3v5h-5"/>
          </svg>
        </button>
      </div>
    </header>

    <!-- 1. Summary Cards Section -->
    <section v-if="!error">
      <SummaryCards :summary="summaryData" />
    </section>

    <!-- Error State -->
    <BaseErrorMessage 
      v-if="error" 
      :message="error" 
      @retry="fetchData" 
    />

    <!-- Main Content Grid -->
    <div v-if="!error" class="grid grid-cols-1 lg:grid-cols-3 gap-6 sm:gap-8">
      <!-- 2. Device Status Table Section -->
      <section class="lg:col-span-2">
        <DeviceStatusTable 
          :devices="deviceData" 
          :pagination="devicePagination" 
          :loading="isLoading" 
          @change-page="handleDevicePageChange"
        />
      </section>

      <!-- 3. Latest Alerts Card List Section -->
      <section class="lg:col-span-1">
        <LatestAlertsCardList 
          :alerts="alertData" 
          :pagination="alertPagination" 
          :loading="isLoading" 
          @change-page="handleAlertPageChange"
        />
      </section>
    </div>
  </div>
</template>

<style scoped>
/* Additional subtle animations */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

section {
  animation: fadeIn 0.5s ease-out forwards;
}

section:nth-child(2) { animation-delay: 0.1s; }
section:nth-child(3) { animation-delay: 0.2s; }
</style>
