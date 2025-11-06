import apiClient from './axios';
import type { SwapRequest, InitiateSwapRequest, SwappableSlot } from '@/types/swap_types';

// Browse swappable slots from other users
export const getSwappableSlots = () => {
  return apiClient.get<SwappableSlot[]>('/swappable-slots');
};

// Get your sent swap requests
export const getSentSwapRequests = () => {
  return apiClient.get<SwapRequest[]>('/swap-requests?filter=sent');
};

// Get your received swap requests
export const getReceivedSwapRequests = () => {
  return apiClient.get<SwapRequest[]>('/swap-requests?filter=received');
};

// Create a new swap request
export const initiateSwap = (data: InitiateSwapRequest) => {
  return apiClient.post<SwapRequest>('/swap-request', data);
};

// Accept a swap request
export const acceptSwap = (id: number) => {
  return apiClient.post<SwapRequest>(`/swap-request/${id}/accept`);
};

// Reject a swap request
export const rejectSwap = (id: number) => {
  return apiClient.post<SwapRequest>(`/swap-request/${id}/reject`);
};

// Cancel/delete a swap request
export const cancelSwap = (id: number) => {
  return apiClient.delete(`/swap-request/${id}`);
};