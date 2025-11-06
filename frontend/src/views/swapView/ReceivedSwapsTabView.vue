<script setup lang="ts">
import { computed } from 'vue';
import { useSwapsStore } from '../..//stores/swaps_store';
import { SwapStatus } from '../..//types/swap_types';
import { formatDate, formatTime } from '../../utils/date';

const swapsStore = useSwapsStore();

const receivedRequests = computed(() => swapsStore.receivedRequests || []);

const getStatusClass = (status: string) => {
  switch (status) {
    case SwapStatus.PENDING:
      return 'received-swaps__status--pending';
    case SwapStatus.ACCEPTED:
      return 'received-swaps__status--accepted';
    case SwapStatus.REJECTED:
      return 'received-swaps__status--rejected';
    case SwapStatus.CANCELLED:
      return 'received-swaps__status--cancelled';
    default:
      return 'received-swaps__status--default';
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

const handleAcceptSwap = async (id: number) => {
  if (confirm('Accept this swap request? Both events will be exchanged.')) {
    try {
      await swapsStore.acceptSwap(id);
      alert('Swap request accepted! Events have been exchanged.');
    } catch (error: any) {
      alert(error.response?.data?.message || 'Failed to accept swap request');
    }
  }
};

const handleRejectSwap = async (id: number) => {
  if (confirm('Reject this swap request?')) {
    try {
      await swapsStore.rejectSwap(id);
      alert('Swap request rejected.');
    } catch (error: any) {
      alert(error.response?.data?.message || 'Failed to reject swap request');
    }
  }
};
</script>

<template>
  <div>
    <div v-if="swapsStore.isLoading" class="received-swaps__loading">
      <div class="received-swaps__spinner"></div>
      <div class="received-swaps__loading-text">Loading received requests...</div>
    </div>

    <div v-else-if="receivedRequests.length === 0" class="received-swaps__empty">
      <div class="received-swaps__empty-icon">📥</div>
      <div class="received-swaps__empty-title">No received swap requests yet</div>
      <p class="received-swaps__empty-text">When someone requests to swap with you, they'll appear here.</p>
    </div>

    <div v-else>
      <div class="received-swaps__count">
        {{ receivedRequests.length }} received request{{ receivedRequests.length !== 1 ? 's' : '' }}
      </div>

      <div class="received-swaps__list">
        <div
          v-for="(request, index) in receivedRequests"
          :key="request?.id || `request-${index}`"
          class="received-swaps__card"
        >
          <div class="received-swaps__header">
            <h3 class="received-swaps__title">Swap Request Received</h3>
            <span v-if="request?.status" :class="['received-swaps__status', getStatusClass(request.status)]">
              {{ getStatusIcon(request.status) }} {{ request.status }}
            </span>
          </div>

          <div v-if="request?.requester_event && request?.receiver_event" class="received-swaps__events">
            <!-- Their Event -->
            <div class="received-swaps__event received-swaps__event--theirs">
              <div class="received-swaps__event-label">They Offer</div>
              <h4 class="received-swaps__event-title">{{ request.requester_event.title || 'Untitled Event' }}</h4>
              <div v-if="request.requester_event.start_time" class="received-swaps__event-details">
                <div>📅 {{ formatDate(request.requester_event.start_time) }}</div>
                <div>🕐 {{ formatTime(request.requester_event.start_time) }}</div>
              </div>
            </div>

            <!-- Your Event -->
            <div class="received-swaps__event received-swaps__event--yours">
              <div class="received-swaps__event-label">Your Event</div>
              <h4 class="received-swaps__event-title">{{ request.receiver_event.title || 'Untitled Event' }}</h4>
              <div v-if="request.receiver_event.start_time" class="received-swaps__event-details">
                <div>📅 {{ formatDate(request.receiver_event.start_time) }}</div>
                <div>🕐 {{ formatTime(request.receiver_event.start_time) }}</div>
              </div>
            </div>
          </div>

          <div class="received-swaps__footer">
            <div v-if="request.created_at" class="received-swaps__date">
              Received on {{ formatDate(request.created_at) }}
            </div>
            <div v-if="request.status === SwapStatus.PENDING && request?.id" class="received-swaps__actions">
              <button
                @click="handleRejectSwap(request.id)"
                class="received-swaps__reject-btn"
              >
                Reject
              </button>
              <button
                @click="handleAcceptSwap(request.id)"
                class="received-swaps__accept-btn"
              >
                Accept Swap
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.received-swaps__loading {
  text-align: center;
  padding: var(--space-3xl) 0;
}

.received-swaps__spinner {
  display: inline-block;
  width: 2rem;
  height: 2rem;
  border: 3px solid var(--color-gray-200);
  border-radius: var(--radius-full);
  border-top-color: var(--color-success);
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.received-swaps__loading-text {
  color: var(--color-text-secondary);
  margin-top: var(--space-md);
}

.received-swaps__empty {
  text-align: center;
  padding: var(--space-3xl) 0;
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
}

.received-swaps__empty-icon {
  font-size: 4rem;
  margin-bottom: var(--space-md);
}

.received-swaps__empty-title {
  color: var(--color-text-primary);
  font-size: var(--text-lg);
  font-weight: var(--weight-semibold);
}

.received-swaps__empty-text {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin-top: var(--space-sm);
}

.received-swaps__count {
  margin-bottom: var(--space-md);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  font-weight: var(--weight-medium);
}

.received-swaps__list {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.received-swaps__card {
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-lg);
  transition: var(--transition-all);
  border: 1px solid var(--color-gray-200);
}

.received-swaps__card:hover {
  box-shadow: var(--shadow-xl);
  transform: translateY(-2px);
}

.received-swaps__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-md);
}

.received-swaps__title {
  font-weight: var(--weight-semibold);
  font-size: var(--text-lg);
  color: var(--color-text-primary);
}

.received-swaps__status {
  padding: var(--space-xs) var(--space-md);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: var(--weight-semibold);
  border: 2px solid;
}

.received-swaps__status--pending {
  background: var(--status-pending-bg);
  color: var(--status-pending-text);
  border-color: var(--status-pending-border);
}

.received-swaps__status--accepted {
  background: var(--status-swappable-bg);
  color: var(--status-swappable-text);
  border-color: var(--status-swappable-border);
}

.received-swaps__status--rejected {
  background: var(--status-busy-bg);
  color: var(--status-busy-text);
  border-color: var(--status-busy-border);
}

.received-swaps__status--cancelled,
.received-swaps__status--default {
  background: var(--color-gray-100);
  color: var(--color-gray-700);
  border-color: var(--color-gray-300);
}

.received-swaps__events {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: var(--space-lg);
  margin-bottom: var(--space-md);
}

@media (min-width: 768px) {
  .received-swaps__events {
    grid-template-columns: repeat(2, 1fr);
  }
}

.received-swaps__event {
  border-left: 4px solid;
  padding-left: var(--space-md);
  border-radius: 0 var(--radius-lg) var(--radius-lg) 0;
  padding: var(--space-md);
}

.received-swaps__event--theirs {
  border-left-color: var(--color-primary);
  background: linear-gradient(135deg, rgba(255, 107, 107, 0.1), transparent);
}

.received-swaps__event--yours {
  border-left-color: var(--color-secondary);
  background: linear-gradient(135deg, rgba(78, 205, 196, 0.1), transparent);
}

.received-swaps__event-label {
  font-size: var(--text-xs);
  font-weight: var(--weight-semibold);
  text-transform: uppercase;
  margin-bottom: var(--space-xs);
}

.received-swaps__event--theirs .received-swaps__event-label {
  color: var(--color-primary);
}

.received-swaps__event--yours .received-swaps__event-label {
  color: var(--color-secondary);
}

.received-swaps__event-title {
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  font-size: var(--text-lg);
}

.received-swaps__event-details {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin-top: var(--space-sm);
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.received-swaps__footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--space-md);
  border-top: 1px solid var(--color-gray-200);
}

.received-swaps__date {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.received-swaps__actions {
  display: flex;
  gap: var(--space-sm);
}

.received-swaps__reject-btn {
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

.received-swaps__reject-btn:hover {
  background: var(--color-error);
  color: white;
  transform: translateY(-2px);
}

.received-swaps__accept-btn {
  padding: var(--space-sm) var(--space-md);
  font-size: var(--text-sm);
  font-weight: var(--weight-semibold);
  color: white;
  background: var(--gradient-secondary-btn);
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: var(--transition-all);
  box-shadow: var(--shadow-success);
}

.received-swaps__accept-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}
</style>