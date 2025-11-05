package dto

import "time"

type InitiateSwapRequest struct {
	RequesterEventID uint `json:"requester_event_id" binding:"required"`
	ReceiverEventID  uint `json:"receiver_event_id" binding:"required"`
}

type SwapRequestResponse struct {
	ID               uint          `json:"id"`
	RequesterEventID uint          `json:"requester_event_id"`
	ReceiverEventID  uint          `json:"receiver_event_id"`
	RequesterID      uint          `json:"requester_id"`
	ReceiverID       uint          `json:"receiver_id"`
	Status           string        `json:"status"`
	RequesterEvent   EventResponse `json:"requester_event"`
	ReceiverEvent    EventResponse `json:"receiver_event"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type SwappableSlotResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	UserID    uint      `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	User      UserInfo  `json:"user"`
}
