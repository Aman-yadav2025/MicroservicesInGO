package repository

import (
	"auth-service/internal/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect_DB() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("env can't be loaded")
	}
	dsn := os.Getenv("DB_DSN")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database with error %v", err)
	}

	err = database.AutoMigrate(&model.User{})

	if err != nil {
		log.Fatalf("Failed to migrate database with error %v",err)
	}

	DB = database

	log.Println("Successfully connect to database and migrated the tables")
}
