<script setup lang="ts">
import type { MonitoringLatestStateOutput, EngineTelemetryOutput } from '@/types/monitoring'

const props = defineProps<{
  latestState: MonitoringLatestStateOutput
  engineTelemetry: EngineTelemetryOutput | null
}>()

const fmt = (v: number | undefined | null, decimals = 1): string =>
  v !== undefined && v !== null ? v.toFixed(decimals) : '—'

const fmtInt = (v: number | undefined | null): string =>
  v !== undefined && v !== null ? Math.round(v).toString() : '—'

const fmtRuntime = (seconds: number | undefined | null): string => {
  if (seconds === undefined || seconds === null) return '—'
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  return `${h}h ${m}m ${s}s`
}

const isRunning = (): boolean => props.latestState?.engine_running ?? false

// ── Sub-group definitions ────────────────────────────────────────────────────

interface Metric {
  label: string
  value: string
  unit: string
  warn?: boolean
  danger?: boolean
}

interface Group {
  title: string
  icon: string
  color: string
  metrics: Metric[]
}

const e = () => props.engineTelemetry
const s = () => props.latestState

const groups = (): Group[] => {
  const et = e()
  const ls = s()

  return [
    {
      title: 'Status & Speed',
      icon: '🔄',
      color: 'indigo',
      metrics: [
        {
          label: 'Engine Speed',
          value: fmtInt(et?.speed ?? ls.speed),
          unit: 'RPM',
          warn: (et?.speed ?? ls.speed ?? 0) > 1800,
          danger: (et?.speed ?? ls.speed ?? 0) > 2000,
        },
        {
          label: 'Rated Speed',
          value: fmtInt(et?.rated_speed),
          unit: 'RPM',
        },
        {
          label: 'Desired Speed',
          value: fmtInt(et?.desired_operating_speed),
          unit: 'RPM',
        },
        {
          label: 'Rated Power',
          value: fmt(et?.rated_power),
          unit: 'kW',
        },
        {
          label: 'Battery Volt',
          value: fmt(et?.batt_volt ?? ls.batt_volt),
          unit: 'V',
          warn: (et?.batt_volt ?? ls.batt_volt ?? 99) < 22,
          danger: (et?.batt_volt ?? ls.batt_volt ?? 99) < 11,
        },
        {
          label: 'Keyswitch Batt',
          value: fmt(et?.keyswitch_batt_potential),
          unit: 'V',
        },
        {
          label: 'Run Time',
          value: fmtRuntime(et?.run_time),
          unit: '',
        },
      ],
    },
    {
      title: 'Fuel System',
      icon: '⛽',
      color: 'amber',
      metrics: [
        {
          label: 'Fuel Level Top',
          value: fmt(et?.fuel_level_top ?? ls.fuel_level),
          unit: '%',
          warn: (et?.fuel_level_top ?? ls.fuel_level ?? 100) < 30,
          danger: (et?.fuel_level_top ?? ls.fuel_level ?? 100) < 15,
        },
        {
          label: 'Fuel Level Bottom',
          value: fmt(et?.fuel_level_bottom),
          unit: '%',
        },
        {
          label: 'Fuel Pressure 1',
          value: fmt(et?.fuel_level_pressure_1),
          unit: 'kPa',
        },
        {
          label: 'Fuel Pressure 2',
          value: fmt(et?.fuel_level_pressure_2),
          unit: 'kPa',
        },
        {
          label: 'Fuel Rate',
          value: fmt(et?.fuel_rate),
          unit: 'L/h',
        },
        {
          label: 'Avg Fuel Rate',
          value: fmt(et?.avg_fuel_rate),
          unit: 'L/h',
        },
        {
          label: 'Total Fuel',
          value: fmt(et?.total_fuel, 2),
          unit: 'L',
        },
        {
          label: 'Trip Fuel',
          value: fmt(et?.trip_fuel, 2),
          unit: 'L',
        },
      ],
    },
    {
      title: 'Pressures & Temperatures',
      icon: '🌡',
      color: 'rose',
      metrics: [
        {
          label: 'Oil Pressure',
          value: fmt(et?.oil_pressure ?? ls.oil_pressure),
          unit: 'Bar',
          warn: (et?.oil_pressure ?? ls.oil_pressure ?? 99) < 2,
          danger: (et?.oil_pressure ?? ls.oil_pressure ?? 99) < 1,
        },
        {
          label: 'Oil Filter Pressure',
          value: fmt(et?.oil_filter_out_pressure),
          unit: 'Bar',
        },
        {
          label: 'Coolant Temp',
          value: fmt(et?.coolant_temperature ?? ls.coolant_temperature),
          unit: '°C',
          warn: (et?.coolant_temperature ?? ls.coolant_temperature ?? 0) > 95,
          danger: (et?.coolant_temperature ?? ls.coolant_temperature ?? 0) > 105,
        },
        {
          label: 'ECU Temp',
          value: fmt(et?.ecu_temperature),
          unit: '°C',
          warn: (et?.ecu_temperature ?? 0) > 80,
        },
        {
          label: 'Turbo Pressure',
          value: fmt(et?.turbo_pressure),
          unit: 'kPa',
        },
        {
          label: 'Intake Manifold P.',
          value: fmt(et?.intake_manifold_pressure),
          unit: 'kPa',
        },
        {
          label: 'Intake Manifold T.',
          value: fmt(et?.intake_manifold_temperature),
          unit: '°C',
        },
      ],
    },
  ]
}

const colorMap: Record<string, { header: string; card: string; value: string }> = {
  indigo: {
    header: 'text-indigo-600 dark:text-indigo-400',
    card: 'border-indigo-100 dark:border-indigo-900/30',
    value: 'text-indigo-700 dark:text-indigo-300',
  },
  amber: {
    header: 'text-amber-600 dark:text-amber-400',
    card: 'border-amber-100 dark:border-amber-900/30',
    value: 'text-amber-700 dark:text-amber-300',
  },
  rose: {
    header: 'text-rose-600 dark:text-rose-400',
    card: 'border-rose-100 dark:border-rose-900/30',
    value: 'text-rose-700 dark:text-rose-300',
  },
}

const valueClass = (m: Metric, color: string): string => {
  if (m.danger) return 'text-rose-600 dark:text-rose-400'
  if (m.warn) return 'text-amber-500 dark:text-amber-400'
  return colorMap[color]?.value ?? 'text-slate-900 dark:text-white'
}
</script>

<template>
  <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 flex items-center justify-between">
      <h3 class="font-bold text-slate-900 dark:text-slate-100 flex items-center gap-2 text-sm">
        <svg class="w-4 h-4 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        Engine Telemetry
      </h3>
      <!-- Engine state badge -->
      <span
        class="px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-widest border transition-all"
        :class="isRunning()
          ? 'bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-900/15 dark:text-blue-400 dark:border-blue-800/40'
          : 'bg-slate-100 text-slate-500 border-slate-200 dark:bg-slate-800/30 dark:text-slate-400 dark:border-slate-700'"
      >
        <span v-if="isRunning()" class="mr-1">⚡</span>
        {{ isRunning() ? 'Running' : 'Stopped' }}
      </span>
    </div>

    <!-- Sub-groups -->
    <div class="p-6 space-y-6">
      <div v-for="group in groups()" :key="group.title">
        <!-- Group title -->
        <div class="flex items-center gap-2 mb-3">
          <span class="text-base">{{ group.icon }}</span>
          <h4 class="text-xs font-bold uppercase tracking-widest" :class="colorMap[group.color]?.header">
            {{ group.title }}
          </h4>
          <div class="flex-1 h-px bg-slate-100 dark:bg-slate-800"></div>
        </div>

        <!-- Metric grid -->
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
          <div
            v-for="metric in group.metrics" :key="metric.label"
            class="bg-slate-50 dark:bg-slate-950 border rounded-xl p-3 flex flex-col gap-0.5 transition-colors"
            :class="colorMap[group.color]?.card"
          >
            <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">{{ metric.label }}</span>
            <div class="flex items-baseline gap-1 mt-1">
              <span class="text-lg font-mono font-bold leading-none transition-colors" :class="valueClass(metric, group.color)">
                {{ metric.value }}
              </span>
              <span v-if="metric.unit" class="text-[10px] text-slate-400 font-mono">{{ metric.unit }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Telemetry recorded at -->
      <p v-if="engineTelemetry?.created_at" class="text-[11px] text-slate-400 font-mono text-right">
        Recorded: {{ new Date(engineTelemetry.created_at).toLocaleString('en-GB', { dateStyle: 'short', timeStyle: 'medium' }) }}
      </p>
    </div>
  </div>
</template>
