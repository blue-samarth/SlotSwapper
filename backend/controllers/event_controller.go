package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"SlotSwapper/config"
	"SlotSwapper/models"
	"SlotSwapper/dto"
)

func GetEvents(c *gin.Context) {
	userID := c.GetUint("user_id")

	var events []models.Event
	db := config.GetDB()
	if err := db.Where("user_id = ? AND is_deleted = false", userID).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	var responseEvents []dto.EventResponse
	for _, event := range events {
		responseEvents = append(responseEvents, dto.EventResponse{
			ID:        event.ID,
			UserID:    event.UserID,
			Title:     event.Title,
			StartTime: event.StartTime,
			EndTime:   event.EndTime,
			Status:    event.Status,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, responseEvents)
}

func CreateEvent(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event := models.Event{
		Title:     req.Title,
		UserID:    userID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Status:    models.StatusBusy,
	}

	db := config.GetDB()
	if err := db.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, dto.EventResponse{
		ID:        event.ID,
		Title:     event.Title,
		UserID:    event.UserID,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		Status:    event.Status,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	})
}

func UpdateEventStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var req dto.UpdateEventStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	var event models.Event

	if err := db.Where("id = ? AND user_id = ? AND is_deleted = false", eventID, userID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch event"})
		return
	}

	if err := event.CanChangeStatus(req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.Status = req.Status
	event.UpdatedAt = time.Now()
	if err := db.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event status"})
		return
	}

	c.JSON(http.StatusOK, dto.EventResponse{
		ID:        event.ID,
		Title:     event.Title,
		UserID:    event.UserID,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		Status:    event.Status,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	})
}

func DeleteEvent(c *gin.Context) {
	// This function soft deletes an event by setting IsDeleted to true
	userID := c.GetUint("user_id")
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	db := config.GetDB()
	var event models.Event

	if err := db.Where("id = ? AND user_id = ? AND is_deleted = false", eventID, userID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch event"})
			return
		}
	}

    if event.Status == models.StatusSwapPending {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete event with pending swap"})
        return
    }
	event.IsDeleted = true
	event.UpdatedAt = time.Now()

	if err := db.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}