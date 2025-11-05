package dto

import (
	"time"

	"SlotSwapper/models"
)

type CreateEventRequest struct {
	Title     string    `json:"title" binding:"required,min=1,max=100"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

type UpdateEventStatusRequest struct {
	Status models.EventStatus `json:"status" binding:"required,oneof=BUSY SWAPPABLE"`
}

type EventResponse struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	UserID    uint               `json:"user_id"`
	StartTime time.Time          `json:"start_time"`
	EndTime   time.Time          `json:"end_time"`
	Status    models.EventStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
