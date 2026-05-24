<script setup lang="ts">
import type { ElectricalTelemetryOutput, MonitoringLatestStateOutput } from '@/types/monitoring'

const props = defineProps<{
  latestState: MonitoringLatestStateOutput
  electricalTelemetry: ElectricalTelemetryOutput | null
}>()

const fmt = (v: number | undefined | null, decimals = 1): string =>
  v !== undefined && v !== null ? v.toFixed(decimals) : '—'

const fmt3 = (v: number | undefined | null): string => fmt(v, 3)

// ── Row / group definitions ──────────────────────────────────────────────────

interface Row {
  label: string
  value: string
  unit: string
}

interface Group {
  title: string
  color: string
  rows: Row[]
}

const e = () => props.electricalTelemetry
const s = () => props.latestState

const groups = (): Group[] => {
  const et = e()
  const ls = s()
  return [
    {
      title: 'Voltage',
      color: 'blue',
      rows: [
        { label: 'L1-N',   value: fmt(et?.l1_n_volt),  unit: 'V' },
        { label: 'L2-N',   value: fmt(et?.l2_n_volt),  unit: 'V' },
        { label: 'L3-N',   value: fmt(et?.l3_n_volt),  unit: 'V' },
        { label: 'L1-L2',  value: fmt(et?.l1_l2_volt), unit: 'V' },
        { label: 'L2-L3',  value: fmt(et?.l2_l3_volt), unit: 'V' },
        { label: 'L3-L1',  value: fmt(et?.l3_l1_volt), unit: 'V' },
        { label: 'Chrg Alt', value: fmt(et?.charge_alt_volt), unit: 'V' },
      ],
    },
    {
      title: 'Current',
      color: 'emerald',
      rows: [
        { label: 'L1',    value: fmt(et?.l1_curr),   unit: 'A' },
        { label: 'L2',    value: fmt(et?.l2_curr),   unit: 'A' },
        { label: 'L3',    value: fmt(et?.l3_curr),   unit: 'A' },
        { label: 'Earth', value: fmt(et?.earth_curr), unit: 'A' },
      ],
    },
    {
      title: 'Apparent Power (VA)',
      color: 'violet',
      rows: [
        { label: 'L1 VA',    value: fmt(et?.l1_va),   unit: 'VA' },
        { label: 'L2 VA',    value: fmt(et?.l2_va),   unit: 'VA' },
        { label: 'L3 VA',    value: fmt(et?.l3_va),   unit: 'VA' },
        { label: 'Total VA', value: fmt((et?.total_va ?? ls.total_va) / 1000, 2), unit: 'kVA' },
      ],
    },
    {
      title: 'Reactive Power (VAR)',
      color: 'cyan',
      rows: [
        { label: 'L1 VAR',    value: fmt(et?.l1_var),   unit: 'VAR' },
        { label: 'L2 VAR',    value: fmt(et?.l2_var),   unit: 'VAR' },
        { label: 'L3 VAR',    value: fmt(et?.l3_var),   unit: 'VAR' },
        { label: 'Total VAR', value: fmt(et?.total_var), unit: 'VAR' },
      ],
    },
    {
      title: 'Power Factor & Frequency',
      color: 'amber',
      rows: [
        { label: 'PF L1',  value: fmt3(et?.pf_l1),                       unit: '' },
        { label: 'PF L2',  value: fmt3(et?.pf_l2),                       unit: '' },
        { label: 'PF L3',  value: fmt3(et?.pf_l3),                       unit: '' },
        { label: 'PF Avg', value: fmt3(et?.pf_avg ?? ls.pf_avg),         unit: '' },
        { label: 'Frequency', value: fmt(et?.frequency ?? ls.frequency, 2), unit: 'Hz' },
        { label: 'Load %V',   value: fmt(et?.percent_fv),                unit: '%' },
        { label: 'Load %P',   value: fmt(et?.percent_fp),                unit: '%' },
      ],
    },
  ]
}

const colorMap: Record<string, { header: string; pill: string; value: string }> = {
  blue: {
    header: 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/10 border-blue-200 dark:border-blue-800/40',
    pill:   'bg-blue-100 dark:bg-blue-900/20',
    value:  'text-blue-700 dark:text-blue-300',
  },
  emerald: {
    header: 'text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-900/10 border-emerald-200 dark:border-emerald-800/40',
    pill:   'bg-emerald-100 dark:bg-emerald-900/20',
    value:  'text-emerald-700 dark:text-emerald-300',
  },
  violet: {
    header: 'text-violet-600 dark:text-violet-400 bg-violet-50 dark:bg-violet-900/10 border-violet-200 dark:border-violet-800/40',
    pill:   'bg-violet-100 dark:bg-violet-900/20',
    value:  'text-violet-700 dark:text-violet-300',
  },
  cyan: {
    header: 'text-cyan-600 dark:text-cyan-400 bg-cyan-50 dark:bg-cyan-900/10 border-cyan-200 dark:border-cyan-800/40',
    pill:   'bg-cyan-100 dark:bg-cyan-900/20',
    value:  'text-cyan-700 dark:text-cyan-300',
  },
  amber: {
    header: 'text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/10 border-amber-200 dark:border-amber-800/40',
    pill:   'bg-amber-100 dark:bg-amber-900/20',
    value:  'text-amber-700 dark:text-amber-300',
  },
}
</script>

<template>
  <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 flex items-center justify-between">
      <h3 class="font-bold text-slate-900 dark:text-slate-100 flex items-center gap-2 text-sm">
        <svg class="w-4 h-4 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
        Electrical Telemetry
      </h3>
      <p v-if="electricalTelemetry?.created_at" class="text-[11px] text-slate-400 font-mono hidden sm:block">
        {{ new Date(electricalTelemetry.created_at).toLocaleString('en-GB', { dateStyle: 'short', timeStyle: 'medium' }) }}
      </p>
    </div>

    <!-- Groups -->
    <div class="p-6 space-y-6">
      <div v-for="group in groups()" :key="group.title">

        <!-- Group pill header -->
        <div class="inline-flex items-center px-3 py-1 rounded-lg border text-[10px] font-bold uppercase tracking-widest mb-3 transition-all" :class="colorMap[group.color]?.header">
          {{ group.title }}
        </div>

        <!-- Row grid -->
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-2.5">
          <div
            v-for="row in group.rows" :key="row.label"
            class="flex items-center justify-between px-3 py-2.5 rounded-xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800"
          >
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider shrink-0">{{ row.label }}</span>
            <div class="flex items-baseline gap-0.5 ml-2">
              <span class="text-sm font-mono font-bold tabular-nums" :class="colorMap[group.color]?.value">{{ row.value }}</span>
              <span v-if="row.unit" class="text-[10px] text-slate-400 font-mono">{{ row.unit }}</span>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>
