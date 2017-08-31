package hsts

import "net/http"

// NewHandler wraps the given HTTP handler with a handler that sets the HSTS header.
func NewHandler(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 31536000 seconds = 1 year.
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		hf.ServeHTTP(w, r)
	})
}
