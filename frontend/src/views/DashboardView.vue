<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useEventsStore } from '../stores/event_store';
import { useSwapsStore } from '../stores/swaps_store';

const router = useRouter();
const authStore = useAuthStore();
const eventsStore = useEventsStore();
const swapsStore = useSwapsStore();

const isLoading = ref(true);
const loadError = ref<string | null>(null);

const handleLogout = () => {
  authStore.logout();
  router.push('/login');
};

onMounted(async () => {
  console.log('DashboardView mounted');
  console.log('User:', authStore.user);
  
  // Load initial data
  try {
    await Promise.all([
      eventsStore.fetchEvents(),
      swapsStore.fetchSentRequests(),
      swapsStore.fetchReceivedRequests(),
    ]);
    console.log('Dashboard data loaded successfully');
  } catch (error: any) {
    console.error('Failed to load dashboard data:', error);
    loadError.value = error.message || 'Failed to load dashboard data';
  } finally {
    isLoading.value = false;
  }
});
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Simple Header -->
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-900">SlotSwapper</h1>
        <div class="flex items-center gap-4">
          <span class="text-sm text-gray-600">{{ authStore.user?.username }}</span>
          <button
            @click="handleLogout"
            class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
          >
            Logout
          </button>
        </div>
      </div>
    </header>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex items-center justify-center min-h-[60vh]">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <p class="mt-4 text-gray-600">Loading dashboard...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="loadError" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
        <p class="text-red-800 font-medium">{{ loadError }}</p>
        <button 
          @click="router.go(0)" 
          class="mt-4 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        >
          Retry
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <main v-else class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="mb-8">
        <h2 class="text-3xl font-bold text-gray-900">Dashboard</h2>
        <p class="text-gray-600 mt-2">Welcome back, {{ authStore.user?.username }}!</p>
      </div>

      <!-- Quick Stats -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-sm font-medium text-gray-500">Total Events</div>
          <div class="mt-2 text-3xl font-bold text-gray-900">
            {{ eventsStore.events?.length || 0 }}
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-sm font-medium text-gray-500">Sent Swap Requests</div>
          <div class="mt-2 text-3xl font-bold text-blue-600">
            {{ swapsStore.sentRequests?.length || 0 }}
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-sm font-medium text-gray-500">Received Swap Requests</div>
          <div class="mt-2 text-3xl font-bold text-green-600">
            {{ swapsStore.receivedRequests?.length || 0 }}
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="bg-white rounded-lg shadow p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <router-link
            to="/events"
            class="p-4 border-2 border-gray-200 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors text-center"
          >
            <div class="text-2xl mb-2">📅</div>
            <div class="font-medium text-gray-900">Manage Events</div>
          </router-link>
          <router-link
            to="/swaps"
            class="p-4 border-2 border-gray-200 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors text-center"
          >
            <div class="text-2xl mb-2">🔄</div>
            <div class="font-medium text-gray-900">View Swaps</div>
          </router-link>
          <router-link
            to="/profile"
            class="p-4 border-2 border-gray-200 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors text-center"
          >
            <div class="text-2xl mb-2">👤</div>
            <div class="font-medium text-gray-900">My Profile</div>
          </router-link>
        </div>
      </div>
    </main>
  </div>
</template>