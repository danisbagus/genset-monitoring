<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useTimeAgo, useTimestamp } from '@vueuse/core'
import SummaryCards from '../../components/SummaryCards.vue'
import DeviceStatusTable from '../../components/DeviceStatusTable.vue'
import LatestAlertsTable from '../../components/LatestAlertsTable.vue'
import type { DashboardSummary, DeviceStatus, Alert, Pagination } from '../../types/dashboard'

// WebSocket & Update Status
const wsStatus = ref<'online' | 'reconnecting' | 'offline'>('online')
const lastUpdate = ref(new Date())
const lastUpdateText = useTimeAgo(lastUpdate)

// Dummy Data Initialization
const summaryData = ref<DashboardSummary>({
  critical_alerts: 2,
  offline_devices: 1,
  online_devices: 4,
  running_engines: 3,
  total_devices: 5,
  warning_alerts: 5
})

const deviceData = ref<DeviceStatus[]>([
  {
    device_id: "GS-001",
    device_name: "Genset Alpha-1",
    device_online: true,
    engine_running: true,
    fuel_level: 85,
    coolant_temperature: 78,
    last_seen_at: new Date().toISOString()
  },
  {
    device_id: "GS-002",
    device_name: "Genset Bravo-2",
    device_online: true,
    engine_running: false,
    fuel_level: 64,
    coolant_temperature: 42,
    last_seen_at: new Date().toISOString()
  },
  {
    device_id: "GS-003",
    device_name: "Genset Gamma-3",
    device_online: false,
    engine_running: false,
    fuel_level: 12,
    coolant_temperature: 25,
    last_seen_at: new Date(Date.now() - 3600000).toISOString()
  },
  {
    device_id: "GS-004",
    device_name: "Genset Delta-4",
    device_online: true,
    engine_running: true,
    fuel_level: 92,
    coolant_temperature: 95,
    last_seen_at: new Date().toISOString()
  },
  {
    device_id: "GS-005",
    device_name: "Genset Epsilon-5",
    device_online: true,
    engine_running: true,
    fuel_level: 45,
    coolant_temperature: 82,
    last_seen_at: new Date().toISOString()
  }
])

const alertData = ref<Alert[]>([
  {
    alert_id: "AL-101",
    device_id: "GS-004",
    device_name: "Genset Delta-4",
    severity: "critical",
    message: "High Coolant Temperature detected (95°C)",
    acknowledged: false,
    created_at: new Date(Date.now() - 600000).toISOString()
  },
  {
    alert_id: "AL-102",
    device_id: "GS-003",
    device_name: "Genset Gamma-3",
    severity: "warning",
    message: "Low Fuel Level (12%)",
    acknowledged: true,
    created_at: new Date(Date.now() - 1800000).toISOString()
  },
  {
    alert_id: "AL-103",
    device_id: "GS-001",
    device_name: "Genset Alpha-1",
    severity: "info",
    message: "Scheduled maintenance in 24 hours",
    acknowledged: false,
    created_at: new Date(Date.now() - 3600000).toISOString()
  }
])

const devicePagination = ref<Pagination>({
  limit: 10,
  page: 1,
  total: 5
})

const alertPagination = ref<Pagination>({
  limit: 10,
  page: 1,
  total: 3
})

const isLoading = ref(false)

const refreshData = () => {
  lastUpdate.value = new Date()
  // Simulate data update
  if (deviceData.value[0]) {
    deviceData.value[0].last_seen_at = new Date().toISOString()
  }
}

onMounted(() => {
  // Simulate initial loading
  isLoading.value = true
  setTimeout(() => {
    isLoading.value = false
    refreshData()
  }, 800)

  // Simulate periodic WebSocket updates
  setInterval(() => {
    refreshData()
  }, 10000)

  // Simulate random WS status changes for demo
  setTimeout(() => {
    // wsStatus.value = 'reconnecting'
  }, 15000)
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
            <p class="text-xs font-medium text-slate-400 dark:text-slate-500 italic">Last updated {{ lastUpdateText }}</p>
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
    <section>
      <SummaryCards :summary="summaryData" />
    </section>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 xl:grid-cols-2 gap-6 sm:gap-8">
      <!-- 2. Device Status Table Section -->
      <section>
        <DeviceStatusTable 
          :devices="deviceData" 
          :pagination="devicePagination" 
          :loading="isLoading" 
        />
      </section>

      <!-- 3. Latest Alerts Table Section -->
      <section>
        <LatestAlertsTable 
          :alerts="alertData" 
          :pagination="alertPagination" 
          :loading="isLoading" 
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
