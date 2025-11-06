import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User } from '@/types/auth_types';
import * as authApi from '@/api/auth_api';

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null);
  const token = ref<string | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value);

  // Actions
  function setAuth(authData: { token: string; user: User }) {
    token.value = authData.token;
    user.value = authData.user;
    localStorage.setItem('token', authData.token);
    localStorage.setItem('user', JSON.stringify(authData.user));
  }

  function clearAuth() {
    token.value = null;
    user.value = null;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }

  function loadFromStorage() {
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');
    
    if (storedToken && storedUser) {
      try {
        token.value = storedToken;
        user.value = JSON.parse(storedUser);
      } catch (e) {
        console.error('Failed to parse stored user data', e);
        clearAuth();
      }
    }
  }

  async function login(email: string, password: string) {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await authApi.login({ email, password });
      setAuth(response.data);
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Login failed';
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  async function signup(username: string, email: string, password: string) {
    isLoading.value = true;
    error.value = null;
    
    try {
      const response = await authApi.signup({ username, email, password });
      setAuth(response.data);
    } catch (e: any) {
      error.value = e.response?.data?.message || 'Signup failed';
      throw e;
    } finally {
      isLoading.value = false;
    }
  }

  function logout() {
    clearAuth();
  }

  return {
    // State
    user,
    token,
    isLoading,
    error,
    // Getters
    isAuthenticated,
    // Actions
    login,
    signup,
    logout,
    loadFromStorage,
  };
});