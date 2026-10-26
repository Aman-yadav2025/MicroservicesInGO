package service

import (
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte(os.Getenv("jwtSecretKey"))
var jwtRefreshKey = []byte(os.Getenv("jwtRefreshKey"))

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
		//fmt.Println("DEBUG: User not found in database!")
		return "", "", errors.New("invalid email or password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		//fmt.Println("DEBUG: Bcrypt password mismatch!")
		return "", "", errors.New("invalid email or password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecretKey)

	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(jwtRefreshKey)

	if err != nil {
		return "", "", err
	}

	//repository.DB.Where("user_id = ?",user.Id).Delete(&model.RefreshToken{})
	dbToken := model.RefreshToken{
		UserId:     user.Id,
		Token:      refreshTokenString,
		ExpiryTime: time.Now().Add(time.Hour * 24 * 7),
	}

	if err = repository.DB.Create(&dbToken).Error; err != nil {
		return "", "", errors.New("failed to save token in the database")
	}

	return tokenString, refreshTokenString, nil
}

func RefreshAccessToken(refreshTokenString string) (string, string, error) {
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

	userId := claims["user_id"].(float64)

	var dbToken model.RefreshToken

	result := repository.DB.Where("token = ?", refreshTokenString).First(&dbToken)

	if result.Error != nil {
		// token not found security breach

		err := repository.DB.Where("user_id = ?", uint(userId)).Delete(&model.RefreshToken{}).Error

		if err != nil {
			log.Println("CRITICAL: Failed to delete sessions during breach:", err)
		}
		return "", "", errors.New("Security breach deteched: token being reused, all sesssion deleted")
	}
	//delete old token assign new one
	repository.DB.Delete(&dbToken)

	newAccessTokken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims["user_id"],
		"email":   claims["email"],
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	newAccessString, err := newAccessTokken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims["user_id"],
		"email":   claims["email"],
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	newRefreshString, err := newRefreshToken.SignedString(jwtRefreshKey)
	if err != nil {
		return "", "", err
	}

	newDbToken := model.RefreshToken{
		UserId:     uint(userId),
		Token:      newRefreshString,
		ExpiryTime: time.Now().Add(time.Hour * 24 * 7),
	}

	repository.DB.Create(&newDbToken)
	return newAccessString, newRefreshString, nil
}

func LogoutUser(refreshTokenString string) error {

	fmt.Println("DEBUG: Attempting to delete token:", refreshTokenString)

	result := repository.DB.Where("token = ?", refreshTokenString).Delete(&model.RefreshToken{})

	return result.Error
}

func RequestPasswordReset(email string) error {
	var user model.User
	if err := repository.DB.Where("email = ?", email).First(&user).Error; err != nil {
		//return nil if email not found to preserve privacy of user
		//even if email is found do not send opt send because if again reveals that user exists on this plateform
		return nil
	}
	repository.DB.Where("email = ?", email).Delete(&model.PasswordReset{})
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	resetData := model.PasswordReset{
		Email:  email,
		OTP:    otp,
		Expiry: time.Now().Add(5 * time.Minute),
	}

	repository.DB.Create(&resetData)
	//---------------------------------------------
	//smtp code here//curently temperary testing code
	log.Printf("MOCK EMAIL TO: %s\n", email)
	log.Printf("SUBJECT: Your Password Reset Code\n")
	log.Printf("BODY: Your OTP is %s. It expires in 5 minutes.\n", otp)
	//---------------------------------------------
	return nil
}

func ResetPassword(email string, otp string, password string) error {
	var resetData model.PasswordReset

	result := repository.DB.Where("email = ? AND otp = ?", email, otp).First(&resetData)

	if result.Error != nil {
		return errors.New("Invalid OTP")
	}

	if time.Now().After(resetData.Expiry) {
		repository.DB.Delete(&resetData)
		return errors.New("OTP expired")
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	repository.DB.Model(&model.User{}).Where("email = ?", email).Update("password", string(hashedpassword))

	repository.DB.Where("email = ?", email).Delete(&model.PasswordReset{})

	var user model.User
	if err := repository.DB.Where("email = ?", email).First(&user).Error; err == nil {
		repository.DB.Where("user_id = ?", user.Id).Delete(&model.RefreshToken{})
	}

	return nil
}
