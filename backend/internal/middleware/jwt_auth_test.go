package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/moLIart/gomoku-backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func dummyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("dummy"))
}

func TestJWTAuth_MissingHeader(t *testing.T) {
	jwtSvc := services.NewJWTService("testsecret")
	handler := JWTAuth(dummyHandler, jwtSvc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	ps := httprouter.Params{}

	handler(rr, req, ps)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing or invalid Authorization header")
}

func TestJWTAuth_InvalidPrefix(t *testing.T) {
	jwtSvc := services.NewJWTService("testsecret")
	handler := JWTAuth(dummyHandler, jwtSvc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token something")
	rr := httptest.NewRecorder()
	ps := httprouter.Params{}

	handler(rr, req, ps)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing or invalid Authorization header")
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	jwtSvc := services.NewJWTService("testsecret")
	handler := JWTAuth(dummyHandler, jwtSvc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rr := httptest.NewRecorder()
	ps := httprouter.Params{}

	handler(rr, req, ps)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid token")
}

func TestJWTAuth_ValidToken(t *testing.T) {
	jwtSvc := services.NewJWTService("testsecret")
	token, err := jwtSvc.Sign("testuser")
	assert.NoError(t, err)

	handler := JWTAuth(dummyHandler, jwtSvc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	ps := httprouter.Params{}

	handler(rr, req, ps)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "dummy", strings.TrimSpace(rr.Body.String()))
}
