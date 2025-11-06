<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth_store';
import { useToast } from '../../composables/useToast';

const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const username = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');

const handleSignup = async () => {
  try {
    if (password.value !== confirmPassword.value) {
      toast.error('Passwords do not match');
      return;
    }

    if (password.value.length < 6) {
      toast.error('Password must be at least 6 characters');
      return;
    }

    if (username.value.length < 3) {
      toast.error('Username must be at least 3 characters');
      return;
    }

    await authStore.signup(username.value, email.value, password.value);
    toast.success('Account created successfully! Welcome to SlotSwapper.');
    router.push('/dashboard');
  } catch (error: any) {
    const errorMessage = error.response?.data?.message || 'Signup failed';
    toast.error(errorMessage);
  }
};
</script>

<template>
  <div class="signup-view">
    <div class="signup-view__card">
      <h2 class="signup-view__title">Create Account</h2>
      
      <form @submit.prevent="handleSignup" class="signup-view__form">
        <div class="signup-view__field">
          <label class="signup-view__label">Username</label>
          <input
            v-model="username"
            type="text"
            required
            minlength="3"
            class="signup-view__input"
            placeholder="johndoe"
          />
        </div>

        <div class="signup-view__field">
          <label class="signup-view__label">Email</label>
          <input
            v-model="email"
            type="email"
            required
            class="signup-view__input"
            placeholder="your@email.com"
          />
        </div>

        <div class="signup-view__field">
          <label class="signup-view__label">Password</label>
          <input
            v-model="password"
            type="password"
            required
            minlength="6"
            class="signup-view__input"
            placeholder="••••••••"
          />
        </div>

        <div class="signup-view__field">
          <label class="signup-view__label">Confirm Password</label>
          <input
            v-model="confirmPassword"
            type="password"
            required
            class="signup-view__input"
            placeholder="••••••••"
          />
        </div>

        <button
          type="submit"
          :disabled="authStore.isLoading"
          class="signup-view__submit"
        >
          {{ authStore.isLoading ? 'Creating account...' : 'Sign Up' }}
        </button>
      </form>

      <p class="signup-view__footer">
        Already have an account?
        <router-link to="/login" class="signup-view__link">
          Login
        </router-link>
      </p>
    </div>
  </div>
</template>

<style scoped>
.signup-view {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gradient-hero);
  padding: var(--space-lg);
}

.signup-view__card {
  max-width: 28rem;
  width: 100%;
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-2xl);
  padding: var(--space-2xl);
  border: 1px solid var(--color-gray-200);
}

.signup-view__title {
  font-size: var(--text-2xl);
  font-weight: var(--weight-bold);
  text-align: center;
  margin-bottom: var(--space-xl);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.signup-view__form {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.signup-view__field {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.signup-view__label {
  font-size: var(--text-sm);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
}

.signup-view__input {
  width: 100%;
  padding: var(--space-sm) var(--space-md);
  border: 2px solid var(--color-gray-300);
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  transition: var(--transition-all);
}

.signup-view__input:focus {
  outline: none;
  border-color: var(--color-secondary);
  box-shadow: 0 0 0 3px rgba(78, 205, 196, 0.15);
  transform: translateY(-2px);
}

.signup-view__submit {
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

.signup-view__submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

.signup-view__submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.signup-view__footer {
  text-align: center;
  margin-top: var(--space-md);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.signup-view__link {
  color: var(--color-secondary);
  font-weight: var(--weight-semibold);
  text-decoration: none;
  transition: var(--transition-all);
}

.signup-view__link:hover {
  color: var(--color-primary);
  text-decoration: underline;
}
</style>