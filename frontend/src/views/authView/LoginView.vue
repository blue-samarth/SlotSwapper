<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth_store';
import { useToast } from '../../composables/useToast';

const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const email = ref('');
const password = ref('');

const handleLogin = async () => {
  try {
    await authStore.login(email.value, password.value);
    toast.success('Login successful! Welcome back.');
    router.push('/dashboard');
  } catch (error: any) {
    const errorMessage = error.response?.data?.message || 'Login failed';
    toast.error(errorMessage);
  }
};
</script>

<template>
  <div class="login-view">
    <div class="login-view__card">
      <h2 class="login-view__title">Login to SlotSwapper</h2>
      
      <form @submit.prevent="handleLogin" class="login-view__form">
        <div class="login-view__field">
          <label class="login-view__label">Email</label>
          <input
            v-model="email"
            type="email"
            required
            class="login-view__input"
            placeholder="your@email.com"
          />
        </div>

        <div class="login-view__field">
          <label class="login-view__label">Password</label>
          <input
            v-model="password"
            type="password"
            required
            class="login-view__input"
            placeholder="••••••••"
          />
        </div>

        <button
          type="submit"
          :disabled="authStore.isLoading"
          class="login-view__submit"
        >
          {{ authStore.isLoading ? 'Logging in...' : 'Login' }}
        </button>
      </form>

      <p class="login-view__footer">
        Don't have an account?
        <router-link to="/signup" class="login-view__link">
          Sign up
        </router-link>
      </p>
    </div>
  </div>
</template>

<style scoped>
.login-view {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gradient-hero);
  padding: var(--space-lg);
}

.login-view__card {
  max-width: 28rem;
  width: 100%;
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-2xl);
  padding: var(--space-2xl);
  border: 1px solid var(--color-gray-200);
}

.login-view__title {
  font-size: var(--text-2xl);
  font-weight: var(--weight-bold);
  text-align: center;
  margin-bottom: var(--space-xl);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.login-view__form {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.login-view__field {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.login-view__label {
  font-size: var(--text-sm);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
}

.login-view__input {
  width: 100%;
  padding: var(--space-sm) var(--space-md);
  border: 2px solid var(--color-gray-300);
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  transition: var(--transition-all);
}

.login-view__input:focus {
  outline: none;
  border-color: var(--color-secondary);
  box-shadow: 0 0 0 3px rgba(78, 205, 196, 0.15);
  transform: translateY(-2px);
}

.login-view__submit {
  width: 100%;
  padding: var(--space-md);
  background: var(--gradient-primary-btn);
  color: var(--color-text-inverse);
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  font-weight: var(--weight-semibold);
  cursor: pointer;
  transition: var(--transition-all);
  box-shadow: var(--shadow-primary);
  margin-top: var(--space-sm);
}

.login-view__submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

.login-view__submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.login-view__footer {
  text-align: center;
  margin-top: var(--space-md);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.login-view__link {
  color: var(--color-secondary);
  font-weight: var(--weight-semibold);
  text-decoration: none;
  transition: var(--transition-all);
}

.login-view__link:hover {
  color: var(--color-primary);
  text-decoration: underline;
}
</style>