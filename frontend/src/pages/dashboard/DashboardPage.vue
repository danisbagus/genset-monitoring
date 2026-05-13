<script setup lang="ts">
import { ref } from 'vue'

const devices = ref([
  { id: 1, name: 'Genset Alpha-1', status: 'Running', load: '75%', fuel: '82%', temp: '85°C' },
  { id: 2, name: 'Genset Bravo-2', status: 'Running', load: '45%', fuel: '64%', temp: '72°C' },
  { id: 3, name: 'Genset Gamma-3', status: 'Standby', load: '0%', fuel: '95%', temp: '24°C' },
  { id: 4, name: 'Genset Delta-4', status: 'Warning', load: '92%', fuel: '12%', temp: '98°C' },
])

const stats = [
  { label: 'Total Power', value: '450 kW', icon: '⚡' },
  { label: 'Avg Fuel Level', value: '63%', icon: '⛽' },
  { label: 'Active Devices', value: '3/4', icon: '🔧' },
  { label: 'System Health', value: 'Normal', icon: '🛡️' },
]
</script>

<template>
  <div class="space-y-8">
    <!-- Top Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div 
        v-for="stat in stats" 
        :key="stat.label"
        class="bg-slate-900 border border-slate-800 p-6 rounded-xl hover:border-blue-500/50 transition-all duration-300"
      >
        <div class="flex items-center justify-between">
          <span class="text-slate-400 text-sm font-medium uppercase tracking-wider">{{ stat.label }}</span>
          <span class="text-2xl">{{ stat.icon }}</span>
        </div>
        <div class="mt-4 text-3xl font-bold text-white">{{ stat.value }}</div>
      </div>
    </div>

    <!-- Realtime Monitor Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <div class="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden">
        <div class="px-6 py-4 border-b border-slate-800 flex justify-between items-center bg-slate-900/50">
          <h2 class="font-semibold text-slate-300">Active Fleet Status</h2>
          <button class="text-xs text-blue-400 hover:underline">View All</button>
        </div>
        <div class="p-0">
          <table class="w-full text-left">
            <thead>
              <tr class="text-slate-500 text-xs uppercase border-b border-slate-800">
                <th class="px-6 py-3 font-medium">Device Name</th>
                <th class="px-6 py-3 font-medium">Status</th>
                <th class="px-6 py-3 font-medium">Load</th>
                <th class="px-6 py-3 font-medium text-right">Fuel</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800">
              <tr 
                v-for="device in devices" 
                :key="device.id"
                class="hover:bg-slate-800/30 cursor-pointer transition-colors"
                @click="$router.push(`/devices/${device.id}`)"
              >
                <td class="px-6 py-4 font-medium text-slate-200">{{ device.name }}</td>
                <td class="px-6 py-4">
                  <span 
                    class="px-2 py-1 rounded text-[10px] font-bold uppercase"
                    :class="{
                      'bg-green-500/10 text-green-500': device.status === 'Running',
                      'bg-blue-500/10 text-blue-500': device.status === 'Standby',
                      'bg-orange-500/10 text-orange-500': device.status === 'Warning',
                    }"
                  >
                    {{ device.status }}
                  </span>
                </td>
                <td class="px-6 py-4 text-slate-400">{{ device.load }}</td>
                <td class="px-6 py-4 text-right">
                  <div class="flex items-center justify-end space-x-2">
                    <div class="w-16 bg-slate-800 h-1.5 rounded-full overflow-hidden">
                      <div 
                        class="h-full rounded-full"
                        :class="parseInt(device.fuel) < 20 ? 'bg-red-500' : 'bg-blue-500'"
                        :style="{ width: device.fuel }"
                      ></div>
                    </div>
                    <span class="text-xs text-slate-300 w-8">{{ device.fuel }}</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Placeholder Chart -->
      <div class="bg-slate-900 border border-slate-800 rounded-xl flex flex-col">
        <div class="px-6 py-4 border-b border-slate-800 flex justify-between items-center bg-slate-900/50">
          <h2 class="font-semibold text-slate-300">Total Power Output (24h)</h2>
        </div>
        <div class="flex-1 flex items-center justify-center p-8">
          <div class="w-full h-48 bg-slate-800/50 rounded-lg flex items-center justify-center border border-dashed border-slate-700">
            <span class="text-slate-500 text-sm italic">Realtime ECharts Visualization Placeholder</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
