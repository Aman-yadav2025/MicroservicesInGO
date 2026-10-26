package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testJWTSecret = "test-secret"

func TestRequestAuthRejectsMissingAndInvalidTokens(t *testing.T) {
	t.Setenv("jwtSecretKey", testJWTSecret)

	tests := []struct {
		name  string
		token string
	}{
		{name: "missing token"},
		{name: "wrong authorization scheme", token: "Basic token"},
		{name: "malformed token", token: "Bearer not-a-jwt"},
		{name: "invalid signature", token: "Bearer " + signedTestToken(t, "other-secret", jwt.MapClaims{"user_id": 42, "exp": time.Now().Add(time.Hour).Unix()})},
		{name: "expired token", token: "Bearer " + signedTestToken(t, testJWTSecret, jwt.MapClaims{"user_id": 42, "exp": time.Now().Add(-time.Hour).Unix()})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				called = true
			})
			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			rec := httptest.NewRecorder()

			RequestAuth(next).ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
			}
			if called {
				t.Error("protected handler was called for an invalid authorization header")
			}
		})
	}
}

func TestRequestAuthPassesValidTokenUserID(t *testing.T) {
	t.Setenv("jwtSecretKey", testJWTSecret)
	token := signedTestToken(t, testJWTSecret, jwt.MapClaims{
		"user_id": 42,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Context().Value("user_id"); got != float64(42) {
			t.Errorf("user_id in context = %v, want 42", got)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	RequestAuth(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func signedTestToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}
	return signed
}
