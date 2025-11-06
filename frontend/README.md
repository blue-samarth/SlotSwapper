# 🎨 SlotSwapper Frontend

A modern Vue 3 + TypeScript application for managing and swapping time slots.

## 🚀 Tech Stack

- **Framework**: Vue 3 (Composition API with `<script setup>`)
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4
- **State Management**: Pinia
- **Routing**: Vue Router
- **HTTP Client**: Axios
- **Build Tool**: Vite
- **Date Utilities**: date-fns

## 📁 Project Structure

```
frontend/
├── src/
│   ├── api/                    # API client modules
│   │   ├── axios.ts           # Axios configuration with interceptors
│   │   ├── auth_api.ts        # Authentication endpoints
│   │   ├── events_api.ts      # Events endpoints
│   │   ├── swap_api.ts        # Swap requests endpoints
│   │   └── users_api.ts       # User management endpoints
│   ├── assets/
│   │   └── styles/
│   │       └── main.css       # Global styles + Tailwind imports
│   ├── components/
│   │   ├── common/            # Reusable UI components
│   │   │   ├── AppButton.vue
│   │   │   ├── AppHeader.vue
│   │   │   ├── AppModal.vue
│   │   │   ├── AppToast.vue
│   │   │   ├── ConfirmDialog.vue
│   │   │   └── LoadingSpinner.vue
│   │   ├── events/            # Event-specific components
│   │   │   ├── EventCard.vue
│   │   │   └── EventForm.vue
│   │   └── swaps/             # Swap-specific components
│   ├── composables/           # Reusable composition functions
│   │   └── useToast.ts        # Toast notification system
│   ├── router/                # Vue Router configuration
│   │   ├── index_router.ts    # Main router setup
│   │   ├── main_router.ts     # Route definitions
│   │   └── guards_router.ts   # Navigation guards (auth)
│   ├── stores/                # Pinia state management
│   │   ├── auth_store.ts      # Authentication state
│   │   ├── event_store.ts     # Events state
│   │   ├── swaps_store.ts     # Swap requests state
│   │   └── ui_store.ts        # UI state (modals, loading, etc.)
│   ├── types/                 # TypeScript type definitions
│   │   ├── api_types.ts       # Generic API types
│   │   ├── auth_types.ts      # Auth-related types
│   │   ├── event_types.ts     # Event types
│   │   ├── swap_types.ts      # Swap request types
│   │   └── user_types.ts      # User types
│   ├── utils/                 # Utility functions
│   │   ├── date.ts            # Date formatting utilities
│   │   └── validation.ts      # Form validation utilities
│   ├── views/                 # Page components
│   │   ├── authView/
│   │   │   ├── LoginView.vue
│   │   │   └── SignupView.vue
│   │   ├── swapView/
│   │   │   ├── BrowseSlotsTabView.vue
│   │   │   ├── ReceivedSwapsTabView.vue
│   │   │   └── SentSwapsTabView.vue
│   │   ├── DashboardView.vue
│   │   ├── EventsView.vue
│   │   ├── ProfileView.vue
│   │   └── SwapsView.vue
│   ├── App.vue                # Root component
│   ├── main.ts                # Application entry point
│   └── style.css              # Additional global styles
├── public/                    # Static assets
├── Dockerfile                 # Multi-stage production build
├── nginx.conf                 # Nginx configuration for production
├── package.json               # Dependencies and scripts
├── postcss.config.js          # PostCSS configuration
├── tailwind.config.js         # Tailwind CSS configuration
├── tsconfig.json              # TypeScript configuration
└── vite.config.ts             # Vite build configuration
```

## 🛠️ Development Setup

### Prerequisites

- Node.js 18+ and npm
- Backend API running (default: http://localhost:8080)

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

The app will be available at **http://localhost:5173**

### Environment Variables

Create a `.env` file:

```env
VITE_API_BASE_URL=http://localhost:8080/api
```

## 📦 Build & Deploy

### Development Build

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

### Docker Build (Production)

```bash
# Build Docker image
docker build -t slotswapper-frontend .

# Run container
docker run -p 80:80 slotswapper-frontend
```

### Docker Compose (Full Stack)

```bash
# From project root
docker-compose up --build
```

Access at **http://localhost**

## 🎯 Key Features

### Components

#### Common Components
- **AppToast** - Global toast notification system (success, error, warning, info)
- **AppModal** - Reusable modal dialog
- **AppButton** - Consistent button styling (primary, secondary, danger variants)
- **AppHeader** - Navigation header with routing
- **ConfirmDialog** - Confirmation dialog for destructive actions
- **LoadingSpinner** - Loading states (3 sizes, fullscreen mode)

#### Event Components
- **EventCard** - Event card with status badges and actions
- **EventForm** - Form for creating/editing events

### State Management (Pinia Stores)

#### Auth Store
- User authentication state
- Login/logout/signup actions
- Token management (localStorage)
- Auto-load on app start

#### Events Store
- Events CRUD operations
- Filter by status (BUSY/SWAPPABLE)
- Toggle event swappable status

#### Swaps Store
- Browse swappable slots
- Manage sent/received swap requests
- Accept/reject/cancel swaps

### Utilities

#### Date Utilities (`utils/date.ts`)
- `formatDate()` - Format date (e.g., "Mon, Nov 6, 2025")
- `formatTime()` - Format time (e.g., "2:30 PM")
- `formatDateTime()` - Combined date and time
- `calculateDuration()` - Calculate time difference
- `getRelativeTime()` - Human-readable relative time
- `isPast()` - Check if date is in the past
- `isToday()` - Check if date is today
- `getDayOfWeek()` - Get day name

#### Validation Utilities (`utils/validation.ts`)
- `validateEmail()` - Email format validation
- `validatePassword()` - Password strength validation
- `validateUsername()` - Username validation
- `validateRequired()` - Required field validation
- `validateMinLength()` / `validateMaxLength()` - Length validation
- `validateFutureDate()` - Date must be in future
- `validateDateFormat()` - Date format validation
- `validateTimeFormat()` - Time format validation

### Toast Notification System

Global toast notifications using provide/inject pattern:

```typescript
import { useToast } from '@/composables/useToast';

const toast = useToast();

// Show notifications
toast.success('Event created successfully!');
toast.error('Failed to delete event');
toast.warning('Event status changed');
toast.info('Loading data...');
```

Features:
- 4 notification types with icons
- Auto-dismiss after 5 seconds
- Click to dismiss
- Smooth animations
- Stack multiple toasts

## 🎨 Styling

### Tailwind CSS v4

Using the latest Tailwind CSS v4 with PostCSS plugin:

```css
/* src/assets/styles/main.css */
@import "tailwindcss";
```

### Custom Styles

- Component-scoped styles with `<style scoped>`
- Global styles in `src/assets/styles/main.css`
- Utility-first approach with Tailwind classes

### Responsive Design

Mobile-first responsive design:
- `sm:` - 640px and up
- `md:` - 768px and up
- `lg:` - 1024px and up
- `xl:` - 1280px and up

## 🔒 Authentication

### JWT Token Flow

1. User logs in → Backend returns JWT token
2. Token stored in `localStorage`
3. Axios interceptor adds token to all requests
4. On 401 error → Auto logout and redirect to login

### Protected Routes

Routes with `meta: { requiresAuth: true }` require authentication:
- `/dashboard`
- `/events`
- `/swaps`
- `/profile`

Routes with `meta: { hideForAuth: true }` redirect authenticated users:
- `/login`
- `/signup`

## 🐛 Debugging

### Console Logs

Debug features are enabled in development:

```typescript
// Check toast system
window.__toast // Access toast methods in console

// Store debugging
const authStore = useAuthStore();
console.log(authStore.user);
```

### Vue DevTools

Install [Vue DevTools](https://devtools.vuejs.org/) for component inspection and state debugging.

## 📝 Code Style

- **Composition API** with `<script setup>`
- **TypeScript** for type safety
- **ESLint** for linting (when configured)
- **Prettier** for code formatting (when configured)

### Component Template

```vue
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { User } from '@/types/auth_types';

// Props
const props = defineProps<{
  user: User;
}>();

// Emits
const emit = defineEmits<{
  update: [id: number];
}>();

// State
const loading = ref(false);

// Computed
const userName = computed(() => props.user.username);

// Methods
const handleUpdate = () => {
  emit('update', props.user.id);
};

// Lifecycle
onMounted(() => {
  console.log('Component mounted');
});
</script>

<template>
  <div>
    <!-- Template -->
  </div>
</template>

<style scoped>
/* Scoped styles */
</style>
```

## 🚢 Production Deployment

### Build Optimizations

- Multi-stage Docker build (build stage + nginx stage)
- Tree-shaking with Vite
- Code splitting for routes
- Asset minification
- Gzip compression in nginx

### Nginx Configuration

Production nginx serves:
- Static assets with 1-year caching
- API proxy to backend (avoids CORS)
- SPA fallback for Vue Router
- Security headers
- Gzip compression

## 🔧 Troubleshooting

### Build Errors

**TypeScript errors in production build:**
- The Dockerfile uses `vite build` directly (skips `vue-tsc`)
- This is intentional for faster Docker builds
- Run `npm run build` locally to check type errors

**CSS not loading:**
- Ensure PostCSS is configured: `postcss.config.js`
- Check Tailwind import in `main.css`
- Restart dev server after config changes

### CORS Errors

**In Docker:**
- Frontend nginx proxies `/api/*` to backend
- API calls use relative paths (`/api/...`)

**In Development:**
- Frontend calls `http://localhost:8080/api`
- Backend CORS must allow `http://localhost:5173`

### Invalid Date Errors

- Always add null checks: `v-if="event.start_time"`
- Use utility functions: `formatDate(date)` handles null safely
- Add fallbacks: `event.user?.username || 'Unknown'`

## 📚 Additional Resources

- [Vue 3 Documentation](https://vuejs.org/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)

## 🤝 Contributing

1. Create feature branch
2. Make changes with proper TypeScript types
3. Test locally with `npm run dev`
4. Test production build with `npm run build`
5. Submit pull request

## 📄 License

[Add your license here]
