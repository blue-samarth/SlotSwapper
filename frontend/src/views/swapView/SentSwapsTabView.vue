<script setup lang="ts">
import { computed } from 'vue';
import { useSwapsStore } from '../..//stores/swaps_store';
import { SwapStatus } from '../..//types/swap_types';
import { formatDate, formatTime } from '../../utils/date';

const swapsStore = useSwapsStore();

const sentRequests = computed(() => swapsStore.sentRequests || []);

const getStatusColor = (status: string) => {
  switch (status) {
    case SwapStatus.PENDING:
      return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    case SwapStatus.ACCEPTED:
      return 'bg-green-100 text-green-800 border-green-200';
    case SwapStatus.REJECTED:
      return 'bg-red-100 text-red-800 border-red-200';
    case SwapStatus.CANCELLED:
      return 'bg-gray-100 text-gray-800 border-gray-200';
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200';
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
    <div v-if="swapsStore.isLoading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <div class="text-gray-600 mt-4">Loading sent requests...</div>
    </div>

    <div v-else-if="sentRequests.length === 0" class="text-center py-12 bg-white rounded-lg shadow">
      <div class="text-6xl mb-4">📤</div>
      <div class="text-gray-500 text-lg font-medium">No sent swap requests yet</div>
      <p class="text-sm text-gray-400 mt-2">Browse available slots to initiate your first swap!</p>
    </div>

    <div v-else>
      <div class="mb-4 text-sm text-gray-600">
        {{ sentRequests.length }} sent request{{ sentRequests.length !== 1 ? 's' : '' }}
      </div>

      <div class="space-y-4">
        <div
          v-for="(request, index) in sentRequests"
          :key="request?.id || `request-${index}`"
          class="bg-white rounded-lg shadow hover:shadow-md transition-shadow p-6"
        >
          <div class="flex justify-between items-start mb-4">
            <h3 class="font-semibold text-lg text-gray-900">Swap Request</h3>
            <span v-if="request?.status" :class="['px-3 py-1 rounded-full text-xs font-semibold border', getStatusColor(request.status)]">
              {{ getStatusIcon(request.status) }} {{ request.status }}
            </span>
          </div>

          <div v-if="request?.requester_event && request?.receiver_event" class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-4">
            <!-- Your Event -->
            <div class="border-l-4 border-blue-500 pl-4 bg-blue-50 rounded-r-lg p-4">
              <div class="text-xs font-semibold text-blue-600 uppercase mb-2">Your Event</div>
              <h4 class="font-semibold text-gray-900 text-lg">{{ request.requester_event.title || 'Untitled Event' }}</h4>
              <div v-if="request.requester_event.start_time" class="text-sm text-gray-600 mt-2 space-y-1">
                <div>📅 {{ formatDate(request.requester_event.start_time) }}</div>
                <div>🕐 {{ formatTime(request.requester_event.start_time) }}</div>
              </div>
            </div>

            <!-- Their Event -->
            <div class="border-l-4 border-green-500 pl-4 bg-green-50 rounded-r-lg p-4">
              <div class="text-xs font-semibold text-green-600 uppercase mb-2">Their Event</div>
              <h4 class="font-semibold text-gray-900 text-lg">{{ request.receiver_event.title || 'Untitled Event' }}</h4>
              <div v-if="request.receiver_event.start_time" class="text-sm text-gray-600 mt-2 space-y-1">
                <div>📅 {{ formatDate(request.receiver_event.start_time) }}</div>
                <div>🕐 {{ formatTime(request.receiver_event.start_time) }}</div>
              </div>
            </div>
          </div>

          <div class="flex justify-between items-center pt-4 border-t">
            <div v-if="request.created_at" class="text-sm text-gray-500">
              Sent on {{ formatDate(request.created_at) }}
            </div>
            <button
              v-if="request.status === SwapStatus.PENDING && request?.id"
              @click="handleCancelSwap(request.id)"
              class="px-4 py-2 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors"
            >
              Cancel Request
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>