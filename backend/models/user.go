package models

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"not null;uniqueIndex;size:255;" json:"username"`
	Email     string         `gorm:"unique;not null;size:255;" json:"email"`
	Password  string         `gorm:"not null;size:75;" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
	// Role	  string         `gorm:"size:50;default:'user'" json:"role"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	return nil
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
