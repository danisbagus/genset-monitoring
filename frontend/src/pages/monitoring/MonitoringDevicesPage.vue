<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useMonitoringDevices } from '@/composables/monitoring/useMonitoringDevices'
import MonitoringDeviceFilters from '@/components/monitoring/MonitoringDeviceFilters.vue'
import MonitoringDeviceTable from '@/components/monitoring/MonitoringDeviceTable.vue'
import BaseErrorMessage from '@/components/BaseErrorMessage.vue'

const {
  devices,
  pagination,
  filters,
  isLoading,
  error,
  currentPage,
  fetchDevices,
  applyFilters,
  resetFilters,
  goToPage
} = useMonitoringDevices()

onMounted(() => fetchDevices())
</script>

<template>
  <div class="space-y-5">

    <!-- Page Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2.5">
          <span class="w-8 h-8 rounded-xl bg-blue-600 flex items-center justify-center text-white shadow-sm shadow-blue-600/30">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
            </svg>
          </span>
          Device Monitoring
        </h1>
        <p class="text-sm text-slate-500 dark:text-slate-400 mt-0.5">
          {{ isLoading ? 'Loading…' : `${pagination.total} device${pagination.total !== 1 ? 's' : ''} total` }}
        </p>
      </div>

      <!-- Refresh button -->
      <button
        @click="fetchDevices()"
        :disabled="isLoading"
        class="self-start sm:self-auto flex items-center gap-2 px-4 py-2 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl text-sm font-bold text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800 transition-all disabled:opacity-50 shadow-sm"
      >
        <svg class="w-4 h-4 transition-transform" :class="{ 'animate-spin': isLoading }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M21 3v5h-5" />
        </svg>
        Refresh
      </button>
    </header>

    <!-- Filters -->
    <MonitoringDeviceFilters
      v-model="filters"
      :loading="isLoading"
      @apply="applyFilters"
      @reset="resetFilters"
    />

    <!-- Error state -->
    <BaseErrorMessage
      v-if="error && !isLoading"
      :message="error"
      @retry="fetchDevices()"
    />

    <!-- Devices table -->
    <MonitoringDeviceTable
      v-if="!error"
      :devices="devices"
      :pagination="{ ...pagination, offset: filters.offset }"
      :loading="isLoading"
      @page-change="goToPage"
    />

  </div>
</template>
