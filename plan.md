# Complete Technical Implementation Plan: Go Gin + Vue Stack

## **EXECUTIVE SUMMARY**

**Total Estimated Time:** 24 hours  
**Architecture:** Monorepo with separate backend/frontend  
**Database:** PostgreSQL (can use SQLite for faster setup)  
**Authentication:** JWT with Bearer tokens  
**Key Challenge:** Transactional swap logic with race condition handling  

---

## **PART 1: BACKEND ARCHITECTURE (Go + Gin)**

### **1.1 Technology Stack & Dependencies**

#### Core Dependencies
```
github.com/gin-gonic/gin v1.9.1          // Web framework
gorm.io/gorm v1.25.5                     // ORM
gorm.io/driver/postgres v1.5.4           // PostgreSQL driver
github.com/golang-jwt/jwt/v5 v5.2.0      // JWT authentication
golang.org/x/crypto v0.17.0              // Password hashing (bcrypt)
github.com/gin-contrib/cors v1.5.0       // CORS middleware
github.com/joho/godotenv v1.5.1          // Environment variables
```

#### Optional but Recommended
```
github.com/go-playground/validator/v10   // Request validation
github.com/stretchr/testify v1.8.4       // Testing
github.com/rs/zerolog v1.31.0            // Structured logging
```

---

### **1.2 Project Directory Structure**

```
backend/
├── main.go                          // Application entry point
├── .env                             // Environment variables
├── go.mod                           // Go modules
├── go.sum                           // Dependency checksums
│
├── config/
│   ├── database.go                  // DB connection & migration
│   └── env.go                       // Environment config loader
│
├── models/
│   ├── user.go                      // User model + validation
│   ├── event.go                     // Event model + status enum
│   ├── swap_request.go              // SwapRequest model + status enum
│   └── base.go                      // Shared model fields (optional)
│
├── controllers/
│   ├── auth_controller.go           // Signup, Login
│   ├── event_controller.go          // CRUD for events
│   ├── swap_controller.go           // Swap operations
│   └── user_controller.go           // Get current user info (optional)
│
├── services/
│   ├── swap_service.go              // Core swap transaction logic
│   ├── event_service.go             // Event business logic (optional)
│   └── auth_service.go              // Auth business logic (optional)
│
├── middleware/
│   ├── auth.go                      // JWT validation middleware
│   ├── error_handler.go             // Global error handler
│   └── logger.go                    // Request logging (optional)
│
├── utils/
│   ├── jwt.go                       // Token generation & validation
│   ├── response.go                  // Standardized API responses
│   └── validator.go                 // Custom validation rules
│
├── routes/
│   └── routes.go                    // Route definitions
│
├── dto/                             // Data Transfer Objects
│   ├── auth_dto.go                  // Login/Signup request/response
│   ├── event_dto.go                 // Event request/response
│   └── swap_dto.go                  // Swap request/response
│
└── tests/
    ├── auth_test.go
    ├── event_test.go
    └── swap_test.go
```

---

### **1.3 Database Schema Design**

#### **Table: users**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,  -- bcrypt hashed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster email lookups
CREATE INDEX idx_users_email ON users(email);
```

#### **Table: events**
```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(20) DEFAULT 'BUSY' CHECK (status IN ('BUSY', 'SWAPPABLE', 'SWAP_PENDING')),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_start_time ON events(start_time);

-- Composite index for marketplace query
CREATE INDEX idx_events_status_user ON events(status, user_id);
```

#### **Table: swap_requests**
```sql
CREATE TABLE swap_requests (
    id SERIAL PRIMARY KEY,
    requester_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    my_slot_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    their_slot_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for queries
CREATE INDEX idx_swap_requests_requester ON swap_requests(requester_id);
CREATE INDEX idx_swap_requests_receiver ON swap_requests(receiver_id);
CREATE INDEX idx_swap_requests_status ON swap_requests(status);
```

---

### **1.4 GORM Models Implementation Details**

#### **User Model (`models/user.go`)**
```
Fields:
- ID: uint (primary key, auto-increment)
- Name: string (not null, max 255 chars)
- Email: string (unique, not null, lowercase before save)
- Password: string (not null, never returned in JSON using json:"-" tag)
- CreatedAt: time.Time (auto-managed by GORM)
- UpdatedAt: time.Time (auto-managed by GORM)

Methods:
- BeforeSave(tx *gorm.DB) error: Hook to lowercase email
- HashPassword(password string) error: Bcrypt with cost 12
- CheckPassword(password string) bool: Compare hashed password

Validation Tags:
- Name: binding:"required,min=2,max=255"
- Email: binding:"required,email"
- Password: binding:"required,min=6"
```

#### **Event Model (`models/event.go`)**
```
Enums:
- EventStatus: BUSY | SWAPPABLE | SWAP_PENDING

Fields:
- ID: uint (primary key)
- Title: string (not null, max 255)
- StartTime: time.Time (not null, must be future during creation)
- EndTime: time.Time (not null, must be after StartTime)
- Status: EventStatus (default: BUSY)
- UserID: uint (foreign key, indexed, not null)
- User: User (lazy-loaded relationship)
- CreatedAt: time.Time
- UpdatedAt: time.Time

Validations:
- Title: binding:"required,min=1,max=255"
- StartTime: binding:"required"
- EndTime: binding:"required,gtfield=StartTime"

Business Rules:
- Can only change BUSY → SWAPPABLE (user action)
- System changes SWAPPABLE → SWAP_PENDING (on swap request)
- System changes SWAP_PENDING → BUSY/SWAPPABLE (on accept/reject)
```

#### **SwapRequest Model (`models/swap_request.go`)**
```
Enums:
- SwapStatus: PENDING | ACCEPTED | REJECTED

Fields:
- ID: uint (primary key)
- RequesterID: uint (foreign key to users, indexed)
- Requester: User (relationship)
- ReceiverID: uint (foreign key to users, indexed)
- Receiver: User (relationship)
- MySlotID: uint (requester's slot, foreign key to events)
- MySlot: Event (relationship with Preload)
- TheirSlotID: uint (receiver's slot, foreign key to events)
- TheirSlot: Event (relationship with Preload)
- Status: SwapStatus (default: PENDING)
- CreatedAt: time.Time
- UpdatedAt: time.Time

Constraints:
- RequesterID != ReceiverID (check in controller)
- MySlotID != TheirSlotID (check in controller)
- Both slots must exist and be SWAPPABLE at creation time
```

---

### **1.5 Authentication Implementation**

#### **JWT Structure**
```
Header:
{
  "alg": "HS256",
  "typ": "JWT"
}

Payload (Claims):
{
  "user_id": 123,
  "email": "user@example.com",
  "exp": 1735257600,      // 24 hours from issue
  "iat": 1735171200       // issued at
}

Signature:
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  SECRET_KEY
)
```

#### **Token Flow**
```
1. User sends POST /api/auth/signup or /api/auth/login
2. Server validates credentials
3. Server generates JWT with user_id and email
4. Server returns: { "token": "...", "user": {...} }
5. Client stores token in localStorage
6. Client sends token in header: "Authorization: Bearer <token>"
7. Middleware validates token on protected routes
8. Middleware extracts user_id and sets in gin.Context
```

#### **Password Security**
```
Algorithm: bcrypt
Cost: 12 (good balance of security and performance)
Salt: Automatically handled by bcrypt

Process:
1. User submits password (min 6 chars)
2. Server validates strength (optional: regex for complexity)
3. Server hashes: bcrypt.GenerateFromPassword([]byte(password), 12)
4. Server stores hash in database (never plaintext)
5. Login: bcrypt.CompareHashAndPassword(stored, submitted)
```

#### **Middleware Logic**
```
Gin Middleware Execution Flow:
1. Extract "Authorization" header
2. Validate format: "Bearer <token>"
3. Parse and validate JWT:
   - Check signature
   - Check expiration
   - Check claims structure
4. If valid:
   - Extract user_id and email from claims
   - Set in context: c.Set("user_id", claims.UserID)
   - Call c.Next() to continue to handler
5. If invalid:
   - Return 401 Unauthorized
   - Call c.Abort() to stop execution
```

---

### **1.6 API Endpoint Specifications**

#### **Authentication Endpoints**

**POST /api/auth/signup**
```
Request Body:
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepass123"
}

Success Response (201):
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-01-15T10:30:00Z"
  }
}

Error Response (400):
{
  "error": "Email already exists"
}

Validations:
- Email must be unique
- Password min 6 characters
- Name required
```

**POST /api/auth/login**
```
Request Body:
{
  "email": "john@example.com",
  "password": "securepass123"
}

Success Response (200):
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}

Error Response (401):
{
  "error": "Invalid credentials"
}
```

---

#### **Event Management Endpoints**

**GET /api/events**
```
Headers:
Authorization: Bearer <token>

Success Response (200):
[
  {
    "id": 1,
    "title": "Team Meeting",
    "start_time": "2025-01-20T10:00:00Z",
    "end_time": "2025-01-20T11:00:00Z",
    "status": "SWAPPABLE",
    "user_id": 1,
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
]

Notes:
- Returns only events belonging to authenticated user
- Ordered by start_time ASC
- Includes all statuses (BUSY, SWAPPABLE, SWAP_PENDING)
```

**POST /api/events**
```
Headers:
Authorization: Bearer <token>

Request Body:
{
  "title": "Focus Block",
  "start_time": "2025-01-21T14:00:00Z",
  "end_time": "2025-01-21T15:00:00Z"
}

Success Response (201):
{
  "id": 2,
  "title": "Focus Block",
  "start_time": "2025-01-21T14:00:00Z",
  "end_time": "2025-01-21T15:00:00Z",
  "status": "BUSY",
  "user_id": 1,
  "created_at": "2025-01-15T11:00:00Z",
  "updated_at": "2025-01-15T11:00:00Z"
}

Validations:
- end_time must be after start_time
- title required (1-255 chars)
- Automatically sets user_id from JWT
- Default status: BUSY
```

**PATCH /api/events/:id/status**
```
Headers:
Authorization: Bearer <token>

Request Body:
{
  "status": "SWAPPABLE"
}

Success Response (200):
{
  "id": 1,
  "title": "Team Meeting",
  "status": "SWAPPABLE",
  ...
}

Error Responses:
404: Event not found
403: Not authorized (not your event)
400: Invalid status transition

Allowed Transitions:
- BUSY → SWAPPABLE (user can mark as swappable)
- SWAPPABLE → BUSY (user can unmark)
- SWAP_PENDING → (read-only, only system can change)
```

**DELETE /api/events/:id** (Optional)
```
Headers:
Authorization: Bearer <token>

Success Response (204): No content

Error Responses:
404: Event not found
403: Not authorized
400: Cannot delete event in SWAP_PENDING status
```

---

#### **Swap Marketplace Endpoints**

**GET /api/swappable-slots**
```
Headers:
Authorization: Bearer <token>

Success Response (200):
[
  {
    "id": 5,
    "title": "Client Call",
    "start_time": "2025-01-22T09:00:00Z",
    "end_time": "2025-01-22T10:00:00Z",
    "status": "SWAPPABLE",
    "user_id": 3,
    "user": {
      "id": 3,
      "name": "Jane Smith",
      "email": "jane@example.com"
    }
  }
]

Query Logic (SQL):
SELECT * FROM events 
WHERE status = 'SWAPPABLE' 
  AND user_id != <current_user_id>
ORDER BY start_time ASC;

Performance Optimization:
- Use composite index on (status, user_id)
- Preload User data to avoid N+1 queries
- Consider pagination for large datasets
```

**POST /api/swap-request**
```
Headers:
Authorization: Bearer <token>

Request Body:
{
  "my_slot_id": 2,
  "their_slot_id": 5
}

Success Response (201):
{
  "id": 1,
  "requester_id": 1,
  "receiver_id": 3,
  "my_slot_id": 2,
  "their_slot_id": 5,
  "status": "PENDING",
  "created_at": "2025-01-15T12:00:00Z"
}

Error Responses:
404: Slot not found
400: Slot not swappable
403: Not your slot (my_slot_id)
409: Slot already in pending swap

Atomic Transaction Steps:
1. BEGIN TRANSACTION
2. SELECT my_slot FOR UPDATE (row lock)
3. SELECT their_slot FOR UPDATE (row lock)
4. Validate both are SWAPPABLE
5. INSERT INTO swap_requests
6. UPDATE my_slot SET status = 'SWAP_PENDING'
7. UPDATE their_slot SET status = 'SWAP_PENDING'
8. COMMIT
9. If any step fails: ROLLBACK
```

---

#### **Swap Response Endpoints**

**POST /api/swap-response/:id**
```
Headers:
Authorization: Bearer <token>

URL Params:
:id = swap_request_id

Request Body:
{
  "accept": true  // or false
}

Success Response (200):
{
  "message": "Swap accepted successfully",
  "swap_request": {
    "id": 1,
    "status": "ACCEPTED",
    ...
  }
}

Error Responses:
404: Swap request not found
403: Not authorized (not the receiver)
400: Request already processed
400: Slots no longer in valid state

ACCEPT Transaction Logic:
1. BEGIN TRANSACTION
2. SELECT swap_request FOR UPDATE
3. Verify receiver_id matches current user
4. Verify status is PENDING
5. SELECT both slots FOR UPDATE
6. Verify both slots are SWAP_PENDING
7. SWAP OWNERS:
   - temp = my_slot.user_id
   - my_slot.user_id = their_slot.user_id
   - their_slot.user_id = temp
8. UPDATE my_slot SET status = 'BUSY', user_id = <new_owner>
9. UPDATE their_slot SET status = 'BUSY', user_id = <new_owner>
10. UPDATE swap_request SET status = 'ACCEPTED'
11. COMMIT

REJECT Transaction Logic:
1. BEGIN TRANSACTION
2. SELECT swap_request FOR UPDATE
3. Verify receiver_id matches current user
4. Verify status is PENDING
5. SELECT both slots FOR UPDATE
6. UPDATE my_slot SET status = 'SWAPPABLE'
7. UPDATE their_slot SET status = 'SWAPPABLE'
8. UPDATE swap_request SET status = 'REJECTED'
9. COMMIT
```

**GET /api/swap-requests/incoming**
```
Headers:
Authorization: Bearer <token>

Success Response (200):
[
  {
    "id": 1,
    "requester": {
      "id": 2,
      "name": "Alice Brown",
      "email": "alice@example.com"
    },
    "receiver_id": 1,
    "my_slot": {
      "id": 5,
      "title": "Client Call",
      "start_time": "2025-01-22T09:00:00Z",
      "end_time": "2025-01-22T10:00:00Z"
    },
    "their_slot": {
      "id": 3,
      "title": "Focus Time",
      "start_time": "2025-01-23T14:00:00Z",
      "end_time": "2025-01-23T15:00:00Z"
    },
    "status": "PENDING",
    "created_at": "2025-01-15T12:00:00Z"
  }
]

Query Logic:
SELECT * FROM swap_requests 
WHERE receiver_id = <current_user_id>
ORDER BY created_at DESC;

Preload:
- Requester (user)
- MySlot (event)
- TheirSlot (event)
```

**GET /api/swap-requests/outgoing**
```
Headers:
Authorization: Bearer <token>

Success Response (200):
[
  {
    "id": 2,
    "requester_id": 1,
    "receiver": {
      "id": 3,
      "name": "Bob Johnson",
      "email": "bob@example.com"
    },
    "my_slot": { ... },
    "their_slot": { ... },
    "status": "PENDING",
    "created_at": "2025-01-15T13:00:00Z"
  }
]

Query Logic:
SELECT * FROM swap_requests 
WHERE requester_id = <current_user_id>
ORDER BY created_at DESC;
```

---

### **1.7 Critical: Swap Transaction Logic Deep Dive**

#### **Race Condition Scenarios & Solutions**

**Scenario 1: Double Swap Request**
```
Problem:
- User A marks slot X as SWAPPABLE
- User B requests to swap with slot X
- User C also requests to swap with slot X (before B's request updates status)
- Both requests get created, but only one should succeed

Solution:
- Use database transactions with row-level locking
- SELECT ... FOR UPDATE on slot rows
- First transaction locks the row, second waits
- Second transaction will fail validation (status != SWAPPABLE)

Implementation:
tx := db.Begin()
var slot Event
tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&slot, slotID)
// Now slot is locked until COMMIT
```

**Scenario 2: Concurrent Acceptance**
```
Problem:
- Request is PENDING
- Both requester and receiver somehow try to accept
- (Shouldn't happen in UI, but API should prevent it)

Solution:
- Only receiver can accept (verify receiver_id in transaction)
- Use SELECT ... FOR UPDATE on swap_request
- Verify status is still PENDING before proceeding

Implementation:
tx := db.Begin()
var swap SwapRequest
tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&swap, requestID)
if swap.ReceiverID != currentUserID {
    return ErrUnauthorized
}
```

**Scenario 3: Slot Deleted During Swap**
```
Problem:
- Swap request is PENDING
- User deletes one of the slots involved

Solution Option 1 (Recommended):
- Don't allow deletion of slots in SWAP_PENDING status
- Return 400 error if user tries

Solution Option 2:
- Use ON DELETE CASCADE in foreign keys
- Automatically delete swap_request if either slot is deleted
- May confuse users, not recommended
```

**Scenario 4: User Accepts Multiple Swaps for Same Slot**
```
Problem:
- User has slot X in SWAP_PENDING with User A
- User also creates request for slot X with User B
- Should not be possible

Prevention:
- When creating swap request, verify my_slot is SWAPPABLE (not SWAP_PENDING)
- Status SWAP_PENDING prevents slot from being in other requests
```

#### **Transaction Isolation Level**

```
Recommended: READ COMMITTED (PostgreSQL default)

Why not SERIALIZABLE?
- More expensive performance-wise
- Not needed if we use SELECT ... FOR UPDATE properly
- READ COMMITTED + row locks = safe for our use case

GORM Transaction:
tx := db.Begin(&sql.TxOptions{
    Isolation: sql.LevelReadCommitted,
})

Rollback Strategy:
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// Do work
if err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

#### **Error Handling Matrix**

```
| Error Type          | HTTP Code | User Message                    | Action            |
|---------------------|-----------|----------------------------------|-------------------|
| Slot not found      | 404       | "Slot not found"                | Return error      |
| Not slot owner      | 403       | "Unauthorized"                  | Return error      |
| Slot not swappable  | 400       | "Slot is not available"         | Return error      |
| Request not pending | 400       | "Request already processed"     | Return error      |
| DB transaction fail | 500       | "Transaction failed"            | Rollback + log    |
| Validation error    | 400       | Specific field error            | Return error      |
| Unauthorized JWT    | 401       | "Invalid or expired token"      | Reject request    |
```

---

### **1.8 Service Layer Architecture**

#### **Separation of Concerns**

```
Controller Layer:
- HTTP request/response handling
- Input validation (using binding tags)
- Authentication verification
- Call service layer
- Format responses

Service Layer:
- Business logic
- Database transactions
- Complex validations
- Orchestrate multiple operations
- Return domain errors

Repository Layer (Optional):
- Direct database access
- Query building
- Can skip if using GORM directly in services
```

#### **Swap Service Detailed Design**

```go
// services/swap_service.go

type SwapService struct {
    db *gorm.DB
}

func NewSwapService(db *gorm.DB) *SwapService {
    return &SwapService{db: db}
}

// CreateSwapRequest handles the entire swap request flow
func (s *SwapService) CreateSwapRequest(requesterID uint, mySlotID uint, theirSlotID uint) (*SwapRequest, error) {
    // Start transaction
    // Lock and validate my slot
    // Lock and validate their slot
    // Create swap request record
    // Update both slot statuses to SWAP_PENDING
    // Commit or rollback
}

// AcceptSwapRequest handles the acceptance flow
func (s *SwapService) AcceptSwapRequest(requestID uint, receiverID uint) error {
    // Start transaction
    // Lock swap request
    // Verify receiver
    // Lock both slots
    // Swap owners
    // Update statuses to BUSY
    // Update request status to ACCEPTED
    // Commit or rollback
}

// RejectSwapRequest handles the rejection flow
func (s *SwapService) RejectSwapRequest(requestID uint, receiverID uint) error {
    // Start transaction
    // Lock swap request
    // Verify receiver
    // Lock both slots
    // Update statuses back to SWAPPABLE
    // Update request status to REJECTED
    // Commit or rollback
}

// ValidateSlotForSwap checks business rules
func (s *SwapService) ValidateSlotForSwap(slot *Event, ownerID uint) error {
    // Check ownership
    // Check status is SWAPPABLE
    // Optional: Check time not in past
    // Optional: Check no time conflicts
}
```

---

### **1.9 Error Handling Strategy**

#### **Custom Error Types**

```go
// utils/errors.go

type AppError struct {
    Code    int    // HTTP status code
    Message string // User-facing message
    Err     error  // Internal error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

// Predefined errors
var (
    ErrNotFound         = &AppError{Code: 404, Message: "Resource not found"}
    ErrUnauthorized     = &AppError{Code: 401, Message: "Unauthorized"}
    ErrForbidden        = &AppError{Code: 403, Message: "Forbidden"}
    ErrBadRequest       = &AppError{Code: 400, Message: "Bad request"}
    ErrInternalServer   = &AppError{Code: 500, Message: "Internal server error"}
    ErrSlotNotSwappable = &AppError{Code: 400, Message: "Slot is not swappable"}
    ErrRequestProcessed = &AppError{Code: 400, Message: "Request already processed"}
)
```

#### **Global Error Handler Middleware**

```go
// middleware/error_handler.go

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // Check if there were any errors
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err

            switch e := err.(type) {
            case *AppError:
                c.JSON(e.Code, gin.H{"error": e.Message})
            case validator.ValidationErrors:
                c.JSON(400, gin.H{"error": formatValidationErrors(e)})
            default:
                c.JSON(500, gin.H{"error": "Internal server error"})
            }
        }
    }
}
```

---

### **1.10 Logging Strategy**

#### **Structured Logging**

```go
// Use zerolog for structured, performant logging

import "github.com/rs/zerolog/log"

// In service layer:
log.Info().
    Uint("user_id", userID).
    Uint("slot_id", slotID).
    Msg("Creating swap request")

// Error logging:
log.Error().
    Err(err).
    Uint("request_id", requestID).
    Msg("Failed to accept swap")

// Performance logging:
start := time.Now()
// ... operation ...
log.Debug().
    Dur("duration", time.Since(start)).
    Str("operation", "swap_accept").
    Msg("Operation completed")
```

#### **Request Logging Middleware**

```go
// middleware/logger.go

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()
        
        duration := time.Since(start)
        
        log.Info().
            Str("method", c.Request.Method).
            Str("path", path).
            Int("status", c.Writer.Status()).
            Dur("duration", duration).
            Str("client_ip", c.ClientIP()).
            Msg("Request handled")
    }
}
```

---

### **1.11 Testing Strategy**

#### **Unit Tests**

```
Focus Areas:
1. Swap service transaction logic
2. JWT generation and validation
3. Password hashing and verification
4. Model validation

Tools:
- testing (Go standard library)
- testify/assert (assertions)
- testify/mock (mocking)

Example Test Structure:
func TestAcceptSwapRequest(t *testing.T) {
    // Setup: Create test database, seed data
    // Execute: Call service method
    // Assert: Verify expected outcomes
    // Cleanup: Rollback or delete test data
}

Test Cases for Swap Acceptance:
- Happy path: successful swap
- Error: request not found
- Error: unauthorized user
- Error: request already processed
- Error: slots in invalid state
- Concurrency: multiple simultaneous accepts
```

#### **Integration Tests**

```
Focus Areas:
1. Full API endpoint flows
2. Database transactions
3. Authentication middleware
4. CORS configuration

Tools:
- httptest (Go standard library)
- Test containers (optional, for isolated DB)

Example:
func TestSwapRequestEndToEnd(t *testing.T) {
    // 1. Create test server
    router := SetupTestRouter()
    
    // 2. Create two users
    user1 := createTestUser(t, router)
    user2 := createTestUser(t, router)
    
    // 3. Create events for both
    slot1 := createTestEvent(t, router, user1.Token)
    slot2 := createTestEvent(t, router, user2.Token)
    
    // 4. Mark both as swappable
    updateSlotStatus(t, router, user1.Token, slot1.ID, "SWAPPABLE")
    updateSlotStatus(t, router, user2.Token, slot2.ID, "SWAPPABLE")
    
```
    // 5. User1 requests swap
    swapReq := createSwapRequest(t, router, user1.Token, slot1.ID, slot2.ID)
    
    // 6. Verify slots are now SWAP_PENDING
    assertSlotStatus(t, router, user1.Token, slot1.ID, "SWAP_PENDING")
    assertSlotStatus(t, router, user2.Token, slot2.ID, "SWAP_PENDING")
    
    // 7. User2 accepts swap
    acceptSwap(t, router, user2.Token, swapReq.ID)
    
    // 8. Verify ownership swapped
    events1 := getMyEvents(t, router, user1.Token)
    events2 := getMyEvents(t, router, user2.Token)
    
    assert.Contains(t, events2, slot1) // User2 now owns slot1
    assert.Contains(t, events1, slot2) // User1 now owns slot2
    
    // 9. Verify statuses are BUSY
    assertSlotStatus(t, router, user1.Token, slot2.ID, "BUSY")
    assertSlotStatus(t, router, user2.Token, slot1.ID, "BUSY")
}
```

---

### **1.12 Configuration Management**

#### **.env File Structure**

```bash
# Server Configuration
PORT=8080
GIN_MODE=release  # or debug for development
ENVIRONMENT=development  # development, staging, production

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=slotswapper
DB_SSLMODE=disable  # or require for production

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-in-production-min-32-chars
JWT_EXPIRATION_HOURS=24

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Logging
LOG_LEVEL=info  # debug, info, warn, error

# Optional: Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

#### **Config Loader (`config/env.go`)**

```go
package config

import (
    "log"
    "os"
    "strconv"
    "strings"
    "github.com/joho/godotenv"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
    CORS     CORSConfig
}

type ServerConfig struct {
    Port    string
    GinMode string
    Env     string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type JWTConfig struct {
    Secret          string
    ExpirationHours int
}

type CORSConfig struct {
    AllowedOrigins []string
}

var AppConfig *Config

func LoadConfig() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    AppConfig = &Config{
        Server: ServerConfig{
            Port:    getEnv("PORT", "8080"),
            GinMode: getEnv("GIN_MODE", "debug"),
            Env:     getEnv("ENVIRONMENT", "development"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "postgres"),
            DBName:   getEnv("DB_NAME", "slotswapper"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        JWT: JWTConfig{
            Secret:          getEnv("JWT_SECRET", "default-secret-change-me"),
            ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
        },
        CORS: CORSConfig{
            AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:5173"), ","),
        },
    }

    validateConfig()
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}

func validateConfig() {
    if AppConfig.Server.Env == "production" {
        if AppConfig.JWT.Secret == "default-secret-change-me" {
            log.Fatal("JWT_SECRET must be set in production")
        }
        if len(AppConfig.JWT.Secret) < 32 {
            log.Fatal("JWT_SECRET must be at least 32 characters")
        }
    }
}

// Database connection string builder
func (c *DatabaseConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
    )
}
```

---

### **1.13 CORS Configuration**

#### **CORS Middleware Setup**

```go
// In routes/routes.go

import "github.com/gin-contrib/cors"

func SetupRoutes() *gin.Engine {
    r := gin.Default()

    // CORS configuration
    corsConfig := cors.Config{
        AllowOrigins:     config.AppConfig.CORS.AllowedOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }

    r.Use(cors.New(corsConfig))

    // Routes...
    return r
}
```

#### **Preflight Request Handling**

```
Browser automatically sends OPTIONS request for:
- Cross-origin requests
- Custom headers (like Authorization)
- Methods other than GET/POST

Gin's CORS middleware handles this automatically
Returns 204 No Content with appropriate headers:
- Access-Control-Allow-Origin
- Access-Control-Allow-Methods
- Access-Control-Allow-Headers
- Access-Control-Max-Age
```

---

### **1.14 Database Connection & Migration**

#### **Connection Pool Configuration**

```go
// config/database.go

func ConnectDatabase() {
    dsn := AppConfig.Database.DSN()
    
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC() // Always use UTC
        },
        PrepareStmt: true, // Prepared statement cache
    })
    
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to database")
    }

    // Get underlying SQL DB for connection pool settings
    sqlDB, err := database.DB()
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to get SQL DB")
    }

    // Connection pool configuration
    sqlDB.SetMaxIdleConns(10)           // Max idle connections
    sqlDB.SetMaxOpenConns(100)          // Max open connections
    sqlDB.SetConnMaxLifetime(time.Hour) // Connection lifetime

    DB = database
    log.Info().Msg("Database connected successfully")
}
```

#### **Auto-Migration Strategy**

```go
func MigrateDatabase() {
    log.Info().Msg("Running database migrations...")
    
    err := DB.AutoMigrate(
        &models.User{},
        &models.Event{},
        &models.SwapRequest{},
    )
    
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to run migrations")
    }

    log.Info().Msg("Database migrations completed")
}

// Alternative: Manual migrations for production
// Use tools like golang-migrate or GORM migrator
```

#### **Seeding (Development Only)**

```go
func SeedDatabase() {
    if AppConfig.Server.Env != "development" {
        return
    }

    log.Info().Msg("Seeding database with test data...")

    // Create test users
    users := []models.User{
        {Name: "Alice Johnson", Email: "alice@example.com", Password: hashPassword("password123")},
        {Name: "Bob Smith", Email: "bob@example.com", Password: hashPassword("password123")},
        {Name: "Charlie Brown", Email: "charlie@example.com", Password: hashPassword("password123")},
    }

    for _, user := range users {
        DB.FirstOrCreate(&user, models.User{Email: user.Email})
    }

    // Create test events
    now := time.Now()
    events := []models.Event{
        {
            Title:     "Team Meeting",
            StartTime: now.Add(24 * time.Hour),
            EndTime:   now.Add(25 * time.Hour),
            Status:    models.StatusSwappable,
            UserID:    1,
        },
        {
            Title:     "Focus Block",
            StartTime: now.Add(48 * time.Hour),
            EndTime:   now.Add(49 * time.Hour),
            Status:    models.StatusSwappable,
            UserID:    2,
        },
    }

    for _, event := range events {
        DB.FirstOrCreate(&event, models.Event{
            Title:  event.Title,
            UserID: event.UserID,
        })
    }

    log.Info().Msg("Database seeding completed")
}
```

---

### **1.15 Performance Optimization**

#### **Database Query Optimization**

```go
// BAD: N+1 Query Problem
func GetSwappableSlots() []Event {
    var slots []Event
    DB.Where("status = ?", StatusSwappable).Find(&slots)
    
    // Each slot access triggers another query for user
    for _, slot := range slots {
        fmt.Println(slot.User.Name) // SELECT * FROM users WHERE id = ?
    }
    return slots
}

// GOOD: Eager Loading with Preload
func GetSwappableSlots() []Event {
    var slots []Event
    DB.Preload("User").  // Single JOIN query
        Where("status = ?", StatusSwappable).
        Find(&slots)
    return slots
}

// BETTER: Select only needed fields
func GetSwappableSlots() []Event {
    var slots []Event
    DB.Preload("User", func(db *gorm.DB) *gorm.DB {
        return db.Select("id", "name", "email") // Don't load password
    }).
    Select("id", "title", "start_time", "end_time", "status", "user_id").
    Where("status = ?", StatusSwappable).
    Find(&slots)
    return slots
}
```

#### **Indexing Strategy**

```
Essential Indexes (created in schema):
1. users.email (UNIQUE)
2. events.user_id
3. events.status
4. events (status, user_id) - Composite for marketplace query
5. swap_requests.requester_id
6. swap_requests.receiver_id

Optional Indexes for scale:
7. events.start_time - For time-based queries
8. swap_requests (receiver_id, status) - For filtering incoming pending
```

#### **Pagination Implementation**

```go
// For marketplace with many slots
func GetSwappableSlots(page int, limit int) ([]Event, int64) {
    var slots []Event
    var total int64
    
    offset := (page - 1) * limit
    
    query := DB.Where("status = ? AND user_id != ?", StatusSwappable, currentUserID)
    
    // Count total before pagination
    query.Model(&Event{}).Count(&total)
    
    // Apply pagination
    query.Preload("User").
        Offset(offset).
        Limit(limit).
        Order("start_time ASC").
        Find(&slots)
    
    return slots, total
}

// Response format
{
    "data": [...],
    "pagination": {
        "page": 1,
        "limit": 20,
        "total": 150,
        "pages": 8
    }
}
```

#### **Caching Strategy (Advanced)**

```
For high-traffic scenarios:

1. Cache swappable slots list (Redis)
   - TTL: 30 seconds
   - Invalidate on: new swappable slot, swap request created
   - Key: "swappable_slots:user:{user_id}"

2. Cache user profile
   - TTL: 5 minutes
   - Invalidate on: profile update
   - Key: "user:{user_id}"

Implementation:
- Use go-redis/redis
- Cache-aside pattern
- Consider cache warming for common queries

Note: Skip caching for MVP, add if performance issues arise
```

---

### **1.16 Security Considerations**

#### **Input Validation**

```go
// Use struct tags for automatic validation
type CreateEventRequest struct {
    Title     string `json:"title" binding:"required,min=1,max=255"`
    StartTime string `json:"start_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
    EndTime   string `json:"end_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00,gtfield=StartTime"`
}

// Custom validator for time logic
func validateEventTimes(sl validator.StructLevel) {
    event := sl.Current().Interface().(CreateEventRequest)
    
    start, _ := time.Parse(time.RFC3339, event.StartTime)
    end, _ := time.Parse(time.RFC3339, event.EndTime)
    
    if end.Before(start) {
        sl.ReportError(event.EndTime, "end_time", "EndTime", "must be after start_time", "")
    }
    
    if start.Before(time.Now()) {
        sl.ReportError(event.StartTime, "start_time", "StartTime", "must be in future", "")
    }
}
```

#### **SQL Injection Prevention**

```
GORM automatically prevents SQL injection:
- Uses prepared statements
- Parameterizes all queries
- Escapes special characters

BAD (raw SQL - avoid):
db.Raw("SELECT * FROM users WHERE email = '" + email + "'")

GOOD (parameterized):
db.Where("email = ?", email).Find(&user)

GOOD (named parameters):
db.Where("email = @email", sql.Named("email", email)).Find(&user)
```

#### **Rate Limiting**

```go
// Use github.com/ulule/limiter/v3

import (
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimitMiddleware() gin.HandlerFunc {
    // 100 requests per minute per IP
    rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  100,
    }
    
    store := memory.NewStore()
    instance := limiter.New(store, rate)
    
    return func(c *gin.Context) {
        ip := c.ClientIP()
        context, err := instance.Get(c, ip)
        
        if err != nil {
            c.JSON(500, gin.H{"error": "Rate limiter error"})
            c.Abort()
            return
        }
        
        // Set rate limit headers
        c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
        c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
        c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))
        
        if context.Reached {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Apply to routes
api.Use(RateLimitMiddleware())
```

#### **Password Policy**

```
Minimum Requirements:
- Length: 6 characters (8+ recommended for production)
- No maximum (bcrypt handles long passwords)
- Optional: Require mix of upper, lower, numbers, symbols

Validation:
func validatePassword(password string) error {
    if len(password) < 6 {
        return errors.New("password must be at least 6 characters")
    }
    
    // Optional complexity checks
    var (
        hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
        hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
        hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
        hasSpecial = regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)
    )
    
    if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
        return errors.New("password must contain uppercase, lowercase, number, and special character")
    }
    
    return nil
}
```

#### **HTTPS & Security Headers**

```go
// Use Gin's secure middleware
import "github.com/gin-contrib/secure"

func SecurityMiddleware() gin.HandlerFunc {
    return secure.New(secure.Config{
        SSLRedirect:           true,  // Redirect HTTP to HTTPS
        STSSeconds:            315360000,
        STSIncludeSubdomains:  true,
        FrameDeny:             true,
        ContentTypeNosniff:    true,
        BrowserXssFilter:      true,
        ContentSecurityPolicy: "default-src 'self'",
    })
}

// Apply only in production
if config.AppConfig.Server.Env == "production" {
    r.Use(SecurityMiddleware())
}
```

---

### **1.17 Docker Setup**

#### **Dockerfile (Multi-stage Build)**

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
```

#### **docker-compose.yml**

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: slotswapper-db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: slotswapper
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: slotswapper-api
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: slotswapper
      PORT: 8080
      GIN_MODE: release
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: slotswapper-web
    ports:
      - "5173:80"
    depends_on:
      - backend

volumes:
  postgres_data:
```

#### **Development vs Production**

```bash
# Development (hot reload)
# Use air for auto-reload
go install github.com/cosmtrek/air@latest

# .air.toml configuration
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ."
  bin = "tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["assets", "tmp", "vendor"]
  delay = 1000

# Run with: air

# Production
# Build optimized binary
CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o main .
# Flags: -s -w remove symbol table and debug info (smaller binary)
```

---

### **1.18 API Documentation**

#### **Swagger/OpenAPI Setup**

```go
// Use swaggo/swag for automatic documentation
// go get -u github.com/swaggo/swag/cmd/swag
// go get -u github.com/swaggo/gin-swagger
// go get -u github.com/swaggo/files

import (
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// In main.go
// @title SlotSwapper API
// @version 1.0
// @description Peer-to-peer time-slot scheduling API
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
    // ... setup code ...
    
    // Swagger endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    r.Run()
}

// In controllers, add swagger comments:
// @Summary Create a new event
// @Description Create a new calendar event for the authenticated user
// @Tags events
// @Accept json
// @Produce json
// @Param event body CreateEventRequest true "Event details"
// @Success 201 {object} Event
// @Failure 400 {object} ErrorResponse
// @Security BearerAuth
// @Router /events [post]
func CreateEvent(c *gin.Context) {
    // ...
}

// Generate docs
// $ swag init
// Access at: http://localhost:8080/swagger/index.html
```

---

## **PART 2: FRONTEND ARCHITECTURE (Vue 3)**

### **2.1 Technology Stack**

```
Core:
- Vue 3 (Composition API)
- Vite (build tool)
- Vue Router (routing)
- Pinia (state management)
- Axios (HTTP client)

UI/Styling:
- Tailwind CSS (utility-first CSS)
- Heroicons or Lucide Vue (icons)

Optional:
- VueUse (composition utilities)
- Vue Toastification (notifications)
- Day.js (date manipulation)
```

### **2.2 Project Structure**

```
frontend/
├── public/
│   └── favicon.ico
├── src/
│   ├── main.js                  // Application entry
│   ├── App.vue                  // Root component
│   │
│   ├── router/
│   │   └── index.js             // Route definitions
│   │
│   ├── stores/
│   │   ├── auth.js              // Auth state & actions
│   │   ├── events.js            // Events state
│   │   └── swaps.js             // Swaps state
│   │
│   ├── api/
│   │   ├── axios.js             // Axios instance config
│   │   ├── auth.js              // Auth API calls
│   │   ├── events.js            // Events API calls
│   │   └── swaps.js             // Swaps API calls
│   │
│   ├── views/
│   │   ├── LoginView.vue
│   │   ├── SignupView.vue
│   │   ├── DashboardView.vue    // My calendar
│   │   ├── MarketplaceView.vue  // Browse swappable slots
│   │   └── RequestsView.vue     // Incoming/Outgoing
│   │
│   ├── components/
│   │   ├── layout/
│   │   │   ├── Navbar.vue
│   │   │   └── Layout.vue
│   │   ├── events/
│   │   │   ├── EventCard.vue
│   │   │   ├── EventForm.vue
│   │   │   └── EventList.vue
│   │   ├── swaps/
│   │   │   ├── SwapRequestCard.vue
│   │   │   ├── SwapModal.vue
│   │   │   └── SlotCard.vue
│   │   └── common/
│   │       ├── Button.vue
│   │       ├── Input.vue
│   │       └── LoadingSpinner.vue
│   │
│   ├── composables/
│   │   ├── useAuth.js           // Auth logic
│   │   ├── useEvents.js         // Events logic
│   │   └── useSwaps.js          // Swaps logic
│   │
│   ├── utils/
│   │   ├── dateFormatter.js
│   │   └── validators.js
│   │
│   └── assets/
│       └── styles/
│           └── main.css          // Tailwind imports
│
├── index.html
├── vite.config.js
├── tailwind.config.js
├── package.json
└── .env
```

### **2.3 State Management Architecture (Pinia)**

#### **Auth Store (`stores/auth.js`)**

```javascript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, signup } from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref(localStorage.getItem('token') || null)
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value)

  // Actions
  async function loginUser(credentials) {
    loading.value = true
    error.value = null
    
    try {
      const response = await login(credentials)
      token.value = response.token
      user.value = response.user
      
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
      
      router.push('/dashboard')
    } catch (err) {
      error.value = err.response?.data?.error || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function signupUser(credentials) {
    loading.value = true
    error.value = null
    
    try {
      const response = await signup(credentials)
      token.value = response.token
      user.value = response.user
      
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
      
      router.push('/dashboard')
    } catch (err) {
      error.value = err.response?.data?.error || 'Signup failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/login')
  }

  return {
    token,
    user,
    loading,
    error,
    isAuthenticated,
    loginUser,
    signupUser,
    logout
  }
})
```

#### **Events Store (`stores/events.js`)**

```javascript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getEvents, createEvent, updateEventStatus, deleteEvent } from '@/api/events'

export const useEventsStore = defineStore('events', () => {
  const events = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchEvents() {
    loading.value = true
    error.value = null
    
    try {
      const data = await getEvents()
      events.value = data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch events'
    } finally {
      loading.value = false
    }
  }

  async function addEvent(eventData) {
    loading.value = true
    error.value = null
    
    try {
      const newEvent = await createEvent(eventData)
      events.value.push(newEvent)
      return newEvent
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to create event'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateStatus(eventId, status) {
    loading.value = true
    error.value = null
    
    try {
      const updated = await updateEventStatus(eventId, status)
      const index = events.value.findIndex(e => e.id === eventId)
      if (index !== -1) {
        events.value[index] = updated
      }
      return updated
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to update event'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function removeEvent(eventId) {
    loading.value = true
    error.value = null
    
    try {
      await deleteEvent(eventId)
      events.value = events.value.filter(e => e.id !== eventId)
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to delete event'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    events,
    loading,
    error,
    fetchEvents,
    addEvent,
    updateStatus,
    removeEvent
  }
})
```

#### **Swaps Store (`stores/swaps.js`)**

```javascript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { 
  getSwappableSlots, 
  createSwapRequest, 
  respondToSwap,
  getIncomingRequests,
  getOutgoingRequests 
} from '@/api/swaps'
import { useEventsStore } from './events'

export const useSwapsStore = defineStore('swaps', () => {
  const swappableSlots = ref([])
  const incomingRequests = ref([])
  const outgoingRequests = ref([])
  const loading = ref(false)
  const error = ref(null)

  const pendingIncoming = computed(() => 
    incomingRequests.value.filter(r => r.status === 'PENDING')
  )

  const pendingOutgoing = computed(() => 
    outgoingRequests.value.filter(r => r.status === 'PENDING')
  )

  async function fetchSwappableSlots() {
    loading.value = true
    error.value = null
    error.value = null
    
    try {
      const data = await getSwappableSlots()
      swappableSlots.value = data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch slots'
    } finally {
      loading.value = false
    }
  }

  async function fetchIncomingRequests() {
    loading.value = true
    error.value = null
    
    try {
      const data = await getIncomingRequests()
      incomingRequests.value = data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch requests'
    } finally {
      loading.value = false
    }
  }

  async function fetchOutgoingRequests() {
    loading.value = true
    error.value = null
    
    try {
      const data = await getOutgoingRequests()
      outgoingRequests.value = data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch requests'
    } finally {
      loading.value = false
    }
  }

  async function requestSwap(mySlotId, theirSlotId) {
    loading.value = true
    error.value = null
    
    try {
      const newRequest = await createSwapRequest(mySlotId, theirSlotId)
      outgoingRequests.value.unshift(newRequest)
      
      // Refresh events to reflect SWAP_PENDING status
      const eventsStore = useEventsStore()
      await eventsStore.fetchEvents()
      
      // Remove slot from marketplace
      swappableSlots.value = swappableSlots.value.filter(
        s => s.id !== theirSlotId
      )
      
      return newRequest
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to create swap request'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function acceptSwap(requestId) {
    loading.value = true
    error.value = null
    
    try {
      await respondToSwap(requestId, true)
      
      // Update request status locally
      const request = incomingRequests.value.find(r => r.id === requestId)
      if (request) {
        request.status = 'ACCEPTED'
      }
      
      // Refresh events to reflect ownership change
      const eventsStore = useEventsStore()
      await eventsStore.fetchEvents()
      
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to accept swap'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function rejectSwap(requestId) {
    loading.value = true
    error.value = null
    
    try {
      await respondToSwap(requestId, false)
      
      // Update request status locally
      const request = incomingRequests.value.find(r => r.id === requestId)
      if (request) {
        request.status = 'REJECTED'
      }
      
      // Refresh events to reflect SWAPPABLE status
      const eventsStore = useEventsStore()
      await eventsStore.fetchEvents()
      
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to reject swap'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    swappableSlots,
    incomingRequests,
    outgoingRequests,
    pendingIncoming,
    pendingOutgoing,
    loading,
    error,
    fetchSwappableSlots,
    fetchIncomingRequests,
    fetchOutgoingRequests,
    requestSwap,
    acceptSwap,
    rejectSwap
  }
})
```

---

### **2.4 API Layer Implementation**

#### **Axios Configuration (`api/axios.js`)**

```javascript
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json'
  },
  timeout: 10000
})

// Request interceptor - Add JWT token
apiClient.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor - Handle errors globally
apiClient.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response) {
      // Handle 401 Unauthorized - logout user
      if (error.response.status === 401) {
        const authStore = useAuthStore()
        authStore.logout()
        router.push('/login')
      }
      
      // Handle 403 Forbidden
      if (error.response.status === 403) {
        console.error('Forbidden:', error.response.data)
      }
      
      // Handle 404 Not Found
      if (error.response.status === 404) {
        console.error('Not found:', error.response.data)
      }
      
      // Handle 500 Server Error
      if (error.response.status >= 500) {
        console.error('Server error:', error.response.data)
      }
    } else if (error.request) {
      // Network error
      console.error('Network error:', error.message)
    }
    
    return Promise.reject(error)
  }
)

export default apiClient
```

#### **Auth API (`api/auth.js`)**

```javascript
import apiClient from './axios'

export const signup = async (credentials) => {
  return await apiClient.post('/auth/signup', credentials)
}

export const login = async (credentials) => {
  return await apiClient.post('/auth/login', credentials)
}
```

#### **Events API (`api/events.js`)**

```javascript
import apiClient from './axios'

export const getEvents = async () => {
  return await apiClient.get('/events')
}

export const createEvent = async (eventData) => {
  return await apiClient.post('/events', eventData)
}

export const updateEventStatus = async (eventId, status) => {
  return await apiClient.patch(`/events/${eventId}/status`, { status })
}

export const deleteEvent = async (eventId) => {
  return await apiClient.delete(`/events/${eventId}`)
}
```

#### **Swaps API (`api/swaps.js`)**

```javascript
import apiClient from './axios'

export const getSwappableSlots = async () => {
  return await apiClient.get('/swappable-slots')
}

export const createSwapRequest = async (mySlotId, theirSlotId) => {
  return await apiClient.post('/swap-request', {
    my_slot_id: mySlotId,
    their_slot_id: theirSlotId
  })
}

export const respondToSwap = async (requestId, accept) => {
  return await apiClient.post(`/swap-response/${requestId}`, {
    accept
  })
}

export const getIncomingRequests = async () => {
  return await apiClient.get('/swap-requests/incoming')
}

export const getOutgoingRequests = async () => {
  return await apiClient.get('/swap-requests/outgoing')
}
```

---

### **2.5 Router Configuration**

#### **Router Setup (`router/index.js`)**

```javascript
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { requiresAuth: false, hideForAuth: true }
  },
  {
    path: '/signup',
    name: 'Signup',
    component: () => import('@/views/SignupView.vue'),
    meta: { requiresAuth: false, hideForAuth: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/DashboardView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/marketplace',
    name: 'Marketplace',
    component: () => import('@/views/MarketplaceView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/requests',
    name: 'Requests',
    component: () => import('@/views/RequestsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFoundView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const hideForAuth = to.matched.some(record => record.meta.hideForAuth)
  
  // Redirect authenticated users away from login/signup
  if (hideForAuth && authStore.isAuthenticated) {
    next('/dashboard')
  }
  // Redirect unauthenticated users to login
  else if (requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  }
  else {
    next()
  }
})

export default router
```

---

### **2.6 Key Component Specifications**

#### **LoginView.vue**

```
Purpose: User authentication form

Features:
- Email and password inputs
- Form validation (client-side)
- Loading state during submission
- Error display
- Link to signup page

State Management:
- Uses authStore.loginUser()
- Displays authStore.error
- Shows authStore.loading

Layout:
- Centered card layout
- Responsive (mobile-first)
- Simple, clean design

Validation:
- Email format validation
- Password minimum length (6 chars)
- Required field validation
```

#### **DashboardView.vue (My Calendar)**

```
Purpose: Display user's events with management capabilities

Features:
- List of all user's events (sorted by start_time)
- Create new event button/form
- For each event:
  - Display title, date, time, status
  - Status badge (BUSY, SWAPPABLE, SWAP_PENDING)
  - "Make Swappable" button (if BUSY)
  - "Remove from Marketplace" button (if SWAPPABLE)
  - Delete button (if not SWAP_PENDING)
- Filter by status (optional)
- Empty state message

Data Flow:
- onMounted: eventsStore.fetchEvents()
- Watch for store changes (reactive)
- Status updates: eventsStore.updateStatus()

UI Components:
- EventList component
- EventCard component
- EventForm modal/component
- Status badges with colors:
  - BUSY: gray
  - SWAPPABLE: green
  - SWAP_PENDING: yellow
```

#### **MarketplaceView.vue**

```
Purpose: Browse and request swaps for available slots

Features:
- Grid/list of swappable slots from other users
- For each slot:
  - Display title, owner name, date, time
  - "Request Swap" button
- Clicking "Request Swap" opens modal
- Modal shows:
  - Selected slot details
  - User's SWAPPABLE slots to offer
  - Confirm button
- Empty state if no slots available
- Loading state
- Success toast notification after request

Data Flow:
- onMounted: swapsStore.fetchSwappableSlots()
- User selects slot to request
- Shows modal with user's swappable slots
- User selects their slot to offer
- Calls swapsStore.requestSwap(mySlotId, theirSlotId)
- Shows success notification
- Closes modal

Validation:
- User must have at least one SWAPPABLE slot to request
- If no swappable slots, show message: "You need to mark a slot as swappable first"

UI Components:
- SlotCard component (for marketplace slots)
- SwapModal component
- EmptyState component
```

#### **RequestsView.vue**

```
Purpose: Manage incoming and outgoing swap requests

Layout:
- Two-column layout (or tabs on mobile)
- Left: Incoming Requests
- Right: Outgoing Requests

Incoming Requests Section:
- List of pending requests where user is receiver
- For each request:
  - Requester name
  - "They want:" (their slot details)
  - "In exchange for:" (your slot details)
  - "Accept" button (green)
  - "Reject" button (red)
- Tabs/filter: Pending | Accepted | Rejected
- Empty state if no requests

Outgoing Requests Section:
- List of requests user has sent
- For each request:
  - Receiver name
  - "You offered:" (your slot details)
  - "For:" (their slot details)
  - Status badge: Pending/Accepted/Rejected
- No actions (waiting for receiver response)
- Empty state if no requests

Data Flow:
- onMounted:
  - swapsStore.fetchIncomingRequests()
  - swapsStore.fetchOutgoingRequests()
- Accept: swapsStore.acceptSwap(requestId)
- Reject: swapsStore.rejectSwap(requestId)
- Confirmation dialog before accept/reject
- Success notifications

Real-time Updates:
- Poll every 30 seconds for new requests (setInterval)
- Or implement with WebSockets (bonus feature)

UI Components:
- SwapRequestCard component
- ConfirmationDialog component
- Status badges
```

#### **Navbar Component**

```
Purpose: Global navigation

Features:
- App logo/title
- Navigation links:
  - Dashboard
  - Marketplace
  - Requests (with badge showing pending count)
- User info dropdown:
  - Display user name
  - Logout button
- Responsive (hamburger menu on mobile)

Data:
- Uses authStore.user
- Uses swapsStore.pendingIncoming.length for badge

Styling:
- Sticky/fixed to top
- Consistent across all pages
- Active route highlighting
```

---

### **2.7 Component Communication Patterns**

#### **Props Down, Events Up**

```javascript
// Parent component
<EventCard 
  :event="event" 
  @update-status="handleStatusUpdate"
  @delete="handleDelete"
/>

// Child component (EventCard.vue)
const props = defineProps({
  event: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update-status', 'delete'])

function updateStatus(newStatus) {
  emit('update-status', props.event.id, newStatus)
}
```

#### **Composables for Reusable Logic**

```javascript
// composables/useEvents.js
import { useEventsStore } from '@/stores/events'
import { storeToRefs } from 'pinia'

export function useEvents() {
  const eventsStore = useEventsStore()
  const { events, loading, error } = storeToRefs(eventsStore)

  const busyEvents = computed(() => 
    events.value.filter(e => e.status === 'BUSY')
  )

  const swappableEvents = computed(() => 
    events.value.filter(e => e.status === 'SWAPPABLE')
  )

  const pendingEvents = computed(() => 
    events.value.filter(e => e.status === 'SWAP_PENDING')
  )

  return {
    events,
    loading,
    error,
    busyEvents,
    swappableEvents,
    pendingEvents,
    fetchEvents: eventsStore.fetchEvents,
    addEvent: eventsStore.addEvent,
    updateStatus: eventsStore.updateStatus,
    removeEvent: eventsStore.removeEvent
  }
}
```

---

### **2.8 Form Handling & Validation**

#### **Event Creation Form**

```javascript
// EventForm.vue
<script setup>
import { ref, computed } from 'vue'
import { useEvents } from '@/composables/useEvents'

const { addEvent } = useEvents()

const form = ref({
  title: '',
  startTime: '',
  endTime: ''
})

const errors = ref({})
const isSubmitting = ref(false)

const isValid = computed(() => {
  return form.value.title &&
         form.value.startTime &&
         form.value.endTime &&
         new Date(form.value.endTime) > new Date(form.value.startTime)
})

function validateForm() {
  errors.value = {}
  
  if (!form.value.title) {
    errors.value.title = 'Title is required'
  }
  
  if (!form.value.startTime) {
    errors.value.startTime = 'Start time is required'
  }
  
  if (!form.value.endTime) {
    errors.value.endTime = 'End time is required'
  }
  
  if (form.value.startTime && form.value.endTime) {
    const start = new Date(form.value.startTime)
    const end = new Date(form.value.endTime)
    
    if (end <= start) {
      errors.value.endTime = 'End time must be after start time'
    }
    
    if (start < new Date()) {
      errors.value.startTime = 'Start time must be in the future'
    }
  }
  
  return Object.keys(errors.value).length === 0
}

async function handleSubmit() {
  if (!validateForm()) return
  
  isSubmitting.value = true
  
  try {
    await addEvent({
      title: form.value.title,
      start_time: new Date(form.value.startTime).toISOString(),
      end_time: new Date(form.value.endTime).toISOString()
    })
    
    // Reset form
    form.value = { title: '', startTime: '', endTime: '' }
    
    // Show success notification
    // Close modal if in modal
    emit('success')
  } catch (error) {
    errors.value.submit = error.response?.data?.error || 'Failed to create event'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <div>
      <label>Title</label>
      <input 
        v-model="form.title" 
        type="text" 
        placeholder="Team Meeting"
        :class="{ 'border-red-500': errors.title }"
      />
      <span v-if="errors.title" class="text-red-500">{{ errors.title }}</span>
    </div>
    
    <div>
      <label>Start Time</label>
      <input 
        v-model="form.startTime" 
        type="datetime-local"
        :class="{ 'border-red-500': errors.startTime }"
      />
      <span v-if="errors.startTime" class="text-red-500">{{ errors.startTime }}</span>
    </div>
    
    <div>
      <label>End Time</label>
      <input 
        v-model="form.endTime" 
        type="datetime-local"
        :class="{ 'border-red-500': errors.endTime }"
      />
      <span v-if="errors.endTime" class="text-red-500">{{ errors.endTime }}</span>
    </div>
    
    <span v-if="errors.submit" class="text-red-500">{{ errors.submit }}</span>
    
    <button 
      type="submit" 
      :disabled="!isValid || isSubmitting"
    >
      {{ isSubmitting ? 'Creating...' : 'Create Event' }}
    </button>
  </form>
</template>
```

---

### **2.9 Date/Time Handling**

#### **Date Formatting Utilities**

```javascript
// utils/dateFormatter.js
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

export function formatDateTime(dateString) {
  return dayjs(dateString).format('MMM D, YYYY h:mm A')
}

export function formatDate(dateString) {
  return dayjs(dateString).format('MMM D, YYYY')
}

export function formatTime(dateString) {
  return dayjs(dateString).format('h:mm A')
}

export function formatTimeRange(startTime, endTime) {
  const start = dayjs(startTime)
  const end = dayjs(endTime)
  
  if (start.isSame(end, 'day')) {
    return `${start.format('MMM D')} • ${start.format('h:mm A')} - ${end.format('h:mm A')}`
  }
  
  return `${start.format('MMM D h:mm A')} - ${end.format('MMM D h:mm A')}`
}

export function timeFromNow(dateString) {
  return dayjs(dateString).fromNow()
}

export function isInPast(dateString) {
  return dayjs(dateString).isBefore(dayjs())
}

export function getDuration(startTime, endTime) {
  const start = dayjs(startTime)
  const end = dayjs(endTime)
  const minutes = end.diff(start, 'minute')
  
  if (minutes < 60) {
    return `${minutes} min`
  }
  
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  
  if (remainingMinutes === 0) {
    return `${hours} hr`
  }
  
  return `${hours} hr ${remainingMinutes} min`
}
```

---

### **2.10 Error Handling & User Feedback**

#### **Toast Notifications**

```javascript
// Install: npm install vue-toastification
// main.js
import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'

app.use(Toast, {
  position: 'top-right',
  timeout: 3000,
  closeOnClick: true,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: false,
  hideProgressBar: false,
  closeButton: 'button',
  icon: true,
  rtl: false
})

// Usage in components
import { useToast } from 'vue-toastification'

const toast = useToast()

// Success
toast.success('Event created successfully!')

// Error
toast.error('Failed to create event')

// Info
toast.info('Swap request sent')

// Warning
toast.warning('This slot is already pending')
```

#### **Loading States**

```javascript
// LoadingSpinner.vue
<template>
  <div v-if="loading" class="flex justify-center items-center p-8">
    <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
  </div>
</template>

// Usage
<LoadingSpinner :loading="eventsStore.loading" />

// Or inline
<div v-if="loading" class="loading-spinner"></div>
<div v-else>
  <!-- Content -->
</div>
```

#### **Error States**

```javascript
// ErrorMessage.vue
<template>
  <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
    <p>{{ error }}</p>
    <button @click="$emit('retry')" class="text-sm underline">Retry</button>
  </div>
</template>

// Usage
<ErrorMessage 
  :error="eventsStore.error" 
  @retry="eventsStore.fetchEvents()"
/>
```

#### **Empty States**

```javascript
// EmptyState.vue
<template>
  <div class="flex flex-col items-center justify-center p-12 text-gray-500">
    <svg class="w-16 h-16 mb-4" /><!-- Icon -->
    <h3 class="text-lg font-medium mb-2">{{ title }}</h3>
    <p class="text-sm text-center mb-4">{{ message }}</p>
    <slot name="action"></slot>
  </div>
</template>

// Usage
<EmptyState 
  v-if="events.length === 0"
  title="No events yet"
  message="Create your first event to get started"
>
  <template #action>
    <button @click="showCreateForm">Create Event</button>
  </template>
</EmptyState>
```

---

### **2.11 Styling with Tailwind CSS**

#### **tailwind.config.js**

```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
        }
      }
    },
  },
  plugins: [],
}
```

#### **Common Utility Classes**

```
Layout:
- Container: max-w-7xl mx-auto px-4 sm:px-6 lg:px-8
- Card: bg-white rounded-lg shadow p-6
- Grid: grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6

Buttons:
- Primary: bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded
- Secondary: bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded
- Danger: bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded
- Disabled: opacity-50 cursor-not-allowed

Status Badges:
- BUSY: bg-gray-100 text-gray-800 px-2 py-1 rounded text-sm
- SWAPPABLE: bg-green-100 text-green-800 px-2 py-1 rounded text-sm
- SWAP_PENDING: bg-yellow-100 text-yellow-800 px-2 py-1 rounded text-sm

Forms:
- Input: border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500
- Label: block text-sm font-medium text-gray-700 mb-1
- Error: text-red-500 text-sm mt-1
```

---

### **2.12 Environment Configuration**

#### **.env File**

```bash
# API Base URL
VITE_API_BASE_URL=http://localhost:8080/api

# App Configuration
VITE_APP_NAME=SlotSwapper
VITE_APP_DESCRIPTION=Peer-to-peer time-slot scheduling

# Feature Flags
VITE_ENABLE_WEBSOCKETS=false
```

#### **vite.config.js**

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

---

## **PART 3: INTEGRATION & DEPLOYMENT**

### **3.1 Development Workflow**

#### **Step-by-Step Development Order**

```
Day 1: Backend Core (12 hours)
├── Hour 0-1: Project setup, database connection
├── Hour 1-3: User model, auth endpoints, JWT
├── Hour 3-5: Event model, CRUD endpoints
├── Hour 5-8: Swap logic (CRITICAL - take time)
├── Hour 8-10: Testing swap transactions
├── Hour 10-12: API documentation, cleanup

Day 2: Frontend & Integration (12 hours)
├── Hour 12-13: Vue project setup, routing
├── Hour 13-15: Auth pages (login/signup)
├── Hour 15-17: Dashboard (events CRUD)
├── Hour 17-19: Marketplace view
├── Hour 19-21: Requests view (accept/reject)
├── Hour 21-22: Polish, error handling
├── Hour 22-23: End-to-end testing
├── Hour 23-24: README, deployment prep
```

#### **Testing Checklist**

```
Backend Tests:
☐ User can signup with valid credentials
☐ User cannot signup with existing email
☐ User can login with correct password
☐ User cannot login with wrong password
☐ JWT token is generated correctly
☐ Protected routes reject requests without token
☐ User can create events
☐ User can update event status to SWAPPABLE
☐ User can view only their own events
☐ Marketplace shows only other users' swappable slots
☐ User can create swap request with valid slots
☐ Swap request sets both slots to SWAP_PENDING
☐ User cannot request swap with non-swappable slot
☐ Receiver can accept swap request
☐ Acceptance swaps ownership correctly
☐ Acceptance sets both slots to BUSY
☐ Receiver can reject swap request
☐ Rejection sets both slots back to SWAPPABLE
☐ Only receiver can respond to swap request
☐ Cannot respond to already-processed request

Frontend Tests:
☐ Login page redirects to dashboard on success
☐ Dashboard shows user's events
☐ User can create new event via form
☐ User can mark event as swappable
☐ Marketplace shows available slots
☐ User can request swap from marketplace
☐ Requests page shows incoming and outgoing
☐ User can accept incoming request
☐ User can reject incoming request
☐ UI updates after accept/reject
☐ Loading states display correctly
☐ Error messages display correctly
☐ Logout clears token and redirects to login
☐ Protected routes redirect to login if not authenticated

Integration Tests:
☐ Complete swap flow from start to finish
☐ Multiple users can swap simultaneously
☐ Race conditions handled correctly
☐ CORS allows frontend to call backend
☐ JWT expiration works correctly
```

---

### **3.2 Deployment Strategy**

#### **Backend Deployment (Render/Railway/Fly.io)**

```bash
# Render.com (Recommended for simplicity)
1. Create account on render.com
2. New > Web Service
3. Connect GitHub repository
4. Configure:
   - Name: slotswapper-api
   - Environment: Go
   - Build Command: go build -o main .
   - Start Command: ./main
   - Environment Variables:
     * PORT (Render sets automatically)
     * DB_HOST (from Render PostgreSQL)
     * DB_PORT
     * DB_USER
     * DB_PASSWORD
     * DB_NAME
     * JWT_SECRET (generate secure key)
     * GIN_MODE=release
5. Add PostgreSQL database
6. Deploy

# Alternative: Railway.app
```bash
# Alternative: Railway.app
1. Create account on railway.app
2. New Project > Deploy from GitHub
3. Add PostgreSQL service
4. Configure environment variables (auto-detected)
5. Railway generates domain automatically

# Alternative: Fly.io (More control)
1. Install flyctl CLI
2. fly auth login
3. fly launch (in backend directory)
4. fly postgres create
5. fly secrets set JWT_SECRET=your-secret
6. fly deploy

# Health Check Endpoint (Add to backend)
// In routes/routes.go
r.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "healthy",
        "timestamp": time.Now().Unix()
    })
})
```

#### **Frontend Deployment (Vercel/Netlify)**

```bash
# Vercel (Recommended)
1. Install Vercel CLI: npm i -g vercel
2. In frontend directory: vercel login
3. vercel
4. Follow prompts
5. Set environment variable:
   - VITE_API_BASE_URL=https://your-backend.render.com/api
6. vercel --prod

# Alternative: Netlify
1. Install Netlify CLI: npm i -g netlify-cli
2. In frontend directory: netlify login
3. netlify init
4. Build command: npm run build
5. Publish directory: dist
6. Environment variables:
   - VITE_API_BASE_URL
7. netlify deploy --prod

# Build Configuration
# netlify.toml (in frontend root)
[build]
  command = "npm run build"
  publish = "dist"

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200
```

#### **Docker Deployment (Optional)**

```dockerfile
# Frontend Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

```nginx
# nginx.conf (for Vue Router history mode)
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

### **3.3 Environment-Specific Configuration**

#### **Backend Environment Variables**

```bash
# Development (.env.development)
PORT=8080
GIN_MODE=debug
ENVIRONMENT=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=slotswapper_dev
DB_SSLMODE=disable
JWT_SECRET=dev-secret-key-not-for-production
JWT_EXPIRATION_HOURS=24
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
LOG_LEVEL=debug

# Production (.env.production - set in hosting platform)
PORT=8080
GIN_MODE=release
ENVIRONMENT=production
DB_HOST=<production-db-host>
DB_PORT=5432
DB_USER=<production-db-user>
DB_PASSWORD=<production-db-password>
DB_NAME=slotswapper
DB_SSLMODE=require
JWT_SECRET=<generate-32-char-random-string>
JWT_EXPIRATION_HOURS=24
ALLOWED_ORIGINS=https://your-frontend.vercel.app
LOG_LEVEL=info
```

#### **Frontend Environment Variables**

```bash
# .env.development
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_NAME=SlotSwapper (Dev)
VITE_ENABLE_DEBUG=true

# .env.production
VITE_API_BASE_URL=https://your-backend.render.com/api
VITE_APP_NAME=SlotSwapper
VITE_ENABLE_DEBUG=false
```

---

### **3.4 Monitoring & Debugging**

#### **Backend Logging**

```go
// Enhanced logging for production
import "github.com/rs/zerolog/log"

// In swap service
func (s *SwapService) AcceptSwapRequest(requestID uint, receiverID uint) error {
    log.Info().
        Uint("request_id", requestID).
        Uint("receiver_id", receiverID).
        Msg("Starting swap acceptance")
    
    err := s.db.Transaction(func(tx *gorm.DB) error {
        // ... transaction logic ...
    })
    
    if err != nil {
        log.Error().
            Err(err).
            Uint("request_id", requestID).
            Msg("Failed to accept swap")
        return err
    }
    
    log.Info().
        Uint("request_id", requestID).
        Msg("Swap accepted successfully")
    
    return nil
}
```

#### **Frontend Error Tracking**

```javascript
// Add Sentry (optional)
// npm install @sentry/vue

// main.js
import * as Sentry from "@sentry/vue"

if (import.meta.env.PROD) {
  Sentry.init({
    app,
    dsn: "your-sentry-dsn",
    integrations: [
      new Sentry.BrowserTracing({
        routingInstrumentation: Sentry.vueRouterInstrumentation(router),
      }),
    ],
    tracesSampleRate: 1.0,
  })
}
```

#### **Performance Monitoring**

```go
// Add request timing middleware
func TimingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        
        log.Debug().
            Str("method", c.Request.Method).
            Str("path", c.Request.URL.Path).
            Int("status", c.Writer.Status()).
            Dur("duration", duration).
            Msg("Request completed")
        
        // Alert if request takes > 1 second
        if duration > time.Second {
            log.Warn().
                Str("path", c.Request.URL.Path).
                Dur("duration", duration).
                Msg("Slow request detected")
        }
    }
}
```

---

## **PART 4: BONUS FEATURES (If Time Permits)**

### **4.1 Real-time Notifications (WebSockets)**

#### **Backend WebSocket Implementation**

```go
// Install: go get github.com/gorilla/websocket

// websocket/hub.go
type Hub struct {
    clients    map[uint]*Client  // userID -> Client
    broadcast  chan Message
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

type Client struct {
    hub    *Hub
    conn   *websocket.Conn
    userID uint
    send   chan []byte
}

type Message struct {
    Type     string      `json:"type"`
    UserID   uint        `json:"user_id"`
    Data     interface{} `json:"data"`
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[uint]*Client),
        broadcast:  make(chan Message),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client.userID] = client
            h.mu.Unlock()
            
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client.userID]; ok {
                delete(h.clients, client.userID)
                close(client.send)
            }
            h.mu.Unlock()
            
        case message := <-h.broadcast:
            h.mu.RLock()
            if client, ok := h.clients[message.UserID]; ok {
                select {
                case client.send <- marshalMessage(message):
                default:
                    close(client.send)
                    delete(h.clients, message.UserID)
                }
            }
            h.mu.RUnlock()
        }
    }
}

// In swap service, after accepting swap:
hub.broadcast <- Message{
    Type:   "swap_accepted",
    UserID: swapReq.RequesterID,
    Data: map[string]interface{}{
        "request_id": swapReq.ID,
        "message": "Your swap request was accepted!",
    },
}
```

#### **Frontend WebSocket Client**

```javascript
// composables/useWebSocket.js
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

export function useWebSocket() {
  const authStore = useAuthStore()
  const ws = ref(null)
  const isConnected = ref(false)

  function connect() {
    if (!authStore.token) return

    ws.value = new WebSocket(
      `ws://localhost:8080/ws?token=${authStore.token}`
    )

    ws.value.onopen = () => {
      console.log('WebSocket connected')
      isConnected.value = true
    }

    ws.value.onmessage = (event) => {
      const message = JSON.parse(event.data)
      handleMessage(message)
    }

    ws.value.onclose = () => {
      console.log('WebSocket disconnected')
      isConnected.value = false
      // Reconnect after 5 seconds
      setTimeout(connect, 5000)
    }

    ws.value.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }

  function handleMessage(message) {
    switch (message.type) {
      case 'swap_accepted':
        toast.success('Your swap request was accepted!')
        // Refresh data
        break
      case 'swap_requested':
        toast.info('You have a new swap request')
        // Refresh incoming requests
        break
      case 'swap_rejected':
        toast.warning('Your swap request was rejected')
        break
    }
  }

  function disconnect() {
    if (ws.value) {
      ws.value.close()
    }
  }

  onMounted(connect)
  onUnmounted(disconnect)

  return {
    isConnected,
    connect,
    disconnect
  }
}
```

---

### **4.2 Unit Testing**

#### **Backend Tests (Go)**

```go
// tests/swap_test.go
package tests

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestAcceptSwapRequest(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    defer cleanupTestDB(db)
    
    service := services.NewSwapService(db)
    
    // Create test users
    user1 := createTestUser(db, "user1@test.com")
    user2 := createTestUser(db, "user2@test.com")
    
    // Create test events
    event1 := createTestEvent(db, user1.ID, models.StatusSwappable)
    event2 := createTestEvent(db, user2.ID, models.StatusSwappable)
    
    // Create swap request
    swapReq := createTestSwapRequest(db, user1.ID, user2.ID, event1.ID, event2.ID)
    
    // Test acceptance
    err := service.AcceptSwapRequest(swapReq.ID, user2.ID)
    assert.NoError(t, err)
    
    // Verify owners swapped
    var updatedEvent1 models.Event
    db.First(&updatedEvent1, event1.ID)
    assert.Equal(t, user2.ID, updatedEvent1.UserID)
    
    var updatedEvent2 models.Event
    db.First(&updatedEvent2, event2.ID)
    assert.Equal(t, user1.ID, updatedEvent2.UserID)
    
    // Verify statuses changed to BUSY
    assert.Equal(t, models.StatusBusy, updatedEvent1.Status)
    assert.Equal(t, models.StatusBusy, updatedEvent2.Status)
}

func TestRejectSwapRequest(t *testing.T) {
    // Similar setup...
    
    err := service.RejectSwapRequest(swapReq.ID, user2.ID)
    assert.NoError(t, err)
    
    // Verify statuses back to SWAPPABLE
    var updatedEvent1 models.Event
    db.First(&updatedEvent1, event1.ID)
    assert.Equal(t, models.StatusSwappable, updatedEvent1.Status)
}

func TestConcurrentSwapRequests(t *testing.T) {
    // Test race conditions
    // Create multiple goroutines trying to swap same slot
    // Ensure only one succeeds
}
```

#### **Frontend Tests (Vitest + Vue Test Utils)**

```javascript
// tests/EventCard.test.js
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import EventCard from '@/components/events/EventCard.vue'

describe('EventCard', () => {
  it('renders event details correctly', () => {
    const event = {
      id: 1,
      title: 'Test Meeting',
      start_time: '2025-01-20T10:00:00Z',
      end_time: '2025-01-20T11:00:00Z',
      status: 'BUSY'
    }
    
    const wrapper = mount(EventCard, {
      props: { event }
    })
    
    expect(wrapper.text()).toContain('Test Meeting')
    expect(wrapper.text()).toContain('BUSY')
  })
  
  it('emits update-status event when button clicked', async () => {
    const event = {
      id: 1,
      title: 'Test Meeting',
      status: 'BUSY'
    }
    
    const wrapper = mount(EventCard, {
      props: { event }
    })
    
    await wrapper.find('[data-test="make-swappable"]').trigger('click')
    
    expect(wrapper.emitted('update-status')).toBeTruthy()
    expect(wrapper.emitted('update-status')[0]).toEqual([1, 'SWAPPABLE'])
  })
})
```

---

### **4.3 Advanced Features**

#### **Email Notifications**

```go
// services/email_service.go
import "github.com/sendgrid/sendgrid-go"

type EmailService struct {
    apiKey string
}

func (s *EmailService) SendSwapNotification(recipient string, swapDetails SwapRequest) error {
    message := mail.NewSingleEmail(
        mail.NewEmail("SlotSwapper", "noreply@slotswapper.com"),
        "New Swap Request",
        mail.NewEmail("", recipient),
        fmt.Sprintf("You have a new swap request from %s", swapDetails.Requester.Name),
        fmt.Sprintf("<p>%s wants to swap slots with you!</p>", swapDetails.Requester.Name),
    )
    
    client := sendgrid.NewSendClient(s.apiKey)
    _, err := client.Send(message)
    return err
}
```

#### **Search & Filters**

```javascript
// In MarketplaceView.vue
const searchQuery = ref('')
const dateFilter = ref('all') // all, today, this-week, next-week

const filteredSlots = computed(() => {
  let slots = swappableSlots.value
  
  // Search by title or owner name
  if (searchQuery.value) {
    slots = slots.filter(slot => 
      slot.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      slot.user.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  
  // Filter by date range
  if (dateFilter.value !== 'all') {
    const now = dayjs()
    slots = slots.filter(slot => {
      const start = dayjs(slot.start_time)
      
      switch (dateFilter.value) {
        case 'today':
          return start.isSame(now, 'day')
        case 'this-week':
          return start.isSame(now, 'week')
        case 'next-week':
          return start.isSame(now.add(1, 'week'), 'week')
        default:
          return true
      }
    })
  }
  
  return slots
})
```

#### **Calendar View Integration**

```javascript
// Install: npm install @fullcalendar/vue3 @fullcalendar/daygrid

import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from '@fullcalendar/interaction'

const calendarEvents = computed(() => {
  return events.value.map(event => ({
    id: event.id,
    title: event.title,
    start: event.start_time,
    end: event.end_time,
    backgroundColor: getEventColor(event.status),
    borderColor: getEventColor(event.status)
  }))
})

function getEventColor(status) {
  switch (status) {
    case 'BUSY': return '#6B7280'
    case 'SWAPPABLE': return '#10B981'
    case 'SWAP_PENDING': return '#F59E0B'
    default: return '#6B7280'
  }
}
```

---

## **PART 5: README DOCUMENTATION**

### **5.1 README Structure**

```markdown
# SlotSwapper

A peer-to-peer time-slot scheduling application that allows users to swap calendar events.

## Features

- 🔐 JWT-based authentication
- 📅 Calendar management (Create, Read, Update, Delete events)
- 🔄 Peer-to-peer slot swapping
- 🔔 Real-time swap request notifications
- 📱 Responsive design

## Tech Stack

**Backend:**
- Go 1.21+
- Gin Web Framework
- GORM (PostgreSQL)
- JWT for authentication

**Frontend:**
- Vue 3 (Composition API)
- Vite
- Pinia (State Management)
- Tailwind CSS
- Axios

**Database:**
- PostgreSQL 15+

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- PostgreSQL 15 or higher
- Git

## Local Setup

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/slotswapper.git
cd slotswapper
```

### 2. Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Create .env file
cp .env.example .env

# Edit .env with your database credentials
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=postgres
# DB_NAME=slotswapper

# Create database
createdb slotswapper

# Run the server (auto-migration enabled)
go run main.go

# Server runs on http://localhost:8080
```

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Create .env file
cp .env.example .env

# Edit .env
# VITE_API_BASE_URL=http://localhost:8080/api

# Run development server
npm run dev

# Frontend runs on http://localhost:5173
```

## Docker Setup

```bash
# Run entire stack with docker-compose
docker-compose up

# Backend: http://localhost:8080
# Frontend: http://localhost:5173
# PostgreSQL: localhost:5432
```

## API Documentation

### Authentication

**POST /api/auth/signup**
- Request: `{ "name": "string", "email": "string", "password": "string" }`
- Response: `{ "token": "string", "user": {...} }`

**POST /api/auth/login**
- Request: `{ "email": "string", "password": "string" }`
- Response: `{ "token": "string", "user": {...} }`

### Events

**GET /api/events**
- Headers: `Authorization: Bearer <token>`
- Response: Array of user's events

**POST /api/events**
- Headers: `Authorization: Bearer <token>`
- Request: `{ "title": "string", "start_time": "ISO8601", "end_time": "ISO8601" }`
- Response: Created event object

**PATCH /api/events/:id/status**
- Headers: `Authorization: Bearer <token>`
- Request: `{ "status": "BUSY|SWAPPABLE" }`
- Response: Updated event object

### Swaps

**GET /api/swappable-slots**
- Headers: `Authorization: Bearer <token>`
- Response: Array of swappable slots from other users

**POST /api/swap-request**
- Headers: `Authorization: Bearer <token>`
- Request: `{ "my_slot_id": number, "their_slot_id": number }`
- Response: Created swap request

**POST /api/swap-response/:id**
- Headers: `Authorization: Bearer <token>`
- Request: `{ "accept": boolean }`
- Response: `{ "message": "string" }`

**GET /api/swap-requests/incoming**
- Headers: `Authorization: Bearer <token>`
- Response: Array of incoming swap requests

**GET /api/swap-requests/outgoing**
- Headers: `Authorization: Bearer <token>`
- Response: Array of outgoing swap requests

## Project Structure

```
slotswapper/
├── backend/
│   ├── config/          # Database & environment configuration
│   ├── controllers/     # HTTP request handlers
│   ├── middleware/      # Auth, CORS, logging middleware
│   ├── models/          # Database models
│   ├── routes/          # Route definitions
│   ├── services/        # Business logic
│   ├── utils/           # Helper functions
│   └── main.go          # Application entry point
│
├── frontend/
│   ├── src/
│   │   ├── api/         # API client functions
│   │   ├── components/  # Vue components
│   │   ├── composables/ # Reusable composition functions
│   │   ├── router/      # Vue Router configuration
│   │   ├── stores/      # Pinia stores
│   │   ├── utils/       # Helper functions
│   │   ├── views/       # Page components
│   │   └── main.js      # Application entry point
│   └── package.json
│
└── docker-compose.yml
```

## Design Decisions

### Backend Architecture
- **Layered architecture:** Separation of concerns (controllers, services, models)
- **Transaction safety:** All swap operations use database transactions to prevent race conditions
- **JWT authentication:** Stateless authentication for scalability
- **GORM AutoMigrate:** Automatic schema migration for rapid development

### Frontend Architecture
- **Composition API:** Modern Vue 3 approach for better code organization
- **Pinia stores:** Centralized state management with type safety
- **Route guards:** Protecting authenticated routes at router level
- **Axios interceptors:** Centralized error handling and token management

### Swap Logic
The core swap transaction follows these steps:
1. Lock both slot records (SELECT FOR UPDATE)
2. Verify both slots are in SWAP_PENDING status
3. Swap the user_id fields between the two slots
4. Update both slot statuses to BUSY
5. Update swap request status to ACCEPTED
6. Commit transaction (or rollback on any error)

This ensures atomic ownership transfer and prevents race conditions.

## Known Limitations

- No pagination on marketplace (can be slow with many slots)
- No recurring events support
- No timezone handling (uses UTC)
- No calendar integration (Google Calendar, Outlook)
- Polling for new requests (no WebSockets in base version)

## Future Enhancements

- [ ] Real-time notifications with WebSockets
- [ ] Email notifications
- [ ] Calendar view (FullCalendar integration)
- [ ] Search and filter marketplace
- [ ] Recurring events
- [ ] Timezone support
- [ ] Mobile app (React Native)

## Testing

### Backend Tests
```bash
cd backend
go test ./tests/... -v
```

### Frontend Tests
```bash
cd frontend
npm run test
```

## Deployment

### Backend (Render)
1. Push code to GitHub
2. Create new Web Service on Render
3. Connect repository
4. Add PostgreSQL database
5. Set environment variables
6. Deploy

### Frontend (Vercel)
1. Push code to GitHub
2. Import project on Vercel
3. Set environment variable: `VITE_API_BASE_URL`
4. Deploy

## Troubleshooting

**Database connection fails:**
- Ensure PostgreSQL is running
- Check credentials in .env file
- Verify database exists

**CORS errors:**
- Check ALLOWED_ORIGINS in backend .env
- Ensure frontend URL is whitelisted

**JWT token invalid:**
- Check JWT_SECRET is set and consistent
- Token may be expired (24h default)
- Clear localStorage and login again

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

MIT License - See LICENSE file for details

## Contact

Your Name - [your.email@example.com](mailto:your.email@example.com)

Project Link: [https://github.com/yourusername/slotswapper](https://github.com/yourusername/slotswapper)

Live Demo: [https://slotswapper.vercel.app](https://slotswapper.vercel.app)
```

---

## **FINAL CHECKLIST**

### **Must-Have Before Submission**

```
✓ Backend:
  ✓ All API endpoints working
  ✓ JWT authentication functional
  ✓ Swap transaction logic tested
  ✓ Database migrations applied
  ✓ CORS configured correctly
  ✓ Environment variables documented

✓ Frontend:
  ✓ All views implemented
  ✓ Authentication flow complete
  ✓ State management working
  ✓ API integration successful
  ✓ Responsive design
  ✓ Error handling implemented

✓ Documentation:
  ✓ Comprehensive README
  ✓ Setup instructions clear
  ✓ API endpoints documented
  ✓ Design decisions explained
  ✓ Known limitations listed

✓ Repository:
  ✓ Clean commit history
  ✓ .gitignore properly configured
  ✓ .env.example files included
  ✓ Code commented where necessary

✓ Optional (if time):
  ✓ Deployed and accessible
  ✓ Some tests written
  ✓ Docker setup included
```

---

**END OF TECHNICAL PLAN**

This comprehensive plan covers every technical detail needed to implement SlotSwapper in 24 hours using Go Gin and Vue. The plan prioritizes the core swap logic (the most critical evaluation criteria) while providing clear structure for rapid development. Follow the hour-by-hour timeline strictly to ensure completion within the deadline.