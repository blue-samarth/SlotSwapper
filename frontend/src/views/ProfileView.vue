<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useToast } from '../composables/useToast';
import { formatDate } from '../utils/date';
import { validateUsername } from '../utils/validation';
import AppButton from '../components/common/AppButton.vue';
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
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-900">SlotSwapper</h1>
        <div class="flex items-center gap-4">
          <router-link to="/dashboard" class="text-sm text-gray-600 hover:text-gray-900">Dashboard</router-link>
          <span class="text-sm text-gray-600">{{ authStore.user?.username }}</span>
          <button @click="handleLogout" class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700">
            Logout
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-3xl font-bold text-gray-900">My Profile</h2>
        <AppButton
          v-if="!isEditMode"
          variant="primary"
          @click="enableEditMode"
        >
          Edit Profile
        </AppButton>
      </div>

      <div class="bg-white rounded-lg shadow p-6 space-y-6">
        <!-- Username Field -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Username</label>
          <div v-if="!isEditMode" class="text-lg font-semibold">
            {{ authStore.user?.username }}
          </div>
          <div v-else>
            <input
              v-model="editForm.username"
              type="text"
              class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :class="{ 'border-red-500': errors.username }"
              placeholder="Enter username"
            />
            <p v-if="errors.username" class="mt-1 text-sm text-red-600">
              {{ errors.username }}
            </p>
          </div>
        </div>

        <!-- Email Field (Read-only) -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
          <div class="text-lg">{{ authStore.user?.email }}</div>
          <p class="mt-1 text-xs text-gray-500">Email cannot be changed</p>
        </div>

        <!-- User ID (Read-only) -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">User ID</label>
          <div class="text-sm text-gray-600 font-mono">{{ authStore.user?.id }}</div>
        </div>

        <!-- Member Since (Read-only) -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Member Since</label>
          <div>{{ formatDate(authStore.user?.created_at || '') }}</div>
        </div>

        <!-- Action Buttons (Edit Mode) -->
        <div v-if="isEditMode" class="flex gap-3 pt-4 border-t">
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