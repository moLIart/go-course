package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentTypeMiddleware(t *testing.T) {
	const mime = "application/json"

	// Dummy handler to check if middleware passes control
	handlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	middleware := ContentType(mime)
	wrapped := middleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	assert.True(t, handlerCalled, "Next handler should be called")
	assert.Equal(t, mime, rec.Header().Get("Content-Type"), "Content-Type header should be set")
}
