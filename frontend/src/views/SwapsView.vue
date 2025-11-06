<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useSwapsStore } from '../stores/swaps_store';
import { useEventsStore } from '../stores/event_store';
import BrowseSlotsTab from './swapView/BrowseSlotsTabView.vue';
import SentSwapsTab from './swapView/SentSwapsTabView.vue';
import ReceivedSwapsTab from './swapView/ReceivedSwapsTabView.vue';

const router = useRouter();
const authStore = useAuthStore();
const swapsStore = useSwapsStore();
const eventsStore = useEventsStore();

const activeTab = ref<'browse' | 'sent' | 'received'>('browse');

const handleLogout = () => {
  authStore.logout();
  router.push('/login');
};

onMounted(async () => {
  await Promise.all([
    eventsStore.fetchEvents(), // Need this for browse tab
    swapsStore.fetchSwappableSlots(),
    swapsStore.fetchSentRequests(),
    swapsStore.fetchReceivedRequests(),
  ]);
});
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow-sm sticky top-0 z-10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-900">SlotSwapper</h1>
        <div class="flex items-center gap-4">
          <router-link 
            to="/dashboard" 
            class="text-sm text-gray-600 hover:text-gray-900 font-medium"
          >
            Dashboard
          </router-link>
          <router-link 
            to="/events" 
            class="text-sm text-gray-600 hover:text-gray-900 font-medium"
          >
            Events
          </router-link>
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

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h2 class="text-3xl font-bold text-gray-900 mb-6">Swap Management</h2>

      <!-- Tabs -->
      <div class="border-b border-gray-200 mb-6">
        <nav class="flex gap-8">
          <button
            @click="activeTab = 'browse'"
            :class="{
              'py-4 px-1 border-b-2 font-medium text-sm transition-colors': true,
              'border-blue-500 text-blue-600': activeTab === 'browse',
              'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300': activeTab !== 'browse',
            }"
          >
            🔍 Browse Slots
          </button>
          <button
            @click="activeTab = 'sent'"
            :class="{
              'py-4 px-1 border-b-2 font-medium text-sm transition-colors relative': true,
              'border-blue-500 text-blue-600': activeTab === 'sent',
              'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300': activeTab !== 'sent',
            }"
          >
            📤 Sent 
            <span 
              v-if="swapsStore.sentRequests.length > 0"
              class="ml-1 px-2 py-0.5 text-xs rounded-full bg-blue-100 text-blue-800"
            >
              {{ swapsStore.sentRequests.length }}
            </span>
          </button>
          <button
            @click="activeTab = 'received'"
            :class="{
              'py-4 px-1 border-b-2 font-medium text-sm transition-colors relative': true,
              'border-blue-500 text-blue-600': activeTab === 'received',
              'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300': activeTab !== 'received',
            }"
          >
            📥 Received 
            <span 
              v-if="swapsStore.receivedRequests.length > 0"
              class="ml-1 px-2 py-0.5 text-xs rounded-full bg-green-100 text-green-800"
            >
              {{ swapsStore.receivedRequests.length }}
            </span>
          </button>
        </nav>
      </div>

      <!-- Tab Content -->
      <BrowseSlotsTab v-if="activeTab === 'browse'" />
      <SentSwapsTab v-else-if="activeTab === 'sent'" />
      <ReceivedSwapsTab v-else />
    </main>
  </div>
</template>