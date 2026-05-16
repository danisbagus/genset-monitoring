<script setup lang="ts">
defineProps<{
  title?: string
  subtitle?: string
  noPadding?: boolean
}>()
</script>

<template>
  <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm hover:shadow-xl hover:border-blue-500/30 dark:hover:border-blue-400/30 transition-all duration-300 overflow-hidden flex flex-col">
    <div v-if="title || $slots.header" class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 flex justify-between items-center bg-slate-50/30 dark:bg-slate-900/30">
      <div>
        <h3 v-if="title" class="text-base sm:text-lg font-bold text-slate-900 dark:text-slate-100 uppercase tracking-tight">{{ title }}</h3>
        <p v-if="subtitle" class="text-[10px] sm:text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-widest mt-0.5">{{ subtitle }}</p>
      </div>
      <div v-if="$slots.headerAction" class="shrink-0">
        <slot name="headerAction" />
      </div>
    </div>
    <div class="flex-1" :class="[noPadding ? '' : 'p-5 sm:p-6']">
      <slot />
    </div>
    <div v-if="$slots.footer" class="px-6 py-4 bg-slate-50/50 dark:bg-slate-800/20 border-t border-slate-100 dark:border-slate-800">
      <slot name="footer" />
    </div>
  </div>
</template>
