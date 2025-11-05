package controllers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"SlotSwapper/config"
	"SlotSwapper/models"
	"SlotSwapper/dto"
	"SlotSwapper/utils"
)

func Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	db := config.GetDB()
	var existingUser models.User
	if err := db.Where("email = ?", strings.ToLower(req.Email)).First(&existingUser).Error; err == nil {
		c.JSON(409, gin.H{"error": "Email already registered"})
		return
	}

	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err != gorm.ErrRecordNotFound {
		c.JSON(409, gin.H{"error": "Username already taken"})
		return
	}

	user := models.User{
		Username:  req.Username,
		Email:     strings.ToLower(req.Email),
	}

	if err := user.HashPassword(req.Password); err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "unique") {
			c.JSON(409, gin.H{"error": "User with provided details already exists"})
		} else {
			c.JSON(500, gin.H{"error": "Internal server error"})
		}
		return
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Email, user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	response := dto.AuthResponse{
		Token:    token,
		User: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}

	c.JSON(201, response)
}

func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	db := config.GetDB()
	var user models.User
	if err := db.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(401, gin.H{"error": "Invalid email or password"})
			return
		} else {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Email, user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	response := dto.AuthResponse{
		Token:    token,
		User: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}

	c.JSON(200, response)
}