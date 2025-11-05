package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "SlotSwapper/config"
    "SlotSwapper/dto"
    "SlotSwapper/models"
    "SlotSwapper/services"
)

func GetSwappableSlots(c *gin.Context) {
	userID := c.GetUint("user_id")

	var events []models.Event
	db := config.GetDB()
	if err := db.Where("status = ? AND user_id != ? AND is_deleted = false", models.StatusSwappable, userID).
		Preload("User").
		Order("start_time ASC").
		Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch swappable slots"})
		return
	}

	var responseEvents []dto.SwappableSlotResponse
    for _, event := range events {
        responseEvents = append(responseEvents, dto.SwappableSlotResponse{
            ID:        event.ID,
            Title:     event.Title,
            UserID:    event.UserID,
            StartTime: event.StartTime,
            EndTime:   event.EndTime,
            Status:    string(event.Status),
            User: dto.UserInfo{
                ID:       event.User.ID,
                Username: event.User.Username,
                Email:    event.User.Email,
            },
        })
    }
	c.JSON(http.StatusOK, responseEvents)
}

func InitiateSwap(c *gin.Context) {
    userID := c.GetUint("user_id")

    var req dto.InitiateSwapRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    swapRequest, err := services.InitiateSwap(req.RequesterEventID, req.ReceiverEventID, userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response := dto.SwapRequestResponse{
        ID:               swapRequest.ID,
        RequesterEventID: swapRequest.RequesterEventID,
        ReceiverEventID:  swapRequest.ReceiverEventID,
        RequesterID:      swapRequest.RequesterID,
        ReceiverID:       swapRequest.ReceiverID,
        Status:           string(swapRequest.Status),
		RequesterEvent: dto.EventResponse{
            ID:        swapRequest.RequesterEvent.ID,
            Title:     swapRequest.RequesterEvent.Title,
            UserID:    swapRequest.RequesterEvent.UserID,
            StartTime: swapRequest.RequesterEvent.StartTime,
            EndTime:   swapRequest.RequesterEvent.EndTime,
            Status:    swapRequest.RequesterEvent.Status,
            CreatedAt: swapRequest.RequesterEvent.CreatedAt,
            UpdatedAt: swapRequest.RequesterEvent.UpdatedAt,
        },
        ReceiverEvent: dto.EventResponse{
            ID:        swapRequest.ReceiverEvent.ID,
            Title:     swapRequest.ReceiverEvent.Title,
            UserID:    swapRequest.ReceiverEvent.UserID,
            StartTime: swapRequest.ReceiverEvent.StartTime,
            EndTime:   swapRequest.ReceiverEvent.EndTime,
            Status:    swapRequest.ReceiverEvent.Status,
            CreatedAt: swapRequest.ReceiverEvent.CreatedAt,
            UpdatedAt: swapRequest.ReceiverEvent.UpdatedAt,
        },
        CreatedAt: swapRequest.CreatedAt,
        UpdatedAt: swapRequest.UpdatedAt,
    }
	c.JSON(http.StatusCreated, response)
}

func GetMySwapRequests(c *gin.Context) {
    userID := c.GetUint("user_id")

    var swapRequests []models.SwapRequest
    db := config.GetDB()

    if err := db.Where("(requester_id = ? OR receiver_id = ?) AND is_deleted = false", userID, userID).
        Preload("RequesterEvent").
        Preload("ReceiverEvent").
        Preload("Requester").
        Preload("Receiver").
        Order("created_at DESC").
        Find(&swapRequests).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch swap requests"})
        return
    }

    var response []dto.SwapRequestResponse
    for _, sr := range swapRequests {
        response = append(response, dto.SwapRequestResponse{
            ID:               sr.ID,
            RequesterEventID: sr.RequesterEventID,
            ReceiverEventID:  sr.ReceiverEventID,
			RequesterID:      sr.RequesterID,
            ReceiverID:       sr.ReceiverID,
            Status:           string(sr.Status),
            RequesterEvent: dto.EventResponse{
                ID:        sr.RequesterEvent.ID,
                Title:     sr.RequesterEvent.Title,
                UserID:    sr.RequesterEvent.UserID,
                StartTime: sr.RequesterEvent.StartTime,
                EndTime:   sr.RequesterEvent.EndTime,
                Status:    sr.RequesterEvent.Status,
                CreatedAt: sr.RequesterEvent.CreatedAt,
                UpdatedAt: sr.RequesterEvent.UpdatedAt,
            },
            ReceiverEvent: dto.EventResponse{
                ID:        sr.ReceiverEvent.ID,
                Title:     sr.ReceiverEvent.Title,
                UserID:    sr.ReceiverEvent.UserID,
                StartTime: sr.ReceiverEvent.StartTime,
                EndTime:   sr.ReceiverEvent.EndTime,
                Status:    sr.ReceiverEvent.Status,
                CreatedAt: sr.ReceiverEvent.CreatedAt,
                UpdatedAt: sr.ReceiverEvent.UpdatedAt,
            },
            CreatedAt: sr.CreatedAt,
            UpdatedAt: sr.UpdatedAt,
        })
    }

    c.JSON(http.StatusOK, response)
}

func AcceptSwap(c *gin.Context) {
    userID := c.GetUint("user_id")
    swapRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid swap request ID"})
        return
    }

    if err := services.AcceptSwap(uint(swapRequestID), userID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Swap accepted successfully"})
}

func RejectSwap(c *gin.Context) {
    userID := c.GetUint("user_id")
    swapRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid swap request ID"})
        return
    }

    if err := services.RejectSwap(uint(swapRequestID), userID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Swap rejected successfully"})
}

func CancelSwap(c *gin.Context) {
    userID := c.GetUint("user_id")
    swapRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid swap request ID"})
        return
    }

    if err := services.CancelSwap(uint(swapRequestID), userID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Swap cancelled successfully"})
}