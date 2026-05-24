<script setup lang="ts">
import { ref, watch } from 'vue'
import type { MonitoringDeviceFilters } from '@/types/monitoring'

const props = defineProps<{
  modelValue: MonitoringDeviceFilters
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: MonitoringDeviceFilters): void
  (e: 'apply'): void
  (e: 'reset'): void
}>()

// Local copy for controlled form
const local = ref<MonitoringDeviceFilters>({ ...props.modelValue })

watch(() => props.modelValue, (v) => { local.value = { ...v } }, { deep: true })

// Debounced search
let searchTimer: ReturnType<typeof setTimeout>
const onSearchInput = (e: Event) => {
  local.value.search = (e.target as HTMLInputElement).value
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    emit('update:modelValue', { ...local.value })
    emit('apply')
  }, 400)
}

const onFilterChange = () => {
  emit('update:modelValue', { ...local.value })
  emit('apply')
}

const onReset = () => {
  emit('reset')
}

const statusOptions = [
  { label: 'All Status', value: '' },
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
  { label: 'Maintenance', value: 'maintenance' },
]

const onlineOptions = [
  { label: 'All Connectivity', value: '' },
  { label: 'Online', value: 'true' },
  { label: 'Offline', value: 'false' },
]

const engineOptions = [
  { label: 'All Engine States', value: '' },
  { label: 'Running', value: 'true' },
  { label: 'Stopped', value: 'false' },
]

const sortByOptions = [
  { label: 'Last Seen', value: 'last_seen_at' },
  { label: 'Last Updated', value: 'updated_at' },
  { label: 'Telemetry Time', value: 'telemetry_recorded_at' },
  { label: 'Name', value: 'name' },
]

// Determine if any filter is active (besides defaults)
const hasActiveFilters = () =>
  !!local.value.search ||
  local.value.online !== '' ||
  local.value.engine_running !== '' ||
  local.value.status !== ''
</script>

<template>
  <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 shadow-sm">
    <div class="flex flex-wrap gap-3 items-center">

      <!-- Search -->
      <div class="relative flex-1 min-w-[200px]">
        <div class="absolute inset-y-0 left-3 flex items-center pointer-events-none text-slate-400">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <input
          id="monitoring-search"
          type="text"
          :value="local.search"
          @input="onSearchInput"
          placeholder="Search device name, code, serial..."
          class="w-full pl-9 pr-4 py-2 text-sm bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-900 dark:text-white placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/40 focus:border-blue-500 transition-all"
        />
      </div>

      <!-- Online filter -->
      <select
        id="monitoring-online-filter"
        v-model="local.online"
        @change="onFilterChange"
        class="h-9 px-3 py-1.5 text-sm bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-700 dark:text-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-500/40 focus:border-blue-500 transition-all min-w-[160px] cursor-pointer"
      >
        <option v-for="o in onlineOptions" :key="String(o.value)" :value="o.value">{{ o.label }}</option>
      </select>

      <!-- Engine state filter -->
      <select
        id="monitoring-engine-filter"
        v-model="local.engine_running"
        @change="onFilterChange"
        class="h-9 px-3 py-1.5 text-sm bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-700 dark:text-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-500/40 focus:border-blue-500 transition-all min-w-[160px] cursor-pointer"
      >
        <option v-for="o in engineOptions" :key="String(o.value)" :value="o.value">{{ o.label }}</option>
      </select>

      <!-- Device status filter -->
      <select
        id="monitoring-status-filter"
        v-model="local.status"
        @change="onFilterChange"
        class="h-9 px-3 py-1.5 text-sm bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-700 dark:text-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-500/40 focus:border-blue-500 transition-all min-w-[140px] cursor-pointer"
      >
        <option v-for="o in statusOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
      </select>

      <!-- Sort by -->
      <select
        id="monitoring-sort-filter"
        v-model="local.sort_by"
        @change="onFilterChange"
        class="h-9 px-3 py-1.5 text-sm bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-700 dark:text-slate-300 focus:outline-none focus:ring-2 focus:ring-blue-500/40 focus:border-blue-500 transition-all min-w-[160px] cursor-pointer"
      >
        <option v-for="o in sortByOptions" :key="o.value" :value="o.value">Sort: {{ o.label }}</option>
      </select>

      <!-- Sort direction toggle -->
      <button
        id="monitoring-sort-dir"
        @click="() => { local.sort_order = local.sort_order === 'desc' ? 'asc' : 'desc'; onFilterChange() }"
        class="h-9 w-9 flex items-center justify-center bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-all"
        :title="local.sort_order === 'desc' ? 'Descending' : 'Ascending'"
      >
        <svg class="w-4 h-4 transition-transform" :class="local.sort_order === 'asc' ? 'rotate-180' : ''" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
        </svg>
      </button>

      <!-- Reset filters -->
      <button
        v-if="hasActiveFilters()"
        id="monitoring-reset-filters"
        @click="onReset"
        :disabled="loading"
        class="h-9 px-3 flex items-center gap-1.5 text-xs font-bold text-rose-600 dark:text-rose-400 bg-rose-50 dark:bg-rose-900/10 border border-rose-200 dark:border-rose-800/40 rounded-xl hover:bg-rose-100 dark:hover:bg-rose-900/20 transition-all disabled:opacity-50"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
        </svg>
        Reset
      </button>

    </div>
  </div>
</template>
