<script setup lang="ts">
import LoginForm from '@/components/auth/LoginForm.vue';
</script>

<template>
  <div class="w-full max-w-md">
    <!-- Branding -->
    <div class="flex items-center justify-center gap-3 mb-10">
      <div class="w-12 h-12 bg-blue-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-900/30">
        <svg viewBox="0 0 24 24" class="w-7 h-7 text-white" fill="none" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
        </svg>
      </div>
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-white uppercase leading-none">Genset</h1>
        <span class="text-xs font-bold text-slate-500 uppercase tracking-[0.2em]">Monitoring</span>
      </div>
    </div>

    <div class="bg-slate-900/50 border border-slate-800 p-8 rounded-2xl shadow-2xl backdrop-blur-sm">
      <div class="mb-8 text-center">
        <h2 class="text-3xl font-bold text-white mb-2">Welcome Back</h2>
        <p class="text-slate-400">Please enter your credentials to access the system.</p>
      </div>

      <!-- Login Form Component -->
      <LoginForm />
    </div>

    <!-- Footer -->
    <div class="mt-8 flex flex-col items-center gap-4 text-xs text-slate-600 font-medium">
      <p>© 2026 Industrial Monitoring Systems</p>
      <div class="flex gap-6">
        <a href="#" class="hover:text-slate-400 transition-colors">Privacy Policy</a>
        <a href="#" class="hover:text-slate-400 transition-colors">Terms of Service</a>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Any specific page transitions or animations can go here */
</style>
