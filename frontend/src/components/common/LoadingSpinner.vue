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
  gap: var(--space-md);
}

.spinner-container.spinner--fullscreen {
  position: fixed;
  inset: 0;
  background: var(--color-overlay-light);
  backdrop-filter: blur(8px);
  z-index: 9998;
}

.spinner {
  border: 3px solid var(--color-gray-200);
  border-radius: var(--radius-full);
  border-top-color: var(--color-primary);
  border-right-color: var(--color-secondary);
  animation: spin 0.8s var(--ease-smooth) infinite;
  box-shadow: var(--shadow-primary);
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
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  font-weight: var(--weight-medium);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
