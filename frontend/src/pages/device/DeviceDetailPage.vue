<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
const deviceId = route.params.id

const telemetry = [
  { label: 'Voltage L1-N', value: '230.5 V', status: 'normal' },
  { label: 'Voltage L2-N', value: '231.2 V', status: 'normal' },
  { label: 'Voltage L3-N', value: '229.8 V', status: 'normal' },
  { label: 'Frequency', value: '50.02 Hz', status: 'normal' },
  { label: 'Current L1', value: '145.2 A', status: 'normal' },
  { label: 'Current L2', value: '142.8 A', status: 'normal' },
  { label: 'Current L3', value: '144.1 A', status: 'normal' },
  { label: 'Power Factor', value: '0.85', status: 'warning' },
]

const engineStats = [
  { label: 'Oil Pressure', value: '4.5 Bar', icon: '💧' },
  { label: 'Coolant Temp', value: '88 °C', icon: '🌡️' },
  { label: 'Battery', value: '24.2 V', icon: '🔋' },
  { label: 'Engine Speed', value: '1500 RPM', icon: '🔄' },
]
</script>

<template>
  <div class="space-y-8">
    <!-- Breadcrumbs/Back -->
    <div class="flex items-center space-x-4">
      <router-link to="/" class="text-blue-400 hover:text-blue-300 text-sm">← Back to Dashboard</router-link>
      <div class="h-4 w-px bg-slate-800"></div>
      <div class="text-slate-500 text-sm">Device ID: {{ deviceId }}</div>
    </div>

    <!-- Header Info -->
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center bg-slate-900 border border-slate-800 p-6 rounded-xl">
      <div>
        <h2 class="text-2xl font-bold text-white">Genset Alpha-001</h2>
        <p class="text-slate-400 text-sm mt-1">SN: GS-2024-X991-B</p>
      </div>
      <div class="mt-4 md:mt-0 flex space-x-3">
        <button class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg text-sm font-medium transition-colors">START</button>
        <button class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg text-sm font-medium transition-colors">STOP</button>
        <button class="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg text-sm font-medium transition-colors">MAINTENANCE</button>
      </div>
    </div>

    <!-- Telemetry Grid -->
    <div class="grid grid-cols-1 xl:grid-cols-3 gap-8">
      <!-- Electrical Monitoring -->
      <div class="xl:col-span-2 space-y-6">
        <div class="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden">
          <div class="px-6 py-4 border-b border-slate-800 bg-slate-900/50">
            <h3 class="font-semibold text-slate-300">Electrical Telemetry</h3>
          </div>
          <div class="p-6 grid grid-cols-2 md:grid-cols-4 gap-6">
            <div v-for="item in telemetry" :key="item.label" class="space-y-1">
              <div class="text-[10px] uppercase text-slate-500 font-bold tracking-widest">{{ item.label }}</div>
              <div class="text-xl font-mono font-bold" :class="item.status === 'warning' ? 'text-orange-400' : 'text-blue-400'">
                {{ item.value }}
              </div>
            </div>
          </div>
        </div>

        <!-- Charts Placeholder -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="bg-slate-900 border border-slate-800 rounded-xl p-6 h-64 flex flex-col">
            <h4 class="text-sm font-medium text-slate-400 mb-4">Voltage History (V)</h4>
            <div class="flex-1 bg-slate-950 rounded border border-slate-800 border-dashed flex items-center justify-center">
              <span class="text-slate-600 text-xs italic">Voltage Chart Component</span>
            </div>
          </div>
          <div class="bg-slate-900 border border-slate-800 rounded-xl p-6 h-64 flex flex-col">
            <h4 class="text-sm font-medium text-slate-400 mb-4">Load Analysis (kW)</h4>
            <div class="flex-1 bg-slate-950 rounded border border-slate-800 border-dashed flex items-center justify-center">
              <span class="text-slate-600 text-xs italic">Load Chart Component</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Engine Status & Events -->
      <div class="space-y-6">
        <div class="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden">
          <div class="px-6 py-4 border-b border-slate-800 bg-slate-900/50">
            <h3 class="font-semibold text-slate-300">Engine Parameters</h3>
          </div>
          <div class="p-6 space-y-4">
            <div v-for="stat in engineStats" :key="stat.label" class="flex items-center justify-between p-3 bg-slate-950 rounded-lg border border-slate-800">
              <div class="flex items-center space-x-3">
                <span class="text-xl">{{ stat.icon }}</span>
                <span class="text-sm text-slate-400">{{ stat.label }}</span>
              </div>
              <div class="text-lg font-bold text-white">{{ stat.value }}</div>
            </div>
          </div>
        </div>

        <div class="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden">
          <div class="px-6 py-4 border-b border-slate-800 bg-slate-900/50">
            <h3 class="font-semibold text-slate-300">Recent Events</h3>
          </div>
          <div class="p-4 space-y-3">
            <div class="flex space-x-3 p-2 hover:bg-slate-800/50 rounded transition-colors border-l-2 border-green-500 bg-green-500/5">
              <div class="text-[10px] text-slate-500 mt-0.5">14:22:05</div>
              <div class="text-xs text-slate-300">Engine started successfully via remote command.</div>
            </div>
            <div class="flex space-x-3 p-2 hover:bg-slate-800/50 rounded transition-colors border-l-2 border-blue-500 bg-blue-500/5">
              <div class="text-[10px] text-slate-500 mt-0.5">12:05:12</div>
              <div class="text-xs text-slate-300">Fuel level drop detected (Refill needed soon).</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
