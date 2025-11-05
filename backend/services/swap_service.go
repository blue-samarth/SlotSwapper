package services

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"SlotSwapper/config"
	"SlotSwapper/models"
	"SlotSwapper/utils"
)

func InitiateSwap(requesterEventID, receiverEventID uint, userID uint) (*models.SwapRequest, error) {
	db := config.GetDB()

	var swapRequest models.SwapRequest
	var swapRequestID uint
	err := utils.WithDefaultTransaction(db, func(tx *gorm.DB) error {
		var event1, event2 models.Event
		var eventID1, eventID2 uint
		if requesterEventID < receiverEventID {
			eventID1, eventID2 = requesterEventID, receiverEventID
		} else {
			eventID1, eventID2 = receiverEventID, requesterEventID
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND is_deleted = false", eventID1).
			First(&event1).Error; err != nil {
			return errors.New("first event not found")
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND is_deleted = false", eventID2).
			First(&event2).Error; err != nil {
			return errors.New("second event not found")
		}

		var requesterEvent, receiverEvent *models.Event
		if event1.ID == requesterEventID {
			requesterEvent = &event1
			receiverEvent = &event2
		} else {
			requesterEvent = &event2
			receiverEvent = &event1
		}

		if requesterEvent.UserID != userID {
			return errors.New("you don't own the requester event")
		}

		receiverID := receiverEvent.UserID

		if !requesterEvent.IsSwappable() {
			return errors.New("your event is not swappable")
		}
		if !receiverEvent.IsSwappable() {
			return errors.New("receiver event is not swappable")
		}

		swapRequest = models.SwapRequest{
			RequesterID:      userID,
			ReceiverID:       receiverID,
			RequesterEventID: requesterEventID,
			ReceiverEventID:  receiverEventID,
			Status:           models.StatusPending,
		}

		if err := tx.Create(&swapRequest).Error; err != nil {
			return err
		}
		requesterEvent.Status = models.StatusSwapPending
		if err := tx.Save(requesterEvent).Error; err != nil {
			return err
		}

		receiverEvent.Status = models.StatusSwapPending
		if err := tx.Save(receiverEvent).Error; err != nil {
			return err
		}

		swapRequestID = swapRequest.ID
		return nil
	})

	if err != nil {
		return nil, err
	}
	if err := db.Preload("RequesterEvent").Preload("ReceiverEvent").
		Preload("Requester").Preload("Receiver").
		First(&swapRequest, swapRequestID).Error; err != nil {
		return nil, err
	}

	return &swapRequest, nil
}

func AcceptSwap(swapRequestID, receiverID uint) error {
	db := config.GetDB()

	return utils.WithDefaultTransaction(db, func(tx *gorm.DB) error {
		var swapRequest models.SwapRequest
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND is_deleted = false", swapRequestID).
			First(&swapRequest).Error; err != nil {
			return errors.New("swap request not found")
		}

		if swapRequest.ReceiverID != receiverID {
			return errors.New("you are not the receiver of this swap")
		}

		if swapRequest.Status != models.StatusPending {
			return errors.New("swap request is not pending")
		}

		eventID1, eventID2 := swapRequest.RequesterEventID, swapRequest.ReceiverEventID
		if eventID1 > eventID2 {
			eventID1, eventID2 = eventID2, eventID1
		}

		var event1, event2 models.Event
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND is_deleted = false", eventID1).
			First(&event1).Error; err != nil {
			return errors.New("first event not found")
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND is_deleted = false", eventID2).
			First(&event2).Error; err != nil {
			return errors.New("second event not found")
		}

		var requesterEvent, receiverEvent *models.Event
		if event1.ID == swapRequest.RequesterEventID {
			requesterEvent = &event1
			receiverEvent = &event2
		} else {
			requesterEvent = &event2
			receiverEvent = &event1
		}

		if requesterEvent.Status != models.StatusSwapPending {
			return errors.New("requester event is not in swap pending state")
		}
		if receiverEvent.Status != models.StatusSwapPending {
			return errors.New("receiver event is not in swap pending state")
		}

		// Swap user IDs within transactio
		requesterEvent.UserID, receiverEvent.UserID = receiverEvent.UserID, requesterEvent.UserID

		requesterEvent.Status = models.StatusBusy
		receiverEvent.Status = models.StatusBusy

		if err := tx.Save(requesterEvent).Error; err != nil {
			return err
		}
		if err := tx.Save(receiverEvent).Error; err != nil {
			return err
		}

		swapRequest.Status = models.StatusAccepted
		if err := tx.Save(&swapRequest).Error; err != nil {
			return err
		}

		return nil
	})
}

func RejectSwap(swapRequestID, receiverID uint) error {
	db := config.GetDB()

	return utils.WithDefaultTransaction(db, func(tx *gorm.DB) error {
		var swapRequest models.SwapRequest
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND is_deleted = false", swapRequestID).
			First(&swapRequest).Error; err != nil {
			return errors.New("swap request not found")
		}

		if swapRequest.ReceiverID != receiverID {
			return errors.New("you are not the receiver of this swap")
		}

		if swapRequest.Status != models.StatusPending {
			return errors.New("swap request is not pending")
		}

		// Use UpdateColumn to bypass BeforeSave validation when only updating status
		if err := tx.Model(&models.Event{}).
			Where("id IN (?, ?) AND is_deleted = false", swapRequest.RequesterEventID, swapRequest.ReceiverEventID).
			UpdateColumn("status", models.StatusSwappable).Error; err != nil {
			return err
		}

		swapRequest.Status = models.StatusRejected
		return tx.Save(&swapRequest).Error
	})
}

func CancelSwap(swapRequestID, requesterID uint) error {
	db := config.GetDB()

	return utils.WithDefaultTransaction(db, func(tx *gorm.DB) error {
		var swapRequest models.SwapRequest
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND is_deleted = false", swapRequestID).
			First(&swapRequest).Error; err != nil {
			return errors.New("swap request not found")
		}

		if swapRequest.RequesterID != requesterID {
			return errors.New("you are not the requester of this swap")
		}

		if swapRequest.Status != models.StatusPending {
			return errors.New("swap request is not pending")
		}

		// Use UpdateColumn to bypass BeforeSave validation when only updating status
		if err := tx.Model(&models.Event{}).
			Where("id IN (?, ?) AND is_deleted = false", swapRequest.RequesterEventID, swapRequest.ReceiverEventID).
			UpdateColumn("status", models.StatusSwappable).Error; err != nil {
			return err
		}

		swapRequest.Status = models.StatusCancelled
		return tx.Save(&swapRequest).Error
	})
}
