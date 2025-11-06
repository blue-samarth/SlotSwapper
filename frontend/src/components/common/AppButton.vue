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
  font-weight: var(--weight-semibold);
  border-radius: var(--radius-md);
  transition: var(--transition-all);
  cursor: pointer;
  border: 1px solid transparent;
  font-family: var(--font-primary);
}

.app-button:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(255, 107, 107, 0.3);
}

/* Size variants */
.app-button--sm {
  padding: var(--space-sm) var(--space-md);
  font-size: var(--text-sm);
}

.app-button--md {
  padding: 0.625rem 1.25rem;
  font-size: var(--text-base);
}

.app-button--lg {
  padding: var(--space-md) var(--space-lg);
  font-size: var(--text-lg);
}

/* Color variants with gradients */
.app-button--primary {
  background: var(--gradient-primary-btn);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-primary);
}

.app-button--primary:hover:not(.app-button--disabled):not(.app-button--loading) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

.app-button--secondary {
  background: var(--gradient-secondary-btn);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-secondary);
}

.app-button--secondary:hover:not(.app-button--disabled):not(.app-button--loading) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

.app-button--danger {
  background: var(--color-error);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-primary);
}

.app-button--danger:hover:not(.app-button--disabled):not(.app-button--loading) {
  background: var(--color-error-dark);
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

.app-button--success {
  background: var(--color-success);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-secondary);
}

.app-button--success:hover:not(.app-button--disabled):not(.app-button--loading) {
  background: var(--color-success-dark);
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
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
