package middleware

import "net/http"

// ContentType sets Content-Type header to desired mime type.
func ContentType(mime string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", mime)
			next.ServeHTTP(w, r)
		})
	}
}
