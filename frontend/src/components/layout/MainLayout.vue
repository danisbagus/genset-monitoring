<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth.store'
import { useRouter, useRoute } from 'vue-router'
import { useTheme } from '@/composables/useTheme'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const { mode, toggleTheme } = useTheme()

// Persistent sidebar state
const sidebarOpen = ref(window.innerWidth >= 1024)
const mobileMenuOpen = ref(false)

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
  if (window.innerWidth >= 1024) {
    localStorage.setItem('sidebar_collapsed', (!sidebarOpen.value).toString())
  }
}

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

// Close mobile menu on route change
watch(() => route.path, () => {
  mobileMenuOpen.value = false
})

onMounted(() => {
  const collapsed = localStorage.getItem('sidebar_collapsed') === 'true'
  if (window.innerWidth >= 1024) {
    sidebarOpen.value = !collapsed
  } else {
    sidebarOpen.value = false
  }
})

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}

const navigation = [
  { name: 'Dashboard', href: '/', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' },
  { name: 'Devices', href: '/devices/1', icon: 'M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z' },
]
</script>

<template>
  <div class="min-h-screen flex bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 transition-colors duration-300">
    <!-- Mobile Overlay -->
    <div 
      v-if="mobileMenuOpen" 
      @click="mobileMenuOpen = false"
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-30 lg:hidden"
    ></div>

    <!-- Sidebar -->
    <aside 
      class="fixed inset-y-0 left-0 z-40 w-64 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 transition-transform duration-300 lg:translate-x-0 lg:static lg:inset-0"
      :class="[
        mobileMenuOpen ? 'translate-x-0' : '-translate-x-full',
        sidebarOpen ? 'lg:w-64' : 'lg:w-20'
      ]"
    >
      <div class="h-16 flex items-center px-6 border-b border-slate-200 dark:border-slate-800 shrink-0">
        <div class="w-8 h-8 bg-blue-600 rounded flex items-center justify-center font-bold shrink-0 text-white">G</div>
        <span v-if="sidebarOpen || mobileMenuOpen" class="ml-3 font-semibold tracking-wider uppercase text-sm truncate">Genset Monitor</span>
      </div>
      
      <nav class="flex-1 p-4 space-y-2 overflow-y-auto custom-scrollbar">
        <router-link 
          v-for="item in navigation" 
          :key="item.name"
          :to="item.href"
          class="flex items-center p-3 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors group"
          active-class="bg-blue-50 dark:bg-slate-800 text-blue-600 dark:text-blue-400"
        >
          <svg 
            class="w-6 h-6 shrink-0 transition-colors" 
            :class="(sidebarOpen || mobileMenuOpen) ? 'mr-3' : 'mx-auto'"
            fill="none" 
            viewBox="0 0 24 24" 
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="item.icon" />
          </svg>
          <span v-if="sidebarOpen || mobileMenuOpen" class="truncate font-medium">{{ item.name }}</span>
          
          <!-- Tooltip for collapsed state -->
          <div v-if="!sidebarOpen && !mobileMenuOpen" class="fixed left-20 ml-2 px-2 py-1 bg-slate-800 text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity whitespace-nowrap z-50">
            {{ item.name }}
          </div>
        </router-link>
      </nav>
      
      <div class="p-4 border-t border-slate-200 dark:border-slate-800 hidden lg:block">
        <button 
          @click="toggleSidebar"
          class="w-full p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded flex justify-center items-center transition-colors text-slate-500 dark:text-slate-400 hover:text-blue-600 dark:hover:text-white"
          :title="sidebarOpen ? 'Collapse Sidebar' : 'Expand Sidebar'"
        >
          <svg 
            class="w-6 h-6 transition-transform duration-300" 
            :class="{ 'rotate-180': !sidebarOpen }"
            fill="none" 
            viewBox="0 0 24 24" 
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
          </svg>
          <span v-if="sidebarOpen" class="ml-2 font-medium text-sm">Collapse</span>
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col min-w-0 h-screen overflow-hidden">
      <header class="h-16 border-b border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 flex items-center justify-between px-4 lg:px-6 shrink-0 transition-colors duration-300">
        <div class="flex items-center gap-4">
          <button 
            @click="toggleMobileMenu"
            class="lg:hidden p-2 text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
          >
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          
          <h1 class="text-base lg:text-lg font-semibold text-slate-900 dark:text-white truncate max-w-[150px] sm:max-w-none">
            {{ $route.meta.title || 'Dashboard' }}
          </h1>
          
          <div class="hidden md:flex items-center space-x-2 px-3 py-1 bg-slate-100 dark:bg-slate-950 rounded-full border border-slate-200 dark:border-slate-800">
            <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <span class="text-[10px] font-bold uppercase tracking-wider text-slate-500 dark:text-slate-400">System Live</span>
          </div>
        </div>

        <div class="flex items-center space-x-2 sm:space-x-4">
          <!-- Theme Toggle -->
          <button 
            @click="toggleTheme"
            class="p-2 text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
            :title="mode === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
          >
            <svg v-if="mode === 'dark'" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 9H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707M16.95 16.95l.707.707M7.05 7.05l.707-.707M12 8a4 4 0 100 8 4 4 0 000-8z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          </button>

          <!-- User Profile & Username -->
          <div class="flex items-center space-x-2 sm:space-x-3 pr-2 sm:pr-4 border-r border-slate-200 dark:border-slate-800">
            <div class="text-right hidden sm:block">
              <p class="text-sm font-semibold text-slate-900 dark:text-white leading-none">{{ authStore.user?.username || 'Guest' }}</p>
              <p class="text-[10px] text-slate-500 uppercase tracking-tighter mt-1">{{ authStore.user?.role || 'user' }}</p>
            </div>
            <div class="w-8 h-8 sm:w-10 sm:h-10 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-full flex items-center justify-center text-white font-bold shadow-lg shadow-blue-900/20 border-2 border-white dark:border-slate-800">
              {{ (authStore.user?.username || 'G').charAt(0).toUpperCase() }}
            </div>
          </div>

          <!-- Logout Button -->
          <button 
            @click="handleLogout"
            class="p-2 text-slate-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-400/10 rounded-lg transition-all duration-200 group relative"
            title="Logout"
          >
            <svg class="w-5 h-5 sm:w-6 sm:h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
          </button>
        </div>
      </header>
      
      <section class="flex-1 overflow-y-auto p-4 sm:p-6 lg:p-8 custom-scrollbar bg-slate-50 dark:bg-slate-950 transition-colors duration-300">
        <div class="max-w-7xl mx-auto">
          <slot />
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  @apply bg-transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  @apply bg-slate-300 dark:bg-slate-800 rounded-full;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  @apply bg-slate-400 dark:bg-slate-700;
}

.router-link-active svg {
  @apply text-blue-400;
}
</style>
