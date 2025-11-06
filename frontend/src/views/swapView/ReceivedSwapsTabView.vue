<script setup lang="ts">
import { computed } from 'vue';
import { useSwapsStore } from '../..//stores/swaps_store';
import { SwapStatus } from '../..//types/swap_types';
import { formatDate, formatTime } from '../../utils/date';

const swapsStore = useSwapsStore();

const receivedRequests = computed(() => swapsStore.receivedRequests || []);

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
    <div v-if="swapsStore.isLoading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <div class="text-gray-600 mt-4">Loading received requests...</div>
    </div>

    <div v-else-if="receivedRequests.length === 0" class="text-center py-12 bg-white rounded-lg shadow">
      <div class="text-6xl mb-4">📥</div>
      <div class="text-gray-500 text-lg font-medium">No received swap requests yet</div>
      <p class="text-sm text-gray-400 mt-2">When someone requests to swap with you, they'll appear here.</p>
    </div>

    <div v-else>
      <div class="mb-4 text-sm text-gray-600">
        {{ receivedRequests.length }} received request{{ receivedRequests.length !== 1 ? 's' : '' }}
      </div>

      <div class="space-y-4">
        <div
          v-for="request in receivedRequests"
          :key="request.id"
          class="bg-white rounded-lg shadow hover:shadow-md transition-shadow p-6"
        >
          <div class="flex justify-between items-start mb-4">
            <h3 class="font-semibold text-lg text-gray-900">Swap Request Received</h3>
            <span :class="['px-3 py-1 rounded-full text-xs font-semibold border', getStatusColor(request.status)]">
              {{ getStatusIcon(request.status) }} {{ request.status }}
            </span>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-4">
            <!-- Their Event -->
            <div class="border-l-4 border-purple-500 pl-4 bg-purple-50 rounded-r-lg p-4">
              <div class="text-xs font-semibold text-purple-600 uppercase mb-2">They Offer</div>
              <h4 class="font-semibold text-gray-900 text-lg">{{ request.requester_event.title }}</h4>
              <div v-if="request.requester_event.start_time" class="text-sm text-gray-600 mt-2 space-y-1">
                <div>📅 {{ formatDate(request.requester_event.start_time) }}</div>
                <div>🕐 {{ formatTime(request.requester_event.start_time) }}</div>
              </div>
            </div>

            <!-- Your Event -->
            <div class="border-l-4 border-blue-500 pl-4 bg-blue-50 rounded-r-lg p-4">
              <div class="text-xs font-semibold text-blue-600 uppercase mb-2">Your Event</div>
              <h4 class="font-semibold text-gray-900 text-lg">{{ request.receiver_event.title }}</h4>
              <div v-if="request.receiver_event.start_time" class="text-sm text-gray-600 mt-2 space-y-1">
                <div>📅 {{ formatDate(request.receiver_event.start_time) }}</div>
                <div>🕐 {{ formatTime(request.receiver_event.start_time) }}</div>
              </div>
            </div>
          </div>

          <div class="flex justify-between items-center pt-4 border-t">
            <div v-if="request.created_at" class="text-sm text-gray-500">
              Received on {{ formatDate(request.created_at) }}
            </div>
            <div v-if="request.status === SwapStatus.PENDING" class="flex gap-2">
              <button
                @click="handleRejectSwap(request.id)"
                class="px-4 py-2 text-sm font-medium text-red-600 hover:bg-red-50 border border-red-200 rounded-lg transition-colors"
              >
                Reject
              </button>
              <button
                @click="handleAcceptSwap(request.id)"
                class="px-4 py-2 text-sm font-medium text-white bg-green-600 hover:bg-green-700 rounded-lg transition-colors"
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