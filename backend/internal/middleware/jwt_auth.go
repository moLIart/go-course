package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/moLIart/gomoku-backend/internal/services"
)

type authPlayerNameKey struct{}

var AuthPlayerNameKey = authPlayerNameKey{}

func JWTAuth(jwtSvc *services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwtSvc.Verify(tokenString)
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			val, ok := claims["name"]
			if !ok {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			newContext := context.
				WithValue(r.Context(), AuthPlayerNameKey, val)
			next.ServeHTTP(w, r.WithContext(newContext))
		})
	}
}
