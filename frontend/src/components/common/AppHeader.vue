<template>
  <header class="app-header">
    <div class="app-header__container">
      <div class="app-header__left">
        <router-link to="/dashboard" class="app-header__logo">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="app-header__logo-icon">
            <path fill-rule="evenodd" d="M15.312 11.424a5.5 5.5 0 01-9.201 2.466l-.312-.311h2.433a.75.75 0 000-1.5H3.989a.75.75 0 00-.75.75v4.242a.75.75 0 001.5 0v-2.43l.31.31a7 7 0 0011.712-3.138.75.75 0 00-1.449-.39zm1.23-3.723a.75.75 0 00.219-.53V2.929a.75.75 0 00-1.5 0V5.36l-.31-.31A7 7 0 003.239 8.188a.75.75 0 101.448.389A5.5 5.5 0 0113.89 6.11l.311.31h-2.432a.75.75 0 000 1.5h4.243a.75.75 0 00.53-.219z" clip-rule="evenodd" />
          </svg>
          <span class="app-header__title">SlotSwapper</span>
        </router-link>
      </div>

      <nav class="app-header__nav">
        <router-link 
          to="/dashboard" 
          class="app-header__link"
          active-class="app-header__link--active"
        >
          Dashboard
        </router-link>
        <router-link 
          to="/events" 
          class="app-header__link"
          active-class="app-header__link--active"
        >
          Events
        </router-link>
        <router-link 
          to="/swaps" 
          class="app-header__link"
          active-class="app-header__link--active"
        >
          Swaps
        </router-link>
        <router-link 
          to="/profile" 
          class="app-header__link"
          active-class="app-header__link--active"
        >
          Profile
        </router-link>
      </nav>

      <div class="app-header__right">
        <span class="app-header__username">{{ username }}</span>
        <AppButton
          size="sm"
          variant="danger"
          @click="emit('logout')"
        >
          Logout
        </AppButton>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import AppButton from './AppButton.vue';

interface Props {
  username?: string;
}

withDefaults(defineProps<Props>(), {
  username: 'User',
});

const emit = defineEmits<{
  logout: [];
}>();
</script>

<style scoped>
.app-header {
  background: white;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 40;
}

.app-header__container {
  max-width: 80rem;
  margin: 0 auto;
  padding: 1rem 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
}

.app-header__left {
  display: flex;
  align-items: center;
}

.app-header__logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  text-decoration: none;
  transition: transform 0.2s;
}

.app-header__logo:hover {
  transform: scale(1.05);
}

.app-header__logo-icon {
  width: 2rem;
  height: 2rem;
  color: #3b82f6;
}

.app-header__title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #111827;
}

.app-header__nav {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  flex: 1;
  justify-content: center;
}

.app-header__link {
  font-size: 0.9375rem;
  font-weight: 500;
  color: #6b7280;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  transition: all 0.2s;
}

.app-header__link:hover {
  color: #111827;
  background: #f3f4f6;
}

.app-header__link--active {
  color: #3b82f6;
  background: #eff6ff;
}

.app-header__right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.app-header__username {
  font-size: 0.875rem;
  font-weight: 500;
  color: #6b7280;
}

/* Responsive */
@media (max-width: 768px) {
  .app-header__container {
    flex-wrap: wrap;
    gap: 1rem;
  }

  .app-header__nav {
    order: 3;
    flex-basis: 100%;
    justify-content: space-around;
    border-top: 1px solid #e5e7eb;
    padding-top: 1rem;
    margin-top: 0.5rem;
  }

  .app-header__link {
    padding: 0.375rem 0.75rem;
    font-size: 0.875rem;
  }
}
</style>
