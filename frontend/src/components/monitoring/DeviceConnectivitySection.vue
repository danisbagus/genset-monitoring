<script setup lang="ts">
import type { MonitoringConnectivityOutput } from '@/types/monitoring'

const props = defineProps<{
  connectivity: MonitoringConnectivityOutput
}>()

interface ConnMetric {
  label: string
  key: keyof MonitoringConnectivityOutput
  type: 'bool' | 'signal' | 'time'
}

const metrics: ConnMetric[] = [
  { label: 'Server', key: 'server_connected', type: 'bool' },
  { label: 'GPS', key: 'gps_connected', type: 'bool' },
  { label: 'CAN Bus', key: 'can_connected', type: 'bool' },
  { label: 'RS-485', key: 'rs485_connected', type: 'bool' },
  { label: 'SD Card', key: 'sd_card_ok', type: 'bool' },
]

const signalBars = (gsm: number): number => {
  if (gsm >= 20) return 4
  if (gsm >= 14) return 3
  if (gsm >= 8) return 2
  if (gsm >= 2) return 1
  return 0
}

const signalLabel = (gsm: number): string => {
  const bars = signalBars(gsm)
  const labels = ['No Signal', 'Weak', 'Fair', 'Good', 'Excellent']
  return labels[bars]
}

const formatDate = (iso: string | null | undefined): string => {
  if (!iso) return '—'
  return new Date(iso).toLocaleString('en-GB', { dateStyle: 'short', timeStyle: 'medium' })
}
</script>

<template>
  <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50">
      <h3 class="font-bold text-slate-900 dark:text-slate-100 flex items-center gap-2 text-sm">
        <svg class="w-4 h-4 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0" />
        </svg>
        Connectivity
      </h3>
    </div>

    <div class="p-6 space-y-5">
      <!-- Connection indicators -->
      <div class="grid grid-cols-2 sm:grid-cols-5 gap-3">
        <div
          v-for="m in metrics" :key="m.label"
          class="flex flex-col items-center gap-2 p-3 rounded-xl border transition-all"
          :class="connectivity[m.key]
            ? 'bg-emerald-50 border-emerald-200 dark:bg-emerald-900/10 dark:border-emerald-800/40'
            : 'bg-slate-50 border-slate-200 dark:bg-slate-950 dark:border-slate-800'"
        >
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center transition-all"
            :class="connectivity[m.key]
              ? 'bg-emerald-500 shadow-lg shadow-emerald-500/30'
              : 'bg-slate-300 dark:bg-slate-600'"
          >
            <svg class="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path v-if="connectivity[m.key]" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
              <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
          <span class="text-[11px] font-bold text-center"
            :class="connectivity[m.key] ? 'text-emerald-700 dark:text-emerald-400' : 'text-slate-400'"
          >{{ m.label }}</span>
        </div>
      </div>

      <!-- GSM Signal bar -->
      <div class="flex items-center gap-4 px-4 py-3 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
        <div class="flex items-end gap-1 h-8">
          <div v-for="bar in 4" :key="bar"
            class="w-2.5 rounded-sm transition-all"
            :class="bar <= signalBars(connectivity.gsm_signal) ? 'bg-emerald-500' : 'bg-slate-200 dark:bg-slate-700'"
            :style="`height: ${bar * 25}%`"
          ></div>
        </div>
        <div>
          <p class="text-xs font-bold text-slate-700 dark:text-slate-300">GSM Signal: {{ signalLabel(connectivity.gsm_signal) }}</p>
          <p class="text-[11px] text-slate-400 font-mono">{{ connectivity.gsm_signal }} dBm</p>
        </div>
      </div>

      <!-- Last heartbeat -->
      <div class="flex items-center justify-between px-4 py-3 bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 rounded-xl">
        <div>
          <p class="text-[10px] font-bold uppercase tracking-widest text-slate-400 mb-0.5">Last Heartbeat</p>
          <p class="text-sm font-mono font-semibold text-slate-700 dark:text-slate-300">{{ formatDate(connectivity.last_seen) }}</p>
        </div>
        <svg class="w-5 h-5 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
        </svg>
      </div>
    </div>
  </div>
</template>
