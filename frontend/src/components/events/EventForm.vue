<template>
  <form @submit.prevent="handleSubmit" class="event-form">
    <!-- Title Field -->
    <div class="form-group">
      <label for="title" class="form-label">Title *</label>
      <input
        id="title"
        v-model="formData.title"
        type="text"
        class="form-input"
        :class="{ 'form-input--error': errors.title }"
        placeholder="e.g., CS101 Lecture"
        required
      />
      <span v-if="errors.title" class="form-error">{{ errors.title }}</span>
    </div>

    <!-- Date Field -->
    <div class="form-group">
      <label for="date" class="form-label">Date *</label>
      <input
        id="date"
        v-model="formData.date"
        type="date"
        class="form-input"
        :class="{ 'form-input--error': errors.date }"
        :min="minDate"
        required
      />
      <span v-if="errors.date" class="form-error">{{ errors.date }}</span>
    </div>

    <!-- Time Field -->
    <div class="form-group">
      <label for="time" class="form-label">Time *</label>
      <input
        id="time"
        v-model="formData.time"
        type="time"
        class="form-input"
        :class="{ 'form-input--error': errors.time }"
        required
      />
      <span v-if="errors.time" class="form-error">{{ errors.time }}</span>
    </div>

    <!-- Duration Field -->
    <div class="form-group">
      <label for="duration" class="form-label">Duration (minutes) *</label>
      <input
        id="duration"
        v-model.number="formData.duration"
        type="number"
        class="form-input"
        :class="{ 'form-input--error': errors.duration }"
        placeholder="e.g., 60"
        min="15"
        max="480"
        step="15"
        required
      />
      <span v-if="errors.duration" class="form-error">{{ errors.duration }}</span>
      <span class="form-help">Between 15 and 480 minutes (8 hours)</span>
    </div>

    <!-- Status Info (Read-only for create) -->
    <div v-if="!isEditMode" class="form-group">
      <div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
        <p class="text-sm text-blue-800">
          <strong>Note:</strong> New events are created as <strong>BUSY</strong> by default. You can change the status to SWAPPABLE after creation.
        </p>
      </div>
    </div>

    <!-- Status Field (Edit mode only) -->
    <div v-else class="form-group">
      <label for="status" class="form-label">Status *</label>
      <select
        id="status"
        v-model="formData.status"
        class="form-input"
        :class="{ 'form-input--error': errors.status }"
        required
      >
        <option value="BUSY">Busy (Not available for swap)</option>
        <option value="SWAPPABLE">Swappable (Available for swap)</option>
      </select>
      <span v-if="errors.status" class="form-error">{{ errors.status }}</span>
    </div>

    <!-- Submit Buttons -->
    <div class="form-actions">
      <AppButton
        type="button"
        variant="secondary"
        @click="emit('cancel')"
      >
        Cancel
      </AppButton>
      <AppButton
        type="submit"
        variant="primary"
        :loading="isSubmitting"
      >
        {{ isEditMode ? 'Update Event' : 'Create Event' }}
      </AppButton>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import type { Event, EventStatus } from '../../types/event_types';
import AppButton from '../common/AppButton.vue';

interface Props {
  event?: Event | null;
  isSubmitting?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  event: null,
  isSubmitting: false,
});

const emit = defineEmits<{
  submit: [data: {
    title: string;
    date: string;
    time: string;
    duration: number;
    status: EventStatus;
  }];
  cancel: [];
}>();

// Form data
const formData = ref({
  title: '',
  date: '',
  time: '',
  duration: 60,
  status: 'SWAPPABLE' as EventStatus,
});

// Errors
const errors = ref<Record<string, string>>({});

// Computed
const isEditMode = computed(() => !!props.event);

const minDate = computed(() => {
  const today = new Date();
  return today.toISOString().split('T')[0];
});

// Validation
const validateForm = (): boolean => {
  errors.value = {};

  if (!formData.value.title.trim()) {
    errors.value.title = 'Title is required';
  } else if (formData.value.title.length > 200) {
    errors.value.title = 'Title must be less than 200 characters';
  }

  if (!formData.value.date) {
    errors.value.date = 'Date is required';
  } else {
    const selectedDate = new Date(formData.value.date);
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    if (selectedDate < today) {
      errors.value.date = 'Date cannot be in the past';
    }
  }

  if (!formData.value.time) {
    errors.value.time = 'Time is required';
  }

  if (!formData.value.duration) {
    errors.value.duration = 'Duration is required';
  } else if (formData.value.duration < 15) {
    errors.value.duration = 'Duration must be at least 15 minutes';
  } else if (formData.value.duration > 480) {
    errors.value.duration = 'Duration must be less than 480 minutes';
  }

  if (!formData.value.status) {
    errors.value.status = 'Status is required';
  }

  return Object.keys(errors.value).length === 0;
};

// Handlers
const handleSubmit = () => {
  if (!validateForm()) {
    return;
  }

  emit('submit', {
    title: formData.value.title.trim(),
    date: formData.value.date,
    time: formData.value.time,
    duration: formData.value.duration,
    status: formData.value.status,
  });
};

// Watch for prop changes (edit mode)
watch(() => props.event, (event) => {
  if (event) {
    // Parse the start_time string (e.g., "2024-11-15T14:30:00Z")
    const startTime = new Date(event.start_time);
    const endTime = new Date(event.end_time);
    
    // Calculate duration in minutes
    const durationMs = endTime.getTime() - startTime.getTime();
    const durationMinutes = Math.round(durationMs / 60000);
    
    // Extract date and time from start_time
    const extractedDate = startTime.toISOString().split('T')[0];
    const extractedTime = startTime.toTimeString().slice(0, 5);

    formData.value = {
      title: event.title,
      date: extractedDate || '',
      time: extractedTime || '',
      duration: durationMinutes,
      status: event.status,
    };
  }
}, { immediate: true });
</script>

<style scoped>
.event-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
}

.form-input {
  padding: 0.625rem 0.875rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.9375rem;
  transition: border-color 0.2s;
  font-family: inherit;
  width: 100%;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-input--error {
  border-color: #ef4444;
}

.form-input--error:focus {
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

.form-error {
  font-size: 0.8125rem;
  color: #ef4444;
  margin-top: -0.25rem;
}

.form-help {
  font-size: 0.8125rem;
  color: #6b7280;
  margin-top: -0.25rem;
}

.form-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
  margin-top: 0.5rem;
}
</style>
