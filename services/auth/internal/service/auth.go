package service

import (
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = os.Getenv("jwtSecretKey")
var jwtRefreshKey = os.Getenv("jwtRefreshKey")

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

func LoginUser(email string, password string) (string, string, error) {
	var user model.User

	result := repository.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return "", "", errors.New("invalid email or password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecretKey))

	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(jwtRefreshKey)

	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func RefreshAccessToken(refreshTokenString string) (string,string, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtRefreshKey, nil
	})

	if err != nil {
		return "", "", errors.New("Invalid Refresh Token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", "", errors.New("Invalid token claims")
	}

	newAccessTokken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"user_id": claims["user_id"],
		"email":   claims["email"],
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	newAccessString,err := newAccessTokken.SignedString(jwtSecretKey)
	if err != nil {
		return "","",err
	}

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodES256,jwt.MapClaims{
		"user_id": claims["user_id"],
		"email":   claims["email"],
		"exp":     time.Now().Add(time.Hour * 24* 7).Unix(),
	})

	newRefreshString,err:= newRefreshToken.SignedString(jwtRefreshKey)
	if err != nil {
		return "","",err
	}
	return newAccessString,newRefreshString,nil
}
