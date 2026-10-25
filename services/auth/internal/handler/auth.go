package handler

import (
	"auth-service/internal/service"
	"encoding/json"
	"net/http"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(w http.ResponseWriter, r *http.Request){
	var req SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&req);err != nil {
		http.Error(w,"Invalid request Body",http.StatusBadRequest)
		return
	}

	err := service.RegisterUser(req.Email,req.Password)

	if err != nil {
		http.Error(w,err.Error(),http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func Login(w http.ResponseWriter, r *http.Request){
	var req SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body",http.StatusBadRequest)
		return
	}

	token,err:= service.LoginUser(req.Email,req.Password)
	if err != nil {
		http.Error(w,err.Error(),http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":token,
	})
}
