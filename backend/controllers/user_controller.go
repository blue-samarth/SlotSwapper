package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "SlotSwapper/config"
    "SlotSwapper/dto"
    "SlotSwapper/models"
    "SlotSwapper/utils"
)

func GetCurrentUser(c *gin.Context) {
    userID := c.GetUint("user_id")

    db := config.GetDB()
    var user models.User

    if err := db.Select("id, email, username, created_at, updated_at").
        Where("id = ?", userID).First(&user).Error; err != nil {
        utils.NotFound(c, "User")
        return
    }

    response := dto.UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Username:  user.Username,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }

    utils.RespondWithSuccess(c, http.StatusOK, "", response)
}

func UpdateUserProfile(c *gin.Context) {
    userID := c.GetUint("user_id")

    var req dto.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationError(c, "Invalid request data", map[string]interface{}{
            "validation_error": err.Error(),
        })
        return
    }

    db := config.GetDB()
    var user models.User

    if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
        utils.NotFound(c, "User")
        return
    }

    // Update allowed fields
    if req.Username != "" {
        user.Username = req.Username
    }

    if err := db.Save(&user).Error; err != nil {
        utils.InternalServerError(c, "Failed to update user profile")
        return
    }

    response := dto.UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Username:  user.Username,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }

    utils.RespondWithSuccess(c, http.StatusOK, "Profile updated successfully", response)
}