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
  background: var(--color-bg-card);
  box-shadow: var(--shadow-md);
  position: sticky;
  top: 0;
  z-index: 40;
  border-bottom: 2px solid transparent;
  background-image: linear-gradient(var(--color-bg-card), var(--color-bg-card)), 
                    var(--gradient-primary);
  background-origin: border-box;
  background-clip: padding-box, border-box;
}

.app-header__container {
  max-width: 80rem;
  margin: 0 auto;
  padding: var(--space-md) var(--space-lg);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-xl);
}

.app-header__left {
  display: flex;
  align-items: center;
}

.app-header__logo {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  text-decoration: none;
  transition: var(--transition-all);
}

.app-header__logo:hover {
  transform: scale(1.05) rotate(-2deg);
}

.app-header__logo-icon {
  width: 2rem;
  height: 2rem;
  color: var(--color-secondary);
  filter: drop-shadow(var(--shadow-secondary));
}

.app-header__title {
  font-size: var(--text-2xl);
  font-weight: var(--weight-bold);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.app-header__nav {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  flex: 1;
  justify-content: center;
}

.app-header__link {
  font-size: var(--text-base);
  font-weight: var(--weight-medium);
  color: var(--color-text-secondary);
  text-decoration: none;
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-lg);
  transition: var(--transition-all);
  position: relative;
}

.app-header__link::before {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 0;
  height: 2px;
  background: var(--gradient-primary);
  transform: translateX(-50%);
  transition: var(--transition-all);
}

.app-header__link:hover {
  color: var(--color-primary);
  background: var(--color-gray-100);
  transform: translateY(-2px);
}

.app-header__link:hover::before {
  width: 80%;
}

.app-header__link--active {
  color: var(--color-primary);
  background: linear-gradient(135deg, rgba(255, 107, 107, 0.1), rgba(78, 205, 196, 0.1));
  box-shadow: var(--shadow-sm);
}

.app-header__link--active::before {
  width: 100%;
}

.app-header__right {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.app-header__username {
  font-size: var(--text-sm);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  padding: var(--space-xs) var(--space-md);
  background: var(--color-gray-100);
  border-radius: var(--radius-full);
  border: 1px solid var(--color-gray-200);
}

/* Responsive */
@media (max-width: 768px) {
  .app-header__container {
    flex-wrap: wrap;
    gap: var(--space-md);
  }

  .app-header__nav {
    order: 3;
    flex-basis: 100%;
    justify-content: space-around;
    border-top: 1px solid var(--color-gray-200);
    padding-top: var(--space-md);
    margin-top: var(--space-sm);
  }

  .app-header__link {
    padding: var(--space-xs) var(--space-sm);
    font-size: var(--text-sm);
  }
}
</style>
