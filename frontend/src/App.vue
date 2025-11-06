<script setup lang="ts">
import { RouterView } from 'vue-router';
import AppToast from './components/common/AppToast.vue';
import { setupToast } from './composables/useToast';

// Setup toast notification system
const toast = setupToast();

// Make toast available globally
(window as any).__toast = toast;
</script>

<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <RouterView />
    <AppToast />
  </div>
</template>

<style scoped>
/* Add any global app styles here if needed */
</style>