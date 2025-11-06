<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useEventsStore } from '../stores/event_store';
import { useToast } from '../composables/useToast';
import type { Event, EventStatus } from '../types/event_types';
import AppModal from '../components/common/AppModal.vue';
import AppButton from '../components/common/AppButton.vue';
import EventForm from '../components/events/EventForm.vue';
import EventCard from '../components/events/EventCard.vue';
import ConfirmDialog from '../components/common/ConfirmDialog.vue';
import LoadingSpinner from '../components/common/LoadingSpinner.vue';
import AppHeader from '../components/common/AppHeader.vue';

const router = useRouter();
const authStore = useAuthStore();
const eventsStore = useEventsStore();
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
  <div class="events-view">
    <!-- Header -->
    <AppHeader :username="authStore.user?.username" @logout="handleLogout" />

    <!-- Main Content -->
    <main class="events-view__content">
      <!-- Header with Create Button -->
      <div class="events-view__header">
        <h2 class="events-view__title">My Events</h2>
        <AppButton variant="primary" @click="openCreateModal">
          + Create Event
        </AppButton>
      </div>
      
      <!-- Loading State -->
      <div v-if="eventsStore.isLoading" class="events-view__loading">
        <LoadingSpinner size="lg" message="Loading events..." />
      </div>

      <!-- Empty State -->
      <div v-else-if="eventsStore.events.length === 0" class="events-view__empty">
        <div class="events-view__empty-text">No events yet. Create your first event!</div>
        <AppButton variant="primary" @click="openCreateModal">
          Create Event
        </AppButton>
      </div>

      <!-- Events Grid -->
      <div v-else class="events-view__grid">
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

<style scoped>
.events-view {
  min-height: 100vh;
  background: var(--gradient-bg-subtle);
}

.events-view__content {
  max-width: 80rem;
  margin: 0 auto;
  padding: var(--space-xl) var(--space-lg);
}

.events-view__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-xl);
}

.events-view__title {
  font-size: var(--text-3xl);
  font-weight: var(--weight-bold);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.events-view__loading {
  text-align: center;
  padding: var(--space-3xl) 0;
}

.events-view__empty {
  text-align: center;
  padding: var(--space-3xl) 0;
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
}

.events-view__empty-text {
  color: var(--color-text-secondary);
  margin-bottom: var(--space-md);
  font-size: var(--text-lg);
}

.events-view__grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: var(--space-lg);
}

@media (min-width: 768px) {
  .events-view__grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .events-view__grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>