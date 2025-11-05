package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type EventStatus string
const (
	StatusBusy        EventStatus = "BUSY"
	StatusSwappable   EventStatus = "SWAPPABLE"
	StatusSwapPending EventStatus = "SWAP_PENDING"
)

type Event struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"type:VARCHAR(100);not null" json:"title"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	StartTime time.Time      `gorm:"not null;index" json:"start_time"`
	EndTime   time.Time      `gorm:"not null;index" json:"end_time"`
	Status    EventStatus    `gorm:"type:VARCHAR(20);not null;default:'BUSY'" json:"status"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (e *Event) BeforeSave(tx *gorm.DB) (err error) {
	if !e.EndTime.After(e.StartTime) {
		return errors.New("end_time must be after start_time")
	}
	if e.ID == 0 && e.StartTime.Before(time.Now()) {
		return errors.New("start_time must be in the future")
	}
	return nil
}

func (e *Event) IsSwappable() bool {return e.Status == StatusSwappable}

func (e *Event) CanChangeStatus(newStatus EventStatus) error {
	if newStatus == StatusSwapPending {
		return errors.New("cannot set status to SWAP_PENDING directly")
	}
	if e.Status == StatusSwapPending {
		return errors.New("cannot change status from SWAP_PENDING")
	}

	if (e.Status == StatusBusy && newStatus != StatusSwappable) ||
		(e.Status == StatusSwappable && newStatus != StatusBusy) {
		return errors.New("invalid status transition")
	}
	return nil
}
