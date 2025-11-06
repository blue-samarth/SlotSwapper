<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useEventsStore } from '../stores/event_store';
import { useSwapsStore } from '../stores/swaps_store';
import AppHeader from '../components/common/AppHeader.vue';

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
  <div class="dashboard-view">
    <!-- Header -->
    <AppHeader :username="authStore.user?.username" @logout="handleLogout" />

    <!-- Loading State -->
    <div v-if="isLoading" class="dashboard-view__loading">
      <div class="dashboard-view__spinner"></div>
      <p class="dashboard-view__loading-text">Loading dashboard...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="loadError" class="dashboard-view__error-container">
      <div class="dashboard-view__error">
        <p class="dashboard-view__error-text">{{ loadError }}</p>
        <button @click="router.go(0)" class="dashboard-view__retry-btn">
          Retry
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <main v-else class="dashboard-view__content">
      <div class="dashboard-view__welcome">
        <h2 class="dashboard-view__title">Dashboard</h2>
        <p class="dashboard-view__subtitle">Welcome back, {{ authStore.user?.username }}!</p>
      </div>

      <!-- Quick Stats -->
      <div class="dashboard-view__stats">
        <div class="dashboard-view__stat-card dashboard-view__stat-card--primary">
          <div class="dashboard-view__stat-label">Total Events</div>
          <div class="dashboard-view__stat-value">
            {{ eventsStore.events?.length || 0 }}
          </div>
        </div>
        <div class="dashboard-view__stat-card dashboard-view__stat-card--secondary">
          <div class="dashboard-view__stat-label">Sent Swap Requests</div>
          <div class="dashboard-view__stat-value">
            {{ swapsStore.sentRequests?.length || 0 }}
          </div>
        </div>
        <div class="dashboard-view__stat-card dashboard-view__stat-card--success">
          <div class="dashboard-view__stat-label">Received Swap Requests</div>
          <div class="dashboard-view__stat-value">
            {{ swapsStore.receivedRequests?.length || 0 }}
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="dashboard-view__actions-card">
        <h3 class="dashboard-view__actions-title">Quick Actions</h3>
        <div class="dashboard-view__actions">
          <router-link to="/events" class="dashboard-view__action">
            <div class="dashboard-view__action-icon">📅</div>
            <div class="dashboard-view__action-text">Manage Events</div>
          </router-link>
          <router-link to="/swaps" class="dashboard-view__action">
            <div class="dashboard-view__action-icon">�</div>
            <div class="dashboard-view__action-text">View Swaps</div>
          </router-link>
          <router-link to="/profile" class="dashboard-view__action">
            <div class="dashboard-view__action-icon">👤</div>
            <div class="dashboard-view__action-text">My Profile</div>
          </router-link>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.dashboard-view {
  min-height: 100vh;
  background: var(--gradient-bg-subtle);
}

.dashboard-view__loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
}

.dashboard-view__spinner {
  display: inline-block;
  width: 3rem;
  height: 3rem;
  border: 4px solid var(--color-gray-200);
  border-radius: var(--radius-full);
  border-top-color: var(--color-secondary);
  animation: spin 0.8s linear infinite;
}

.dashboard-view__loading-text {
  margin-top: var(--space-md);
  color: var(--color-text-secondary);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.dashboard-view__error-container {
  max-width: 80rem;
  margin: 0 auto;
  padding: var(--space-xl) var(--space-lg);
}

.dashboard-view__error {
  background: linear-gradient(135deg, rgba(255, 107, 107, 0.1), rgba(255, 107, 107, 0.05));
  border: 2px solid var(--color-error);
  border-radius: var(--radius-lg);
  padding: var(--space-xl);
  text-align: center;
}

.dashboard-view__error-text {
  color: var(--color-error);
  font-weight: var(--weight-semibold);
}

.dashboard-view__retry-btn {
  margin-top: var(--space-md);
  padding: var(--space-sm) var(--space-md);
  background: var(--color-error);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--weight-medium);
  cursor: pointer;
  transition: var(--transition-all);
}

.dashboard-view__retry-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.dashboard-view__content {
  max-width: 80rem;
  margin: 0 auto;
  padding: var(--space-xl) var(--space-lg);
}

.dashboard-view__welcome {
  margin-bottom: var(--space-xl);
}

.dashboard-view__title {
  font-size: var(--text-3xl);
  font-weight: var(--weight-bold);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.dashboard-view__subtitle {
  color: var(--color-text-secondary);
  margin-top: var(--space-sm);
  font-size: var(--text-lg);
}

.dashboard-view__stats {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: var(--space-lg);
  margin-bottom: var(--space-xl);
}

@media (min-width: 768px) {
  .dashboard-view__stats {
    grid-template-columns: repeat(3, 1fr);
  }
}

.dashboard-view__stat-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-lg);
  transition: var(--transition-all);
}

.dashboard-view__stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.dashboard-view__stat-card--primary {
  border-left: 4px solid var(--color-primary);
}

.dashboard-view__stat-card--secondary {
  border-left: 4px solid var(--color-secondary);
}

.dashboard-view__stat-card--success {
  border-left: 4px solid var(--color-success);
}

.dashboard-view__stat-label {
  font-size: var(--text-sm);
  font-weight: var(--weight-medium);
  color: var(--color-text-secondary);
}

.dashboard-view__stat-value {
  margin-top: var(--space-sm);
  font-size: var(--text-3xl);
  font-weight: var(--weight-bold);
  color: var(--color-text-primary);
}

.dashboard-view__actions-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-lg);
}

.dashboard-view__actions-title {
  font-size: var(--text-lg);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-md);
}

.dashboard-view__actions {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: var(--space-md);
}

@media (min-width: 768px) {
  .dashboard-view__actions {
    grid-template-columns: repeat(3, 1fr);
  }
}

.dashboard-view__action {
  padding: var(--space-md);
  border: 2px solid var(--color-gray-200);
  border-radius: var(--radius-lg);
  text-decoration: none;
  text-align: center;
  transition: var(--transition-all);
}

.dashboard-view__action:hover {
  border-color: var(--color-secondary);
  background: linear-gradient(135deg, rgba(78, 205, 196, 0.1), transparent);
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.dashboard-view__action-icon {
  font-size: var(--text-2xl);
  margin-bottom: var(--space-sm);
}

.dashboard-view__action-text {
  font-weight: var(--weight-medium);
  color: var(--color-text-primary);
}
</style>