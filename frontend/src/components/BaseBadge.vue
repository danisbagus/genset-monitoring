<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  variant?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'neutral'
  size?: 'sm' | 'md' | 'lg'
  dot?: boolean
}>(), {
  variant: 'neutral',
  size: 'md',
  dot: false
})

const variantClasses = computed(() => {
  const base = 'inline-flex items-center font-bold rounded-lg border transition-all duration-300'
  const variants = {
    primary: 'bg-blue-50/50 text-blue-700 border-blue-200 dark:bg-blue-900/10 dark:text-blue-400 dark:border-blue-800/30',
    success: 'bg-emerald-50/50 text-emerald-700 border-emerald-200 dark:bg-emerald-900/10 dark:text-emerald-400 dark:border-emerald-800/30',
    warning: 'bg-amber-50/50 text-amber-700 border-amber-200 dark:bg-amber-900/10 dark:text-amber-400 dark:border-amber-800/30',
    danger: 'bg-rose-50/50 text-rose-700 border-rose-200 dark:bg-rose-900/10 dark:text-rose-400 dark:border-rose-800/30',
    info: 'bg-indigo-50/50 text-indigo-700 border-indigo-200 dark:bg-indigo-900/10 dark:text-indigo-400 dark:border-indigo-800/30',
    neutral: 'bg-slate-100/50 text-slate-700 border-slate-200 dark:bg-slate-800/20 dark:text-slate-300 dark:border-slate-700/50'
  }
  
  const sizes = {
    sm: 'px-2 py-0.5 text-[10px] uppercase tracking-widest',
    md: 'px-3 py-1 text-xs sm:text-sm tracking-wide',
    lg: 'px-4 py-1.5 text-sm sm:text-base tracking-wide'
  }
  
  return `${base} ${variants[props.variant]} ${sizes[props.size]}`
})

const dotClasses = computed(() => {
  const variants = {
    primary: 'bg-blue-500',
    success: 'bg-emerald-500',
    warning: 'bg-amber-500',
    danger: 'bg-rose-500',
    info: 'bg-indigo-500',
    neutral: 'bg-slate-400'
  }
  return `w-1.5 h-1.5 rounded-full mr-2 ${variants[props.variant]}`
})
</script>

<template>
  <span :class="variantClasses">
    <span v-if="dot" :class="dotClasses"></span>
    <slot />
  </span>
</template>
