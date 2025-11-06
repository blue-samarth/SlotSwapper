<script setup lang="ts">
import { computed } from 'vue';
import { useSwapsStore } from '../../stores/swaps_store';
import { useEventsStore } from '../../stores/event_store';
import { formatDate, formatTime } from '../../utils/date';

const swapsStore = useSwapsStore();
const eventsStore = useEventsStore();

const availableSlots = computed(() => swapsStore.swappableSlots || []);
const myEvents = computed(() => eventsStore.events?.filter(e => e.status === 'SWAPPABLE') || []);

const handleRequestSwap = async (theirSlotId: number) => {
  // Debug logging to diagnose the validation error
  console.log('=== SWAP REQUEST DEBUG ===');
  console.log('Their slot ID:', theirSlotId);
  console.log('Their slot ID type:', typeof theirSlotId);
  console.log('My events:', myEvents.value);
  console.log('Selected event:', myEvents.value[0]);
  console.log('Selected event ID:', myEvents.value[0]?.id);
  
  // This will be enhanced later with a modal to select which of your events to swap
  if (myEvents.value.length === 0) {
    alert('You need to have at least one SWAPPABLE event to request a swap');
    return;
  }

  const myEvent = myEvents.value[0];
  if (!myEvent) {
    alert('No swappable event found');
    return;
  }

  const requestData = {
    requester_event_id: myEvent.id,
    receiver_event_id: theirSlotId,
  };
  console.log('Request payload:', requestData);

  try {
    await swapsStore.initiateSwap(requestData);
    
    // Refresh the sent requests list to show the new swap immediately
    await swapsStore.fetchSentRequests();
    
    // Also refresh events to update status to SWAP_PENDING
    await eventsStore.fetchEvents();
    
    alert('Swap request sent successfully!');
  } catch (error: any) {
    console.error('Swap request error:', error);
    console.error('Error response:', error.response?.data);
    alert(error.response?.data?.message || 'Failed to send swap request');
  }
};
</script>

<template>
  <div>
    <div v-if="swapsStore.isLoading" class="browse-slots__loading">
      <div class="browse-slots__spinner"></div>
      <div class="browse-slots__loading-text">Loading available slots...</div>
    </div>

    <div v-else-if="availableSlots.length === 0" class="browse-slots__empty">
      <div class="browse-slots__empty-icon">🔍</div>
      <div class="browse-slots__empty-title">No swappable slots available</div>
      <p class="browse-slots__empty-text">Check back later or ask others to mark their events as swappable.</p>
    </div>

    <div v-else>
      <div class="browse-slots__count">
        Found {{ availableSlots.length }} available slot{{ availableSlots.length !== 1 ? 's' : '' }}
      </div>

      <div class="browse-slots__grid">
        <div
          v-for="(slot, index) in availableSlots"
          :key="slot?.id || `slot-${index}`"
          class="browse-slots__card"
        >
          <div class="browse-slots__card-header">
            <h3 class="browse-slots__card-title">{{ slot.title }}</h3>
            <span class="browse-slots__badge">
              SWAPPABLE
            </span>
          </div>

          <div class="browse-slots__card-details">
            <div class="browse-slots__detail">
              <span class="browse-slots__detail-label">👤 Owner:</span>
              <span class="browse-slots__detail-value">{{ slot.user?.username || 'Unknown' }}</span>
            </div>
            <div v-if="slot.start_time" class="browse-slots__detail">
              <span class="browse-slots__detail-label">📅 Date:</span>
              <span class="browse-slots__detail-value">{{ formatDate(slot.start_time) }}</span>
            </div>
            <div v-if="slot.start_time && slot.end_time" class="browse-slots__detail">
              <span class="browse-slots__detail-label">🕐 Time:</span>
              <span class="browse-slots__detail-value">{{ formatTime(slot.start_time) }} - {{ formatTime(slot.end_time) }}</span>
            </div>
          </div>

          <button
            @click="slot?.id && handleRequestSwap(slot.id)"
            :disabled="myEvents.length === 0 || !slot?.id"
            class="browse-slots__request-btn"
            :title="myEvents.length === 0 ? 'You need swappable events to request a swap' : !slot?.id ? 'Invalid slot' : 'Request swap with this slot'"
          >
            {{ myEvents.length === 0 ? 'No Swappable Events' : 'Request Swap' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.browse-slots__loading {
  text-align: center;
  padding: var(--space-3xl) 0;
}

.browse-slots__spinner {
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

.browse-slots__loading-text {
  color: var(--color-text-secondary);
  margin-top: var(--space-md);
}

.browse-slots__empty {
  text-align: center;
  padding: var(--space-3xl) 0;
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
}

.browse-slots__empty-icon {
  font-size: 4rem;
  margin-bottom: var(--space-md);
}

.browse-slots__empty-title {
  color: var(--color-text-primary);
  font-size: var(--text-lg);
  font-weight: var(--weight-semibold);
}

.browse-slots__empty-text {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin-top: var(--space-sm);
}

.browse-slots__count {
  margin-bottom: var(--space-md);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  font-weight: var(--weight-medium);
}

.browse-slots__grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: var(--space-lg);
}

@media (min-width: 768px) {
  .browse-slots__grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .browse-slots__grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

.browse-slots__card {
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-lg);
  transition: var(--transition-all);
  border: 1px solid var(--color-gray-200);
}

.browse-slots__card:hover {
  box-shadow: var(--shadow-xl);
  transform: translateY(-4px);
  border-color: var(--color-success);
}

.browse-slots__card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-md);
  gap: var(--space-sm);
}

.browse-slots__card-title {
  font-weight: var(--weight-semibold);
  font-size: var(--text-lg);
  color: var(--color-text-primary);
  flex: 1;
}

.browse-slots__badge {
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: var(--weight-semibold);
  background: var(--status-swappable-bg);
  color: var(--status-swappable-text);
  border: 2px solid var(--status-swappable-border);
  flex-shrink: 0;
}

.browse-slots__card-details {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  margin-bottom: var(--space-md);
}

.browse-slots__detail {
  display: flex;
  align-items: center;
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.browse-slots__detail-label {
  font-weight: var(--weight-medium);
  margin-right: var(--space-sm);
}

.browse-slots__detail-value {
  color: var(--color-text-primary);
}

.browse-slots__request-btn {
  width: 100%;
  padding: var(--space-sm) var(--space-md);
  background: var(--gradient-secondary-btn);
  color: var(--color-text-inverse);
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--weight-semibold);
  cursor: pointer;
  transition: var(--transition-all);
  box-shadow: var(--shadow-secondary);
}

.browse-slots__request-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

.browse-slots__request-btn:disabled {
  background: var(--color-gray-300);
  color: var(--color-gray-500);
  cursor: not-allowed;
  box-shadow: none;
}
</style>