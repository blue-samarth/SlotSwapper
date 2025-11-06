<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="buttonClasses"
    @click="handleClick"
  >
    <span v-if="loading" class="button-spinner"></span>
    <slot v-else></slot>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue';

type ButtonVariant = 'primary' | 'secondary' | 'danger' | 'success';
type ButtonSize = 'sm' | 'md' | 'lg';

interface Props {
  variant?: ButtonVariant;
  size?: ButtonSize;
  type?: 'button' | 'submit' | 'reset';
  disabled?: boolean;
  loading?: boolean;
  fullWidth?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false,
  fullWidth: false,
});

const emit = defineEmits<{
  click: [event: MouseEvent];
}>();

const handleClick = (event: MouseEvent) => {
  if (!props.disabled && !props.loading) {
    emit('click', event);
  }
};

const buttonClasses = computed(() => {
  const classes = ['app-button'];
  
  // Variant classes
  classes.push(`app-button--${props.variant}`);
  
  // Size classes
  classes.push(`app-button--${props.size}`);
  
  // State classes
  if (props.disabled) classes.push('app-button--disabled');
  if (props.loading) classes.push('app-button--loading');
  if (props.fullWidth) classes.push('app-button--full-width');
  
  return classes.join(' ');
});
</script>

<style scoped>
/* Base button styles */
.app-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 500;
  border-radius: 0.375rem;
  transition: all 0.2s;
  cursor: pointer;
  border: 1px solid transparent;
  font-family: inherit;
}

.app-button:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.3);
}

/* Size variants */
.app-button--sm {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
}

.app-button--md {
  padding: 0.625rem 1.25rem;
  font-size: 0.9375rem;
}

.app-button--lg {
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
}

/* Color variants */
.app-button--primary {
  background-color: #3b82f6;
  color: white;
}

.app-button--primary:hover:not(.app-button--disabled):not(.app-button--loading) {
  background-color: #2563eb;
}

.app-button--secondary {
  background-color: #6b7280;
  color: white;
}

.app-button--secondary:hover:not(.app-button--disabled):not(.app-button--loading) {
  background-color: #4b5563;
}

.app-button--danger {
  background-color: #ef4444;
  color: white;
}

.app-button--danger:hover:not(.app-button--disabled):not(.app-button--loading) {
  background-color: #dc2626;
}

.app-button--success {
  background-color: #10b981;
  color: white;
}

.app-button--success:hover:not(.app-button--disabled):not(.app-button--loading) {
  background-color: #059669;
}

/* State variants */
.app-button--disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.app-button--loading {
  cursor: not-allowed;
  position: relative;
}

.app-button--full-width {
  width: 100%;
}

/* Loading spinner */
.button-spinner {
  display: inline-block;
  width: 1rem;
  height: 1rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
