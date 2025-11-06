<template>
  <div :class="['spinner-container', { 'spinner--fullscreen': fullscreen }]">
    <div :class="['spinner', `spinner--${size}`]"></div>
    <p v-if="message" class="spinner-message">{{ message }}</p>
  </div>
</template>

<script setup lang="ts">
type SpinnerSize = 'sm' | 'md' | 'lg';

interface Props {
  size?: SpinnerSize;
  message?: string;
  fullscreen?: boolean;
}

withDefaults(defineProps<Props>(), {
  size: 'md',
  message: '',
  fullscreen: false,
});
</script>

<style scoped>
.spinner-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}

.spinner-container.spinner--fullscreen {
  position: fixed;
  inset: 0;
  background: rgba(255, 255, 255, 0.9);
  z-index: 9998;
}

.spinner {
  border: 3px solid rgba(59, 130, 246, 0.2);
  border-radius: 50%;
  border-top-color: #3b82f6;
  animation: spin 0.8s linear infinite;
}

.spinner--sm {
  width: 1.5rem;
  height: 1.5rem;
  border-width: 2px;
}

.spinner--md {
  width: 2.5rem;
  height: 2.5rem;
  border-width: 3px;
}

.spinner--lg {
  width: 4rem;
  height: 4rem;
  border-width: 4px;
}

.spinner-message {
  margin: 0;
  font-size: 0.875rem;
  color: #6b7280;
  font-weight: 500;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
