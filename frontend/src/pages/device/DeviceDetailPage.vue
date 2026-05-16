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
  <div class="space-y-6 sm:space-y-8">
    <!-- Breadcrumbs/Back -->
    <div class="flex items-center space-x-3 sm:space-x-4">
      <router-link to="/" class="text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 text-xs sm:text-sm font-bold flex items-center transition-colors">
        <svg class="w-4 h-4 mr-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        Back to Dashboard
      </router-link>
      <div class="h-4 w-px bg-slate-200 dark:bg-slate-800"></div>
      <div class="text-slate-500 text-[10px] sm:text-xs font-mono">DEVICE ID: {{ deviceId }}</div>
    </div>

    <!-- Header Info -->
    <div class="flex flex-col lg:flex-row justify-between items-start lg:items-center bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 p-4 sm:p-6 rounded-2xl shadow-sm transition-all duration-300 gap-4 sm:gap-6">
      <div class="flex items-center gap-4">
        <div class="w-12 h-12 sm:w-16 sm:h-16 bg-blue-600 rounded-2xl flex items-center justify-center text-white shadow-lg shadow-blue-600/20">
          <svg class="w-6 h-6 sm:w-8 sm:h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
          </svg>
        </div>
        <div>
          <h2 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">Genset Alpha-001</h2>
          <div class="flex items-center gap-2 mt-1">
            <span class="text-slate-500 text-xs sm:text-sm">SN: GS-2024-X991-B</span>
            <span class="w-1 h-1 rounded-full bg-slate-300 dark:bg-slate-700"></span>
            <span class="text-emerald-600 dark:text-emerald-400 text-[10px] sm:text-xs font-bold uppercase tracking-widest">Active Now</span>
          </div>
        </div>
      </div>
      <div class="w-full lg:w-auto grid grid-cols-2 sm:flex sm:flex-row gap-3">
        <button class="px-4 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs sm:text-sm font-bold transition-all shadow-lg shadow-emerald-600/20 flex items-center justify-center gap-2">
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd" /></svg>
          START
        </button>
        <button class="px-4 py-2.5 bg-rose-600 hover:bg-rose-700 text-white rounded-xl text-xs sm:text-sm font-bold transition-all shadow-lg shadow-rose-600/20 flex items-center justify-center gap-2">
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8 7a1 1 0 00-1 1v4a1 1 0 001 1h4a1 1 0 001-1V8a1 1 0 00-1-1H8z" clip-rule="evenodd" /></svg>
          STOP
        </button>
        <button class="col-span-2 sm:col-auto px-4 py-2.5 bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-700 dark:text-white rounded-xl text-xs sm:text-sm font-bold transition-all border border-slate-200 dark:border-slate-700 flex items-center justify-center gap-2">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
          MAINTENANCE
        </button>
      </div>
    </div>

    <!-- Telemetry Grid -->
    <div class="grid grid-cols-1 xl:grid-cols-3 gap-6 sm:gap-8">
      <!-- Electrical Monitoring -->
      <div class="xl:col-span-2 space-y-6 sm:space-y-8">
        <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden transition-all duration-300">
          <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 flex items-center justify-between">
            <h3 class="font-bold text-slate-900 dark:text-slate-100 flex items-center gap-2 text-sm sm:text-base">
              <svg class="w-5 h-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
              Electrical Telemetry
            </h3>
            <span class="text-[10px] font-bold uppercase text-slate-400 tracking-tighter hidden sm:block">Realtime Stream</span>
          </div>
          <div class="p-4 sm:p-6 grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-8">
            <div v-for="item in telemetry" :key="item.label" class="space-y-1 group">
              <div class="text-[10px] uppercase text-slate-500 font-bold tracking-widest">{{ item.label }}</div>
              <div class="text-xl sm:text-2xl font-mono font-bold transition-all duration-300" :class="item.status === 'warning' ? 'text-amber-500' : 'text-blue-600 dark:text-blue-400 group-hover:translate-x-1'">
                {{ item.value }}
              </div>
            </div>
          </div>
        </div>

        <!-- Charts Placeholder -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 sm:gap-8">
          <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-6 h-72 flex flex-col shadow-sm transition-all duration-300">
            <div class="flex justify-between items-center mb-6">
              <h4 class="text-xs sm:text-sm font-bold text-slate-900 dark:text-slate-100 uppercase tracking-wider">Voltage History (V)</h4>
              <div class="flex gap-1">
                <div class="w-1.5 h-1.5 rounded-full bg-blue-500"></div>
                <div class="w-1.5 h-1.5 rounded-full bg-blue-500/50"></div>
                <div class="w-1.5 h-1.5 rounded-full bg-blue-500/20"></div>
              </div>
            </div>
            <div class="flex-1 bg-slate-50 dark:bg-slate-950 rounded-xl border border-slate-200 dark:border-slate-800 border-dashed flex flex-col items-center justify-center gap-3">
              <svg class="w-10 h-10 text-slate-300 dark:text-slate-800" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" /></svg>
              <span class="text-slate-400 dark:text-slate-600 text-[10px] font-bold uppercase tracking-[0.2em]">Voltage Chart Active</span>
            </div>
          </div>
          <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-6 h-72 flex flex-col shadow-sm transition-all duration-300">
            <div class="flex justify-between items-center mb-6">
              <h4 class="text-xs sm:text-sm font-bold text-slate-900 dark:text-slate-100 uppercase tracking-wider">Load Analysis (kW)</h4>
              <div class="flex gap-1">
                <div class="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
                <div class="w-1.5 h-1.5 rounded-full bg-emerald-500/50"></div>
                <div class="w-1.5 h-1.5 rounded-full bg-emerald-500/20"></div>
              </div>
            </div>
            <div class="flex-1 bg-slate-50 dark:bg-slate-950 rounded-xl border border-slate-200 dark:border-slate-800 border-dashed flex flex-col items-center justify-center gap-3">
              <svg class="w-10 h-10 text-slate-300 dark:text-slate-800" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" /></svg>
              <span class="text-slate-400 dark:text-slate-600 text-[10px] font-bold uppercase tracking-[0.2em]">Load Analysis Chart Active</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Engine Status & Events -->
      <div class="space-y-6 sm:space-y-8">
        <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden transition-all duration-300">
          <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50">
            <h3 class="font-bold text-slate-900 dark:text-slate-100 flex items-center gap-2 text-sm sm:text-base">
              <svg class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /></svg>
              Engine Parameters
            </h3>
          </div>
          <div class="p-6 space-y-4">
            <div v-for="stat in engineStats" :key="stat.label" class="flex items-center justify-between p-4 bg-slate-50 dark:bg-slate-950 rounded-2xl border border-slate-200 dark:border-slate-800 transition-all hover:border-indigo-500/50 group">
              <div class="flex items-center space-x-3">
                <span class="text-xl group-hover:scale-125 transition-transform">{{ stat.icon }}</span>
                <span class="text-sm font-bold text-slate-500 dark:text-slate-400 uppercase tracking-widest">{{ stat.label }}</span>
              </div>
              <div class="text-lg font-mono font-bold text-slate-900 dark:text-white">{{ stat.value }}</div>
            </div>
          </div>
        </div>

        <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden transition-all duration-300">
          <div class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50">
            <h3 class="font-bold text-slate-900 dark:text-slate-100 flex items-center gap-2 text-sm sm:text-base">
              <svg class="w-5 h-5 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
              Recent Events
            </h3>
          </div>
          <div class="p-4 space-y-3">
            <div class="flex space-x-3 p-3 hover:bg-slate-50 dark:hover:bg-slate-800/50 rounded-xl transition-all border-l-4 border-emerald-500 bg-emerald-50 dark:bg-emerald-500/5 group">
              <div class="text-[10px] text-slate-400 mt-0.5 font-bold font-mono">14:22:05</div>
              <div class="text-xs sm:text-sm text-slate-700 dark:text-slate-300 group-hover:text-slate-900 dark:group-hover:text-white transition-colors">Engine started successfully via remote command.</div>
            </div>
            <div class="flex space-x-3 p-3 hover:bg-slate-50 dark:hover:bg-slate-800/50 rounded-xl transition-all border-l-4 border-blue-500 bg-blue-50 dark:bg-blue-500/5 group">
              <div class="text-[10px] text-slate-400 mt-0.5 font-bold font-mono">12:05:12</div>
              <div class="text-xs sm:text-sm text-slate-700 dark:text-slate-300 group-hover:text-slate-900 dark:group-hover:text-white transition-colors">Fuel level drop detected (Refill needed soon).</div>
            </div>
          </div>
          <div class="px-6 py-4 bg-slate-50 dark:bg-slate-900/80 border-t border-slate-100 dark:border-slate-800 text-center rounded-b-2xl">
            <button class="text-[10px] font-bold text-blue-600 dark:text-blue-400 uppercase tracking-widest hover:underline">View All Event Logs</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
