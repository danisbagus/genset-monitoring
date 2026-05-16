<script setup lang="ts">
import { computed } from 'vue'
import BaseCard from './BaseCard.vue'
import type { DashboardSummary } from '../types/dashboard'

const props = defineProps<{
  summary: DashboardSummary
}>()

const stats = computed(() => [
  {
    title: 'Total Devices',
    value: props.summary.total_devices,
    icon: 'devices',
    color: 'primary',
    bg: 'bg-blue-500/10',
    text: 'text-blue-500'
  },
  {
    title: 'Online Devices',
    value: props.summary.online_devices,
    icon: 'online',
    color: 'success',
    bg: 'bg-emerald-500/10',
    text: 'text-emerald-500'
  },
  {
    title: 'Offline Devices',
    value: props.summary.offline_devices,
    icon: 'offline',
    color: 'danger',
    bg: 'bg-rose-500/10',
    text: 'text-rose-500'
  },
  {
    title: 'Running Engines',
    value: props.summary.running_engines,
    icon: 'engine',
    color: 'info',
    bg: 'bg-indigo-500/10',
    text: 'text-indigo-500'
  },
  {
    title: 'Warning Alerts',
    value: props.summary.warning_alerts,
    icon: 'warning',
    color: 'warning',
    bg: 'bg-amber-500/10',
    text: 'text-amber-500'
  },
  {
    title: 'Critical Alerts',
    value: props.summary.critical_alerts,
    icon: 'critical',
    color: 'danger',
    bg: 'bg-red-500/10',
    text: 'text-red-600'
  }
])
</script>

<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4">
    <BaseCard v-for="stat in stats" :key="stat.title" class="!p-4">
      <div class="flex items-center justify-between mb-3">
        <div :class="['p-2 rounded-xl', stat.bg, stat.text]">
          <!-- Icons mapping -->
          <svg v-if="stat.icon === 'devices'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="16" height="10" x="2" y="4" rx="2"/><line x1="10" x2="10" y1="14" y2="20"/><line x1="15" x2="15" y1="14" y2="20"/><line x1="2" x2="22" y1="20" y2="20"/></svg>
          <svg v-else-if="stat.icon === 'online'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12s2.5 4.5 7 4.5 7-4.5 7-4.5"/><path d="M12 5V3"/><path d="M12 21v-2"/><path d="m17 7 1.5-1.5"/><path d="m6.5 17.5 1.5-1.5"/><path d="m17 17 1.5 1.5"/><path d="m6.5 6.5 1.5 1.5"/></svg>
          <svg v-else-if="stat.icon === 'offline'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 1l22 22"/><path d="M16.7 11.3c.1.2.2.4.3.7 0 4.5-7 4.5-7 4.5"/><path d="M12 5V3"/><path d="M12 21v-2"/><path d="m17 7 1.5-1.5"/><path d="m6.5 17.5 1.5-1.5"/><path d="m17 17 1.5 1.5"/><path d="m6.5 6.5 1.5 1.5"/></svg>
          <svg v-else-if="stat.icon === 'engine'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v4"/><path d="m4.93 4.93 2.83 2.83"/><path d="M2 12h4"/><path d="m4.93 19.07 2.83-2.83"/><path d="M12 22v-4"/><path d="m19.07 19.07-2.83-2.83"/><path d="M22 12h-4"/><path d="m19.07 4.93-2.83 2.83"/><circle cx="12" cy="12" r="4"/></svg>
          <svg v-else-if="stat.icon === 'warning'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" x2="12" y1="9" y2="13"/><line x1="12" x2="12.01" y1="17" y2="17"/></svg>
          <svg v-else-if="stat.icon === 'critical'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" x2="12" y1="8" y2="12"/><line x1="12" x2="12.01" y1="16" y2="16"/></svg>
        </div>
        <div :class="['w-2 h-2 rounded-full', stat.value > 0 ? 'animate-pulse ' + stat.bg.replace('/10', '') : 'bg-slate-300 dark:bg-slate-700']"></div>
      </div>
      <div>
        <p class="text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider">{{ stat.title }}</p>
        <h4 class="text-2xl font-bold text-slate-800 dark:text-slate-100 mt-1">{{ stat.value }}</h4>
      </div>
    </BaseCard>
  </div>
</template>
