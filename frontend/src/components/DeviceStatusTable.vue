<script setup lang="ts">
import { useTimeAgo } from '@vueuse/core'
import BaseCard from './BaseCard.vue'
import BaseBadge from './BaseBadge.vue'
import type { DeviceStatus, Pagination } from '../types/dashboard'

defineProps<{
  devices: DeviceStatus[]
  pagination: Pagination
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'change-page', page: number): void
}>()

const formatRelativeTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return useTimeAgo(new Date(dateStr))
}
</script>

<template>
  <BaseCard title="Device Status" no-padding>
    <template #headerAction>
      <button class="text-sm font-medium text-blue-600 dark:text-blue-400 hover:underline">View All</button>
    </template>

    <!-- Desktop Table View -->
    <div class="hidden md:block overflow-x-auto">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-slate-50 dark:bg-slate-800/50 text-slate-500 dark:text-slate-400 text-xs uppercase tracking-wider">
            <th class="px-6 py-4 font-semibold">Device Name</th>
            <th class="px-6 py-4 font-semibold">Status</th>
            <th class="px-6 py-4 font-semibold">Engine</th>
            <th class="px-6 py-4 font-semibold">Fuel Level</th>
            <th class="px-6 py-4 font-semibold">Temperature</th>
            <th class="px-6 py-4 font-semibold text-right">Last Seen</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
          <template v-if="loading">
            <tr v-for="i in 5" :key="i" class="animate-pulse">
              <td v-for="j in 6" :key="j" class="px-6 py-4">
                <div class="h-4 bg-slate-200 dark:bg-slate-700 rounded w-full"></div>
              </td>
            </tr>
          </template>
          
          <template v-else-if="devices.length === 0">
            <tr>
              <td colspan="6" class="px-6 py-12 text-center text-slate-500 dark:text-slate-400">
                <div class="flex flex-col items-center">
                  <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="mb-3 opacity-20"><rect width="16" height="10" x="2" y="4" rx="2"/><line x1="10" x2="10" y1="14" y2="20"/><line x1="15" x2="15" y1="14" y2="20"/><line x1="2" x2="22" y1="20" y2="20"/></svg>
                  <p>No devices found</p>
                </div>
              </td>
            </tr>
          </template>

          <tr 
            v-for="device in devices" 
            :key="device.device_id"
            class="hover:bg-slate-50 dark:hover:bg-slate-800/30 transition-colors group cursor-pointer"
            @click="$router.push(`/devices/${device.device_id}`)"
          >
            <td class="px-6 py-4">
              <div class="font-medium text-slate-900 dark:text-slate-100 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                {{ device.device_name }}
              </div>
              <div class="text-xs text-slate-500 font-mono">{{ device.device_id }}</div>
            </td>
            <td class="px-6 py-4">
              <BaseBadge :variant="device.device_online ? 'success' : 'danger'" dot>
                {{ device.device_online ? 'Online' : 'Offline' }}
              </BaseBadge>
            </td>
            <td class="px-6 py-4">
              <BaseBadge :variant="device.engine_running ? 'primary' : 'neutral'">
                {{ device.engine_running ? 'Running' : 'Stopped' }}
              </BaseBadge>
            </td>
            <td class="px-6 py-4">
              <div class="flex items-center space-x-2">
                <div class="w-12 bg-slate-100 dark:bg-slate-800 h-1.5 rounded-full overflow-hidden">
                  <div 
                    class="h-full rounded-full transition-all duration-500"
                    :class="device.fuel_level < 20 ? 'bg-rose-500 shadow-[0_0_8px_rgba(244,63,94,0.5)]' : 'bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.5)]'"
                    :style="{ width: `${device.fuel_level}%` }"
                  ></div>
                </div>
                <span class="text-sm font-medium text-slate-700 dark:text-slate-300">{{ device.fuel_level }}%</span>
              </div>
            </td>
            <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400">
              <span :class="device.coolant_temperature > 90 ? 'text-rose-600 font-bold' : ''">
                {{ device.coolant_temperature }}°C
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-slate-500 dark:text-slate-400 text-right">
              {{ formatRelativeTime(device.last_seen_at) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Mobile Card View -->
    <div class="md:hidden divide-y divide-slate-100 dark:divide-slate-800">
      <div v-if="loading" class="p-4 space-y-4">
        <div v-for="i in 3" :key="i" class="animate-pulse space-y-3">
          <div class="h-5 bg-slate-200 dark:bg-slate-700 rounded w-1/2"></div>
          <div class="h-4 bg-slate-100 dark:bg-slate-800 rounded w-3/4"></div>
        </div>
      </div>
      
      <div 
        v-for="device in devices" 
        :key="device.device_id"
        class="p-4 hover:bg-slate-50 dark:hover:bg-slate-800/30 transition-colors"
        @click="$router.push(`/devices/${device.device_id}`)"
      >
        <div class="flex justify-between items-start mb-3">
          <div>
            <div class="font-bold text-slate-900 dark:text-slate-100">{{ device.device_name }}</div>
            <div class="text-xs text-slate-500 font-mono">{{ device.device_id }}</div>
          </div>
          <BaseBadge :variant="device.device_online ? 'success' : 'danger'" dot>
            {{ device.device_online ? 'Online' : 'Offline' }}
          </BaseBadge>
        </div>
        
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <p class="text-[10px] uppercase tracking-wider text-slate-500 mb-1">Engine</p>
            <BaseBadge :variant="device.engine_running ? 'primary' : 'neutral'" size="sm">
              {{ device.engine_running ? 'Running' : 'Stopped' }}
            </BaseBadge>
          </div>
          <div>
            <p class="text-[10px] uppercase tracking-wider text-slate-500 mb-1">Fuel</p>
            <div class="flex items-center space-x-2">
              <div class="w-full max-w-[60px] bg-slate-100 dark:bg-slate-800 h-1.5 rounded-full overflow-hidden">
                <div 
                  class="h-full rounded-full"
                  :class="device.fuel_level < 20 ? 'bg-rose-500' : 'bg-blue-500'"
                  :style="{ width: `${device.fuel_level}%` }"
                ></div>
              </div>
              <span class="font-medium">{{ device.fuel_level }}%</span>
            </div>
          </div>
          <div>
            <p class="text-[10px] uppercase tracking-wider text-slate-500 mb-1">Temp</p>
            <span class="font-medium" :class="device.coolant_temperature > 90 ? 'text-rose-600' : ''">
              {{ device.coolant_temperature }}°C
            </span>
          </div>
          <div class="text-right">
            <p class="text-[10px] uppercase tracking-wider text-slate-500 mb-1">Last Seen</p>
            <span class="text-slate-500">{{ formatRelativeTime(device.last_seen_at) }}</span>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <span class="text-xs sm:text-sm text-slate-500">Showing {{ devices.length }} of {{ pagination.total }} devices</span>
        <div class="flex space-x-2">
          <button 
            class="px-3 py-1.5 border border-slate-200 dark:border-slate-700 rounded-lg text-xs sm:text-sm font-medium hover:bg-slate-50 dark:hover:bg-slate-800 disabled:opacity-50 transition-colors" 
            :disabled="pagination.page === 1 || loading"
            @click="emit('change-page', pagination.page - 1)"
          >
            Previous
          </button>
          <button 
            class="px-3 py-1.5 border border-slate-200 dark:border-slate-700 rounded-lg text-xs sm:text-sm font-medium hover:bg-slate-50 dark:hover:bg-slate-800 disabled:opacity-50 transition-colors" 
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
