<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useToast } from '../composables/useToast';
import { formatDate } from '../utils/date';
import { validateUsername } from '../utils/validation';
import AppButton from '../components/common/AppButton.vue';
import AppHeader from '../components/common/AppHeader.vue';
import * as usersApi from '../api/users_api';

const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

// Edit mode state
const isEditMode = ref(false);
const isSubmitting = ref(false);
const editForm = ref({
  username: authStore.user?.username || '',
});
const errors = ref<Record<string, string>>({});

// Computed
const hasChanges = computed(() => {
  return editForm.value.username !== authStore.user?.username;
});

const handleLogout = () => {
  authStore.logout();
  router.push('/login');
};

const enableEditMode = () => {
  editForm.value.username = authStore.user?.username || '';
  errors.value = {};
  isEditMode.value = true;
};

const cancelEdit = () => {
  isEditMode.value = false;
  editForm.value.username = authStore.user?.username || '';
  errors.value = {};
};

const validateForm = (): boolean => {
  errors.value = {};

  const usernameValidation = validateUsername(editForm.value.username);
  if (!usernameValidation.valid) {
    errors.value.username = usernameValidation.error || 'Invalid username';
    return false;
  }

  return true;
};

const handleSaveProfile = async () => {
  if (!validateForm()) {
    return;
  }

  if (!hasChanges.value) {
    toast.info('No changes to save');
    isEditMode.value = false;
    return;
  }

  isSubmitting.value = true;
  try {
    const response = await usersApi.updateUsername(editForm.value.username);
    
    // Update auth store with new user data
    if (authStore.user) {
      authStore.user.username = response.data.username;
    }
    
    isEditMode.value = false;
    toast.success('Profile updated successfully!');
  } catch (error: any) {
    console.error('Failed to update profile:', error);
    const errorMessage = error.response?.data?.message || 'Failed to update profile';
    toast.error(errorMessage);
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <div class="profile-view">
    <!-- Header -->
    <AppHeader :username="authStore.user?.username" @logout="handleLogout" />

    <!-- Main Content -->
    <main class="profile-view__content">
      <div class="profile-view__header">
        <h2 class="profile-view__title">My Profile</h2>
        <AppButton
          v-if="!isEditMode"
          variant="primary"
          @click="enableEditMode"
        >
          Edit Profile
        </AppButton>
      </div>

      <div class="profile-view__card">
        <!-- Username Field -->
        <div class="profile-view__field">
          <label class="profile-view__label">Username</label>
          <div v-if="!isEditMode" class="profile-view__value">
            {{ authStore.user?.username }}
          </div>
          <div v-else>
            <input
              v-model="editForm.username"
              type="text"
              class="profile-view__input"
              :class="{ 'profile-view__input--error': errors.username }"
              placeholder="Enter username"
            />
            <p v-if="errors.username" class="profile-view__error">
              {{ errors.username }}
            </p>
          </div>
        </div>

        <!-- Email Field (Read-only) -->
        <div class="profile-view__field">
          <label class="profile-view__label">Email</label>
          <div class="profile-view__value">{{ authStore.user?.email }}</div>
          <p class="profile-view__hint">Email cannot be changed</p>
        </div>

        <!-- User ID (Read-only) -->
        <div class="profile-view__field">
          <label class="profile-view__label">User ID</label>
          <div class="profile-view__value profile-view__value--mono">{{ authStore.user?.id }}</div>
        </div>

        <!-- Member Since (Read-only) -->
        <div class="profile-view__field">
          <label class="profile-view__label">Member Since</label>
          <div class="profile-view__value">{{ formatDate(authStore.user?.created_at || '') }}</div>
        </div>

        <!-- Action Buttons (Edit Mode) -->
        <div v-if="isEditMode" class="profile-view__actions">
          <AppButton
            variant="secondary"
            :disabled="isSubmitting"
            @click="cancelEdit"
            full-width
          >
            Cancel
          </AppButton>
          <AppButton
            variant="primary"
            :loading="isSubmitting"
            :disabled="!hasChanges"
            @click="handleSaveProfile"
            full-width
          >
            Save Changes
          </AppButton>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.profile-view {
  min-height: 100vh;
  background: var(--gradient-bg-subtle);
}

.profile-view__content {
  max-width: 48rem;
  margin: 0 auto;
  padding: var(--space-xl) var(--space-lg);
}

.profile-view__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-xl);
}

.profile-view__title {
  font-size: var(--text-3xl);
  font-weight: var(--weight-bold);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.profile-view__card {
  background: var(--color-bg-card);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
  padding: var(--space-xl);
  display: flex;
  flex-direction: column;
  gap: var(--space-xl);
}

.profile-view__field {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.profile-view__label {
  font-size: var(--text-sm);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
}

.profile-view__value {
  font-size: var(--text-lg);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
}

.profile-view__value--mono {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.profile-view__input {
  width: 100%;
  padding: var(--space-sm) var(--space-md);
  border: 2px solid var(--color-gray-300);
  border-radius: var(--radius-md);
  font-size: var(--text-base);
  transition: var(--transition-all);
}

.profile-view__input:focus {
  outline: none;
  border-color: var(--color-secondary);
  box-shadow: 0 0 0 3px rgba(78, 205, 196, 0.15);
}

.profile-view__input--error {
  border-color: var(--color-error);
}

.profile-view__error {
  margin-top: calc(var(--space-xs) * -1);
  font-size: var(--text-sm);
  color: var(--color-error);
}

.profile-view__hint {
  margin-top: calc(var(--space-xs) * -1);
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  font-style: italic;
}

.profile-view__actions {
  display: flex;
  gap: var(--space-sm);
  padding-top: var(--space-md);
  border-top: 1px solid var(--color-gray-200);
}
</style>