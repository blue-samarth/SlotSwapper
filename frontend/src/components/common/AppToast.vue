<template>
  <Teleport to="body">
    <div class="toast-container">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="['toast', `toast--${toast.type}`]"
          @click="handleRemove(toast.id)"
        >
          <div class="toast-icon">
            <svg v-if="toast.type === 'success'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z" clip-rule="evenodd" />
            </svg>
            <svg v-else-if="toast.type === 'error'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" clip-rule="evenodd" />
            </svg>
            <svg v-else-if="toast.type === 'warning'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495zM10 5a.75.75 0 01.75.75v3.5a.75.75 0 01-1.5 0v-3.5A.75.75 0 0110 5zm0 9a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a.75.75 0 000 1.5h.253a.25.25 0 01.244.304l-.459 2.066A1.75 1.75 0 0010.747 15H11a.75.75 0 000-1.5h-.253a.25.25 0 01-.244-.304l.459-2.066A1.75 1.75 0 009.253 9H9z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="toast-content">
            <p class="toast-message">{{ toast.message }}</p>
          </div>
          <button class="toast-close" @click.stop="handleRemove(toast.id)">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
            </svg>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { inject, type Ref } from 'vue';
import type { Toast } from '../../composables/useToast';

const toasts = inject<Ref<Toast[]>>('toasts');
const removeToast = inject<(id: string) => void>('removeToast');

// Type guard to ensure functions exist
const handleRemove = (id: string) => {
  if (removeToast) {
    removeToast(id);
  }
};
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: var(--space-md);
  right: var(--space-md);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  max-width: 24rem;
  width: 100%;
  pointer-events: none;
}

.toast {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  padding: var(--space-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  pointer-events: auto;
  cursor: pointer;
  transition: var(--transition-all);
  border-left: 4px solid;
}

.toast:hover {
  transform: translateX(-8px) scale(1.02);
  box-shadow: var(--shadow-xl);
}

.toast--success {
  border-left-color: var(--color-success);
  background: linear-gradient(to right, rgba(149, 225, 211, 0.1), var(--color-bg-card));
  box-shadow: var(--shadow-success);
}

.toast--error {
  border-left-color: var(--color-error);
  background: linear-gradient(to right, rgba(255, 107, 107, 0.1), var(--color-bg-card));
  box-shadow: var(--shadow-primary);
}

.toast--warning {
  border-left-color: var(--color-warning);
  background: linear-gradient(to right, rgba(243, 129, 129, 0.1), var(--color-bg-card));
  box-shadow: var(--shadow-warning);
}

.toast--info {
  border-left-color: var(--color-secondary);
  background: linear-gradient(to right, rgba(78, 205, 196, 0.1), var(--color-bg-card));
  box-shadow: var(--shadow-secondary);
}

.toast-icon {
  flex-shrink: 0;
  width: 1.5rem;
  height: 1.5rem;
}

.toast--success .toast-icon {
  color: var(--color-success);
}

.toast--error .toast-icon {
  color: var(--color-error);
}

.toast--warning .toast-icon {
  color: var(--color-warning);
}

.toast--info .toast-icon {
  color: var(--color-secondary);
}

.toast-content {
  flex: 1;
  min-width: 0;
}

.toast-message {
  margin: 0;
  font-size: var(--text-sm);
  font-weight: var(--weight-medium);
  color: var(--color-text-primary);
  word-wrap: break-word;
}

.toast-close {
  flex-shrink: 0;
  width: 1.25rem;
  height: 1.25rem;
  padding: 0;
  border: none;
  background: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: var(--transition-all);
  border-radius: var(--radius-sm);
}

.toast-close:hover {
  color: var(--color-text-primary);
  background: var(--color-gray-100);
  transform: rotate(90deg);
}

/* Toast animations */
.toast-enter-active,
.toast-leave-active {
  transition: all var(--duration-slow) var(--ease-bounce);
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%) scale(0.8);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%) scale(0.8);
}

.toast-move {
  transition: transform var(--duration-normal) var(--ease-smooth);
}
</style>
