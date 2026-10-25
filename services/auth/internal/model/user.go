package model

import "time"

type User struct {
	Id         uint   	`gorm:"primaryKey"`
	Email      string 	`gorm:"uniqueIndex;not null"`
	Password   string 	`gorm:"not null"`
	Create_At  time.Time
	Updated_at time.Time
}
