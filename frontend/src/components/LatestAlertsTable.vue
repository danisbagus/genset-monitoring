<script setup lang="ts">
import BaseCard from './BaseCard.vue'
import BaseBadge from './BaseBadge.vue'
import type { Alert, Pagination } from '../types/dashboard'

defineProps<{
  alerts: Alert[]
  pagination: Pagination
  loading?: boolean
}>()

const getSeverityVariant = (severity: string) => {
  switch (severity.toLowerCase()) {
    case 'critical': return 'danger'
    case 'warning': return 'warning'
    case 'info': return 'info'
    default: return 'neutral'
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}
</script>

<template>
  <BaseCard title="Latest Alerts" no-padding>
    <div class="overflow-x-auto">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-slate-50 dark:bg-slate-800/50 text-slate-500 dark:text-slate-400 text-xs uppercase tracking-wider">
            <th class="px-6 py-4 font-semibold">Device</th>
            <th class="px-6 py-4 font-semibold">Severity</th>
            <th class="px-6 py-4 font-semibold">Message</th>
            <th class="px-6 py-4 font-semibold">Status</th>
            <th class="px-6 py-4 font-semibold">Time</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
          <template v-if="loading">
            <tr v-for="i in 5" :key="i" class="animate-pulse">
              <td v-for="j in 5" :key="j" class="px-6 py-4">
                <div class="h-4 bg-slate-200 dark:bg-slate-700 rounded w-full"></div>
              </td>
            </tr>
          </template>

          <template v-else-if="alerts.length === 0">
            <tr>
              <td colspan="5" class="px-6 py-12 text-center text-slate-500 dark:text-slate-400">
                <div class="flex flex-col items-center">
                  <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="mb-3 opacity-20"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" x2="12" y1="9" y2="13"/><line x1="12" x2="12.01" y1="17" y2="17"/></svg>
                  <p>No alerts recorded</p>
                </div>
              </td>
            </tr>
          </template>

          <tr 
            v-for="alert in alerts" 
            :key="alert.alert_id"
            class="hover:bg-slate-50 dark:hover:bg-slate-800/30 transition-colors group"
            :class="{'bg-rose-50/30 dark:bg-rose-900/10': alert.severity === 'critical' && !alert.acknowledged}"
          >
            <td class="px-6 py-4 font-medium text-slate-800 dark:text-slate-200">
              {{ alert.device_name }}
            </td>
            <td class="px-6 py-4">
              <BaseBadge :variant="getSeverityVariant(alert.severity)">
                {{ alert.severity }}
              </BaseBadge>
            </td>
            <td class="px-6 py-4 text-sm text-slate-700 dark:text-slate-300">
              {{ alert.message }}
            </td>
            <td class="px-6 py-4">
              <span 
                class="text-xs font-medium"
                :class="alert.acknowledged ? 'text-emerald-600 dark:text-emerald-400' : 'text-amber-600 dark:text-amber-400'"
              >
                {{ alert.acknowledged ? 'Acknowledged' : 'Pending' }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-500 dark:text-slate-400">
              {{ formatDate(alert.created_at) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <div class="flex items-center justify-between">
        <span class="text-sm text-slate-500">Showing {{ alerts.length }} alerts</span>
        <button class="text-sm font-medium text-blue-600 dark:text-blue-400 hover:underline">View All History</button>
      </div>
    </template>
  </BaseCard>
</template>
