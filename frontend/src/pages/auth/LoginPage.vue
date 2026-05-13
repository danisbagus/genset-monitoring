<script setup lang="ts">
import LoginHero from '@/components/auth/LoginHero.vue';
import LoginForm from '@/components/auth/LoginForm.vue';
</script>

<template>
  <div class="min-h-screen flex bg-slate-950 font-sans">
    <!-- Left Section: Visual/Info (Hidden on small screens) -->
    <LoginHero class="w-1/2" />

    <!-- Right Section: Login Form -->
    <div class="flex-1 flex flex-col items-center justify-center p-6 sm:p-12 lg:p-24 overflow-y-auto">
      <div class="w-full max-w-md">
        <!-- Mobile Branding (Only visible on small screens) -->
        <div class="lg:hidden flex items-center gap-3 mb-10">
          <div class="w-10 h-10 bg-indigo-600 rounded flex items-center justify-center">
            <svg viewBox="0 0 24 24" class="w-6 h-6 text-white" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
            </svg>
          </div>
          <h1 class="text-xl font-bold tracking-tight text-white uppercase">Genset <span class="text-slate-400">Monitoring</span></h1>
        </div>

        <div class="mb-10">
          <h2 class="text-3xl font-bold text-white mb-2">Access Portal</h2>
          <p class="text-slate-500 font-medium">Please enter your credentials to access the system.</p>
        </div>

        <!-- Login Form Component -->
        <LoginForm />

        <!-- Footer -->
        <div class="mt-12 flex flex-col sm:flex-row justify-between items-center gap-4 text-xs text-slate-600 font-medium">
          <p>© 2026 Industrial Monitoring Systems</p>
          <div class="flex gap-6">
            <a href="#" class="hover:text-slate-400 transition-colors">Privacy Policy</a>
            <a href="#" class="hover:text-slate-400 transition-colors">Terms of Service</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Ensure the page takes full height even with scroll */
:global(html, body) {
  height: 100%;
}
</style>
