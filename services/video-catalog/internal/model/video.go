package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Title string `grom:"not null"`
	Description string
	VideoUrl string `gorm:"not null"`
	ThrumnailUrl string
	UserId uint `gorm:"not null"`
	Views int `grom:"default:0"`
}
