# SlotSwapper

A full-stack web application for managing and swapping time slots between users. Built with Go backend and Vue.js frontend, SlotSwapper allows users to mark their calendar events as swappable and request swaps with other users.

## Table of Contents

- [Overview](#overview)
- [Design Choices](#design-choices)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Local Setup Instructions](#local-setup-instructions)
- [API Documentation](#api-documentation)
- [Assumptions](#assumptions)
- [Challenges Faced](#challenges-faced)
- [Technology Stack](#technology-stack)

## Overview

SlotSwapper is a time slot management and exchange platform that enables users to:

- Create and manage calendar events with specific time slots
- Mark events as swappable or busy
- Browse available swappable slots from other users
- Send swap requests to exchange time slots
- Accept or reject incoming swap requests
- Track the status of all swap requests (pending, accepted, rejected, cancelled)

The application uses JWT authentication for secure user sessions and provides a clean, responsive user interface built with Vue 3 and Tailwind CSS.

## Design Choices

### Backend Architecture

**Framework**: Go with Gin web framework for high performance and built-in middleware support.

**Database**: PostgreSQL for relational data with GORM as the ORM for type-safe database operations.

**Authentication**: JWT-based stateless authentication with tokens stored in HTTP-only cookies (when applicable) and localStorage for SPA usage.

**API Design**: RESTful API design with consistent response structure wrapping all responses in `{message, data}` format for predictable client-side handling.

**Error Handling**: Centralized error handling middleware that translates all errors into consistent JSON responses with appropriate HTTP status codes.

**Transaction Management**: Database transactions for operations that modify multiple tables (e.g., swap acceptance updates both events and swap request status).

### Frontend Architecture

**Framework**: Vue 3 with Composition API and TypeScript for type safety and modern reactive patterns.

**State Management**: Pinia stores for global state management with separate stores for authentication, events, swaps, and UI state.

**Styling**: Tailwind CSS v4 for utility-first styling with custom PostCSS configuration.

**Routing**: Vue Router with navigation guards for authentication and authorization checks.

**API Communication**: Axios with centralized configuration and request/response interceptors for token injection and error handling.

**Component Design**: Reusable component library including buttons, modals, toasts, loading spinners, and form components for consistency across the application.

### Database Schema

**Users Table**: Stores user credentials with bcrypt-hashed passwords.

**Events Table**: Stores time slot information with status (BUSY or SWAPPABLE) and foreign key to user.

**Swap Requests Table**: Stores swap request information linking requester event and receiver event with status tracking (PENDING, ACCEPTED, REJECTED, CANCELLED).

### Security Considerations

- Passwords are hashed using bcrypt before storage
- JWT tokens expire after 24 hours
- CORS configuration allows only specified frontend origins
- Rate limiting middleware prevents abuse (100 requests per minute per IP)
- Input validation on all API endpoints
- SQL injection prevention through GORM parameterized queries

## Project Structure

```
SlotSwapper/
├── backend/
│   ├── config/          # Database configuration
│   ├── controllers/     # HTTP request handlers
│   ├── middleware/      # Auth, CORS, logging, rate limiting
│   ├── models/          # Database models
│   ├── routes/          # Route definitions
│   ├── services/        # Business logic layer
│   ├── utils/           # Helper functions (JWT, errors)
│   ├── tests/           # Unit and integration tests
│   ├── go.mod           # Go dependencies
│   └── main.go          # Application entry point
├── frontend/
│   ├── src/
│   │   ├── api/         # API client modules
│   │   ├── components/  # Reusable Vue components
│   │   ├── composables/ # Vue composables
│   │   ├── router/      # Vue Router configuration
│   │   ├── stores/      # Pinia state stores
│   │   ├── types/       # TypeScript type definitions
│   │   ├── utils/       # Utility functions
│   │   └── views/       # Page components
│   ├── package.json     # Node dependencies
│   └── vite.config.ts   # Vite configuration
└── docker-compose.yaml  # Docker orchestration
```

## Prerequisites

### For Local Development

- Go 1.24 or higher
- Node.js 20.x or higher
- PostgreSQL 16.x or higher
- npm or yarn package manager

### For Docker Deployment

- Docker 20.x or higher
- Docker Compose 2.x or higher

## Local Setup Instructions

### Step 1: Clone the Repository

```bash
git clone https://github.com/blue-samarth/SlotSwapper.git
cd SlotSwapper
```

### Step 2: Database Setup

Install and start PostgreSQL, then create the database:

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE slotswapper;

# Exit psql
\q
```

### Step 3: Backend Setup

Navigate to the backend directory and install dependencies:

```bash
cd backend
go mod download
```

Create a `.env` file in the backend directory (optional, defaults will be used):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=slotswapper
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
```

Run database migrations (auto-migration runs on startup):

```bash
go run main.go
```

The backend will start on `http://localhost:8080` and automatically create the required database tables.

### Step 4: Frontend Setup

Open a new terminal, navigate to the frontend directory, and install dependencies:

```bash
cd frontend
npm install
```

Create a `.env` file in the frontend directory (optional):

```env
VITE_API_BASE_URL=http://localhost:8080/api
```

Start the development server:

```bash
npm run dev
```

The frontend will start on `http://localhost:5173`.

### Step 5: Access the Application

Open your browser and navigate to `http://localhost:5173`. You can now:

1. Sign up for a new account
2. Log in with your credentials
3. Create events with time slots
4. Mark events as swappable
5. Browse and request swaps from other users

### Running Tests

Backend tests:

```bash
cd backend
go test ./tests/... -v
```

Frontend tests (if configured):

```bash
cd frontend
npm run test
```

## Docker Setup (Alternative)

If you prefer Docker, you can run the entire stack with Docker Compose:

```bash
# Build and start all services
docker-compose up --build

# Access the application at http://localhost:80
```

This will start:
- PostgreSQL database on port 5432
- Backend API on port 8080
- Frontend on port 80 (nginx)

For more details, see `DOCKER.md`.

## API Documentation

### Base URL

```
http://localhost:8080/api
```

### Response Format

All API responses follow this structure:

```json
{
  "message": "Success message or error description",
  "data": "Response data (object, array, or null)"
}
```

### Authentication

Include JWT token in Authorization header for protected endpoints:

```
Authorization: Bearer <token>
```

### Endpoints

#### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/signup` | Create new user account | No |
| POST | `/auth/login` | Login and receive JWT token | No |
| POST | `/auth/logout` | Logout current user | Yes |
| GET | `/auth/me` | Get current user info | Yes |

**Signup Request:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Login Request:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Login Response:**
```json
{
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### Events

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/events` | Get all events for current user | Yes |
| POST | `/events` | Create new event | Yes |
| GET | `/events/:id` | Get event by ID | Yes |
| PUT | `/events/:id` | Update event | Yes |
| DELETE | `/events/:id` | Delete event | Yes |
| PUT | `/events/:id/status` | Toggle event status (BUSY/SWAPPABLE) | Yes |

**Create Event Request:**
```json
{
  "title": "Team Meeting",
  "description": "Weekly team sync",
  "start_time": "2025-11-10T14:00:00Z",
  "end_time": "2025-11-10T15:00:00Z",
  "location": "Conference Room A"
}
```

**Update Event Status Request:**
```json
{
  "status": "SWAPPABLE"
}
```

#### Swaps

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/swaps/available` | Get all swappable slots from other users | Yes |
| POST | `/swaps/initiate` | Initiate swap request | Yes |
| GET | `/swaps/sent` | Get swap requests sent by current user | Yes |
| GET | `/swaps/received` | Get swap requests received by current user | Yes |
| PUT | `/swaps/:id/accept` | Accept swap request | Yes |
| PUT | `/swaps/:id/reject` | Reject swap request | Yes |
| DELETE | `/swaps/:id/cancel` | Cancel swap request | Yes |

**Initiate Swap Request:**
```json
{
  "requester_event_id": 5,
  "receiver_event_id": 12
}
```

**Available Slots Response:**
```json
{
  "message": "Swappable slots fetched successfully",
  "data": [
    {
      "id": 12,
      "title": "Office Hours",
      "start_time": "2025-11-11T10:00:00Z",
      "end_time": "2025-11-11T11:00:00Z",
      "status": "SWAPPABLE",
      "user": {
        "id": 2,
        "username": "jane_smith",
        "email": "jane@example.com"
      }
    }
  ]
}
```

#### Users

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/users/:id` | Get user by ID | Yes |
| PUT | `/users/:id` | Update user profile | Yes |

**Update User Request:**
```json
{
  "username": "new_username"
}
```

#### Health Check

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Check API health status | No |

## Assumptions

1. **Single Time Zone**: All times are stored and handled in UTC. The frontend is responsible for displaying times in the user's local timezone.

2. **Event Ownership**: Users can only modify or delete their own events. Swap requests can only be initiated for events marked as SWAPPABLE.

3. **Swap Request Lifecycle**: A swap request can only be in one of four states: PENDING, ACCEPTED, REJECTED, or CANCELLED. Once accepted or rejected, the request cannot be modified.

4. **Event Status Constraints**: Events with pending swap requests cannot be deleted to maintain referential integrity. Users must cancel the swap request first.

5. **One-to-One Swaps**: Each swap request involves exactly two events: one from the requester and one from the receiver. Group swaps or multi-party swaps are not supported.

6. **User Authentication**: JWT tokens expire after 24 hours and must be refreshed by logging in again. No refresh token mechanism is implemented.

7. **Email Uniqueness**: Email addresses must be unique across all users. Usernames must also be unique.

8. **Password Requirements**: Minimum password length is 8 characters. No maximum length or complexity requirements are enforced at the API level (though the frontend validates for better UX).

9. **Concurrent Requests**: The system does not handle race conditions where two users might try to request a swap for the same event simultaneously. The first request to complete will succeed.

10. **Event Overlap**: The system does not prevent users from creating overlapping events. Users are responsible for managing their own schedules.

## Challenges Faced

### 1. Backend Response Structure Mismatch

**Challenge**: Initially, the backend returned data directly, but later was changed to wrap responses in `{message, data}` structure. This broke the frontend which expected `response.data` to contain the actual data.

**Solution**: Updated all frontend API calls to access `response.data.data` and added null safety checks to handle cases where the backend returns `{data: null}`.

### 2. CORS Issues in Docker Environment

**Challenge**: When running in Docker containers, the frontend could not communicate with the backend due to CORS restrictions and different network contexts.

**Solution**: Implemented nginx reverse proxy in the frontend container to proxy `/api/` requests to the backend container, eliminating cross-origin requests entirely.

### 3. Event Status Behavior

**Challenge**: The backend always creates events with status BUSY regardless of what the frontend sends, but this was not documented or communicated to the frontend.

**Solution**: Removed the status field from the event creation form and added an informational note explaining that events are created as BUSY by default and can be toggled to SWAPPABLE afterwards.

### 4. Null Safety in Frontend

**Challenge**: The backend sometimes returns null or undefined values for related entities (e.g., `user` field in events), causing "Cannot read properties of null" errors in the frontend.

**Solution**: Added comprehensive null safety checks throughout the frontend using optional chaining (`?.`), v-if guards before rendering, and fallback values (`|| 'Unknown'`, `|| []`) for all data access.

### 5. Swap Request Immediate Visibility

**Challenge**: After initiating a swap request, it would not appear in the "Sent" tab until the page was manually refreshed.

**Solution**: Added automatic data refresh calls (`fetchSentRequests()` and `fetchEvents()`) after successful swap initiation to update the UI immediately.

### 6. Event Deletion with Pending Swaps

**Challenge**: The backend returns a 400 error when trying to delete an event that has pending swap requests, but no user-friendly error message was shown.

**Solution**: Added frontend validation to check event status before deletion and display a helpful toast warning message explaining that events with pending swaps cannot be deleted.

### 7. Git Merge Conflicts

**Challenge**: Multiple developers working on the same files led to merge conflicts with duplicate code sections and conflict markers appearing in the files.

**Solution**: Carefully resolved conflicts by keeping the version with null safety improvements and bug fixes, then tested thoroughly to ensure no functionality was lost.

### 8. Tailwind CSS v4 Configuration

**Challenge**: Tailwind CSS v4 requires a PostCSS plugin configuration, which was not set up initially, causing no styles to load in the application.

**Solution**: Created `postcss.config.js` with the `@tailwindcss/postcss` plugin and installed the required npm package.

### 9. Testing Database State

**Challenge**: Running tests would pollute the development database with test data, making it difficult to maintain clean test and development environments.

**Solution**: Implemented test helpers with transaction rollback for each test case, ensuring database state is not persisted between tests.

### 10. Rate Limiting Configuration

**Challenge**: During development, the rate limiter (100 requests per minute) would sometimes block legitimate development requests when rapidly testing features.

**Solution**: Made rate limiter configurable through environment variables and increased the limit during development. Production deployments use the stricter limit.

## Technology Stack

### Backend

- **Language**: Go 1.24
- **Web Framework**: Gin 1.11
- **Database**: PostgreSQL 16
- **ORM**: GORM 1.31
- **Authentication**: JWT (golang-jwt/jwt v5.3)
- **Password Hashing**: bcrypt
- **Environment Variables**: godotenv 1.5
- **Testing**: Go standard testing package

### Frontend

- **Framework**: Vue 3.5
- **Language**: TypeScript 5.9
- **Build Tool**: Vite 7.1
- **State Management**: Pinia 3.0
- **Routing**: Vue Router 4.6
- **HTTP Client**: Axios 1.13
- **Styling**: Tailwind CSS 4.1
- **CSS Processing**: PostCSS 8.5 with Tailwind plugin

### DevOps

- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Web Server**: Nginx (for frontend in production)
- **Database**: PostgreSQL 16 Alpine

### Development Tools

- **Version Control**: Git
- **Code Editor**: VS Code (recommended)
- **API Testing**: curl, Postman (optional)
- **Database Client**: psql, pgAdmin (optional)

## License

This project is private and proprietary. All rights reserved.

## Contact

For questions or issues, please contact the development team or create an issue in the repository.
