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
    <div v-if="swapsStore.isLoading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-b`lue-600"></div>
      <div class="text-gray-600 mt-4">Loading available slots...</div>
    </div>

    <div v-else-if="availableSlots.length === 0" class="text-center py-12 bg-white rounded-lg shadow">
      <div class="text-6xl mb-4">🔍</div>
      <div class="text-gray-500 text-lg font-medium">No swappable slots available</div>
      <p class="text-sm text-gray-400 mt-2">Check back later or ask others to mark their events as swappable.</p>
    </div>

    <div v-else>
      <div class="mb-4 text-sm text-gray-600">
        Found {{ availableSlots.length }} available slot{{ availableSlots.length !== 1 ? 's' : '' }}
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="slot in availableSlots"
          :key="slot.id"
          class="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-6"
        >
          <div class="flex justify-between items-start mb-4">
            <h3 class="font-semibold text-lg text-gray-900 flex-1">{{ slot.title }}</h3>
            <span class="px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 ml-2">
              SWAPPABLE
            </span>
          </div>

          <div class="space-y-2 text-sm text-gray-600 mb-4">
            <div class="flex items-center">
              <span class="font-medium mr-2">👤 Owner:</span>
              <span class="text-gray-900">{{ slot.user?.username || 'Unknown' }}</span>
            </div>
            <div v-if="slot.start_time" class="flex items-center">
              <span class="font-medium mr-2">📅 Date:</span>
              <span>{{ formatDate(slot.start_time) }}</span>
            </div>
            <div v-if="slot.start_time && slot.end_time" class="flex items-center">
              <span class="font-medium mr-2">🕐 Time:</span>
              <span>{{ formatTime(slot.start_time) }} - {{ formatTime(slot.end_time) }}</span>
            </div>
          </div>

          <button
            @click="handleRequestSwap(slot.id)"
            :disabled="myEvents.length === 0"
            class="w-full bg-blue-600 text-white py-2.5 rounded-lg hover:bg-blue-700 transition-colors font-medium disabled:bg-gray-300 disabled:cursor-not-allowed"
            :title="myEvents.length === 0 ? 'You need swappable events to request a swap' : 'Request swap with this slot'"
          >
            {{ myEvents.length === 0 ? 'No Swappable Events' : 'Request Swap' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>