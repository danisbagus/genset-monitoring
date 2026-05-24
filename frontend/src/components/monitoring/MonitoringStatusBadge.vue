<script setup lang="ts">
/**
 * MonitoringStatusBadge
 * Unified status badge for online/offline, engine state, and device status.
 */

type BadgeType = 'online' | 'offline' | 'engine_running' | 'engine_stopped' | 'device_status'

const props = withDefaults(defineProps<{
  type: BadgeType
  value?: boolean | string
  pulse?: boolean
}>(), {
  pulse: true
})

// Derive label and color class from type + value
const config = (() => {
  switch (props.type) {
    case 'online':
      return props.value
        ? { label: 'Online', color: 'emerald', icon: '●' }
        : { label: 'Offline', color: 'slate', icon: '○' }
    case 'offline':
      return { label: 'Offline', color: 'slate', icon: '○' }
    case 'engine_running':
      return props.value
        ? { label: 'Running', color: 'blue', icon: '⚡' }
        : { label: 'Stopped', color: 'slate', icon: '—' }
    case 'engine_stopped':
      return { label: 'Stopped', color: 'slate', icon: '—' }
    case 'device_status': {
      const status = (props.value as string)?.toLowerCase()
      if (status === 'active') return { label: 'Active', color: 'emerald', icon: null }
      if (status === 'maintenance') return { label: 'Maintenance', color: 'amber', icon: null }
      return { label: 'Inactive', color: 'rose', icon: null }
    }
    default:
      return { label: String(props.value ?? '—'), color: 'slate', icon: null }
  }
})()

const colorMap: Record<string, { badge: string; dot: string }> = {
  emerald: {
    badge: 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-900/15 dark:text-emerald-400 dark:border-emerald-800/40',
    dot: 'bg-emerald-500'
  },
  blue: {
    badge: 'bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-900/15 dark:text-blue-400 dark:border-blue-800/40',
    dot: 'bg-blue-500'
  },
  amber: {
    badge: 'bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-900/15 dark:text-amber-400 dark:border-amber-800/40',
    dot: 'bg-amber-500'
  },
  rose: {
    badge: 'bg-rose-50 text-rose-700 border-rose-200 dark:bg-rose-900/15 dark:text-rose-400 dark:border-rose-800/40',
    dot: 'bg-rose-500'
  },
  slate: {
    badge: 'bg-slate-100 text-slate-500 border-slate-200 dark:bg-slate-800/30 dark:text-slate-400 dark:border-slate-700/50',
    dot: 'bg-slate-400'
  }
}

const colorCfg = colorMap[config.color] ?? colorMap.slate
const isPulsing = props.pulse && (props.type === 'online' && props.value === true)
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-md border text-[10px] font-bold uppercase tracking-widest transition-all duration-300"
    :class="colorCfg.badge"
  >
    <!-- Pulse dot for online -->
    <span v-if="isPulsing" class="relative flex h-2 w-2">
      <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75" :class="colorCfg.dot"></span>
      <span class="relative inline-flex rounded-full h-2 w-2" :class="colorCfg.dot"></span>
    </span>
    <!-- Static dot otherwise -->
    <span v-else class="inline-flex h-1.5 w-1.5 rounded-full" :class="colorCfg.dot"></span>

    {{ config.label }}
  </span>
</template>
