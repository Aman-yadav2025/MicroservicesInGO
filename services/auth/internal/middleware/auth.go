package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)
var jwtSecretKey = os.Getenv("jwtSecretKey")
var jwtRefreshKey = os.Getenv("jwtRefreshKey")

func RequestAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		//takes two thing one is raw token string second is a function which return key
		//first strip away the payload and header and then compare
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token or expired Token", http.StatusUnauthorized)
			return
		}
		if claim, ok := token.Claims.(jwt.MapClaims); ok {
			//extract user_id form the token and pack it in the
			//context so that next file can see it directly
			ctx := context.WithValue(r.Context(), "user_id", claim["user_id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
		}
	})
}
