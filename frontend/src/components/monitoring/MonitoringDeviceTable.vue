<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import MonitoringStatusBadge from './MonitoringStatusBadge.vue'
import type { MonitoringDeviceItem, MonitoringPagination } from '@/types/monitoring'

const props = defineProps<{
  devices: MonitoringDeviceItem[]
  pagination: MonitoringPagination
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'page-change', page: number): void
}>()

const router = useRouter()

const currentPage = computed(() =>
  Math.floor((props.pagination.page - 1)) + 1
)

const totalPages = computed(() =>
  Math.ceil(props.pagination.total / props.pagination.limit) || 1
)

const goTo = (deviceId: string) => {
  router.push(`/monitoring/devices/${deviceId}`)
}

const isOnline = (device: MonitoringDeviceItem): boolean => {
  if (!device.last_seen_at) return false
  return (Date.now() - new Date(device.last_seen_at).getTime()) < 5 * 60 * 1000
}

const formatRelativeTime = (iso: string | null | undefined): string => {
  if (!iso) return '—'
  const diff = Math.floor((Date.now() - new Date(iso).getTime()) / 1000)
  if (diff < 60) return `${diff}s ago`
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return `${Math.floor(diff / 86400)}d ago`
}

const formatNumber = (v: number | undefined, unit = '', decimals = 1): string => {
  if (v === undefined || v === null) return '—'
  return `${v.toFixed(decimals)}${unit}`
}

const signalBars = (gsm: number): number => {
  if (gsm >= 20) return 4
  if (gsm >= 14) return 3
  if (gsm >= 8) return 2
  if (gsm >= 2) return 1
  return 0
}
</script>

<template>
  <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden">

    <!-- Loading skeleton -->
    <template v-if="loading">
      <div class="divide-y divide-slate-100 dark:divide-slate-800">
        <div v-for="i in 6" :key="i" class="flex items-center gap-4 px-5 py-4 animate-pulse">
          <div class="w-9 h-9 rounded-xl bg-slate-200 dark:bg-slate-800 shrink-0"></div>
          <div class="flex-1 space-y-2">
            <div class="h-3 w-40 bg-slate-200 dark:bg-slate-800 rounded-full"></div>
            <div class="h-2.5 w-24 bg-slate-100 dark:bg-slate-700 rounded-full"></div>
          </div>
          <div class="hidden md:flex gap-3">
            <div class="h-6 w-16 bg-slate-100 dark:bg-slate-700 rounded-lg"></div>
            <div class="h-6 w-16 bg-slate-100 dark:bg-slate-700 rounded-lg"></div>
          </div>
          <div class="hidden lg:flex gap-6">
            <div class="h-3 w-20 bg-slate-100 dark:bg-slate-700 rounded-full"></div>
            <div class="h-3 w-20 bg-slate-100 dark:bg-slate-700 rounded-full"></div>
          </div>
        </div>
      </div>
    </template>

    <!-- Empty state -->
    <template v-else-if="!loading && devices.length === 0">
      <div class="py-20 flex flex-col items-center gap-4 text-center">
        <div class="w-16 h-16 rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-slate-400">
          <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714c0 .597.237 1.17.659 1.591L19.8 15.3M14.25 3.104c.251.023.501.05.75.082M19.8 15.3l-1.57.393A9.065 9.065 0 0112 15a9.065 9.065 0 00-6.23-.693L5 14.5m14.8.8l1.402 1.402c1 1 .03 2.798-1.414 2.798H4.213c-1.444 0-2.414-1.798-1.414-2.798L4.2 15.3" />
          </svg>
        </div>
        <div>
          <p class="font-bold text-slate-700 dark:text-slate-300 text-base">No devices found</p>
          <p class="text-sm text-slate-400 mt-1">Try adjusting your filters or search query</p>
        </div>
      </div>
    </template>

    <!-- Table -->
    <template v-else>
      <!-- Table header -->
      <div class="px-5 py-3 border-b border-slate-100 dark:border-slate-800 bg-slate-50/70 dark:bg-slate-900/70 hidden md:grid grid-cols-12 gap-4 text-[10px] font-bold uppercase tracking-widest text-slate-400">
        <div class="col-span-4">Device</div>
        <div class="col-span-2">Status</div>
        <div class="col-span-2">Engine</div>
        <div class="col-span-2">Telemetry</div>
        <div class="col-span-2 text-right">Last Seen</div>
      </div>

      <!-- Rows -->
      <div class="divide-y divide-slate-100 dark:divide-slate-800">
        <div
          v-for="device in devices"
          :key="device.device_id"
          @click="goTo(device.device_id)"
          class="px-5 py-4 cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-all duration-150 group md:grid md:grid-cols-12 md:gap-4 md:items-center flex flex-col gap-3"
          :id="`device-row-${device.device_id}`"
        >
          <!-- Col 1: Device info (4/12) -->
          <div class="md:col-span-4 flex items-center gap-3 min-w-0">
            <!-- Icon with online indicator -->
            <div class="relative shrink-0">
              <div
                class="w-9 h-9 rounded-xl flex items-center justify-center text-white shadow-sm transition-all duration-300"
                :class="isOnline(device) ? 'bg-blue-600 shadow-blue-600/20' : 'bg-slate-400 dark:bg-slate-600'"
              >
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
                </svg>
              </div>
              <!-- Online pulse dot -->
              <span v-if="isOnline(device)" class="absolute -top-0.5 -right-0.5 flex h-3 w-3">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500 border-2 border-white dark:border-slate-900"></span>
              </span>
            </div>

            <div class="min-w-0">
              <p class="font-bold text-sm text-slate-900 dark:text-white truncate group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                {{ device.device_name }}
              </p>
              <p class="text-[11px] text-slate-400 font-mono truncate">{{ device.device_code }} · {{ device.serial_number }}</p>
            </div>
          </div>

          <!-- Col 2: Online / Device status (2/12) -->
          <div class="md:col-span-2 flex flex-wrap gap-1.5">
            <MonitoringStatusBadge type="online" :value="isOnline(device)" />
            <MonitoringStatusBadge type="device_status" :value="device.device_status" />
          </div>

          <!-- Col 3: Engine state + RPM (2/12) -->
          <div class="md:col-span-2 flex flex-col gap-1">
            <MonitoringStatusBadge type="engine_running" :value="device.engine_running" />
            <span v-if="device.engine_running && device.speed" class="text-[11px] font-mono text-slate-500 dark:text-slate-400">
              {{ device.speed }} RPM
            </span>
          </div>

          <!-- Col 4: Key telemetry snapshot (2/12) -->
          <div class="md:col-span-2 hidden lg:flex flex-col gap-1 text-[11px] font-mono text-slate-500 dark:text-slate-400">
            <span>⚡ {{ formatNumber(device.total_va / 1000, ' kVA') }}</span>
            <span>🌡 {{ formatNumber(device.coolant_temperature, ' °C') }}</span>
            <span>⛽ {{ formatNumber(device.fuel_level, '%') }}</span>
          </div>

          <!-- Col 5: Last seen (2/12) -->
          <div class="md:col-span-2 flex items-center justify-between md:justify-end gap-2">
            <span class="text-xs text-slate-400 font-mono tabular-nums">{{ formatRelativeTime(device.last_seen_at) }}</span>
            <!-- GSM signal -->
            <div class="flex items-end gap-0.5 h-4 shrink-0" :title="`GSM Signal: ${device.gsm_signal}`">
              <div v-for="bar in 4" :key="bar"
                class="w-1 rounded-sm transition-all"
                :class="bar <= signalBars(device.gsm_signal) ? 'bg-emerald-500' : 'bg-slate-200 dark:bg-slate-700'"
                :style="`height: ${bar * 25}%`"
              ></div>
            </div>
            <svg class="w-4 h-4 text-slate-300 dark:text-slate-600 group-hover:text-blue-400 transition-colors shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div
        v-if="pagination.total > pagination.limit"
        class="px-5 py-3 border-t border-slate-100 dark:border-slate-800 bg-slate-50/70 dark:bg-slate-900/70 flex items-center justify-between gap-4"
      >
        <span class="text-xs text-slate-400">
          Showing {{ Math.min(pagination.offset ?? 0, pagination.total) + 1 }}–{{ Math.min((pagination.offset ?? 0) + pagination.limit, pagination.total) }} of {{ pagination.total }} devices
        </span>
        <div class="flex gap-1.5">
          <button
            :disabled="currentPage <= 1"
            @click="emit('page-change', currentPage - 1)"
            class="h-8 w-8 flex items-center justify-center rounded-lg border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 disabled:opacity-30 disabled:cursor-not-allowed transition-all text-sm"
          >‹</button>

          <button
            v-for="p in totalPages" :key="p"
            v-show="Math.abs(p - currentPage) <= 2"
            @click="emit('page-change', p)"
            class="h-8 w-8 flex items-center justify-center rounded-lg border text-xs font-bold transition-all"
            :class="p === currentPage
              ? 'bg-blue-600 border-blue-600 text-white shadow-sm shadow-blue-600/20'
              : 'border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800'"
          >{{ p }}</button>

          <button
            :disabled="currentPage >= totalPages"
            @click="emit('page-change', currentPage + 1)"
            class="h-8 w-8 flex items-center justify-center rounded-lg border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 disabled:opacity-30 disabled:cursor-not-allowed transition-all text-sm"
          >›</button>
        </div>
      </div>
    </template>

  </div>
</template>
