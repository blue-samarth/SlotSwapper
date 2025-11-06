import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';

/**
 * Navigation guard to protect routes that require authentication
 */
export function authGuard(
  to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const authStore = useAuthStore();
  
  // Load auth state from localStorage on first load
  authStore.loadFromStorage();

  const isAuthenticated = authStore.isAuthenticated;
  const requiresAuth = to.meta.requiresAuth;
  const hideForAuth = to.meta.hideForAuth;

  // If route requires auth and user is not authenticated
  if (requiresAuth && !isAuthenticated) {
    next('/login');
  }
  // If user is authenticated and tries to access login/signup
  else if (hideForAuth && isAuthenticated) {
    next('/dashboard');
  }
  // Allow navigation
  else {
    next();
  }
}