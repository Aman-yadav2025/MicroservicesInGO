package repository

import (
	"log"
	"os"
	"video-catalog/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect_DB() {

	dsn := os.Getenv("Video_catalog_db")


	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the video database:", err)
	}
	err = database.AutoMigrate(&model.Video{})

	if err != nil {
		log.Fatal("Failed to migrate video database:", err)
	}

	DB = database

	log.Println("Successfully connected to Video PostgreSQL and migrated tables!")
}
