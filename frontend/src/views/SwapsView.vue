<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useSwapsStore } from '../stores/swaps_store';
import { useEventsStore } from '../stores/event_store';
import BrowseSlotsTab from './swapView/BrowseSlotsTabView.vue';
import SentSwapsTab from './swapView/SentSwapsTabView.vue';
import ReceivedSwapsTab from './swapView/ReceivedSwapsTabView.vue';
import AppHeader from '../components/common/AppHeader.vue';

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
  <div class="swaps-view">
    <!-- Header -->
    <AppHeader :username="authStore.user?.username" @logout="handleLogout" />

    <!-- Main Content -->
    <main class="swaps-view__content">
      <h2 class="swaps-view__title">Swap Management</h2>

      <!-- Tabs -->
      <div class="swaps-view__tabs">
        <nav class="swaps-view__tabs-nav">
          <button
            @click="activeTab = 'browse'"
            :class="['swaps-view__tab', { 'swaps-view__tab--active': activeTab === 'browse' }]"
          >
            🔍 Browse Slots
          </button>
          <button
            @click="activeTab = 'sent'"
            :class="['swaps-view__tab', { 'swaps-view__tab--active': activeTab === 'sent' }]"
          >
            📤 Sent 
            <span 
              v-if="swapsStore.sentRequests.length > 0"
              class="swaps-view__badge swaps-view__badge--secondary"
            >
              {{ swapsStore.sentRequests.length }}
            </span>
          </button>
          <button
            @click="activeTab = 'received'"
            :class="['swaps-view__tab', { 'swaps-view__tab--active': activeTab === 'received' }]"
          >
            📥 Received 
            <span 
              v-if="swapsStore.receivedRequests.length > 0"
              class="swaps-view__badge swaps-view__badge--success"
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

<style scoped>
.swaps-view {
  min-height: 100vh;
  background: var(--gradient-bg-subtle);
}

.swaps-view__content {
  max-width: 80rem;
  margin: 0 auto;
  padding: var(--space-xl) var(--space-lg);
}

.swaps-view__title {
  font-size: var(--text-3xl);
  font-weight: var(--weight-bold);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: var(--space-xl);
}

.swaps-view__tabs {
  border-bottom: 2px solid var(--color-gray-200);
  margin-bottom: var(--space-xl);
}

.swaps-view__tabs-nav {
  display: flex;
  gap: var(--space-xl);
}

.swaps-view__tab {
  padding: var(--space-md) var(--space-sm);
  border-bottom: 3px solid transparent;
  font-weight: var(--weight-medium);
  font-size: var(--text-base);
  transition: var(--transition-all);
  color: var(--color-text-secondary);
  background: none;
  border-top: none;
  border-left: none;
  border-right: none;
  cursor: pointer;
  position: relative;
  display: flex;
  align-items: center;
  gap: var(--space-xs);
}

.swaps-view__tab:hover {
  color: var(--color-primary);
  border-bottom-color: var(--color-gray-300);
  transform: translateY(-2px);
}

.swaps-view__tab--active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
  font-weight: var(--weight-semibold);
}

.swaps-view__badge {
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: var(--weight-semibold);
  margin-left: var(--space-xs);
}

.swaps-view__badge--secondary {
  background: linear-gradient(135deg, rgba(78, 205, 196, 0.2), rgba(78, 205, 196, 0.1));
  color: var(--color-secondary);
  border: 1px solid var(--color-secondary);
}

.swaps-view__badge--success {
  background: linear-gradient(135deg, rgba(149, 225, 211, 0.2), rgba(149, 225, 211, 0.1));
  color: var(--color-success);
  border: 1px solid var(--color-success);
}
</style>