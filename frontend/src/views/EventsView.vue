<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useEventsStore } from '../stores/event_store';
import { useSwapsStore } from '../stores/swaps_store';
import { useToast } from '../composables/useToast';
import type { Event, EventStatus } from '../types/event_types';
import AppModal from '../components/common/AppModal.vue';
import AppButton from '../components/common/AppButton.vue';
import EventForm from '../components/events/EventForm.vue';
import EventCard from '../components/events/EventCard.vue';
import ConfirmDialog from '../components/common/ConfirmDialog.vue';
import LoadingSpinner from '../components/common/LoadingSpinner.vue';

const router = useRouter();
const authStore = useAuthStore();
const eventsStore = useEventsStore();
const swapsStore = useSwapsStore();
const toast = useToast();

// Modal state
const showCreateModal = ref(false);
const showDeleteModal = ref(false);
const selectedEvent = ref<Event | null>(null);
const isSubmitting = ref(false);

const handleLogout = () => {
  authStore.logout();
  router.push('/login');
};

const openCreateModal = () => {
  selectedEvent.value = null;
  showCreateModal.value = true;
};

const openDeleteModal = (event: Event) => {
  selectedEvent.value = event;
  showDeleteModal.value = true;
};

const handleCreateSubmit = async (data: {
  title: string;
  date: string;
  time: string;
  duration: number;
  status: EventStatus;
}) => {
  isSubmitting.value = true;
  try {
    // Combine date and time into ISO string
    const startTime = new Date(`${data.date}T${data.time}`).toISOString();
    
    // Calculate end time
    const endTime = new Date(new Date(startTime).getTime() + data.duration * 60000).toISOString();
    
    await eventsStore.createEvent({
      title: data.title,
      start_time: startTime,
      end_time: endTime,
    });
    
    showCreateModal.value = false;
    toast.success('Event created successfully!');
  } catch (error) {
    console.error('Failed to create event:', error);
    toast.error('Failed to create event. Please try again.');
  } finally {
    isSubmitting.value = false;
  }
};

const handleDelete = async () => {
  if (!selectedEvent.value) return;
  
  // Check if event has pending swaps
  if (selectedEvent.value.status === 'SWAP_PENDING') {
    toast.warning('Cannot delete event with pending swap. Please cancel the swap request first.');
    showDeleteModal.value = false;
    return;
  }
  isSubmitting.value = true;
  try {
    await eventsStore.deleteEvent(selectedEvent.value.id);
    showDeleteModal.value = false;
    toast.success('Event deleted successfully!');
  } catch (error: any) {
    console.error('Failed to delete event:', error);
    const errorMessage = error.response?.data?.error || 'Failed to delete event. Please try again.';
    toast.error(errorMessage);
    showDeleteModal.value = false;
  } finally {
    isSubmitting.value = false;
  }
};

const handleToggleStatus = async (event: Event) => {
  // Only toggle between BUSY and SWAPPABLE
  if (event.status === 'SWAP_PENDING') {
    toast.warning('Cannot change status of events with pending swaps');
    return;
  }
  
  try {
    const newStatus = event.status === 'BUSY' ? 'SWAPPABLE' : 'BUSY';
    await eventsStore.updateEventStatus(event.id, { status: newStatus });
    toast.success(`Event status changed to ${newStatus}`);
  } catch (error) {
    console.error('Failed to toggle status:', error);
    toast.error('Failed to toggle status. Please try again.');
  }
};

onMounted(async () => {
  await eventsStore.fetchEvents();
});
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-900">SlotSwapper</h1>
        <div class="flex items-center gap-4">
          <router-link to="/dashboard" class="text-sm text-gray-600 hover:text-gray-900">Dashboard</router-link>
          <span class="text-sm text-gray-600">{{ authStore.user?.username }}</span>
          <button @click="handleLogout" class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700">
            Logout
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header with Create Button -->
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold text-gray-900">My Events</h2>
        <AppButton variant="primary" @click="openCreateModal">
          + Create Event
        </AppButton>
      </div>
      
      <!-- Loading State -->
      <div v-if="eventsStore.isLoading" class="text-center py-12">
        <LoadingSpinner size="lg" message="Loading events..." />
      </div>

      <!-- Empty State -->
      <div v-else-if="eventsStore.events.length === 0" class="text-center py-12 bg-white rounded-lg shadow">
        <div class="text-gray-500 mb-4">No events yet. Create your first event!</div>
        <AppButton variant="primary" @click="openCreateModal">
          Create Event
        </AppButton>
      </div>

      <!-- Events Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <EventCard
          v-for="event in eventsStore.events"
          :key="event.id"
          :event="event"
          @toggle-status="handleToggleStatus"
          @delete="openDeleteModal"
        />
      </div>
    </main>

    <!-- Create Event Modal -->
    <AppModal
      v-model="showCreateModal"
      title="Create New Event"
    >
      <EventForm
        :is-submitting="isSubmitting"
        @submit="handleCreateSubmit"
        @cancel="showCreateModal = false"
      />
    </AppModal>

    <!-- Delete Confirmation Dialog -->
    <ConfirmDialog
      v-model="showDeleteModal"
      title="Delete Event"
      :message="`Are you sure you want to delete '${selectedEvent?.title}'?`"
      sub-message="This action cannot be undone."
      type="danger"
      confirm-text="Delete Event"
      :is-loading="isSubmitting"
      @confirm="handleDelete"
      @cancel="showDeleteModal = false"
    />
  </div>
</template>