<script setup lang="ts">
import MonitoringStatusBadge from './MonitoringStatusBadge.vue'
import type { MonitoringDeviceInfoOutput, MonitoringLatestStateOutput } from '@/types/monitoring'

const props = defineProps<{
  deviceInfo: MonitoringDeviceInfoOutput
  latestState: MonitoringLatestStateOutput
}>()

const isOnline = (): boolean => {
  if (!props.latestState?.last_seen_at) return false
  return (Date.now() - new Date(props.latestState.last_seen_at).getTime()) < 5 * 60 * 1000
}

const formatDate = (iso: string | null | undefined): string => {
  if (!iso) return '—'
  return new Date(iso).toLocaleString('en-GB', { dateStyle: 'short', timeStyle: 'short' })
}

const fields: Array<{ label: string; key: keyof MonitoringDeviceInfoOutput }> = [
  { label: 'Device Name', key: 'name' },
  { label: 'Device Code', key: 'device_code' },
  { label: 'Serial Number', key: 'serial_number' },
  { label: 'Engine ID', key: 'engine_id' },
  { label: 'Firmware', key: 'firmware_version' },
  { label: 'GSM Number', key: 'gsm_number' },
]
</script>

<template>
  <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 flex items-center justify-between">
      <h3 class="font-bold text-slate-900 dark:text-slate-100 flex items-center gap-2 text-sm">
        <svg class="w-4 h-4 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
        </svg>
        Device Information
      </h3>
      <div class="flex items-center gap-2">
        <MonitoringStatusBadge type="online" :value="isOnline()" />
        <MonitoringStatusBadge type="device_status" :value="deviceInfo.status" />
      </div>
    </div>

    <!-- Body -->
    <div class="p-6 grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4">
      <div v-for="field in fields" :key="field.key">
        <p class="text-[10px] font-bold uppercase tracking-widest text-slate-400 mb-0.5">{{ field.label }}</p>
        <p class="text-sm font-semibold text-slate-800 dark:text-slate-200 font-mono">
          {{ (deviceInfo[field.key] as string) || '—' }}
        </p>
      </div>

      <!-- Last updated -->
      <div>
        <p class="text-[10px] font-bold uppercase tracking-widest text-slate-400 mb-0.5">Last Seen</p>
        <p class="text-sm font-semibold text-slate-800 dark:text-slate-200 font-mono">{{ formatDate(latestState.last_seen_at) }}</p>
      </div>

      <div>
        <p class="text-[10px] font-bold uppercase tracking-widest text-slate-400 mb-0.5">Last Online</p>
        <p class="text-sm font-semibold text-slate-800 dark:text-slate-200 font-mono">{{ formatDate(latestState.last_online_at) }}</p>
      </div>

      <div>
        <p class="text-[10px] font-bold uppercase tracking-widest text-slate-400 mb-0.5">Telemetry Recorded At</p>
        <p class="text-sm font-semibold text-slate-800 dark:text-slate-200 font-mono">{{ formatDate(latestState.telemetry_recorded_at) }}</p>
      </div>
    </div>
  </div>
</template>
