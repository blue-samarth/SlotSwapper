<script setup lang="ts">
import { computed } from 'vue';
import { useSwapsStore } from '../..//stores/swaps_store';
import { SwapStatus } from '../..//types/swap_types';
import { formatDate, formatTime } from '../../utils/date';

const swapsStore = useSwapsStore();

const sentRequests = computed(() => swapsStore.sentRequests || []);

const getStatusClass = (status: string) => {
  switch (status) {
    case SwapStatus.PENDING:
      return 'sent-swaps__status--pending';
    case SwapStatus.ACCEPTED:
      return 'sent-swaps__status--accepted';
    case SwapStatus.REJECTED:
      return 'sent-swaps__status--rejected';
    case SwapStatus.CANCELLED:
      return 'sent-swaps__status--cancelled';
    default:
      return 'sent-swaps__status--default';
  }
};

const getStatusIcon = (status: string) => {
  switch (status) {
    case SwapStatus.PENDING:
      return '⏳';
    case SwapStatus.ACCEPTED:
      return '✅';
    case SwapStatus.REJECTED:
      return '❌';
    case SwapStatus.CANCELLED:
      return '🚫';
    default:
      return '❓';
  }
};

const handleCancelSwap = async (id: number) => {
  if (confirm('Are you sure you want to cancel this swap request?')) {
    try {
      await swapsStore.cancelSwap(id);
    } catch (error: any) {
      alert(error.response?.data?.message || 'Failed to cancel swap request');
    }
  }
};
</script>

<template>
  <div>
    <div v-if="swapsStore.isLoading" class="sent-swaps__loading">
      <div class="sent-swaps__spinner"></div>
      <div class="sent-swaps__loading-text">Loading sent requests...</div>
    </div>

    <div v-else-if="sentRequests.length === 0" class="sent-swaps__empty">
      <div class="sent-swaps__empty-icon">📤</div>
      <div class="sent-swaps__empty-title">No sent swap requests yet</div>
      <p class="sent-swaps__empty-text">Browse available slots to initiate your first swap!</p>
    </div>

    <div v-else>
      <div class="sent-swaps__count">
        {{ sentRequests.length }} sent request{{ sentRequests.length !== 1 ? 's' : '' }}
      </div>

      <div class="sent-swaps__list">
        <div
          v-for="(request, index) in sentRequests"
          :key="request?.id || `request-${index}`"
          class="sent-swaps__card"
        >
          <div class="sent-swaps__header">
            <h3 class="sent-swaps__title">Swap Request</h3>
            <span v-if="request?.status" :class="['sent-swaps__status', getStatusClass(request.status)]">
              {{ getStatusIcon(request.status) }} {{ request.status }}
            </span>
          </div>

          <div v-if="request?.requester_event && request?.receiver_event" class="sent-swaps__events">
            <!-- Your Event -->
            <div class="sent-swaps__event sent-swaps__event--yours">
              <div class="sent-swaps__event-label">Your Event</div>
              <h4 class="sent-swaps__event-title">{{ request.requester_event.title || 'Untitled Event' }}</h4>
              <div v-if="request.requester_event.start_time" class="sent-swaps__event-details">
                <div>📅 {{ formatDate(request.requester_event.start_time) }}</div>
                <div>🕐 {{ formatTime(request.requester_event.start_time) }}</div>
              </div>
            </div>

            <!-- Their Event -->
            <div class="sent-swaps__event sent-swaps__event--theirs">
              <div class="sent-swaps__event-label">Their Event</div>
              <h4 class="sent-swaps__event-title">{{ request.receiver_event.title || 'Untitled Event' }}</h4>
              <div v-if="request.receiver_event.start_time" class="sent-swaps__event-details">
                <div>📅 {{ formatDate(request.receiver_event.start_time) }}</div>
                <div>🕐 {{ formatTime(request.receiver_event.start_time) }}</div>
              </div>
            </div>
          </div>

          <div class="sent-swaps__footer">
            <div v-if="request.created_at" class="sent-swaps__date">
              Sent on {{ formatDate(request.created_at) }}
            </div>
            <button
              v-if="request.status === SwapStatus.PENDING && request?.id"
              @click="handleCancelSwap(request.id)"
              class="sent-swaps__cancel-btn"
            >
              Cancel Request
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sent-swaps__loading {
  text-align: center;
  padding: var(--space-3xl) 0;
}

.sent-swaps__spinner {
  display: inline-block;
  width: 2rem;
  height: 2rem;
  border: 3px solid var(--color-gray-200);
  border-radius: var(--radius-full);
  border-top-color: var(--color-secondary);
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.sent-swaps__loading-text {
  color: var(--color-text-secondary);
  margin-top: var(--space-md);
}

.sent-swaps__empty {
  text-align: center;
  padding: var(--space-3xl) 0;
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
}

.sent-swaps__empty-icon {
  font-size: 4rem;
  margin-bottom: var(--space-md);
}

.sent-swaps__empty-title {
  color: var(--color-text-primary);
  font-size: var(--text-lg);
  font-weight: var(--weight-semibold);
}

.sent-swaps__empty-text {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin-top: var(--space-sm);
}

.sent-swaps__count {
  margin-bottom: var(--space-md);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  font-weight: var(--weight-medium);
}

.sent-swaps__list {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.sent-swaps__card {
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-lg);
  transition: var(--transition-all);
  border: 1px solid var(--color-gray-200);
}

.sent-swaps__card:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}

.sent-swaps__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-md);
}

.sent-swaps__title {
  font-weight: var(--weight-semibold);
  font-size: var(--text-lg);
  color: var(--color-text-primary);
}

.sent-swaps__status {
  padding: var(--space-xs) var(--space-md);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: var(--weight-semibold);
  border: 2px solid;
}

.sent-swaps__status--pending {
  background: var(--status-pending-bg);
  color: var(--status-pending-text);
  border-color: var(--status-pending-border);
}

.sent-swaps__status--accepted {
  background: var(--status-swappable-bg);
  color: var(--status-swappable-text);
  border-color: var(--status-swappable-border);
}

.sent-swaps__status--rejected {
  background: var(--status-busy-bg);
  color: var(--status-busy-text);
  border-color: var(--status-busy-border);
}

.sent-swaps__status--cancelled,
.sent-swaps__status--default {
  background: var(--color-gray-100);
  color: var(--color-gray-700);
  border-color: var(--color-gray-300);
}

.sent-swaps__events {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: var(--space-lg);
  margin-bottom: var(--space-md);
}

@media (min-width: 768px) {
  .sent-swaps__events {
    grid-template-columns: repeat(2, 1fr);
  }
}

.sent-swaps__event {
  border-left: 4px solid;
  padding-left: var(--space-md);
  border-radius: 0 var(--radius-lg) var(--radius-lg) 0;
  padding: var(--space-md);
}

.sent-swaps__event--yours {
  border-left-color: var(--color-secondary);
  background: linear-gradient(135deg, rgba(78, 205, 196, 0.1), transparent);
}

.sent-swaps__event--theirs {
  border-left-color: var(--color-success);
  background: linear-gradient(135deg, rgba(149, 225, 211, 0.1), transparent);
}

.sent-swaps__event-label {
  font-size: var(--text-xs);
  font-weight: var(--weight-semibold);
  text-transform: uppercase;
  margin-bottom: var(--space-xs);
}

.sent-swaps__event--yours .sent-swaps__event-label {
  color: var(--color-secondary);
}

.sent-swaps__event--theirs .sent-swaps__event-label {
  color: var(--color-success);
}

.sent-swaps__event-title {
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  font-size: var(--text-lg);
}

.sent-swaps__event-details {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin-top: var(--space-sm);
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.sent-swaps__footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--space-md);
  border-top: 1px solid var(--color-gray-200);
}

.sent-swaps__date {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.sent-swaps__cancel-btn {
  padding: var(--space-sm) var(--space-md);
  font-size: var(--text-sm);
  font-weight: var(--weight-medium);
  color: var(--color-error);
  background: transparent;
  border: 1px solid var(--color-error);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: var(--transition-all);
}

.sent-swaps__cancel-btn:hover {
  background: var(--color-error);
  color: white;
  transform: translateY(-2px);
}
</style>