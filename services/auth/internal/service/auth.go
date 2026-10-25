package service

import (
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user := model.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	result := repository.DB.Create(&user)

	if result.Error != nil {
		return errors.New("user already exists or database error")
	}
	return nil
}

func LoginUser(email string, password string) (string, error) {
	var user model.User

	result := repository.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return "", errors.New("invalid email or password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	if err := godotenv.Load();err != nil {
		return "",errors.New("Failed to load Env")
	}
	jwtSecretKey := os.Getenv("jwtSecretKey")
	tokenString, err := token.SignedString(jwtSecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
