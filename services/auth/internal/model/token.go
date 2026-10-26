package model

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserId     uint      `gorm:"not null"`
	Token      string    `gorm:"uniqueIndex;not null"`
	ExpiryTime time.Time `grom:"not null"`
}

type PasswordReset struct {
	Email  string    `gorm:"not null"`
	OTP    string    `gorm:"not null"`
	Expiry time.Time `grom:"not null"`
}
