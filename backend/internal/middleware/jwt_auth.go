package middleware

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/moLIart/gomoku-backend/internal/services"
)

func JWTAuth(handler httprouter.Handle, jwtSvc *services.JWTService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

		handler(w, r, ps)
	}
}
