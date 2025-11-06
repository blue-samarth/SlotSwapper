import { createRouter, createWebHistory } from 'vue-router';
import { authGuard } from './guards_router';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/authView/LoginView.vue'),
      meta: { requiresAuth: false, hideForAuth: true },
    },
    {
      path: '/signup',
      name: 'signup',
      component: () => import('@/views/authView/SignupView.vue'),
      meta: { requiresAuth: false, hideForAuth: true },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/events',
      name: 'events',
      component: () => import('@/views/EventsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/swaps',
      name: 'swaps',
      component: () => import('@/views/SwapsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      redirect: '/dashboard',
    },
  ],
});

// Apply authentication guard
router.beforeEach(authGuard);

export default router;