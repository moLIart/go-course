package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moLIart/gomoku-backend/internal/services"
)

func TestJWTAuth_MissingAuthHeader(t *testing.T) {
	jwtSvc := &services.JWTService{}
	handler := JWTAuth(jwtSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing or invalid Authorization header")
}

func TestJWTAuth_InvalidAuthHeader(t *testing.T) {
	jwtSvc := &services.JWTService{}
	handler := JWTAuth(jwtSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "InvalidToken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing or invalid Authorization header")
}

func TestJWTAuth_VerifyError(t *testing.T) {
	jwtSvc := &services.JWTService{}

	handler := JWTAuth(jwtSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer sometoken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid token")
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	jwtSvc := &services.JWTService{}
	handler := JWTAuth(jwtSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer sometoken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid token")
}

func TestJWTAuth_InvalidClaimsType(t *testing.T) {
	jwtSvc := &services.JWTService{}
	handler := JWTAuth(jwtSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer sometoken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid token")
}

func TestJWTAuth_MissingNameClaim(t *testing.T) {
	jwtSvc := &services.JWTService{}
	handler := JWTAuth(jwtSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer sometoken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid token")
}

func TestJWTAuth_Success(t *testing.T) {
	playerName := "testuser"
	jwtSvc := services.NewJWTService("ec23e6b12163c54ddf3bb814ab1a2a6b0169df44ab78d2a2daf6e6914203dc36f18454c0348689e1183202363a89614faab2e49b1ae8ef558dfb03032c489072")

	called := false
	handler := JWTAuth(jwtSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		val := r.Context().Value(AuthPlayerNameKey)
		require.Equal(t, playerName, val)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQyNDI5Nzc5MDUsIm5hbWUiOiJ0ZXN0dXNlciJ9.fs-DtT-034WmcOnJU4sdlwYfoa9_1ij_CWEdY4A-LtI")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.True(t, called, "next handler should be called")
	assert.Equal(t, http.StatusOK, rr.Code)
}
