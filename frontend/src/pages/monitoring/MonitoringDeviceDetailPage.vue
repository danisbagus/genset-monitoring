<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMonitoringDeviceDetail } from '@/composables/monitoring/useMonitoringDeviceDetail'
import DeviceInfoSection from '@/components/monitoring/DeviceInfoSection.vue'
import DeviceEngineSection from '@/components/monitoring/DeviceEngineSection.vue'
import DeviceElectricalSection from '@/components/monitoring/DeviceElectricalSection.vue'
import DeviceConnectivitySection from '@/components/monitoring/DeviceConnectivitySection.vue'
import BaseErrorMessage from '@/components/BaseErrorMessage.vue'

const route = useRoute()
const router = useRouter()
const deviceId = route.params.deviceId as string

const {
  detail,
  isLoading,
  error,
  notFound,
  fetchDetail
} = useMonitoringDeviceDetail()

onMounted(() => fetchDetail(deviceId))

const goBack = () => router.push('/monitoring/devices')
</script>

<template>
  <div class="space-y-6">

    <!-- Breadcrumb / back -->
    <div class="flex items-center gap-3">
      <button
        @click="goBack"
        class="flex items-center gap-1.5 text-sm font-bold text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 transition-colors"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        Device Monitoring
      </button>
      <span class="h-4 w-px bg-slate-200 dark:bg-slate-700"></span>
      <span class="text-xs font-mono text-slate-400 uppercase tracking-widest truncate max-w-[200px]">{{ deviceId }}</span>
    </div>

    <!-- Loading skeleton -->
    <template v-if="isLoading">
      <div class="animate-pulse space-y-6">
        <!-- Header skeleton -->
        <div class="h-24 bg-slate-100 dark:bg-slate-800 rounded-2xl"></div>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div class="h-64 bg-slate-100 dark:bg-slate-800 rounded-2xl"></div>
          <div class="h-64 bg-slate-100 dark:bg-slate-800 rounded-2xl"></div>
        </div>
        <div class="h-72 bg-slate-100 dark:bg-slate-800 rounded-2xl"></div>
        <div class="h-52 bg-slate-100 dark:bg-slate-800 rounded-2xl"></div>
      </div>
    </template>

    <!-- Not found -->
    <template v-else-if="notFound">
      <div class="py-20 flex flex-col items-center gap-4 text-center bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl">
        <div class="w-16 h-16 rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-slate-400">
          <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div>
          <p class="font-bold text-slate-700 dark:text-slate-300 text-base">Device not found</p>
          <p class="text-sm text-slate-400 mt-1">The device with ID <code class="font-mono bg-slate-100 dark:bg-slate-800 px-1 rounded">{{ deviceId }}</code> does not exist.</p>
        </div>
        <button
          @click="goBack"
          class="mt-2 px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-bold rounded-xl transition-all shadow-sm shadow-blue-600/20"
        >
          Back to Device List
        </button>
      </div>
    </template>

    <!-- API error -->
    <template v-else-if="error">
      <BaseErrorMessage :message="error" @retry="fetchDetail(deviceId)" />
    </template>

    <!-- Detail content -->
    <template v-else-if="detail">
      <!-- Page header with device name -->
      <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl px-6 py-4 shadow-sm">
        <div class="flex items-center gap-4">
          <div
            class="w-12 h-12 rounded-2xl flex items-center justify-center text-white shadow-lg transition-all"
            :class="detail.latest_state?.engine_running ? 'bg-blue-600 shadow-blue-600/25' : 'bg-slate-500 shadow-slate-500/20'"
          >
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
            </svg>
          </div>
          <div>
            <h1 class="text-lg sm:text-xl font-bold text-slate-900 dark:text-white">{{ detail.device_info.name }}</h1>
            <p class="text-xs font-mono text-slate-400">{{ detail.device_info.device_code }} · {{ detail.device_info.serial_number }}</p>
          </div>
        </div>

        <!-- Refresh -->
        <button
          @click="fetchDetail(deviceId)"
          class="self-start sm:self-auto flex items-center gap-2 px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-bold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700 transition-all"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M21 3v5h-5" />
          </svg>
          Refresh
        </button>
      </header>

      <!-- Sections grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Device Info -->
        <DeviceInfoSection
          :device-info="detail.device_info"
          :latest-state="detail.latest_state"
        />

        <!-- Connectivity -->
        <DeviceConnectivitySection
          :connectivity="detail.connectivity"
        />
      </div>

      <!-- Engine Telemetry (full width) -->
      <DeviceEngineSection
        :latest-state="detail.latest_state"
        :engine-telemetry="detail.engine_telemetry"
      />

      <!-- Electrical Telemetry (full width) -->
      <DeviceElectricalSection
        :latest-state="detail.latest_state"
        :electrical-telemetry="detail.electrical_telemetry"
      />
    </template>

  </div>
</template>
