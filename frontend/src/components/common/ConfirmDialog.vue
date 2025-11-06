<template>
  <AppModal
    v-model="isOpen"
    :title="title"
    :close-on-overlay="false"
  >
    <div class="confirm-dialog">
      <div class="confirm-icon" :class="`confirm-icon--${type}`">
        <svg v-if="type === 'danger'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495zM10 5a.75.75 0 01.75.75v3.5a.75.75 0 01-1.5 0v-3.5A.75.75 0 0110 5zm0 9a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a.75.75 0 000 1.5h.253a.25.25 0 01.244.304l-.459 2.066A1.75 1.75 0 0010.747 15H11a.75.75 0 000-1.5h-.253a.25.25 0 01-.244-.304l.459-2.066A1.75 1.75 0 009.253 9H9z" clip-rule="evenodd" />
        </svg>
      </div>
      
      <p class="confirm-message">{{ message }}</p>
      
      <p v-if="subMessage" class="confirm-submessage">{{ subMessage }}</p>
    </div>

    <template #footer>
      <AppButton
        variant="secondary"
        @click="handleCancel"
      >
        {{ cancelText }}
      </AppButton>
      <AppButton
        :variant="type === 'danger' ? 'danger' : 'primary'"
        :loading="isLoading"
        @click="handleConfirm"
      >
        {{ confirmText }}
      </AppButton>
    </template>
  </AppModal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import AppModal from './AppModal.vue';
import AppButton from './AppButton.vue';

type ConfirmType = 'danger' | 'warning' | 'info';

interface Props {
  modelValue: boolean;
  title?: string;
  message: string;
  subMessage?: string;
  type?: ConfirmType;
  confirmText?: string;
  cancelText?: string;
  isLoading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Confirm Action',
  subMessage: '',
  type: 'danger',
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  isLoading: false,
});

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  confirm: [];
  cancel: [];
}>();

const isOpen = ref(props.modelValue);

watch(() => props.modelValue, (value) => {
  isOpen.value = value;
});

watch(isOpen, (value) => {
  emit('update:modelValue', value);
});

const handleConfirm = () => {
  emit('confirm');
};

const handleCancel = () => {
  isOpen.value = false;
  emit('cancel');
};
</script>

<style scoped>
.confirm-dialog {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-md);
  text-align: center;
  padding: var(--space-md) 0;
}

.confirm-icon {
  width: 3.5rem;
  height: 3.5rem;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-md);
  animation: pulse 2s infinite;
}

.confirm-icon--danger {
  background: linear-gradient(135deg, rgba(255, 107, 107, 0.2), rgba(255, 107, 107, 0.1));
  color: var(--color-error);
  border: 2px solid var(--color-error);
}

.confirm-icon--warning {
  background: linear-gradient(135deg, rgba(243, 129, 129, 0.2), rgba(243, 129, 129, 0.1));
  color: var(--color-warning);
  border: 2px solid var(--color-warning);
}

.confirm-icon--info {
  background: linear-gradient(135deg, rgba(78, 205, 196, 0.2), rgba(78, 205, 196, 0.1));
  color: var(--color-secondary);
  border: 2px solid var(--color-secondary);
}

.confirm-icon svg {
  width: 2rem;
  height: 2rem;
}

.confirm-message {
  margin: 0;
  font-size: var(--text-lg);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  line-height: var(--leading-normal);
}

.confirm-submessage {
  margin: 0;
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  line-height: var(--leading-relaxed);
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}
</style>
