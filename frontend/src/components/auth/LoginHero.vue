<script setup lang="ts">
import { ref } from 'vue';
import type { TelemetryPreview } from '@/types/auth';

const telemetryData = ref<TelemetryPreview[]>([
  { label: 'RPM', value: 1500, unit: 'rpm', trend: 'stable', status: 'normal' },
  { label: 'Oil Pressure', value: 4.2, unit: 'bar', trend: 'stable', status: 'normal' },
  { label: 'Coolant Temp', value: 85, unit: '°C', trend: 'up', status: 'warning' },
  { label: 'Voltage', value: 402, unit: 'V', trend: 'stable', status: 'normal' },
  { label: 'Fuel Level', value: 65, unit: '%', trend: 'down', status: 'normal' },
]);

const getStatusColor = (status: TelemetryPreview['status']) => {
  switch (status) {
    case 'warning': return 'text-amber-500 border-amber-500/20 bg-amber-500/5';
    case 'critical': return 'text-red-500 border-red-500/20 bg-red-500/5';
    default: return 'text-emerald-500 border-emerald-500/20 bg-emerald-500/5';
  }
};
</script>

<template>
  <div class="relative hidden lg:flex flex-col justify-between p-12 bg-slate-900 border-r border-slate-800 h-full overflow-hidden">
    <!-- Grid Pattern Overlay -->
    <div class="absolute inset-0 opacity-10 pointer-events-none" 
      style="background-image: radial-gradient(circle at 2px 2px, #475569 1px, transparent 0); background-size: 24px 24px;">
    </div>

    <!-- Top Branding -->
    <div class="relative z-10">
      <div class="flex items-center gap-3 mb-8">
        <div class="w-10 h-10 bg-indigo-600 rounded flex items-center justify-center shadow-lg shadow-indigo-600/20">
          <svg viewBox="0 0 24 24" class="w-6 h-6 text-white" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
          </svg>
        </div>
        <div>
          <h1 class="text-xl font-bold tracking-tight text-white uppercase letter">Genset <span class="text-slate-400">Monitoring</span></h1>
          <p class="text-xs text-indigo-400 font-mono tracking-widest uppercase">Industrial System v2.4</p>
        </div>
      </div>

      <h2 class="text-4xl font-semibold text-white leading-tight mb-4">
        Real-time Performance <br />
        <span class="text-slate-500 italic">Optimization & Control</span>
      </h2>
      <p class="text-slate-400 max-w-md leading-relaxed">
        Advanced telemetry monitoring for critical power infrastructure. Secure access to engine diagnostics, electrical parameters, and predictive maintenance alerts.
      </p>
    </div>

    <!-- Telemetry Cards -->
    <div class="relative z-10 grid grid-cols-2 gap-4">
      <div 
        v-for="item in telemetryData" 
        :key="item.label"
        class="p-4 rounded border border-slate-800 bg-slate-950/50 backdrop-blur-sm group hover:border-slate-700 transition-colors"
      >
        <div class="flex justify-between items-start mb-2">
          <span class="text-[10px] font-bold uppercase tracking-wider text-slate-500">{{ item.label }}</span>
          <div :class="['w-1.5 h-1.5 rounded-full animate-pulse', 
            item.status === 'normal' ? 'bg-emerald-500' : item.status === 'warning' ? 'bg-amber-500' : 'bg-red-500']">
          </div>
        </div>
        <div class="flex items-baseline gap-1">
          <span class="text-2xl font-mono font-medium text-white">{{ item.value }}</span>
          <span class="text-xs text-slate-500 font-mono uppercase">{{ item.unit }}</span>
        </div>
        <div class="mt-2 h-1 w-full bg-slate-800 rounded-full overflow-hidden">
          <div 
            class="h-full transition-all duration-1000" 
            :class="item.status === 'normal' ? 'bg-emerald-500' : item.status === 'warning' ? 'bg-amber-500' : 'bg-red-500'"
            :style="{ width: `${(item.value / (item.label === 'RPM' ? 2000 : 500)) * 100}%` }"
          ></div>
        </div>
      </div>

      <!-- Industrial SVG/Graphic Placeholder -->
      <div class="col-span-2 mt-4 p-4 border border-slate-800 bg-slate-950/50 rounded flex items-center justify-center">
        <div class="w-full h-24 flex items-center justify-around">
          <div v-for="i in 12" :key="i" 
            class="w-1 bg-slate-800 rounded-full transition-all duration-500"
            :style="{ height: `${Math.random() * 80 + 20}%` }">
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Metadata -->
    <div class="relative z-10 flex items-center justify-between text-[10px] font-mono text-slate-500 uppercase tracking-[0.2em]">
      <div class="flex gap-4">
        <span>Node: GS-HQ-01</span>
        <span>Lat: -6.2088 | Lon: 106.8456</span>
      </div>
      <div>
        Est. 2024 © Industrial Systems
      </div>
    </div>
  </div>
</template>

<style scoped>
.letter {
  letter-spacing: -0.02em;
}
</style>
