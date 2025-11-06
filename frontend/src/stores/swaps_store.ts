import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { SwapRequest, InitiateSwapRequest, SwappableSlot } from '@/types/swap_types';
import * as swapsApi from '@/api/swap_api';

export const useSwapsStore = defineStore('swaps', () => {
  // State
  const swappableSlots = ref<SwappableSlot[]>([]);
  const sentRequests = ref<SwapRequest[]>([]);
  const receivedRequests = ref<SwapRequest[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // Actions
  async function fetchSwappableSlots() {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await swapsApi.getSwappableSlots();      
      // Backend wraps response in { message, data } structure
      const data = (response.data as any).data;
      
      // Handle null/undefined data from backend
      if (data === null || data === undefined) {
        swappableSlots.value = [];
      } else if (Array.isArray(data)) {
        swappableSlots.value = data;
      } else {
        swappableSlots.value = [];
      }
      console.log('=== FETCH SWAPPABLE SLOTS DEBUG ===');
      console.log('Full Axios response:', response);
      console.log('response.data:', response.data);
      console.log('response.data.data:', (response.data as any).data);
      
      // Backend wraps response in { message, data } structure
      swappableSlots.value = (response.data as any).data || response.data || [];
      console.log('Stored swappableSlots:', swappableSlots.value);
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to fetch swappable slots';
      swappableSlots.value = [];
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchSentRequests() {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await swapsApi.getSentSwapRequests();
      const data = (response.data as any).data;
      sentRequests.value = Array.isArray(data) ? data : [];
      sentRequests.value = response.data || [];
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to fetch sent requests';
      sentRequests.value = [];
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchReceivedRequests() {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await swapsApi.getReceivedSwapRequests();
      const data = (response.data as any).data;
      receivedRequests.value = Array.isArray(data) ? data : [];
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to fetch received requests';
      receivedRequests.value = [];
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function initiateSwap(data: InitiateSwapRequest) {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await swapsApi.initiateSwap(data);
      sentRequests.value.push(response.data);
      return response.data;
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to initiate swap';
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function acceptSwap(id: number) {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await swapsApi.acceptSwap(id);
      const index = receivedRequests.value.findIndex(r => r.id === id);
      if (index !== -1) {
        receivedRequests.value[index] = response.data;
      }
      return response.data;
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to accept swap';
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function rejectSwap(id: number) {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await swapsApi.rejectSwap(id);
      const index = receivedRequests.value.findIndex(r => r.id === id);
      if (index !== -1) {
        receivedRequests.value[index] = response.data;
      }
      return response.data;
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to reject swap';
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function cancelSwap(id: number) {
    isLoading.value = true;
    error.value = null;
    
    try {
      await swapsApi.cancelSwap(id);
      sentRequests.value = sentRequests.value.filter(r => r.id !== id);
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to cancel swap';
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  return {
    // State
    swappableSlots,
    sentRequests,
    receivedRequests,
    isLoading,
    error,
    // Actions
    fetchSwappableSlots,
    fetchSentRequests,
    fetchReceivedRequests,
    initiateSwap,
    acceptSwap,
    rejectSwap,
    cancelSwap,
  };
});