package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"SlotSwapper/config"
	"SlotSwapper/dto"
	"SlotSwapper/models"
	"SlotSwapper/services"
	"SlotSwapper/utils"
)

// GetMySwapRequests returns swap requests where user is requester or receiver
func GetMySwapRequests(c *gin.Context) {
	userID := c.GetUint("user_id")
	filter := c.Query("filter") // "sent", "received", or "all" (default)

	db := config.GetDB()
	var swapRequests []models.SwapRequest

	query := db.Preload("RequesterEvent").
		Preload("ReceiverEvent").
		Preload("Requester").
		Preload("Receiver").
		Where("is_deleted = false")

	switch filter {
	case "sent":
		query = query.Where("requester_id = ?", userID)
	case "received":
		query = query.Where("receiver_id = ?", userID)
	default:
		query = query.Where("requester_id = ? OR receiver_id = ?", userID, userID)
	}

	if err := query.Order("created_at DESC").Find(&swapRequests).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch swap requests")
		return
	}

	var response []dto.SwapRequestResponse
	for _, sr := range swapRequests {
		response = append(response, mapSwapRequestToResponse(sr))
	}

	utils.RespondWithSuccess(c, http.StatusOK, "", response)
}

// GetSwappableSlots returns events that can be swapped
func GetSwappableSlots(c *gin.Context) {
	userID := c.GetUint("user_id")

	db := config.GetDB()
	var events []models.Event

	// Get events that are swappable and not owned by current user
	if err := db.Where("status = ? AND user_id != ? AND is_deleted = false",
		models.StatusSwappable, userID).
		Preload("User").
		Order("start_time ASC").
		Find(&events).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch swappable slots")
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

	utils.RespondWithSuccess(c, http.StatusOK, "", responseEvents)
}

// InitiateSwap creates a new swap request
func InitiateSwap(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.InitiateSwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", map[string]interface{}{
			"validation_error": err.Error(),
		})
		return
	}

	// Validate that user owns the requester event
	db := config.GetDB()
	var requesterEvent models.Event
	if err := db.Where("id = ? AND user_id = ? AND is_deleted = false",
		req.RequesterEventID, userID).First(&requesterEvent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SwapLogicError(c, "You don't own the requester event", map[string]interface{}{
				"requester_event_id": req.RequesterEventID,
			})
			return
		}
		utils.InternalServerError(c, "Failed to validate requester event")
		return
	}

	// Call service
	swapRequest, err := services.InitiateSwap(req.RequesterEventID, req.ReceiverEventID, userID)
	if err != nil {
		utils.SwapLogicError(c, err.Error(), map[string]interface{}{
			"requester_event_id": req.RequesterEventID,
			"receiver_event_id":  req.ReceiverEventID,
		})
		return
	}

	response := mapSwapRequestToResponse(*swapRequest)
	utils.RespondWithSuccess(c, http.StatusCreated, "Swap request initiated successfully", response)
}

// AcceptSwap accepts a swap request
func AcceptSwap(c *gin.Context) {
	userID := c.GetUint("user_id")
	swapRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid swap request ID", nil)
		return
	}

	// Verify user is the receiver
	db := config.GetDB()
	var swapRequest models.SwapRequest
	if err := db.Where("id = ? AND is_deleted = false", swapRequestID).First(&swapRequest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Swap request")
			return
		}
		utils.InternalServerError(c, "Failed to fetch swap request")
		return
	}

	if swapRequest.ReceiverID != userID {
		utils.Forbidden(c, "You are not authorized to accept this swap request")
		return
	}

	if err := services.AcceptSwap(uint(swapRequestID), userID); err != nil {
		utils.SwapLogicError(c, err.Error(), map[string]interface{}{
			"swap_request_id": swapRequestID,
		})
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, "Swap request accepted successfully", nil)
}

// RejectSwap rejects a swap request
func RejectSwap(c *gin.Context) {
	userID := c.GetUint("user_id")
	swapRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid swap request ID", nil)
		return
	}

	// Verify user is the receiver
	db := config.GetDB()
	var swapRequest models.SwapRequest
	if err := db.Where("id = ? AND is_deleted = false", swapRequestID).First(&swapRequest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Swap request")
			return
		}
		utils.InternalServerError(c, "Failed to fetch swap request")
		return
	}

	if swapRequest.ReceiverID != userID {
		utils.Forbidden(c, "You are not authorized to reject this swap request")
		return
	}

	if err := services.RejectSwap(uint(swapRequestID), userID); err != nil {
		utils.SwapLogicError(c, err.Error(), map[string]interface{}{
			"swap_request_id": swapRequestID,
		})
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, "Swap request rejected successfully", nil)
}

// CancelSwap cancels a swap request (requester only)
func CancelSwap(c *gin.Context) {
	userID := c.GetUint("user_id")
	swapRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid swap request ID", nil)
		return
	}

	// Verify user is the requester
	db := config.GetDB()
	var swapRequest models.SwapRequest
	if err := db.Where("id = ? AND is_deleted = false", swapRequestID).First(&swapRequest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Swap request")
			return
		}
		utils.InternalServerError(c, "Failed to fetch swap request")
		return
	}

	if swapRequest.RequesterID != userID {
		utils.Forbidden(c, "You are not authorized to cancel this swap request")
		return
	}

	if err := services.CancelSwap(uint(swapRequestID), userID); err != nil {
		utils.SwapLogicError(c, err.Error(), map[string]interface{}{
			"swap_request_id": swapRequestID,
		})
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, "Swap request cancelled successfully", nil)
}

// Helper function to map swap request to response
func mapSwapRequestToResponse(sr models.SwapRequest) dto.SwapRequestResponse {
	return dto.SwapRequestResponse{
		ID:               sr.ID,
		RequesterID:      sr.RequesterID,
		ReceiverID:       sr.ReceiverID,
		RequesterEventID: sr.RequesterEventID,
		ReceiverEventID:  sr.ReceiverEventID,
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
	}
}
