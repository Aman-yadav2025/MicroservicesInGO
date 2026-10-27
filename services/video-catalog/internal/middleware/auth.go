package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(os.Getenv("jwtSecretKey"))

func RequestAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		userID := uint(claims["user_id"].(float64))

		ctx := context.WithValue(r.Context(),"user_id",userID)

		next.ServeHTTP(w,r.WithContext(ctx))
	})
}
