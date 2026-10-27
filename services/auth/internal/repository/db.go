package repository

import (
	"auth-service/internal/model"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect_DB() {
	dsn := os.Getenv("DB_DSN")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database with error %v", err)
	}

	err = database.AutoMigrate(&model.User{}, &model.RefreshToken{}, &model.PasswordReset{})

	if err != nil {
		log.Fatalf("Failed to migrate database with error %v", err)
	}

	DB = database

	log.Println("Successfully connect to database and migrated the tables")
}
