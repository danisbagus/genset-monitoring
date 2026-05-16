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
  const base = 'inline-flex items-center font-medium rounded-full'
  const variants = {
    primary: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
    success: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400',
    warning: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400',
    danger: 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-400',
    info: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400',
    neutral: 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-400'
  }
  
  const sizes = {
    sm: 'px-2 py-0.5 text-xs',
    md: 'px-2.5 py-0.5 text-sm',
    lg: 'px-3 py-1 text-base'
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
    neutral: 'bg-slate-500'
  }
  return `w-1.5 h-1.5 rounded-full mr-1.5 ${variants[props.variant]}`
})
</script>

<template>
  <span :class="variantClasses">
    <span v-if="dot" :class="dotClasses"></span>
    <slot />
  </span>
</template>
