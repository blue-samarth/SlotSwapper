# SlotSwapper Backend API

A robust Go-based REST API for managing event swaps between users. Built with Gin framework, GORM ORM, and PostgreSQL.

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Configuration](#configuration)
- [Testing](#testing)
- [Rate Limiting](#rate-limiting)
- [Error Handling](#error-handling)
- [Health Checks](#health-checks)
- [Development](#development)

## ✨ Features

### Core Features
- 🔐 **JWT Authentication** - Secure user authentication and authorization
- 📅 **Event Management** - Create, update, and manage events with status tracking
- 🔄 **Swap System** - Request, accept, reject, and cancel event swaps
- 👤 **User Profiles** - User profile management and preferences
- 🛡️ **Rate Limiting** - Configurable rate limiting per endpoint type
- 🏥 **Health Monitoring** - Multiple health check endpoints for production monitoring
- ⚠️ **Standardized Errors** - Consistent error responses across all endpoints

### Advanced Features
- 🔒 **Concurrent Safety** - Handles concurrent swap requests with proper locking
- 🎯 **Smart Filtering** - Filter swap requests by sent/received/all
- 📊 **System Metrics** - Detailed health checks with memory, goroutines, and DB stats
- 🚦 **Authorization Checks** - Comprehensive ownership and permission validation
- 🧹 **Soft Deletes** - Preserve data integrity with GORM soft delete support

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Web Framework:** Gin v1.10.0
- **ORM:** GORM v1.25.12
- **Database:** PostgreSQL 16+ (SQLite for testing)
- **Authentication:** JWT (golang-jwt/jwt v5.2.1)
- **Validation:** go-playground/validator v10.22.1
- **Environment:** godotenv v1.5.1

## 📁 Project Structure

```
backend/
├── config/          # Database and app configuration
│   └── database.go
├── controllers/     # HTTP request handlers
│   ├── auth_controller.go
│   ├── event_controller.go
│   ├── health_controller.go
│   ├── swap_controller.go
│   └── user_controller.go
├── dto/            # Data Transfer Objects
│   ├── auth_dto.go
│   ├── event_dto.go
│   ├── swap_dto.go
│   └── user_dto.go
├── middleware/     # HTTP middleware
│   ├── auth_middleware.go
│   └── rate_limiter.go
├── models/         # Database models
│   ├── event.go
│   ├── swap_request.go
│   └── user.go
├── routes/         # Route definitions
│   └── route.go
├── services/       # Business logic
│   ├── auth_service.go
│   ├── event_service.go
│   └── swap_service.go
├── tests/          # Test suites
│   ├── auth_test.go
│   ├── errors_test.go
│   ├── event_test.go
│   ├── health_test.go
│   ├── rate_limiter_test.go
│   ├── swap_filter_test.go
│   ├── swap_test.go
│   ├── test_helper.go
│   └── user_test.go
├── utils/          # Utility functions
│   └── errors.go
├── .env.example    # Environment variables template
├── Dockerfile      # Docker configuration
├── go.mod          # Go module dependencies
├── go.sum          # Dependency checksums
└── main.go         # Application entry point
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 16 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/blue-samarth/SlotSwapper.git
   cd SlotSwapper/backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your configuration:
   ```env
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=slotswapper_db
   DB_SSLMODE=disable

   # JWT Configuration
   JWT_SECRET=your_secure_random_secret_key_minimum_32_characters

   # Application Configuration
   APP_ENV=development
   APP_VERSION=1.0.0
   PORT=8080
   ```

4. **Create the database**
   ```bash
   createdb slotswapper_db
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

   The API will be available at `http://localhost:8080`

### Using Docker

```bash
# Build the image
docker build -t slotswapper-backend .

# Run the container
docker run -p 8080:8080 --env-file .env slotswapper-backend
```

### Using Docker Compose (Recommended)

```bash
# From the root directory
docker-compose up
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080
```

### Authentication

#### Register
```http
POST /api/auth/signup
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "username"
    }
  }
}
```

### User Profile

#### Get Current User
```http
GET /api/users/me
Authorization: Bearer <token>
```

#### Update User Profile
```http
PATCH /api/users/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "newusername"
}
```

### Events

#### Create Event
```http
POST /api/events
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Meeting with Team",
  "description": "Weekly team sync",
  "start_time": "2025-11-10T10:00:00Z",
  "end_time": "2025-11-10T11:00:00Z",
  "location": "Conference Room A",
  "status": "BUSY"
}
```

#### Get User's Events
```http
GET /api/events
Authorization: Bearer <token>
```

#### Update Event
```http
PUT /api/events/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Meeting",
  "status": "SWAPPABLE"
}
```

#### Delete Event
```http
DELETE /api/events/:id
Authorization: Bearer <token>
```

### Swap Requests

#### Initiate Swap
```http
POST /api/swap-requests
Authorization: Bearer <token>
Content-Type: application/json

{
  "requester_event_id": 1,
  "receiver_event_id": 2
}
```

#### Get Swap Requests
```http
GET /api/swap-requests?filter=sent
GET /api/swap-requests?filter=received
GET /api/swap-requests?filter=all
Authorization: Bearer <token>
```

#### Accept Swap
```http
POST /api/swap-requests/:id/accept
Authorization: Bearer <token>
```

#### Reject Swap
```http
POST /api/swap-requests/:id/reject
Authorization: Bearer <token>
```

#### Cancel Swap
```http
POST /api/swap-requests/:id/cancel
Authorization: Bearer <token>
```

### Health Checks

#### Basic Health Check
```http
GET /health

Response:
{
  "status": "healthy",
  "timestamp": "2025-11-06T12:00:00Z"
}
```

#### Detailed Health Check
```http
GET /health/detailed

Response:
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": "2h30m15s",
  "database": {
    "status": "connected",
    "max_open_connections": 25,
    "open_connections": 3,
    "in_use": 1,
    "idle": 2
  },
  "system": {
    "memory_alloc": "5.2 MB",
    "total_alloc": "12.8 MB",
    "sys": "18.4 MB",
    "num_gc": 12,
    "goroutines": 8,
    "num_cpu": 8
  }
}
```

#### Kubernetes Readiness Probe
```http
GET /health/ready
```

#### Kubernetes Liveness Probe
```http
GET /health/live
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | PostgreSQL host | localhost | Yes |
| `DB_PORT` | PostgreSQL port | 5432 | Yes |
| `DB_USER` | Database user | postgres | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | slotswapper_db | Yes |
| `DB_SSLMODE` | SSL mode | disable | No |
| `JWT_SECRET` | JWT signing key (min 32 chars) | - | Yes |
| `APP_ENV` | Environment (development/production) | development | No |
| `APP_VERSION` | Application version | 1.0.0 | No |
| `PORT` | Server port | 8080 | No |

### Database Models

#### User
```go
type User struct {
    ID        uint
    Email     string    // Unique, normalized to lowercase
    Username  string    // Normalized to lowercase
    Password  string    // Bcrypt hashed
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
}
```

#### Event
```go
type Event struct {
    ID          uint
    UserID      uint
    Title       string
    Description string
    StartTime   time.Time
    EndTime     time.Time
    Location    string
    Status      string    // BUSY, SWAPPABLE, SWAP_PENDING
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt
}
```

#### SwapRequest
```go
type SwapRequest struct {
    ID               uint
    RequesterID      uint
    RequesterEventID uint
    ReceiverID       uint
    ReceiverEventID  uint
    Status           string    // PENDING, ACCEPTED, REJECTED, CANCELLED
    CreatedAt        time.Time
    UpdatedAt        time.Time
    DeletedAt        gorm.DeletedAt
}
```

## 🧪 Testing

### Run All Tests
```bash
go test ./tests/... -v
```

### Run Specific Test Suite
```bash
# Authentication tests
go test ./tests/... -v -run TestUser
go test ./tests/... -v -run TestJWT

# Event tests
go test ./tests/... -v -run TestEvent

# Swap tests
go test ./tests/... -v -run TestSwap
go test ./tests/... -v -run TestConcurrent

# Health check tests
go test ./tests/... -v -run TestHealth

# Rate limiter tests
go test ./tests/... -v -run TestRateLimiter

# Error handling tests
go test ./tests/... -v -run TestError
```

### Test Coverage
```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Statistics
- **Total Tests:** 70
- **Passing:** 69 (98.6%)
- **Skipped:** 1 (JWT expiration - requires time mocking)
- **Coverage:** High (covers all critical paths)

### Test Categories
- ✅ Authentication & JWT (11 tests)
- ✅ Event Management (13 tests)
- ✅ Swap Operations (14 tests)
- ✅ Error Handling (9 tests)
- ✅ Health Checks (4 tests)
- ✅ Rate Limiting (10 tests)
- ✅ Swap Filters (7 tests)
- ✅ User Management (6 tests)

## 🚦 Rate Limiting

Rate limiting is implemented per IP address with configurable limits:

### Rate Limit Configuration

| Endpoint Type | Requests/Minute | Purpose |
|--------------|-----------------|---------|
| **Global** | 100 | All endpoints |
| **Auth** | 5 | Signup/Login (brute force prevention) |
| **Swap Operations** | 20 | Swap create/accept/reject/cancel |

### Rate Limit Headers
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1699276800
```

### Rate Limit Error Response
```json
{
  "error": "RATE_LIMIT_EXCEEDED",
  "message": "Too many requests. Please try again later.",
  "status_code": 429,
  "timestamp": "2025-11-06T12:00:00Z"
}
```

### Cleanup
- Stale visitor records are automatically cleaned up every 5 minutes
- No manual maintenance required

## ⚠️ Error Handling

All API errors follow a standardized format:

### Error Response Structure
```json
{
  "error": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": {
    "field": "additional context"
  },
  "status_code": 400,
  "timestamp": "2025-11-06T12:00:00Z"
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Invalid input data |
| `BAD_REQUEST` | 400 | Malformed request |
| `AUTHENTICATION_ERROR` | 401 | Invalid or missing token |
| `AUTHORIZATION_ERROR` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource conflict (e.g., duplicate email) |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `SWAP_LOGIC_ERROR` | 422 | Swap business logic violation |
| `INTERNAL_SERVER_ERROR` | 500 | Server error |

### Success Response Structure
```json
{
  "message": "Operation successful",
  "data": {
    "id": 1,
    "field": "value"
  }
}
```

## 🏥 Health Checks

### Endpoints

#### 1. Basic Health Check
```http
GET /health
```
- Returns: Service status and timestamp
- Use: Simple uptime monitoring

#### 2. Detailed Health Check
```http
GET /health/detailed
```
- Returns: Service status, version, uptime, database stats, system metrics
- Use: Comprehensive monitoring dashboards

#### 3. Readiness Check
```http
GET /health/ready
```
- Returns: Database connectivity status
- Use: Kubernetes readiness probe

#### 4. Liveness Check
```http
GET /health/live
```
- Returns: Service alive status
- Use: Kubernetes liveness probe

### Kubernetes Configuration Example
```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

## 💻 Development

### Code Organization

- **Controllers:** Handle HTTP requests and responses
- **Services:** Contain business logic
- **Models:** Define database schemas
- **DTOs:** Define request/response structures
- **Middleware:** Process requests before controllers
- **Utils:** Shared utility functions

### Best Practices

1. **Authentication:** All protected routes require JWT token
2. **Authorization:** Endpoints verify ownership before operations
3. **Validation:** Input validation using struct tags
4. **Error Handling:** Always use standardized error responses
5. **Testing:** Write tests for new features
6. **Database:** Use transactions for multi-step operations

### Adding New Features

1. **Define Model** (if needed) in `models/`
2. **Create DTO** in `dto/`
3. **Implement Service** in `services/`
4. **Create Controller** in `controllers/`
5. **Add Routes** in `routes/route.go`
6. **Write Tests** in `tests/`

### Database Migrations

GORM AutoMigrate runs on startup:
```go
config.DB.AutoMigrate(&models.User{}, &models.Event{}, &models.SwapRequest{})
```

For production, consider using a migration tool like [golang-migrate](https://github.com/golang-migrate/migrate).

### Code Style

Follow standard Go conventions:
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Vet code
go vet ./...
```

## 🐛 Debugging

### Enable Debug Mode
```go
// In main.go
if os.Getenv("APP_ENV") == "development" {
    gin.SetMode(gin.DebugMode)
}
```

### Database Query Logging
```go
// In config/database.go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

### Common Issues

1. **JWT Secret Too Short**
   - Error: "JWT_SECRET must be at least 32 characters"
   - Solution: Use a longer secret in `.env`

2. **Database Connection Failed**
   - Check PostgreSQL is running
   - Verify credentials in `.env`
   - Ensure database exists

3. **Port Already in Use**
   - Change `PORT` in `.env`
   - Or kill process: `lsof -ti:8080 | xargs kill`

## 📝 License

MIT License - see LICENSE file for details

## 👥 Contributors

- Samarth - Initial development

## 🔗 Links

- Repository: https://github.com/blue-samarth/SlotSwapper
- Issues: https://github.com/blue-samarth/SlotSwapper/issues

## 📞 Support

For questions or issues, please open an issue on GitHub.

---

**Status:** ✅ Production-ready | **Tests:** 69/70 passing | **Coverage:** High
