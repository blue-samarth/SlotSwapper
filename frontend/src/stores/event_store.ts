import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Event, CreateEventRequest, UpdateEventStatusRequest } from '@/types/event_types';
import * as eventsApi from '@/api/events_api';

export const useEventsStore = defineStore('events', () => {
  // State
  const events = ref<Event[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // Actions
  async function fetchEvents() {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await eventsApi.getEvents();
      events.value = response.data || [];
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to fetch events';
      events.value = [];
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function createEvent(data: CreateEventRequest) {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await eventsApi.createEvent(data);
      events.value.push(response.data);
      return response.data;
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to create event';
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function updateEventStatus(id: number, data: UpdateEventStatusRequest) {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await eventsApi.updateEventStatus(id, data);
      const index = events.value.findIndex(e => e.id === id);
      if (index !== -1) {
        events.value[index] = response.data;
      }
      return response.data;
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to update event';
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function deleteEvent(id: number) {
    isLoading.value = true;
    error.value = null;
    
    try {
      await eventsApi.deleteEvent(id);
      events.value = events.value.filter(e => e.id !== id);
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Failed to delete event';
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  return {
    // State
    events,
    isLoading,
    error,
    // Actions
    fetchEvents,
    createEvent,
    updateEventStatus,
    deleteEvent,
  };
});