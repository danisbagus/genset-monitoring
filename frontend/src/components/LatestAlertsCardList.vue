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

const emit = defineEmits<{
  (e: 'change-page', page: number): void
}>()

// Vanilla relative time helper to run perfectly inside loops without violating composable rules
const formatRelativeTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHr = Math.floor(diffMin / 60)
  const diffDays = Math.floor(diffHr / 24)

  if (diffSec < 5) return 'just now'
  if (diffSec < 60) return `${diffSec}s ago`
  if (diffMin < 60) return `${diffMin}m ago`
  if (diffHr < 24) return `${diffHr}h ago`
  return `${diffDays}d ago`
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

// Parses messages to separate alert type/title from the actual details
const getAlertTitleAndDesc = (alert: Alert) => {
  const parts = alert.message.split(':')
  if (parts.length > 1) {
    return {
      title: parts[0].trim(),
      desc: parts.slice(1).join(':').trim()
    }
  }

  // Fallbacks if no colon is present
  let title = 'Genset Alarm'
  if (alert.severity === 'critical') {
    title = 'Critical Alarm'
  } else if (alert.severity === 'warning') {
    title = 'System Warning'
  } else if (alert.severity === 'info') {
    title = 'System Notification'
  }

  return {
    title,
    desc: alert.message
  }
}
</script>

<template>
  <BaseCard title="Latest Alerts" no-padding>
    <!-- Scrollable Vertical Cards Container -->
    <div class="p-4 sm:p-5 space-y-3.5 max-h-[600px] overflow-y-auto custom-scrollbar">
      <!-- Loading Skeleton State -->
      <template v-if="loading">
        <div 
          v-for="i in 4" 
          :key="i" 
          class="p-4 rounded-xl border border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 animate-pulse space-y-3"
        >
          <div class="flex justify-between items-center">
            <div class="h-4 bg-slate-200 dark:bg-slate-700 rounded w-1/3"></div>
            <div class="h-5 bg-slate-200 dark:bg-slate-700 rounded w-16"></div>
          </div>
          <div class="h-4 bg-slate-100 dark:bg-slate-800 rounded w-full"></div>
          <div class="flex justify-between items-center pt-1">
            <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-20"></div>
            <div class="h-3 bg-slate-200 dark:bg-slate-700 rounded w-24"></div>
          </div>
        </div>
      </template>

      <!-- Empty State -->
      <template v-else-if="alerts.length === 0">
        <div class="flex flex-col items-center justify-center py-16 text-center text-slate-500 dark:text-slate-400">
          <div class="p-4 rounded-full bg-emerald-50 dark:bg-emerald-950/20 text-emerald-500 mb-4 animate-bounce duration-1000">
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
              <polyline points="22 4 12 14.01 9 11.01" />
            </svg>
          </div>
          <h4 class="font-bold text-slate-900 dark:text-slate-100 mb-1">All Systems Normal</h4>
          <p class="text-xs max-w-[200px] mx-auto text-slate-500 dark:text-slate-400">No active alerts or warnings recorded at this time.</p>
        </div>
      </template>

      <!-- Active Alerts Card List -->
      <template v-else>
        <div 
          v-for="alert in alerts" 
          :key="alert.alert_id"
          class="relative flex flex-col p-4 bg-white dark:bg-slate-900 border-l-4 rounded-r-xl border-y border-r border-slate-200 dark:border-slate-800 shadow-sm hover:shadow-md transition-all duration-300 gap-3 group cursor-pointer hover:-translate-y-0.5"
          :class="{
            'border-l-rose-500 hover:border-l-rose-600 bg-rose-50/5 dark:bg-rose-950/5': alert.severity === 'critical',
            'border-l-amber-500 hover:border-l-amber-600 bg-amber-50/5 dark:bg-amber-950/5': alert.severity === 'warning',
            'border-l-blue-500 hover:border-l-blue-600 bg-blue-50/5 dark:bg-blue-950/5': alert.severity === 'info'
          }"
        >
          <!-- Card Header: Title/Type and Severity Badge -->
          <div class="flex items-start justify-between gap-2">
            <div class="flex items-start gap-2.5">
              <!-- Severity Icon Indicator -->
              <div class="mt-0.5 shrink-0">
                <span v-if="alert.severity === 'critical'">
                  <svg class="h-5 w-5 text-rose-500 dark:text-rose-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                </span>
                <span v-else-if="alert.severity === 'warning'">
                  <svg class="h-5 w-5 text-amber-500 dark:text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </span>
                <span v-else>
                  <svg class="h-5 w-5 text-blue-500 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </span>
              </div>
              <div>
                <h4 class="font-bold text-sm text-slate-800 dark:text-slate-200 tracking-tight group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                  {{ getAlertTitleAndDesc(alert).title }}
                </h4>
              </div>
            </div>
            
            <BaseBadge :variant="getSeverityVariant(alert.severity)" size="sm">
              {{ alert.severity }}
            </BaseBadge>
          </div>

          <!-- Message / Description -->
          <p class="text-xs text-slate-600 dark:text-slate-400 leading-relaxed font-normal pl-7">
            {{ getAlertTitleAndDesc(alert).desc }}
          </p>

          <!-- Card Footer: Device Info, Status and Time -->
          <div class="flex items-center justify-between border-t border-slate-100 dark:border-slate-800/80 pt-2.5 mt-1 text-[11px] font-medium text-slate-500 dark:text-slate-400">
            <div class="flex items-center gap-1.5 hover:text-slate-800 dark:hover:text-slate-200 transition-colors">
              <svg class="h-3.5 w-3.5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
              <span>{{ alert.device_name }}</span>
            </div>
            
            <div class="flex items-center gap-2">
              <!-- Realtime Status Dot Indicator -->
              <span class="flex items-center gap-1">
                <span 
                  class="h-1.5 w-1.5 rounded-full" 
                  :class="[
                    alert.acknowledged 
                      ? 'bg-emerald-500' 
                      : 'bg-amber-500 animate-pulse shadow-[0_0_8px_rgba(245,158,11,0.7)]'
                  ]"
                ></span>
                <span class="text-[9px] text-slate-400 dark:text-slate-500 font-bold tracking-wider">
                  {{ alert.acknowledged ? 'ACK' : 'NEW' }}
                </span>
              </span>
              
              <span class="h-3 w-px bg-slate-200 dark:bg-slate-800"></span>

              <!-- Relative Timestamp (Hover for Absolute) -->
              <div 
                class="flex items-center gap-1.5 cursor-help hover:text-slate-700 dark:hover:text-slate-300" 
                :title="formatDate(alert.created_at)"
              >
                <svg class="h-3.5 w-3.5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span class="font-semibold">{{ formatRelativeTime(alert.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- Pagination Footer -->
    <template #footer>
      <div class="flex items-center justify-between">
        <span class="text-xs sm:text-sm text-slate-500 dark:text-slate-400 font-medium">
          Showing {{ alerts.length }} of {{ pagination.total }} alerts
        </span>
        <div class="flex space-x-2">
          <button 
            class="px-3 py-1.5 border border-slate-200 dark:border-slate-700 rounded-lg text-xs sm:text-sm font-medium hover:bg-slate-50 dark:hover:bg-slate-800 dark:text-slate-300 disabled:opacity-50 disabled:hover:bg-transparent transition-colors" 
            :disabled="pagination.page === 1 || loading"
            @click="emit('change-page', pagination.page - 1)"
          >
            Previous
          </button>
          <button 
            class="px-3 py-1.5 border border-slate-200 dark:border-slate-700 rounded-lg text-xs sm:text-sm font-medium hover:bg-slate-50 dark:hover:bg-slate-800 dark:text-slate-300 disabled:opacity-50 disabled:hover:bg-transparent transition-colors" 
            :disabled="pagination.page * pagination.limit >= pagination.total || loading"
            @click="emit('change-page', pagination.page + 1)"
          >
            Next
          </button>
        </div>
      </div>
    </template>
  </BaseCard>
</template>

<style scoped>
/* High-quality scrollbar experience matching dashboard theme */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}
.dark .custom-scrollbar::-webkit-scrollbar-thumb {
  background: #334155;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
.dark .custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #475569;
}
</style>
