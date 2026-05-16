<script setup lang="ts">
defineProps<{
  title?: string
  subtitle?: string
  noPadding?: boolean
}>()
</script>

<template>
  <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm hover:shadow-md transition-shadow duration-300 overflow-hidden">
    <div v-if="title || $slots.header" class="px-6 py-4 border-b border-slate-100 dark:border-slate-800 flex justify-between items-center">
      <div>
        <h3 v-if="title" class="text-lg font-semibold text-slate-800 dark:text-slate-100">{{ title }}</h3>
        <p v-if="subtitle" class="text-sm text-slate-500 dark:text-slate-400">{{ subtitle }}</p>
      </div>
      <div v-if="$slots.headerAction">
        <slot name="headerAction" />
      </div>
    </div>
    <div :class="[noPadding ? '' : 'p-6']">
      <slot />
    </div>
    <div v-if="$slots.footer" class="px-6 py-4 bg-slate-50 dark:bg-slate-800/50 border-t border-slate-100 dark:border-slate-800">
      <slot name="footer" />
    </div>
  </div>
</template>
