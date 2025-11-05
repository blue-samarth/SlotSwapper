# SlotSwapper Backend Implementation Plan - Complete Enhanced Edition

## Overview
**Total Time: 14 hours** (extended from 12 to accommodate critical enhancements)  
**Stack: Go + Gin + PostgreSQL + GORM + JWT**  
**Critical Focus: Transactional swap logic with comprehensive race condition handling**

---

## Hour-by-Hour Breakdown (Updated)

### **Hours 0-1.5: Project Foundation & Enhanced Database Setup**

**Tasks:**

**1. Project Initialization (20 min)**
- Initialize Go module: `go mod init github.com/yourusername/slotswapper`
- Create directory structure:
  ```
  backend/
  ├── config/         # Database & environment config
  ├── controllers/    # HTTP handlers
  ├── middleware/     # Auth, CORS, rate limiting
  ├── models/         # Database models
  ├── routes/         # Route definitions
  ├── services/       # Business logic
  ├── utils/          # Helpers & transaction wrapper
  ├── dto/            # Data Transfer Objects
  ├── tests/          # Test files
  └── main.go
  ```

**2. Install Dependencies (10 min)**
```bash
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/gin-contrib/cors
go get -u github.com/joho/godotenv
go get -u github.com/rs/zerolog
go get -u github.com/ulule/limiter/v3

**3. Environment Configuration (15 min)**
```bash
# Create .env file
PORT=8080
GIN_MODE=debug
ENVIRONMENT=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=slotswapper_dev
DB_SSLMODE=disable
JWT_SECRET=your-dev-secret-min-32-chars-please-change
JWT_EXPIRATION_HOURS=24
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
LOG_LEVEL=debug
TRANSACTION_TIMEOUT=5
RATE_LIMIT_ENABLED=true
RATE_LIMIT_SWAP_REQUESTS=10
RATE_LIMIT_GENERAL_REQUESTS=100
```

**4. Database Connection with Connection Pooling (20 min)**
```
File: config/database.go

Setup:
- Connect to PostgreSQL using GORM
- Configure connection pool:
  * MaxIdleConns: 10
  * MaxOpenConns: 100
  * ConnMaxLifetime: 1 hour
- Enable prepared statement cache
- Use UTC timezone for consistency
- Test connection before proceeding

Key Configuration:
db.SetMaxIdleConns(10)
db.SetMaxOpenConns(100)
db.SetConnMaxLifetime(time.Hour)
```

**5. Transaction Wrapper Utility (25 min)**
```
File: utils/transaction.go

Create reusable transaction wrapper with:
- Context timeout (5 seconds default)
- Automatic rollback on panic
- PostgreSQL lock timeout setting
- Isolation level configuration (READ COMMITTED)
- Convenience wrapper for default options

Benefits:
- Reduces boilerplate by ~10 lines per transaction
- Consistent error handling
- Prevents deadlocks from hanging forever
- Centralized transaction management

Usage pattern:
return utils.WithDefaultTransaction(db, func(tx *gorm.DB) error {
    // Transaction logic here
    return nil
})
```

**Deliverables:**
- ✅ Project skeleton with all directories
- ✅ All dependencies installed
- ✅ Database connected with connection pooling
- ✅ Environment configuration loaded
- ✅ Transaction wrapper utility ready
- ✅ Basic logging configured

**Time-Saving Tips:**
- Copy dependency list and install all at once
- Use `createdb slotswapper_dev` to create database
- Test database connection with simple query before proceeding
- Keep database client (pgAdmin/TablePlus) open for verification

**Testing Checkpoint:**
```bash
# Verify setup
go run main.go
# Should see: "Database connected successfully"
# Should see: "Server starting on :8080"
```

---

### **Hours 1.5-3.5: Authentication System with Enhanced Security**

**Tasks:**

**1. User Model with Database Constraints (35 min)**
```
File: models/user.go

Struct Definition:
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null;size:255"`
    Email     string    `gorm:"unique;not null;index;size:255"`  // ← UNIQUE constraint
    Password  string    `gorm:"not null" json:"-"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

Key Features:
- UNIQUE constraint on email at database level
- Email is indexed for fast lookups
- Password never returned in JSON (json:"-" tag)
- BeforeSave hook to lowercase email
- HashPassword method (bcrypt cost 12)
- CheckPassword method for validation

Why Email UNIQUE at DB Level:
- Prevents race condition where two concurrent signups succeed
- Database guarantees atomicity
- Returns clear constraint violation error
- Application-level validation alone is insufficient

Implementation Details:
func (u *User) BeforeSave(tx *gorm.DB) error {
    u.Email = strings.ToLower(strings.TrimSpace(u.Email))
    return nil
}

func (u *User) HashPassword(password string) error {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }
    u.Password = string(hashed)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}

Run AutoMigrate:
db.AutoMigrate(&User{})
```

**2. JWT Utilities (30 min)**
```
File: utils/jwt.go

Token Structure:
Claims:
- user_id: uint
- email: string
- exp: 24 hours from issue
- iat: issued at timestamp

Functions:
1. GenerateToken(userID uint, email string) (string, error)
   - Create JWT with HS256 algorithm
   - Set expiration to 24 hours
   - Sign with SECRET from env

2. ValidateToken(tokenString string) (*Claims, error)
   - Parse JWT token
   - Verify signature
   - Check expiration
   - Return claims or error

Security Considerations:
- Use strong secret (min 32 characters)
- Never log tokens
- Short expiration for security
- Use Bearer token format
```

**3. Auth Controllers with Enhanced Error Handling (50 min)**
```
File: controllers/auth_controller.go

Signup Endpoint: POST /api/auth/signup
Request:
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepass123"
}

Logic:
1. Bind JSON to DTO with validation
2. Check if email already exists (fast check before hashing)
3. Hash password (expensive operation)
4. Create user record
5. Handle duplicate email error (409 Conflict)
6. Generate JWT token
7. Return token + user info (without password)

Response (201):
{
    "token": "eyJhbG...",
    "user": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "created_at": "2025-01-15T10:30:00Z"
    }
}

Error Responses:
- 400: Validation error (weak password, invalid email)
- 409: Email already registered (UNIQUE constraint)
- 500: Internal server error

Login Endpoint: POST /api/auth/login
Request:
{
    "email": "john@example.com",
    "password": "securepass123"
}

Logic:
1. Bind JSON to DTO
2. Find user by email (case-insensitive)
3. Check password with bcrypt
4. Generate JWT token
5. Return token + user info

Response (200): Same as signup

Error Responses:
- 401: Invalid credentials (don't specify which field)
- 500: Internal server error

Security Best Practices:
- Generic error message ("Invalid credentials")
- Rate limit login attempts (handled by middleware)
- Constant-time password comparison (bcrypt does this)
- No information leakage about user existence
```

**4. Auth Middleware (20 min)**
```
File: middleware/auth.go

Function: AuthMiddleware() gin.HandlerFunc

Execution Flow:
1. Extract "Authorization" header
2. Validate format: "Bearer <token>"
3. Extract token (strip "Bearer " prefix)
4. Call utils.ValidateToken(token)
5. If valid:
   - Set user_id in Gin context: c.Set("user_id", claims.UserID)
   - Set email in context: c.Set("email", claims.Email)
   - Call c.Next() to proceed
6. If invalid:
   - Return 401 Unauthorized
   - Call c.Abort() to stop execution

Error Cases:
- Missing Authorization header → 401
- Invalid format (not "Bearer X") → 401
- Expired token → 401
- Invalid signature → 401
- Malformed token → 401

Usage in routes:
protected := r.Group("/api")
protected.Use(middleware.AuthMiddleware())
```

**5. Route Setup for Auth (15 min)**
```
File: routes/routes.go

Setup:
func SetupRoutes(db *gorm.DB) *gin.Engine {
    r := gin.Default()
    
    // CORS middleware
    r.Use(cors.New(corsConfig))
    
    // Public routes
    auth := r.Group("/api/auth")
    {
        auth.POST("/signup", controllers.Signup)
        auth.POST("/login", controllers.Login)
    }
    
    return r
}
```

**Deliverables:**
- ✅ Users table with UNIQUE email constraint
- ✅ Secure password hashing (bcrypt cost 12)
- ✅ JWT generation and validation
- ✅ Signup with duplicate email prevention
- ✅ Login with secure credential checking
- ✅ Auth middleware protecting routes
- ✅ Proper error responses (400, 401, 409)

**Testing Checkpoints:**
```bash
# Test Signup
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@test.com","password":"test123"}'
# Expected: 201 with token

# Test Duplicate Email
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob","email":"alice@test.com","password":"test123"}'
# Expected: 409 Conflict

# Test Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","password":"test123"}'
# Expected: 200 with token

# Test Wrong Password
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","password":"wrong"}'
# Expected: 401

# Test Protected Route Without Token
curl http://localhost:8080/api/events
# Expected: 401
```

---

### **Hours 3.5-5.5: Events Management with Validation**

**Tasks:**

**1. Event Model with Constraints (35 min)**
```
File: models/event.go

Status Enum:
type EventStatus string

const (
    StatusBusy        EventStatus = "BUSY"
    StatusSwappable   EventStatus = "SWAPPABLE"
    StatusSwapPending EventStatus = "SWAP_PENDING"
)

Struct Definition:
type Event struct {
    ID        uint        `gorm:"primaryKey"`
    Title     string      `gorm:"not null;size:255"`
    StartTime time.Time   `gorm:"not null;index"`
    EndTime   time.Time   `gorm:"not null"`
    Status    EventStatus `gorm:"default:'BUSY';index;type:varchar(20)"`
    UserID    uint        `gorm:"not null;index"`
    User      User        `gorm:"foreignKey:UserID"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

Database Constraints:
1. CHECK constraint: end_time > start_time
   Add in migration:
   db.Exec(`
       ALTER TABLE events 
       ADD CONSTRAINT check_end_after_start 
       CHECK (end_time > start_time)
   `)

2. Composite index for marketplace query:
   db.Exec(`
       CREATE INDEX IF NOT EXISTS idx_events_status_user 
       ON events(status, user_id)
   `)

BeforeSave Hook (Application-level validation):
func (e *Event) BeforeSave(tx *gorm.DB) error {
    // Validate time range
    if !e.EndTime.After(e.StartTime) {
        return errors.New("end_time must be after start_time")
    }
    
    // Validate not in past (only for new records)
    if e.ID == 0 && e.StartTime.Before(time.Now()) {
        return errors.New("start_time must be in the future")
    }
    
    return nil
}

Why Both DB and Application Validation:
- DB constraint: Last line of defense, prevents bad data at storage
- GORM hook: Immediate feedback, better error messages
- Defense in depth approach

Allowed Status Transitions:
User can change:
- BUSY → SWAPPABLE (mark as available)
- SWAPPABLE → BUSY (unmark)

System changes (during swap):
- SWAPPABLE → SWAP_PENDING (on request created)
- SWAP_PENDING → BUSY (on accept)
- SWAP_PENDING → SWAPPABLE (on reject/cancel)

Run AutoMigrate:
db.AutoMigrate(&Event{})
```

**2. Event Controllers (70 min)**
```
File: controllers/event_controller.go

GET /api/events - List user's events
Logic:
1. Extract user_id from JWT context
2. Query events WHERE user_id = current_user
3. Order by start_time ASC
4. Return array of events

Response (200):
[
    {
        "id": 1,
        "title": "Team Meeting",
        "start_time": "2025-01-20T10:00:00Z",
        "end_time": "2025-01-20T11:00:00Z",
        "status": "SWAPPABLE",
        "user_id": 1,
        "created_at": "...",
        "updated_at": "..."
    }
]

POST /api/events - Create event
Request:
{
    "title": "Focus Block",
    "start_time": "2025-01-21T14:00:00Z",
    "end_time": "2025-01-21T15:00:00Z"
}

Logic:
1. Bind JSON to DTO with validation
2. Extract user_id from JWT context
3. Create event with default status BUSY
4. Auto-set user_id from token
5. Validate times (end > start, start in future)
6. Save to database
7. Return created event

Response (201): Created event object

Error Responses:
- 400: Validation error (invalid times, missing fields)
- 401: Unauthorized (invalid/missing token)
- 500: Internal server error

PATCH /api/events/:id/status - Update status
Request:
{
    "status": "SWAPPABLE"  // or "BUSY"
}

Logic:
1. Extract event ID from URL params
2. Extract user_id from JWT context
3. Find event by ID
4. Verify ownership (event.user_id == current_user_id)
5. Validate status transition:
   - Can only change BUSY ↔ SWAPPABLE
   - Cannot change SWAP_PENDING (system-managed)
6. Update status
7. Return updated event

Response (200): Updated event

Error Responses:
- 400: Invalid transition (trying to set SWAP_PENDING)
- 403: Not authorized (not your event)
- 404: Event not found
- 500: Internal server error

DELETE /api/events/:id - Delete event
Logic:
1. Extract event ID from URL params
2. Extract user_id from JWT context
3. Find event by ID
4. Verify ownership
5. Check status is not SWAP_PENDING
6. Delete event
7. Return 204 No Content

Error Responses:
- 400: Cannot delete SWAP_PENDING event
- 403: Not authorized
- 404: Event not found
- 500: Internal server error

Status Validation Helper:
func validateStatusTransition(currentStatus, newStatus EventStatus) error {
    // User cannot set SWAP_PENDING
    if newStatus == StatusSwapPending {
        return errors.New("cannot manually set SWAP_PENDING status")
    }
    
    // Can only transition BUSY ↔ SWAPPABLE
    if currentStatus == StatusSwapPending {
        return errors.New("cannot change status of pending swap")
    }
    
    // Valid transitions
    validTransitions := map[EventStatus][]EventStatus{
        StatusBusy:      {StatusSwappable},
        StatusSwappable: {StatusBusy},
    }
    
    allowed := validTransitions[currentStatus]
    for _, status := range allowed {
        if status == newStatus {
            return nil
        }
    }
    
    return errors.New("invalid status transition")
}
```

**3. Routes Configuration (15 min)**
```
File: routes/routes.go

CORS Configuration:
corsConfig := cors.Config{
    AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}

Protected Routes:
protected := r.Group("/api")
protected.Use(middleware.AuthMiddleware())
{
    // Events
    protected.GET("/events", controllers.GetEvents)
    protected.POST("/events", controllers.CreateEvent)
    protected.PATCH("/events/:id/status", controllers.UpdateEventStatus)
    protected.DELETE("/events/:id", controllers.DeleteEvent)
}
```

**Deliverables:**
- ✅ Events table with time validation constraint
- ✅ Status enum with proper transitions
- ✅ Full CRUD for events
- ✅ Ownership validation on all operations
- ✅ Status change restrictions enforced
- ✅ SWAP_PENDING status protected from user changes
- ✅ Database indexes for performance

**Testing Checkpoints:**
```bash
TOKEN="<token_from_login>"

# Create Event
curl -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Meeting","start_time":"2025-01-25T10:00:00Z","end_time":"2025-01-25T11:00:00Z"}'
# Expected: 201 with event

# Invalid Times (end before start)
curl -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Meeting","start_time":"2025-01-25T11:00:00Z","end_time":"2025-01-25T10:00:00Z"}'
# Expected: 400 error

# Get Events
curl http://localhost:8080/api/events \
  -H "Authorization: Bearer $TOKEN"
# Expected: 200 with array

# Mark as Swappable
curl -X PATCH http://localhost:8080/api/events/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"SWAPPABLE"}'
# Expected: 200 with updated event

# Try to set SWAP_PENDING manually
curl -X PATCH http://localhost:8080/api/events/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"SWAP_PENDING"}'
# Expected: 400 error
```

---

### **Hours 5.5-9: Complete Swap Logic with Race Condition Handling (CRITICAL)**

**⚠️ This is the most important section - 50% of evaluation weight!**

**Tasks:**

**1. SwapRequest Model with Enhanced Constraints (25 min)**
```
File: models/swap_request.go

Status Enum:
type SwapStatus string

const (
    SwapStatusPending   SwapStatus = "PENDING"
    SwapStatusAccepted  SwapStatus = "ACCEPTED"
    SwapStatusRejected  SwapStatus = "REJECTED"
    SwapStatusCancelled SwapStatus = "CANCELLED"  // NEW
)

Struct Definition:
type SwapRequest struct {
    ID          uint       `gorm:"primaryKey"`
    RequesterID uint       `gorm:"not null;index:idx_requester"`
    Requester   User       `gorm:"foreignKey:RequesterID"`
    ReceiverID  uint       `gorm:"not null;index:idx_receiver"`
    Receiver    User       `gorm:"foreignKey:ReceiverID"`
    MySlotID    uint       `gorm:"not null;index:idx_swap_slots"`
    MySlot      Event      `gorm:"foreignKey:MySlotID"`
    TheirSlotID uint       `gorm:"not null;index:idx_swap_slots"`
    TheirSlot   Event      `gorm:"foreignKey:TheirSlotID"`
    Status      SwapStatus `gorm:"default:'PENDING';index;type:varchar(20)"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

Optional: Partial Unique Index (prevents duplicate pending requests)
db.Exec(`
    CREATE UNIQUE INDEX IF NOT EXISTS idx_swap_unique_pending 
    ON swap_requests(requester_id, receiver_id, my_slot_id, their_slot_id) 
    WHERE status = 'PENDING'
`)

Benefits:
- Prevents same user requesting same swap twice
- Only enforced for PENDING status
- Allows historical records (ACCEPTED/REJECTED)

Run AutoMigrate:
db.AutoMigrate(&SwapRequest{})
```

**2. Marketplace Endpoint (20 min)**
```
File: controllers/swap_controller.go

GET /api/swappable-slots - Browse available slots

Logic:
1. Extract user_id from JWT context
2. Query events:
   WHERE status = 'SWAPPABLE' 
   AND user_id != current_user_id
3. Preload User relationship (avoid N+1)
4. Order by start_time ASC
5. Return array

Optimization:
- Use composite index (status, user_id)
- Select only needed fields
- Limit result size (pagination optional)

Response (200):
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

Query Implementation:
var slots []Event
err := db.Preload("User", func(db *gorm.DB) *gorm.DB {
    return db.Select("id", "name", "email")
}).
Where("status = ? AND user_id != ?", StatusSwappable, currentUserID).
Order("start_time ASC").
Find(&slots).Error
```

**3. Swap Service - Create Request (70 min)**
```
File: services/swap_service.go

Function: CreateSwapRequest(
    requesterID uint, 
    mySlotID uint, 
    theirSlotID uint
) (*SwapRequest, error)

Critical Requirements:
- MUST use database transaction
- MUST use SELECT FOR UPDATE for row-level locking
- MUST validate all constraints before committing
- MUST handle race conditions properly

Implementation Steps:

return utils.WithDefaultTransaction(s.db, func(tx *gorm.DB) error {
    // Step 1: Lock my_slot (requester's slot)
    var mySlot Event
    err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&mySlot, mySlotID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Your slot not found"}
    }
    
    // Step 2: Validate my_slot ownership
    if mySlot.UserID != requesterID {
        return &AppError{Code: 403, Message: "Not your slot"}
    }
    
    // Step 3: Validate my_slot is SWAPPABLE
    if mySlot.Status != StatusSwappable {
        return &AppError{Code: 400, Message: "Your slot is not swappable"}
    }
    
    // Step 4: Lock their_slot
    var theirSlot Event
    err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&theirSlot, theirSlotID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Their slot not found"}
    }
    
    // Step 5: Validate their_slot is SWAPPABLE
    if theirSlot.Status != StatusSwappable {
        return &AppError{Code: 400, Message: "Their slot is no longer available"}
    }
    
    // Step 6: Validate not swapping with yourself
    if theirSlot.UserID == requesterID {
        return &AppError{Code: 400, Message: "Cannot swap with yourself"}
    }
    
    // Step 7: Validate not same slot
    if mySlotID == theirSlotID {
        return &AppError{Code: 400, Message: "Cannot swap same slot"}
    }
    
    // Step 8: Create swap request
    swapReq := &SwapRequest{
        RequesterID: requesterID,
        ReceiverID:  theirSlot.UserID,
        MySlotID:    mySlotID,
        TheirSlotID: theirSlotID,
        Status:      SwapStatusPending,
    }
    
    if err := tx.Create(swapReq).Error; err != nil {
        // Check for unique constraint violation
        if strings.Contains(err.Error(), "idx_swap_unique_pending") {
            return &AppError{Code: 409, Message: "Swap request already exists"}
        }
        return err
    }
    
    // Step 9: Update my_slot status to SWAP_PENDING
    if err := tx.Model(&mySlot).Update("status", StatusSwapPending).Error; err != nil {
        return err
    }
    
    // Step 10: Update their_slot status to SWAP_PENDING
    if err := tx.Model(&theirSlot).Update("status", StatusSwapPending).Error; err != nil {
        return err
    }
    
    // Transaction commits automatically if no error
    return nil
})

Race Condition Handling:
- SELECT FOR UPDATE locks rows until transaction commits
- First transaction acquires lock, second waits
- Second transaction sees updated status (SWAP_PENDING)
- Second transaction fails validation and rolls back
- Only one swap request succeeds

Error Cases Handled:
- 404: Slot not found
- 403: Not slot owner
- 400: Slot not swappable
- 400: Cannot swap with yourself
- 400: Same slot
- 409: Duplicate request (unique constraint)
- 500: Transaction failure

Logging:
log.Info().
    Uint("requester_id", requesterID).
    Uint("my_slot_id", mySlotID).
    Uint("their_slot_id", theirSlotID).
    Msg("Creating swap request")
```

**4. Swap Service - Accept Request (75 min)**
```
Function: AcceptSwapRequest(
    requestID uint, 
    receiverID uint
) error

Critical Requirements:
- MUST be atomic (all-or-nothing)
- MUST swap ownership correctly
- MUST handle concurrent accept/reject
- MUST prevent double-processing

Implementation Steps:

return utils.WithDefaultTransaction(s.db, func(tx *gorm.DB) error {
    // Step 1: Lock swap request
    var swapReq SwapRequest
    err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&swapReq, requestID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Swap request not found"}
    }
    
    // Step 2: Verify receiver (authorization)
    if swapReq.ReceiverID != receiverID {
        return &AppError{Code: 403, Message: "Not authorized"}
    }
    
    // Step 3: CRITICAL - Re-check status after acquiring lock
    // This prevents simultaneous accept/reject race condition
    if swapReq.Status != SwapStatusPending {
        return &AppError{
            Code: 400, 
            Message: fmt.Sprintf("Request already %s", 
                strings.ToLower(string(swapReq.Status))),
        }
    }
    
    // Step 4: Lock both slots
    var mySlot Event
    err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&mySlot, swapReq.MySlotID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Requester's slot not found"}
    }
    
    var theirSlot Event
    err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&theirSlot, swapReq.TheirSlotID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Your slot not found"}
    }
    
    // Step 5: Verify both slots are still SWAP_PENDING
    if mySlot.Status != StatusSwapPending {
        return &AppError{Code: 400, Message: "Requester's slot no longer pending"}
    }
    if theirSlot.Status != StatusSwapPending {
        return &AppError{Code: 400, Message: "Your slot no longer pending"}
    }
    
    // Step 6: CRITICAL - Swap ownership atomically
    // Store original owners
    originalMySlotOwner := mySlot.UserID
    originalTheirSlotOwner := theirSlot.UserID
    
    // Perform swap
    mySlot.UserID = originalTheirSlotOwner
    theirSlot.UserID = originalMySlotOwner
    
    // Step 7: Update both slots (ownership + status to BUSY)
    updates := map[string]interface{}{
        "user_id": mySlot.UserID,
        "status":  StatusBusy,
    }
    if err := tx.Model(&mySlot).Updates(updates).Error; err != nil {
        return err
    }
    
```
    updates["user_id"] = theirSlot.UserID
    if err := tx.Model(&theirSlot).Updates(updates).Error; err != nil {
        return err
    }
    
    // Step 8: Update swap request status to ACCEPTED
    if err := tx.Model(&swapReq).Update("status", SwapStatusAccepted).Error; err != nil {
        return err
    }
    
    // Transaction commits automatically if no error
    log.Info().
        Uint("request_id", requestID).
        Uint("requester_id", swapReq.RequesterID).
        Uint("receiver_id", receiverID).
        Msg("Swap request accepted successfully")
    
    return nil
})

Key Points:
- Status re-check AFTER lock prevents double-acceptance
- Ownership swap must be atomic (both or neither)
- Both slots become BUSY (not back to SWAPPABLE)
- Transaction ensures all updates succeed or all rollback

Race Condition Protection:
Scenario: User clicks "Accept" twice (double-click)
1. First request acquires lock on swap_request
2. Second request waits for lock
3. First request changes status to ACCEPTED and commits
4. Second request acquires lock, sees status != PENDING
5. Second request returns 400 error "Request already accepted"
6. No duplicate processing occurs

Scenario: Simultaneous Accept + Reject
1. Accept acquires lock on swap_request first
2. Reject waits for lock
3. Accept changes status to ACCEPTED and commits
4. Reject acquires lock, sees status != PENDING
5. Reject returns 400 error "Request already accepted"
6. Only accept takes effect
```

**5. Swap Service - Reject Request (40 min)**
```
Function: RejectSwapRequest(
    requestID uint, 
    receiverID uint
) error

Implementation Steps:

return utils.WithDefaultTransaction(s.db, func(tx *gorm.DB) error {
    // Step 1: Lock swap request
    var swapReq SwapRequest
    err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&swapReq, requestID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Swap request not found"}
    }
    
    // Step 2: Verify receiver authorization
    if swapReq.ReceiverID != receiverID {
        return &AppError{Code: 403, Message: "Not authorized"}
    }
    
    // Step 3: CRITICAL - Re-check status after lock
    if swapReq.Status != SwapStatusPending {
        return &AppError{
            Code: 400, 
            Message: fmt.Sprintf("Request already %s", 
                strings.ToLower(string(swapReq.Status))),
        }
    }
    
    // Step 4: Lock both slots
    var mySlot Event
    err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&mySlot, swapReq.MySlotID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Slot not found"}
    }
    
    var theirSlot Event
    err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&theirSlot, swapReq.TheirSlotID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Slot not found"}
    }
    
    // Step 5: Restore both slots to SWAPPABLE
    if err := tx.Model(&mySlot).Update("status", StatusSwappable).Error; err != nil {
        return err
    }
    
    if err := tx.Model(&theirSlot).Update("status", StatusSwappable).Error; err != nil {
        return err
    }
    
    // Step 6: Update swap request status to REJECTED
    if err := tx.Model(&swapReq).Update("status", SwapStatusRejected).Error; err != nil {
        return err
    }
    
    log.Info().
        Uint("request_id", requestID).
        Uint("receiver_id", receiverID).
        Msg("Swap request rejected")
    
    return nil
})

Key Points:
- Both slots return to SWAPPABLE (available again)
- Ownership does NOT change
- Status check prevents race conditions
- Transaction ensures consistency
```

**6. Swap Service - Cancel Request (NEW - 40 min)**
```
Function: CancelSwapRequest(
    requestID uint, 
    requesterID uint  // NOTE: Requester, not receiver!
) error

Business Rules:
- Only REQUESTER can cancel (not receiver)
- Only PENDING requests can be cancelled
- Both slots restore to SWAPPABLE
- Cannot cancel ACCEPTED/REJECTED requests

Implementation Steps:

return utils.WithDefaultTransaction(s.db, func(tx *gorm.DB) error {
    // Step 1: Lock swap request
    var swapReq SwapRequest
    err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&swapReq, requestID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Swap request not found"}
    }
    
    // Step 2: Verify REQUESTER (not receiver!)
    if swapReq.RequesterID != requesterID {
        return &AppError{Code: 403, Message: "Only requester can cancel"}
    }
    
    // Step 3: Verify status is PENDING
    if swapReq.Status != SwapStatusPending {
        return &AppError{
            Code: 400, 
            Message: fmt.Sprintf("Cannot cancel %s request", 
                strings.ToLower(string(swapReq.Status))),
        }
    }
    
    // Step 4: Lock both slots
    var mySlot Event
    err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&mySlot, swapReq.MySlotID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Slot not found"}
    }
    
    var theirSlot Event
    err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&theirSlot, swapReq.TheirSlotID).Error
    if err != nil {
        return &AppError{Code: 404, Message: "Slot not found"}
    }
    
    // Step 5: Restore both slots to SWAPPABLE
    if err := tx.Model(&mySlot).Update("status", StatusSwappable).Error; err != nil {
        return err
    }
    
    if err := tx.Model(&theirSlot).Update("status", StatusSwappable).Error; err != nil {
        return err
    }
    
    // Step 6: Update swap request status to CANCELLED
    if err := tx.Model(&swapReq).Update("status", SwapStatusCancelled).Error; err != nil {
        return err
    }
    
    log.Info().
        Uint("request_id", requestID).
        Uint("requester_id", requesterID).
        Msg("Swap request cancelled by requester")
    
    return nil
})

Difference from Reject:
- Cancel: Called by REQUESTER
- Reject: Called by RECEIVER
- Both restore slots to SWAPPABLE
- Different status (CANCELLED vs REJECTED)
```

**7. Swap Controllers (40 min)**
```
File: controllers/swap_controller.go

POST /api/swap-request - Create swap request

Handler:
func CreateSwapRequest(c *gin.Context) {
    // Extract user_id from JWT context
    userID, _ := c.Get("user_id")
    requesterID := userID.(uint)
    
    // Bind request body
    var req dto.CreateSwapRequestDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Call service
    swapReq, err := swapService.CreateSwapRequest(
        requesterID, 
        req.MySlotID, 
        req.TheirSlotID,
    )
    
    if err != nil {
        handleError(c, err)
        return
    }
    
    c.JSON(201, swapReq)
}

POST /api/swap-response/:id - Accept or reject

Request Body:
{
    "accept": true  // or false
}

Handler:
func RespondToSwapRequest(c *gin.Context) {
    // Extract receiver_id from JWT
    userID, _ := c.Get("user_id")
    receiverID := userID.(uint)
    
    // Extract request ID from URL
    requestID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    
    // Bind request body
    var req dto.SwapResponseDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    var err error
    if req.Accept {
        err = swapService.AcceptSwapRequest(uint(requestID), receiverID)
    } else {
        err = swapService.RejectSwapRequest(uint(requestID), receiverID)
    }
    
    if err != nil {
        handleError(c, err)
        return
    }
    
    message := "accepted"
    if !req.Accept {
        message = "rejected"
    }
    
    c.JSON(200, gin.H{"message": fmt.Sprintf("Swap request %s successfully", message)})
}

DELETE /api/swap-request/:id - Cancel request (NEW)

Handler:
func CancelSwapRequest(c *gin.Context) {
    // Extract requester_id from JWT
    userID, _ := c.Get("user_id")
    requesterID := userID.(uint)
    
    // Extract request ID from URL
    requestID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    
    // Call service
    err := swapService.CancelSwapRequest(uint(requestID), requesterID)
    
    if err != nil {
        handleError(c, err)
        return
    }
    
    c.JSON(200, gin.H{"message": "Swap request cancelled successfully"})
}

GET /api/swap-requests/incoming - View incoming requests

Handler:
func GetIncomingRequests(c *gin.Context) {
    userID, _ := c.Get("user_id")
    receiverID := userID.(uint)
    
    var requests []SwapRequest
    err := db.Preload("Requester").
        Preload("MySlot").
        Preload("TheirSlot").
        Where("receiver_id = ?", receiverID).
        Order("created_at DESC").
        Find(&requests).Error
    
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch requests"})
        return
    }
    
    c.JSON(200, requests)
}

GET /api/swap-requests/outgoing - View outgoing requests

Handler:
func GetOutgoingRequests(c *gin.Context) {
    userID, _ := c.Get("user_id")
    requesterID := userID.(uint)
    
    var requests []SwapRequest
    err := db.Preload("Receiver").
        Preload("MySlot").
        Preload("TheirSlot").
        Where("requester_id = ?", requesterID).
        Order("created_at DESC").
        Find(&requests).Error
    
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch requests"})
        return
    }
    
    c.JSON(200, requests)
}

Error Handler Helper:
func handleError(c *gin.Context, err error) {
    if appErr, ok := err.(*AppError); ok {
        c.JSON(appErr.Code, gin.H{"error": appErr.Message})
    } else {
        log.Error().Err(err).Msg("Internal error")
        c.JSON(500, gin.H{"error": "Internal server error"})
    }
}
```

**8. Update Routes (10 min)**
```
File: routes/routes.go

Add to protected group:
protected := r.Group("/api")
protected.Use(middleware.AuthMiddleware())
{
    // Existing event routes...
    
    // Swap routes
    protected.GET("/swappable-slots", controllers.GetSwappableSlots)
    protected.POST("/swap-request", controllers.CreateSwapRequest)
    protected.POST("/swap-response/:id", controllers.RespondToSwapRequest)
    protected.DELETE("/swap-request/:id", controllers.CancelSwapRequest)
    protected.GET("/swap-requests/incoming", controllers.GetIncomingRequests)
    protected.GET("/swap-requests/outgoing", controllers.GetOutgoingRequests)
}
```

**Deliverables:**
- ✅ Complete swap request creation with row-level locking
- ✅ Accept flow with atomic ownership swap
- ✅ Reject flow with status restoration
- ✅ Cancel flow for requesters (NEW)
- ✅ Race condition protection on all operations
- ✅ Concurrent accept/reject handled correctly
- ✅ Transaction timeout prevents deadlocks
- ✅ All error cases handled gracefully
- ✅ Comprehensive logging for debugging

**Testing Checkpoints (CRITICAL - Test Thoroughly!):**

```bash
# Setup: Create 2 users and 2 swappable events
USER1_TOKEN="<token1>"
USER2_TOKEN="<token2>"

# 1. Create swap request
curl -X POST http://localhost:8080/api/swap-request \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"my_slot_id":1,"their_slot_id":2}'
# Expected: 201, both slots become SWAP_PENDING

# 2. Verify slots are SWAP_PENDING
curl http://localhost:8080/api/events \
  -H "Authorization: Bearer $USER1_TOKEN"
# Expected: Slot 1 status is SWAP_PENDING

# 3. Accept swap request
curl -X POST http://localhost:8080/api/swap-response/1 \
  -H "Authorization: Bearer $USER2_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"accept":true}'
# Expected: 200, ownership swapped

# 4. Verify ownership swapped
curl http://localhost:8080/api/events \
  -H "Authorization: Bearer $USER1_TOKEN"
# Expected: User1 now owns slot 2 (was slot 1)

curl http://localhost:8080/api/events \
  -H "Authorization: Bearer $USER2_TOKEN"
# Expected: User2 now owns slot 1 (was slot 2)

# 5. Test reject flow (create new request first)
curl -X POST http://localhost:8080/api/swap-response/2 \
  -H "Authorization: Bearer $USER2_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"accept":false}'
# Expected: 200, slots back to SWAPPABLE

# 6. Test cancel flow (create new request first)
curl -X DELETE http://localhost:8080/api/swap-request/3 \
  -H "Authorization: Bearer $USER1_TOKEN"
# Expected: 200, slots back to SWAPPABLE

# 7. Test error: Try to accept already-accepted request
curl -X POST http://localhost:8080/api/swap-response/1 \
  -H "Authorization: Bearer $USER2_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"accept":true}'
# Expected: 400 "Request already accepted"

# 8. Test error: Receiver tries to cancel
curl -X DELETE http://localhost:8080/api/swap-request/3 \
  -H "Authorization: Bearer $USER2_TOKEN"
# Expected: 403 "Only requester can cancel"

# 9. Test error: Request swap with non-swappable slot
# (First mark slot 1 as BUSY)
curl -X POST http://localhost:8080/api/swap-request \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"my_slot_id":1,"their_slot_id":2}'
# Expected: 400 "Your slot is not swappable"

# 10. RACE CONDITION TEST: Simultaneous requests for same slot
# Open 2 terminal windows, execute simultaneously:
# Terminal 1:
curl -X POST http://localhost:8080/api/swap-request \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"my_slot_id":1,"their_slot_id":3}'

# Terminal 2 (execute immediately after):
curl -X POST http://localhost:8080/api/swap-request \
  -H "Authorization: Bearer $USER2_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"my_slot_id":2,"their_slot_id":3}'
# Expected: One succeeds (201), one fails (400 "Their slot is no longer available")

# 11. DOUBLE-CLICK TEST: Accept twice rapidly
# Execute this command twice in quick succession:
curl -X POST http://localhost:8080/api/swap-response/4 \
  -H "Authorization: Bearer $USER2_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"accept":true}'
# Expected: First succeeds (200), second fails (400 "Request already accepted")
```

---

### **Hours 9-11: Enhanced Testing & Rate Limiting (NEW)**

**Tasks:**

**1. Comprehensive Manual Testing (75 min)**
```
Test Categories:

A. Authentication (15 min)
✓ Signup with valid data
✓ Signup with duplicate email (409)
✓ Signup with weak password (400)
✓ Login with correct credentials
✓ Login with wrong password (401)
✓ Login with non-existent email (401)
✓ Access protected route without token (401)
✓ Access protected route with expired token (401)
✓ Access protected route with invalid token (401)

B. Events Management (15 min)
✓ Create event with valid data
✓ Create event with end_time < start_time (400)
✓ Create event with end_time = start_time (400)
✓ Get all events for user
✓ Mark event as SWAPPABLE
✓ Mark SWAPPABLE event as BUSY
✓ Try to mark SWAP_PENDING as BUSY (400)
✓ Try to update another user's event (403)
✓ Delete event (not in swap)
✓ Try to delete SWAP_PENDING event (400)

C. Swap Creation (20 min)
✓ Create valid swap request
✓ Verify both slots become SWAP_PENDING
✓ Try to request swap with own slot (400)
✓ Try to request swap with same slot (400)
✓ Try to request swap with non-swappable slot (400)
✓ Try to create duplicate swap request (409)
✓ Race condition: Two users request same slot simultaneously
✓ Verify marketplace doesn't show SWAP_PENDING slots

D. Swap Response (20 min)
✓ Accept swap request (verify ownership swap)
✓ Reject swap request (verify slots become SWAPPABLE)
✓ Cancel swap request (verify slots become SWAPPABLE)
✓ Try to cancel as receiver (403)
✓ Try to accept as requester (403)
✓ Try to respond to already-processed request (400)
✓ Double-click accept (idempotency test)
✓ Simultaneous accept + reject
✓ Verify historical requests (incoming/outgoing)

E. Edge Cases (5 min)
✓ Delete user with pending swaps (cascade behavior)
✓ Very long event titles
✓ Events with same start/end times as existing
✓ Multiple concurrent operations

Testing Tools:
- Postman/Insomnia collection
- Bash scripts for race condition testing
- Database client for state verification
```

**2. Rate Limiting Implementation (30 min)**
```
File: middleware/rate_limit.go

Install dependency:
go get github.com/ulule/limiter/v3
go get github.com/ulule/limiter/v3/drivers/store/memory

Implementation:

package middleware

import (
    "fmt"
    "strconv"
    "strings"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/store/memory"
)

var (
    // General API rate limit: 100 requests per minute
    generalStore   = memory.NewStore()
    generalLimiter = limiter.New(generalStore, limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  100,
    })
    
    // Swap operations: 10 requests per minute (more restrictive)
    swapStore   = memory.NewStore()
    swapLimiter = limiter.New(swapStore, limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  10,
    })
    
    // Auth operations: 5 login attempts per minute
    authStore   = memory.NewStore()
    authLimiter = limiter.New(authStore, limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  5,
    })
)

func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Choose limiter based on route
        var limit *limiter.Limiter
        path := c.Request.URL.Path
        
        if strings.Contains(path, "/auth/login") {
            limit = authLimiter
        } else if strings.Contains(path, "/swap") {
            limit = swapLimiter
        } else {
            limit = generalLimiter
        }
        
        // Use user ID if authenticated, IP otherwise
        key := c.ClientIP()
        if userID, exists := c.Get("user_id"); exists {
            key = fmt.Sprintf("user:%v", userID)
        }
        
        // Get rate limit context
        context, err := limit.Get(c.Request.Context(), key)
        if err != nil {
            log.Error().Err(err).Msg("Rate limiter error")
            c.JSON(500, gin.H{"error": "Rate limiter error"})
            c.Abort()
            return
        }
        
        // Set rate limit headers
        c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
        c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
        c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))
        
        // Check if limit exceeded
        if context.Reached {
            retryAfter := time.Until(time.Unix(context.Reset, 0))
            c.Header("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
            
            c.JSON(429, gin.H{
                "error": "Rate limit exceeded. Please try again later.",
                "retry_after_seconds": int(retryAfter.Seconds()),
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

Apply to routes:
File: routes/routes.go

// Apply globally
r.Use(middleware.RateLimitMiddleware())

// Or apply selectively
auth := r.Group("/api/auth")
auth.Use(middleware.RateLimitMiddleware())

protected := r.Group("/api")
protected.Use(middleware.AuthMiddleware())
protected.Use(middleware.RateLimitMiddleware())
```

**3. Error Handling Refinement (25 min)**
```
File: utils/errors.go

Custom Error Type:
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

Predefined Errors:
var (
    ErrNotFound          = &AppError{Code: 404, Message: "Resource not found"}
    ErrUnauthorized      = &AppError{Code: 401, Message: "Unauthorized"}
    ErrForbidden         = &AppError{Code: 403, Message: "Forbidden"}
    ErrBadRequest        = &AppError{Code: 400, Message: "Bad request"}
    ErrConflict          = &AppError{Code: 409, Message: "Resource conflict"}
    ErrInternalServer    = &AppError{Code: 500, Message: "Internal server error"}
    ErrSlotNotSwappable  = &AppError{Code: 400, Message: "Slot is not swappable"}
    ErrRequestProcessed  = &AppError{Code: 400, Message: "Request already processed"}
    ErrRateLimitExceeded = &AppError{Code: 429, Message: "Rate limit exceeded"}
)

Global Error Handler Middleware:
File: middleware/error_handler.go

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            switch e := err.(type) {
            case *AppError:
                c.JSON(e.Code, gin.H{"error": e.Message})
            case validator.ValidationErrors:
                c.JSON(400, gin.H{"error": formatValidationErrors(e)})
            default:
                log.Error().Err(err).Msg("Unhandled error")
                c.JSON(500, gin.H{"error": "Internal server error"})
            }
        }
    }
}

func formatValidationErrors(errs validator.ValidationErrors) string {
    var messages []string
    for _, err := range errs {
        messages = append(messages, fmt.Sprintf("%s: %s", err.Field(), err.Tag()))
    }
    return strings.Join(messages, "; ")
}
```

**4. Enhanced Logging (20 min)**
```
File: middleware/logger.go

Request Logging Middleware:
func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method
        
        // Process request
        c.Next()
        
        duration := time.Since(start)
        status := c.Writer.Status()
        
        logEvent := log.Info()
        if status >= 400 {
            logEvent = log.Error()
        }
        
        logEvent.
            Str("method", method).
            Str("path", path).
            Int("status", status).
            Dur("duration", duration).
            Str("client_ip", c.ClientIP()).
            Msg("Request completed")
        
        // Warn on slow requests
        if duration > time.Second {
            log.Warn().
                Str("path", path).
                Dur("duration", duration).
                Msg("Slow request detected")
        }
    }
}

Transaction Logging (in swap service):
log.Info().
    Uint("request_id", requestID).
    Uint("requester_id", requesterID).
    Uint("receiver_id", receiverID).
    Str("operation", "accept_swap").
    Msg("Starting swap acceptance")

// After successful commit
log.Info().
    Uint("request_id", requestID).
    Dur("duration", time.Since(start)).
    Msg("Swap accepted successfully")

// On error
log.Error().
    Err(err).
    Uint("request_id", requestID).
    Msg("Failed to accept swap")
```

**Deliverables:**
- ✅ All endpoints tested with 2-3 user accounts
- ✅ Race conditions tested and verified
- ✅ Error responses consistent and helpful
- ✅ Rate limiting prevents abuse
- ✅ Comprehensive logging for debugging
- ✅ Database state verified after operations

---

### **Hours 11-12: Polish, Documentation & Deployment Prep**

**Tasks:**

**1. Code Cleanup (20 min)**
```
Checklist:
✓ Remove all commented code
✓ Remove debug print statements
✓ Add comments to complex logic:
  - Transaction wrapper usage
  - Swap ownership logic
  - Race condition handling
  - Lock ordering
✓ Ensure consistent naming (camelCase for Go)
✓ Run `go fmt ./...` on all files
✓ Run `go vet ./...` to catch issues
✓ Check for unused imports
✓ Verify all error messages are user-friendly
✓ Remove hardcoded values (use env vars)
```

**2. Environment Configuration (15 min)**
```
Create .env.example:
# Server Configuration
PORT=8080
GIN_MODE=debug
ENVIRONMENT=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=slotswapper_dev
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=change-me-to-random-32-char-string-in-production
JWT_EXPIRATION_HOURS=24

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Logging
LOG_LEVEL=debug

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_SWAP_REQUESTS=10
RATE_LIMIT_GENERAL_REQUESTS=100
RATE_LIMIT_AUTH_REQUESTS=5

# Transaction Configuration
TRANSACTION_TIMEOUT=5

# Database Connection Pool
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=3600

Document in README:
- Which variables are required
- Which have defaults
- Security considerations for production
- How to generate secure JWT_SECRET
```

**3. Health Check Endpoint (10 min)**
```
File: controllers/health_controller.go

func HealthCheck(c *gin.Context) {
    // Check database connection
    sqlDB, err := db.DB()
    if err != nil {
        c.JSON(503, gin.H{
            "status": "unhealthy",
            "database": "disconnected",
            "timestamp": time.Now().Unix(),
        })
        return
    }
    
    if err := sqlDB.Ping(); err != nil {
        c.JSON(503, gin.H{
            "status": "unhealthy",
            "database": "unreachable",
            "timestamp": time.Now().Unix(),
        })
        return
    }
    
    c.JSON(200, gin.H{
        "status": "healthy",
        "database": "connected",
        "timestamp": time.Now().Unix(),
        "version": "1.0.0",
    })
}

Add to routes:
r.GET("/health", controllers.HealthCheck)
r.GET("/api/health", controllers.HealthCheck)
```
