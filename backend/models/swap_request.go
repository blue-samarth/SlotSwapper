package models

import (
	"time"

	"gorm.io/gorm"
)

type SwapStatus string

const (
	StatusPending   SwapStatus = "PENDING"
	StatusAccepted  SwapStatus = "ACCEPTED"
	StatusRejected  SwapStatus = "REJECTED"
	StatusCancelled SwapStatus = "CANCELLED"
)

type SwapRequest struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	RequesterID      uint           `gorm:"not null;index" json:"requester_id"`
	ReceiverID       uint           `gorm:"not null;index" json:"receiver_id"`
	RequesterEventID uint           `gorm:"not null;index" json:"requester_event_id"`
	ReceiverEventID  uint           `gorm:"not null;index" json:"receiver_event_id"`
	Status           SwapStatus     `gorm:"type:VARCHAR(20);not null;default:'PENDING'" json:"status"`
	Requester        User           `gorm:"foreignKey:RequesterID" json:"-"`
	Receiver         User           `gorm:"foreignKey:ReceiverID" json:"-"`
	RequesterEvent   Event          `gorm:"foreignKey:RequesterEventID" json:"-"`
	ReceiverEvent    Event          `gorm:"foreignKey:ReceiverEventID" json:"-"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted        bool           `gorm:"default:false" json:"is_deleted"`
}

func (sr *SwapRequest) Validate(db *gorm.DB) error {
	var requesterEvent, receiverEvent Event
	if err := db.First(&requesterEvent, "id = ? AND is_deleted = false", sr.RequesterEventID).Error; err != nil {
		return err
	}
	if requesterEvent.UserID != sr.RequesterID {
		return gorm.ErrInvalidData
	}
	if err := db.First(&receiverEvent, "id = ? AND is_deleted = false", sr.ReceiverEventID).Error; err != nil {
		return err
	}
	if receiverEvent.UserID != sr.ReceiverID {
		return gorm.ErrInvalidData
	}

	if !requesterEvent.IsSwappable() || !receiverEvent.IsSwappable() {
		return gorm.ErrInvalidData
	}

	return nil
}

func (sr *SwapRequest) CanChangeStatus(newStatus SwapStatus) bool {
	switch sr.Status {
	case StatusPending:
		return newStatus == StatusAccepted || newStatus == StatusRejected || newStatus == StatusCancelled
	case StatusAccepted, StatusRejected, StatusCancelled:
		return false
	default:
		return false
	}
}
