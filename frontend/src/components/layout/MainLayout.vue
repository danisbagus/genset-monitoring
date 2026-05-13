<script setup lang="ts">
import { ref } from 'vue'

const sidebarOpen = ref(true)

const navigation = [
  { name: 'Dashboard', href: '/', icon: 'HomeIcon' },
  { name: 'Devices', href: '/devices/1', icon: 'DeviceIcon' }, // Example ID
]
</script>

<template>
  <div class="min-h-screen flex bg-slate-950 text-slate-100">
    <!-- Sidebar -->
    <aside 
      class="w-64 border-r border-slate-800 bg-slate-900 transition-all duration-300 flex flex-col"
      :class="{ 'w-20': !sidebarOpen }"
    >
      <div class="h-16 flex items-center px-6 border-b border-slate-800">
        <div class="w-8 h-8 bg-blue-600 rounded flex items-center justify-center font-bold">G</div>
        <span v-if="sidebarOpen" class="ml-3 font-semibold tracking-wider uppercase text-sm">Genset Monitor</span>
      </div>
      
      <nav class="flex-1 p-4 space-y-2">
        <router-link 
          v-for="item in navigation" 
          :key="item.name"
          :to="item.href"
          class="flex items-center p-3 rounded-lg hover:bg-slate-800 transition-colors"
          active-class="bg-slate-800 text-blue-400"
        >
          <div class="w-6 h-6 bg-slate-700 rounded mr-3"></div>
          <span v-if="sidebarOpen">{{ item.name }}</span>
        </router-link>
      </nav>
      
      <div class="p-4 border-t border-slate-800">
        <button 
          @click="sidebarOpen = !sidebarOpen"
          class="w-full p-2 hover:bg-slate-800 rounded flex justify-center items-center"
        >
          <span v-if="sidebarOpen">Collapse</span>
          <span v-else>→</span>
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <header class="h-16 border-b border-slate-800 bg-slate-900 flex items-center justify-between px-8">
        <h1 class="text-lg font-medium">{{ $route.meta.title }}</h1>
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-2">
            <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <span class="text-xs text-slate-400">System Live</span>
          </div>
          <div class="w-8 h-8 bg-slate-700 rounded-full"></div>
        </div>
      </header>
      
      <section class="flex-1 overflow-y-auto p-8 custom-scrollbar">
        <router-view />
      </section>
    </main>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  @apply bg-slate-950;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  @apply bg-slate-800 rounded;
}
</style>
