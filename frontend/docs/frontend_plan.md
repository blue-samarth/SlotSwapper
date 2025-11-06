# SlotSwapper Frontend Development Plan

## 📋 Project Overview

**Technology Stack:**
- **Framework:** Vue 3 (Composition API)
- **Language:** TypeScript
- **Build Tool:** Vite
- **State Management:** Pinia
- **Routing:** Vue Router
- **HTTP Client:** Axios
- **UI Framework:** Tailwind CSS / Vuetify / PrimeVue (TBD)
- **Date/Time:** date-fns or Day.js
- **Calendar:** FullCalendar Vue
- **Form Validation:** Vuelidate or VeeValidate
- **Icons:** Heroicons / Material Icons

**Estimated Time:** 12-14 hours

---

## 🎯 Core Features

### 1. Authentication System (2 hours)
- Login page
- Signup page
- JWT token management
- Auth guards for protected routes
- Auto-logout on token expiration

### 2. Dashboard (3 hours)
- Calendar view of events
- Quick stats (total events, swap requests)
- Navigation sidebar/header
- User profile dropdown

### 3. Event Management (3 hours)
- Create event modal/form
- List/grid view of events
- Edit event status (BUSY ↔ SWAPPABLE)
- Delete event confirmation
- Event status badges

### 4. Swap System (3 hours)
- Browse swappable slots (other users)
- Initiate swap request
- View sent swap requests
- View received swap requests
- Accept/Reject/Cancel swap actions
- Swap status badges

### 5. User Profile (1 hour)
- View profile information
- Update username
- Display account details

### 6. UI/UX Enhancements (1-2 hours)
- Loading states
- Error handling & toast notifications
- Responsive design
- Empty states
- Confirmation modals

---

## 📁 Project Structure

```
frontend/
├── public/
│   └── favicon.ico
├── src/
│   ├── api/                    # API service layer
│   │   ├── axios.ts           # Axios instance with interceptors
│   │   ├── auth.api.ts        # Auth endpoints
│   │   ├── events.api.ts      # Events endpoints
│   │   ├── swaps.api.ts       # Swap endpoints
│   │   └── users.api.ts       # User endpoints
│   ├── assets/                # Static assets
│   │   ├── styles/
│   │   │   └── main.css       # Global styles
│   │   └── images/
│   ├── components/            # Reusable components
│   │   ├── common/
│   │   │   ├── AppButton.vue
│   │   │   ├── AppModal.vue
│   │   │   ├── AppToast.vue
│   │   │   ├── LoadingSpinner.vue
│   │   │   └── ConfirmDialog.vue
│   │   ├── layout/
│   │   │   ├── AppHeader.vue
│   │   │   ├── AppSidebar.vue
│   │   │   └── AppLayout.vue
│   │   ├── events/
│   │   │   ├── EventCard.vue
│   │   │   ├── EventList.vue
│   │   │   ├── EventForm.vue
│   │   │   ├── EventCalendar.vue
│   │   │   └── EventStatusBadge.vue
│   │   └── swaps/
│   │       ├── SwapRequestCard.vue
│   │       ├── SwapRequestList.vue
│   │       ├── SwappableSlotCard.vue
│   │       ├── SwapStatusBadge.vue
│   │       └── InitiateSwapModal.vue
│   ├── composables/           # Composable functions
│   │   ├── useAuth.ts
│   │   ├── useEvents.ts
│   │   ├── useSwaps.ts
│   │   ├── useToast.ts
│   │   └── useModal.ts
│   ├── layouts/               # Layout components
│   │   ├── AuthLayout.vue     # For login/signup
│   │   └── MainLayout.vue     # For authenticated pages
│   ├── router/                # Vue Router configuration
│   │   ├── index.ts
│   │   └── guards.ts          # Auth guards
│   ├── stores/                # Pinia stores
│   │   ├── auth.store.ts
│   │   ├── events.store.ts
│   │   ├── swaps.store.ts
│   │   └── ui.store.ts
│   ├── types/                 # TypeScript types
│   │   ├── auth.types.ts
│   │   ├── event.types.ts
│   │   ├── swap.types.ts
│   │   ├── user.types.ts
│   │   └── api.types.ts
│   ├── utils/                 # Utility functions
│   │   ├── date.utils.ts
│   │   ├── validation.utils.ts
│   │   ├── storage.utils.ts
│   │   └── error.utils.ts
│   ├── views/                 # Page components
│   │   ├── auth/
│   │   │   ├── LoginView.vue
│   │   │   └── SignupView.vue
│   │   ├── DashboardView.vue
│   │   ├── EventsView.vue
│   │   ├── SwapsView.vue
│   │   │   ├── SentSwapsTab.vue
│   │   │   ├── ReceivedSwapsTab.vue
│   │   │   └── BrowseSlotsTab.vue
│   │   └── ProfileView.vue
│   ├── App.vue                # Root component
│   ├── main.ts                # Entry point
│   └── vite-env.d.ts          # Vite type declarations
├── .env.development           # Dev environment variables
├── .env.production            # Prod environment variables
├── .eslintrc.js               # ESLint configuration
├── .prettierrc                # Prettier configuration
├── index.html                 # HTML entry point
├── package.json               # Dependencies
├── tailwind.config.js         # Tailwind configuration
├── tsconfig.json              # TypeScript configuration
├── tsconfig.node.json         # TypeScript Node configuration
└── vite.config.ts             # Vite configuration
```

---

## 🔧 Phase 1: Project Setup (30 minutes)

### Tasks:
1. Initialize Vite project with Vue 3 + TypeScript
2. Install dependencies
3. Configure Tailwind CSS
4. Set up ESLint & Prettier
5. Create basic folder structure
6. Configure environment variables

### Commands:
```bash
cd frontend
npm create vite@latest . -- --template vue-ts
npm install
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
npm install vue-router@4 pinia axios
npm install @vueuse/core
npm install date-fns
npm install @heroicons/vue
```

### Environment Variables (.env.development):
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_NAME=SlotSwapper
VITE_APP_VERSION=1.0.0
```

---

## 🎨 Phase 2: Core Setup (1 hour)

### 2.1 API Layer Setup

**File: `src/api/axios.ts`**
```typescript
import axios from 'axios';
import type { AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/stores/auth.store';
import router from '@/router';

const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - Add auth token
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor - Handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore();
      authStore.logout();
      router.push('/login');
    }
    return Promise.reject(error);
  }
);

export default api;
```

### 2.2 TypeScript Types

**File: `src/types/auth.types.ts`**
```typescript
export interface SignupRequest {
  username: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}
```

**File: `src/types/event.types.ts`**
```typescript
export enum EventStatus {
  BUSY = 'BUSY',
  SWAPPABLE = 'SWAPPABLE',
  SWAP_PENDING = 'SWAP_PENDING',
}

export interface Event {
  id: number;
  title: string;
  user_id: number;
  start_time: string;
  end_time: string;
  status: EventStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateEventRequest {
  title: string;
  start_time: string;
  end_time: string;
}

export interface UpdateEventStatusRequest {
  status: EventStatus.BUSY | EventStatus.SWAPPABLE;
}
```

**File: `src/types/swap.types.ts`**
```typescript
import type { Event } from './event.types';
import type { User } from './auth.types';

export enum SwapStatus {
  PENDING = 'PENDING',
  ACCEPTED = 'ACCEPTED',
  REJECTED = 'REJECTED',
  CANCELLED = 'CANCELLED',
}

export interface SwapRequest {
  id: number;
  requester_event_id: number;
  receiver_event_id: number;
  requester_id: number;
  receiver_id: number;
  status: SwapStatus;
  requester_event: Event;
  receiver_event: Event;
  created_at: string;
  updated_at: string;
}

export interface InitiateSwapRequest {
  requester_event_id: number;
  receiver_event_id: number;
}

export interface SwappableSlot extends Event {
  user: User;
}
```

### 2.3 Pinia Stores

**File: `src/stores/auth.store.ts`**
```typescript
import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { User } from '@/types/auth.types';
import * as authApi from '@/api/auth.api';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(null);
  const isAuthenticated = ref(false);

  function setAuth(authData: { token: string; user: User }) {
    token.value = authData.token;
    user.value = authData.user;
    isAuthenticated.value = true;
    localStorage.setItem('token', authData.token);
    localStorage.setItem('user', JSON.stringify(authData.user));
  }

  function logout() {
    token.value = null;
    user.value = null;
    isAuthenticated.value = false;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }

  function loadFromStorage() {
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');
    if (storedToken && storedUser) {
      token.value = storedToken;
      user.value = JSON.parse(storedUser);
      isAuthenticated.value = true;
    }
  }

  async function login(email: string, password: string) {
    const response = await authApi.login({ email, password });
    setAuth(response.data);
  }

  async function signup(username: string, email: string, password: string) {
    const response = await authApi.signup({ username, email, password });
    setAuth(response.data);
  }

  return {
    user,
    token,
    isAuthenticated,
    login,
    signup,
    logout,
    loadFromStorage,
  };
});
```

### 2.4 Router Setup

**File: `src/router/index.ts`**
```typescript
import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth.store';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/signup',
      name: 'signup',
      component: () => import('@/views/auth/SignupView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      redirect: '/dashboard',
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
  ],
});

// Auth guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  authStore.loadFromStorage();

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login');
  } else if (!to.meta.requiresAuth && authStore.isAuthenticated) {
    next('/dashboard');
  } else {
    next();
  }
});

export default router;
```

---

## 🔐 Phase 3: Authentication (2 hours)

### 3.1 Login Page
**File: `src/views/auth/LoginView.vue`**

**Features:**
- Email & password inputs
- Form validation
- Loading state during login
- Error message display
- Link to signup page
- "Remember me" checkbox (optional)

**Key Functions:**
- Handle form submission
- Call `authStore.login()`
- Redirect to dashboard on success
- Display error messages from API

### 3.2 Signup Page
**File: `src/views/auth/SignupView.vue`**

**Features:**
- Username, email, password inputs
- Password confirmation field
- Form validation (min lengths, email format)
- Loading state
- Error display
- Link to login page

**Validation Rules:**
- Username: 3-50 characters
- Email: Valid email format
- Password: 6-100 characters
- Password confirmation must match

---

## 📊 Phase 4: Dashboard (3 hours)

### 4.1 Dashboard View
**File: `src/views/DashboardView.vue`**

**Features:**
- Welcome message with username
- Quick stats cards:
  - Total events
  - Swappable events
  - Pending swap requests (sent)
  - Pending swap requests (received)
- Calendar view of all events
- Quick action buttons:
  - Create new event
  - Browse swappable slots
- Recent swap requests (last 5)

### 4.2 Calendar Component
**File: `src/components/events/EventCalendar.vue`**

**Features:**
- Full calendar integration
- Display all user events
- Color-coded by status:
  - BUSY: Red
  - SWAPPABLE: Green
  - SWAP_PENDING: Orange
- Click event to view details
- Navigate months
- Today button

---

## 📅 Phase 5: Event Management (3 hours)

### 5.1 Events View
**File: `src/views/EventsView.vue`**

**Features:**
- List/grid view toggle
- Filter by status (All, BUSY, SWAPPABLE, SWAP_PENDING)
- Sort by date (ascending/descending)
- Search by title
- Create event button (opens modal)
- Display event cards

### 5.2 Event Card Component
**File: `src/components/events/EventCard.vue`**

**Props:**
- event: Event object

**Features:**
- Display event title, date, time
- Status badge (color-coded)
- Actions dropdown:
  - Toggle status (BUSY ↔ SWAPPABLE)
  - Initiate swap (if SWAPPABLE)
  - Delete event
- Click to expand details

### 5.3 Create/Edit Event Form
**File: `src/components/events/EventForm.vue`**

**Features:**
- Title input
- Date picker (start date)
- Time picker (start time, end time)
- Validation:
  - Title required (1-100 chars)
  - End time > Start time
  - Start time in future
- Submit button with loading state
- Cancel button

### 5.4 Event Status Badge
**File: `src/components/events/EventStatusBadge.vue`**

**Props:**
- status: EventStatus

**Design:**
- BUSY: Red badge
- SWAPPABLE: Green badge
- SWAP_PENDING: Orange badge

---

## 🔄 Phase 6: Swap System (3 hours)

### 6.1 Swaps View
**File: `src/views/SwapsView.vue`**

**Features:**
- Tab navigation:
  1. Browse Slots (swappable events from other users)
  2. Sent Requests (your outgoing swap requests)
  3. Received Requests (incoming swap requests)
- Tab-specific content

### 6.2 Browse Swappable Slots Tab
**File: `src/views/SwapsView/BrowseSlotsTab.vue`**

**Features:**
- List of all swappable slots from other users
- Filter by date range
- Search by username or title
- Display SlotCard with:
  - Event title, date, time
  - Owner username
  - "Request Swap" button
- Click "Request Swap" opens modal

### 6.3 Initiate Swap Modal
**File: `src/components/swaps/InitiateSwapModal.vue`**

**Props:**
- receiverEvent: SwappableSlot

**Features:**
- Display receiver's event details
- Dropdown to select your event to swap
- Only show your SWAPPABLE events
- Confirm button
- Cancel button
- Display error if API fails

### 6.4 Sent Swap Requests Tab
**File: `src/views/SwapsView/SentSwapsTab.vue`**

**Features:**
- List of outgoing swap requests
- Filter by status (All, PENDING, ACCEPTED, REJECTED, CANCELLED)
- Display SwapRequestCard with:
  - Your event details
  - Their event details
  - Status badge
  - Cancel button (if PENDING)
- Empty state if no requests

### 6.5 Received Swap Requests Tab
**File: `src/views/SwapsView/ReceivedSwapsTab.vue`**

**Features:**
- List of incoming swap requests
- Filter by status
- Display SwapRequestCard with:
  - Their event details
  - Your event details
  - Status badge
  - Accept/Reject buttons (if PENDING)
- Empty state if no requests

### 6.6 Swap Request Card
**File: `src/components/swaps/SwapRequestCard.vue`**

**Props:**
- swapRequest: SwapRequest
- type: 'sent' | 'received'

**Features:**
- Display both events side-by-side
- Show requester and receiver usernames
- Status badge
- Action buttons based on type and status
- Created/updated timestamps

---

## 👤 Phase 7: User Profile (1 hour)

### 7.1 Profile View
**File: `src/views/ProfileView.vue`**

**Features:**
- Display current user info:
  - Username
  - Email
  - Account created date
- Edit username form
- Validation (3-50 chars)
- Save button with loading state
- Success/error messages
- Logout button

---

## 🎨 Phase 8: UI Components (1-2 hours)

### 8.1 Common Components

**AppButton.vue:**
- Props: variant (primary, secondary, danger), loading, disabled
- Slot for content
- Consistent styling

**AppModal.vue:**
- Props: isOpen, title, size
- Slots: header, body, footer
- Close button
- Backdrop click to close

**LoadingSpinner.vue:**
- Props: size (small, medium, large)
- Centered spinner animation

**AppToast.vue:**
- Props: type (success, error, warning, info), message, duration
- Auto-dismiss after duration
- Close button

**ConfirmDialog.vue:**
- Props: isOpen, title, message, confirmText, cancelText
- Emit: confirm, cancel
- Used for delete confirmations

### 8.2 Layout Components

**AppHeader.vue:**
- Logo/App name
- Navigation links
- User dropdown:
  - View profile
  - Logout
- Responsive mobile menu

**AppSidebar.vue:**
- Navigation menu:
  - Dashboard
  - Events
  - Swaps
  - Profile
- Active route highlighting
- Collapsible on mobile

**MainLayout.vue:**
- Wrapper with Header + Sidebar + Content
- Responsive grid layout

---

## 🔧 Phase 9: Utilities & Composables (1 hour)

### 9.1 Utilities

**date.utils.ts:**
```typescript
import { format, parseISO, isAfter, isBefore } from 'date-fns';

export const formatDate = (date: string) => {
  return format(parseISO(date), 'MMM dd, yyyy');
};

export const formatTime = (date: string) => {
  return format(parseISO(date), 'hh:mm a');
};

export const formatDateTime = (date: string) => {
  return format(parseISO(date), 'MMM dd, yyyy hh:mm a');
};

export const isDateInFuture = (date: string) => {
  return isAfter(parseISO(date), new Date());
};
```

**error.utils.ts:**
```typescript
import type { AxiosError } from 'axios';

export const getErrorMessage = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    return error.response?.data?.message || 'An error occurred';
  }
  return 'An unexpected error occurred';
};

export const getErrorCode = (error: unknown): string | null => {
  if (axios.isAxiosError(error)) {
    return error.response?.data?.error || null;
  }
  return null;
};
```

### 9.2 Composables

**useToast.ts:**
```typescript
import { ref } from 'vue';

interface Toast {
  id: number;
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
}

const toasts = ref<Toast[]>([]);

export function useToast() {
  const show = (type: Toast['type'], message: string) => {
    const id = Date.now();
    toasts.value.push({ id, type, message });
    setTimeout(() => {
      toasts.value = toasts.value.filter((t) => t.id !== id);
    }, 5000);
  };

  const success = (message: string) => show('success', message);
  const error = (message: string) => show('error', message);
  const warning = (message: string) => show('warning', message);
  const info = (message: string) => show('info', message);

  return { toasts, success, error, warning, info };
}
```

**useModal.ts:**
```typescript
import { ref } from 'vue';

export function useModal() {
  const isOpen = ref(false);

  const open = () => {
    isOpen.value = true;
  };

  const close = () => {
    isOpen.value = false;
  };

  const toggle = () => {
    isOpen.value = !isOpen.value;
  };

  return { isOpen, open, close, toggle };
}
```

---

## 🎯 Phase 10: Polish & Testing (1-2 hours)

### 10.1 Features to Add:
- [ ] Loading states for all async operations
- [ ] Empty states with helpful messages
- [ ] Confirmation dialogs for destructive actions
- [ ] Toast notifications for success/error
- [ ] Responsive design (mobile, tablet, desktop)
- [ ] Keyboard navigation support
- [ ] Form validation error messages
- [ ] Rate limit warning display
- [ ] Optimistic UI updates

### 10.2 Error Handling:
- Network errors (connection lost)
- Authentication errors (auto-redirect to login)
- Validation errors (display field-specific errors)
- Rate limit errors (show countdown)
- Server errors (generic error message)

### 10.3 Performance:
- Lazy load routes
- Debounce search inputs
- Cache API responses (optional)
- Pagination for large lists (if needed)

---

## 📝 Development Checklist

### Setup Phase
- [ ] Initialize Vite project
- [ ] Install all dependencies
- [ ] Configure Tailwind CSS
- [ ] Set up folder structure
- [ ] Create TypeScript types
- [ ] Set up Axios instance
- [ ] Create Pinia stores
- [ ] Configure Vue Router

### Authentication
- [ ] Create login page
- [ ] Create signup page
- [ ] Implement JWT storage
- [ ] Add auth guards
- [ ] Handle token expiration
- [ ] Add logout functionality

### Dashboard
- [ ] Create dashboard layout
- [ ] Add stats cards
- [ ] Integrate calendar component
- [ ] Display recent activity
- [ ] Add quick actions

### Events
- [ ] Create events list view
- [ ] Create event card component
- [ ] Add create event modal
- [ ] Implement event status toggle
- [ ] Add delete event functionality
- [ ] Add filters and search

### Swaps
- [ ] Create swaps view with tabs
- [ ] Display swappable slots
- [ ] Create initiate swap modal
- [ ] Display sent swap requests
- [ ] Display received swap requests
- [ ] Implement accept/reject/cancel actions
- [ ] Add filters

### Profile
- [ ] Display user information
- [ ] Add edit username form
- [ ] Implement save functionality

### UI/UX
- [ ] Create common components (Button, Modal, Toast)
- [ ] Add loading states
- [ ] Add error handling
- [ ] Implement toast notifications
- [ ] Add confirmation dialogs
- [ ] Make responsive
- [ ] Add empty states

### Testing
- [ ] Manual testing of all features
- [ ] Test error scenarios
- [ ] Test on different devices
- [ ] Test with different data states

---

## 🚀 Deployment Checklist

- [ ] Update API base URL for production
- [ ] Build production bundle
- [ ] Test production build locally
- [ ] Configure environment variables
- [ ] Deploy to hosting (Vercel/Netlify/etc.)
- [ ] Test deployed application
- [ ] Set up CI/CD (optional)

---

## 📦 Dependencies Summary

```json
{
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.2.0",
    "pinia": "^2.1.0",
    "axios": "^1.6.0",
    "date-fns": "^3.0.0",
    "@vueuse/core": "^10.7.0",
    "@heroicons/vue": "^2.1.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0",
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "tailwindcss": "^3.4.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0",
    "eslint": "^8.56.0",
    "prettier": "^3.1.0"
  }
}
```

---

## 🎨 Design Guidelines

### Color Scheme:
- **Primary:** Blue (#3B82F6)
- **Secondary:** Gray (#6B7280)
- **Success:** Green (#10B981)
- **Warning:** Orange (#F59E0B)
- **Danger:** Red (#EF4444)
- **Info:** Cyan (#06B6D4)

### Event Status Colors:
- **BUSY:** Red (#EF4444)
- **SWAPPABLE:** Green (#10B981)
- **SWAP_PENDING:** Orange (#F59E0B)

### Swap Status Colors:
- **PENDING:** Yellow (#EAB308)
- **ACCEPTED:** Green (#10B981)
- **REJECTED:** Red (#EF4444)
- **CANCELLED:** Gray (#6B7280)

### Typography:
- **Font:** Inter or System UI
- **Headings:** Bold, larger sizes
- **Body:** Regular weight
- **Small text:** 0.875rem

### Spacing:
- Use Tailwind's spacing scale (4, 8, 16, 24, 32, 48, 64px)
- Consistent padding/margins throughout

---

## 🔒 Security Considerations

1. **Token Storage:** Store JWT in localStorage (or httpOnly cookie for enhanced security)
2. **XSS Prevention:** Sanitize user inputs
3. **CSRF Protection:** Use CORS properly
4. **Auth Guards:** Protect all authenticated routes
5. **API Errors:** Don't expose sensitive error details to users
6. **Rate Limiting:** Display remaining requests to prevent user frustration

---

## 📈 Future Enhancements (Post-MVP)

- Real-time notifications (WebSockets)
- Push notifications
- Email notifications for swap requests
- Advanced calendar features (recurring events)
- Event categories/tags
- User preferences and settings
- Dark mode
- Multi-language support (i18n)
- Export calendar to iCal
- Integration with Google Calendar
- Mobile app (React Native/Flutter)

---

## ✅ Success Criteria

The frontend is complete when:
- ✅ Users can signup and login
- ✅ Users can create, view, and delete events
- ✅ Users can toggle event status (BUSY ↔ SWAPPABLE)
- ✅ Users can browse swappable slots from other users
- ✅ Users can initiate swap requests
- ✅ Users can accept/reject received swap requests
- ✅ Users can cancel sent swap requests
- ✅ Users can view their profile and update username
- ✅ All API errors are handled gracefully
- ✅ Loading states are shown during async operations
- ✅ The app is responsive on mobile and desktop
- ✅ The app works correctly in production

---

**Estimated Total Time:** 12-14 hours

**Start Date:** TBD
**Target Completion:** TBD

---

## 📞 Support

For questions or issues during development, refer to:
- Backend API documentation: `frontend/docs/backend_api.md`
- Vue 3 docs: https://vuejs.org/
- Pinia docs: https://pinia.vuejs.org/
- Vue Router docs: https://router.vuejs.org/
- Tailwind CSS docs: https://tailwindcss.com/
