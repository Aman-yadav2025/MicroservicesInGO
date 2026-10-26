package handler

import (
	//"auth-service/internal/model"
	"auth-service/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email"`
	OTP      string `json:"otp"`
	Password string `json:"password"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	err := service.RegisterUser(req.Email, req.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	AccessToken, RefreshToken, err := service.LoginUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  AccessToken,
		"refresh_token": RefreshToken,
	})
}

func RefrehToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshToken

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}
	newAccessToken, newRefreshToken, err := service.RefreshAccessToken(req.RefreshToken)

	if err != nil {
		http.Error(w, "Invalide Refresh Token", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "json/application")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshToken

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid requset body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
        http.Error(w, "Token is missing from request", http.StatusBadRequest)
        return
    }

	if err := service.LogoutUser(req.RefreshToken); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successfully"))
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}
	err := service.RequestPasswordReset(req.Email)

	if err != nil {
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil, "OTP is sent to email %s", req.Email))
}

func ResetPassword(w http.ResponseWriter, r *http.Request){
	var req ResetPasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req);err != nil {
		http.Error(w,"invalid request body",http.StatusBadRequest)
		return
	}
	if err := service.ResetPassword(req.Email,req.OTP,req.Password);err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password has been reset successfully"))
}
