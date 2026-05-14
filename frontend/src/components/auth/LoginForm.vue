<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuth } from '@/composables/useAuth';

const router = useRouter();
const route = useRoute();
const { login, isLoading } = useAuth();

const username = ref('');
const password = ref('');
const rememberMe = ref(false);
const errors = reactive({
  username: '',
  password: '',
  general: ''
});

const validate = () => {
  let isValid = true;
  errors.username = '';
  errors.password = '';
  errors.general = '';

  if (!username.value) {
    errors.username = 'Username is required';
    isValid = false;
  }

  if (!password.value) {
    errors.password = 'Password is required';
    isValid = false;
  } else if (password.value.length < 6) {
    errors.password = 'Password must be at least 6 characters';
    isValid = false;
  }

  return isValid;
};

const handleSubmit = async () => {
  if (!validate()) return;

  const result = await login({
    username: username.value,
    password: password.value
  });

  if (result.success) {
    const redirectPath = (route.query.redirect as string) || '/';
    router.push(redirectPath);
  } else {
    errors.general = result.message;
  }
};
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-6">
    <!-- General Error Alert -->
    <div v-if="errors.general" class="p-3 bg-red-500/10 border border-red-500/20 rounded text-red-500 text-sm flex items-center gap-2">
      <svg viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4">
        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" clip-rule="evenodd" />
      </svg>
      {{ errors.general }}
    </div>

    <!-- Username Field -->
    <div class="space-y-1.5">
      <label for="username" class="block text-xs font-bold uppercase tracking-widest text-slate-400">
        Username
      </label>
      <div class="relative group">
        <input 
          id="username"
          v-model="username"
          type="text"
          placeholder="admin_genset"
          :disabled="isLoading"
          class="w-full bg-slate-900 border px-4 py-3 rounded text-slate-100 placeholder:text-slate-600 focus:outline-none focus:ring-1 transition-all"
          :class="[
            errors.username 
              ? 'border-red-500/50 focus:border-red-500 focus:ring-red-500/20' 
              : 'border-slate-800 focus:border-indigo-500 focus:ring-indigo-500/20'
          ]"
        />
        <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
          <svg viewBox="0 0 24 24" class="w-5 h-5 text-slate-600" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" />
            <circle cx="12" cy="7" r="4" />
          </svg>
        </div>
      </div>
      <p v-if="errors.username" class="text-red-500 text-[11px] font-medium mt-1">{{ errors.username }}</p>
    </div>

    <!-- Password Field -->
    <div class="space-y-1.5">
      <div class="flex justify-between items-center">
        <label for="password" class="block text-xs font-bold uppercase tracking-widest text-slate-400">
          System Password
        </label>
        <a href="#" class="text-[11px] font-semibold text-indigo-400 hover:text-indigo-300 transition-colors">
          Forgot Password?
        </a>
      </div>
      <div class="relative group">
        <input 
          id="password"
          v-model="password"
          type="password"
          placeholder="••••••••"
          :disabled="isLoading"
          class="w-full bg-slate-900 border px-4 py-3 rounded text-slate-100 placeholder:text-slate-600 focus:outline-none focus:ring-1 transition-all"
          :class="[
            errors.password 
              ? 'border-red-500/50 focus:border-red-500 focus:ring-red-500/20' 
              : 'border-slate-800 focus:border-indigo-500 focus:ring-indigo-500/20'
          ]"
        />
        <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
          <svg viewBox="0 0 24 24" class="w-5 h-5 text-slate-600" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-11V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </div>
      </div>
      <p v-if="errors.password" class="text-red-500 text-[11px] font-medium mt-1">{{ errors.password }}</p>
    </div>

    <!-- Options -->
    <div class="flex items-center">
      <label class="relative flex items-center cursor-pointer group">
        <input 
          type="checkbox" 
          v-model="rememberMe"
          :disabled="isLoading"
          class="sr-only peer"
        />
        <div class="w-4 h-4 bg-slate-900 border border-slate-700 rounded peer-checked:bg-indigo-600 peer-checked:border-indigo-600 transition-all flex items-center justify-center">
          <svg viewBox="0 0 24 24" class="w-3 h-3 text-white scale-0 peer-checked:scale-100 transition-transform" fill="none" stroke="currentColor" stroke-width="3">
            <polyline points="20 6 9 17 4 12" />
          </svg>
        </div>
        <span class="ml-2 text-xs text-slate-400 font-medium group-hover:text-slate-300 transition-colors">Keep me authenticated for 30 days</span>
      </label>
    </div>

    <!-- Submit Button -->
    <button 
      type="submit"
      :disabled="isLoading"
      class="w-full py-3.5 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-70 disabled:hover:bg-indigo-600 text-white font-bold rounded shadow-lg shadow-indigo-600/20 transition-all flex items-center justify-center gap-2 group"
    >
      <template v-if="isLoading">
        <svg class="animate-spin h-4 w-4 text-white" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Authenticating...
      </template>
      <template v-else>
        Sign In to Control Center
        <svg viewBox="0 0 24 24" class="w-4 h-4 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 5l7 7m0 0l-7 7m7-7H3" />
        </svg>
      </template>
    </button>

    <!-- System Info -->
    <div class="pt-6 border-t border-slate-800">
      <div class="flex items-center gap-2 text-[10px] text-slate-500 font-mono tracking-tight">
        <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"></span>
        SYSTEM SECURE | SSL ENCRYPTED | ENDPOINT VALIDATED
      </div>
    </div>
  </form>
</template>
