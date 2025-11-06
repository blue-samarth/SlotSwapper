📋 Backend API Reference for Frontend Development
🌐 Base Configuration
Base URL
Headers Required
Rate Limits (Per IP)
Global: 100 requests/minute (all endpoints)
Auth: 5 requests/minute (signup/login)
Swap Operations: 20 requests/minute
Rate Limit Headers (returned by API)
🔐 Authentication Endpoints
1. Sign Up
Request:

Success Response (200):

Error Response (409 - Email exists):

2. Login
Request:

Success Response (200):

Error Response (401):

👤 User Endpoints (Protected)
3. Get Current User Profile
Success Response (200):

4. Update User Profile
Request:

Success Response (200):

📅 Event Endpoints (Protected)
5. Get User's Events
Success Response (200):

6. Create Event
Request:

Example:

Success Response (201):

Error Response (400 - Invalid times):

7. Update Event Status
Request:

Success Response (200):

Error Response (403 - Not owner):

8. Delete Event
Success Response (200):

Error Response (404):

🔄 Swap Endpoints (Protected + Rate Limited 20/min)
9. Get Swappable Slots (All Users)
Success Response (200):

10. Get My Swap Requests
Query Parameters:

filter (optional): sent, received, or all (default: all)
Success Response (200):

Swap Request Statuses:

PENDING - Awaiting receiver's response
ACCEPTED - Swap completed successfully
REJECTED - Receiver declined the swap
CANCELLED - Requester cancelled before acceptance
11. Initiate Swap Request
Request:

Success Response (201):

Error Response (422 - Event not swappable):

12. Accept Swap Request
Success Response (200):

Error Response (403 - Not receiver):

13. Reject Swap Request
Success Response (200):

14. Cancel Swap Request
Success Response (200):

Error Response (403 - Not requester):

🏥 Health Check Endpoints (Public)
15. Basic Health Check
Response (200):

16. Detailed Health Check
Response (200):

⚠️ Error Codes Reference
Error Code	HTTP Status	When It Occurs
VALIDATION_ERROR	400	Invalid input (missing fields, wrong format)
BAD_REQUEST	400	Malformed request body
AUTHENTICATION_ERROR	401	Invalid/missing token, wrong credentials
AUTHORIZATION_ERROR	403	User doesn't have permission
NOT_FOUND	404	Resource doesn't exist
CONFLICT	409	Duplicate resource (email exists)
SWAP_LOGIC_ERROR	422	Business logic violation (event not swappable)
RATE_LIMIT_EXCEEDED	429	Too many requests
INTERNAL_SERVER_ERROR	500	Server error
🎯 Frontend Implementation Tips
1. Store JWT Token
2. Axios Interceptor for Auth
3. Handle Rate Limiting
4. Error Handling
5. Event Status Colors
6. Swap Request Status Colors
📊 Data Models Summary
Event Statuses
BUSY - Not available for swapping
SWAPPABLE - Available for swap requests
SWAP_PENDING - Currently in a pending swap (auto-set)
Swap Request Statuses
PENDING - Awaiting decision
ACCEPTED - Swap completed
REJECTED - Receiver declined
CANCELLED - Requester cancelled
Key Validation Rules
Username: 3-50 characters
Email: Valid email format
Password: 6-100 characters
Event Title: 1-100 characters
End Time: Must be after start time
Start Time: Must be in the future (for new events)