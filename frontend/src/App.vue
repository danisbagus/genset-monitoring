<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import MainLayout from '@/components/layout/MainLayout.vue'
import AuthLayout from '@/components/layout/AuthLayout.vue'

const authStore = useAuthStore()
const route = useRoute()

const layout = computed(() => {
  return route.meta.layout === 'auth' ? AuthLayout : MainLayout
})

onMounted(() => {
  authStore.setupStorageListener()
})
</script>

<template>
  <component :is="layout">
    <router-view />
  </component>
</template>

<style>
/* Global styles if needed */
</style>
