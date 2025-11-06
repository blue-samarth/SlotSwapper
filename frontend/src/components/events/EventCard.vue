<template>
  <div class="event-card">
    <!-- Event Header -->
    <div class="event-card__header">
      <h3 class="event-card__title">{{ event.title }}</h3>
      <span :class="['event-card__status', `event-card__status--${event.status.toLowerCase()}`]">
        {{ event.status }}
      </span>
    </div>

    <!-- Event Details -->
    <div class="event-card__details">
      <div class="event-card__detail">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="event-card__icon">
          <path d="M5.25 12a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H6a.75.75 0 01-.75-.75V12zM6 13.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V14a.75.75 0 00-.75-.75H6zM7.25 12a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H8a.75.75 0 01-.75-.75V12zM8 13.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V14a.75.75 0 00-.75-.75H8zM9.25 10a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H10a.75.75 0 01-.75-.75V10zM10 11.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V12a.75.75 0 00-.75-.75H10zM9.25 14a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H10a.75.75 0 01-.75-.75V14zM12 9.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V10a.75.75 0 00-.75-.75H12zM11.25 12a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H12a.75.75 0 01-.75-.75V12zM12 13.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V14a.75.75 0 00-.75-.75H12zM13.25 10a.75.75 0 01.75-.75h.01a.75.75 0 01.75.75v.01a.75.75 0 01-.75.75H14a.75.75 0 01-.75-.75V10zM14 11.25a.75.75 0 00-.75.75v.01c0 .414.336.75.75.75h.01a.75.75 0 00.75-.75V12a.75.75 0 00-.75-.75H14z" />
          <path fill-rule="evenodd" d="M5.75 2a.75.75 0 01.75.75V4h7V2.75a.75.75 0 011.5 0V4h.25A2.75 2.75 0 0118 6.75v8.5A2.75 2.75 0 0115.25 18H4.75A2.75 2.75 0 012 15.25v-8.5A2.75 2.75 0 014.75 4H5V2.75A.75.75 0 015.75 2zm-1 5.5c-.69 0-1.25.56-1.25 1.25v6.5c0 .69.56 1.25 1.25 1.25h10.5c.69 0 1.25-.56 1.25-1.25v-6.5c0-.69-.56-1.25-1.25-1.25H4.75z" clip-rule="evenodd" />
        </svg>
        <span><strong>Start:</strong> {{ formatDateTime(event.start_time) }}</span>
      </div>
      <div class="event-card__detail">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="event-card__icon">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z" clip-rule="evenodd" />
        </svg>
        <span><strong>Duration:</strong> {{ calculateDuration(event.start_time, event.end_time) }}</span>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="event-card__actions">
      <AppButton
        v-if="event.status !== 'SWAP_PENDING'"
        size="sm"
        variant="secondary"
        :full-width="true"
        @click="emit('toggle-status', event)"
      >
        {{ event.status === 'BUSY' ? 'Make Swappable' : 'Mark Busy' }}
      </AppButton>
      <AppButton
        size="sm"
        variant="danger"
        @click="emit('delete', event)"
      >
        Delete
      </AppButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Event } from '../../types/event_types';
import { formatDateTime, calculateDuration } from '../../utils/date';
import AppButton from '../common/AppButton.vue';

interface Props {
  event: Event;
}

defineProps<Props>();

const emit = defineEmits<{
  'toggle-status': [event: Event];
  delete: [event: Event];
}>();
</script>

<style scoped>
.event-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-lg);
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
  transition: var(--transition-all);
  border: 1px solid var(--color-gray-200);
}

.event-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
  border-color: var(--color-secondary);
}

.event-card__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--space-md);
  margin-bottom: var(--space-sm);
}

.event-card__title {
  font-size: var(--text-lg);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
  flex: 1;
  line-height: var(--leading-tight);
}

.event-card__status {
  padding: var(--space-xs) var(--space-md);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: var(--weight-semibold);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  white-space: nowrap;
  flex-shrink: 0;
  box-shadow: var(--shadow-sm);
}

.event-card__status--busy {
  background: var(--status-busy-bg);
  color: var(--status-busy-text);
  border: 2px solid var(--status-busy-border);
}

.event-card__status--swappable {
  background: var(--status-swappable-bg);
  color: var(--status-swappable-text);
  border: 2px solid var(--status-swappable-border);
}

.event-card__status--swap_pending {
  background: var(--status-pending-bg);
  color: var(--status-pending-text);
  border: 2px solid var(--status-pending-border);
}

.event-card__details {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.event-card__detail {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.event-card__icon {
  width: 1.125rem;
  height: 1.125rem;
  flex-shrink: 0;
  color: var(--color-secondary);
}

.event-card__actions {
  display: flex;
  gap: var(--space-sm);
  margin-top: var(--space-sm);
}
</style>
