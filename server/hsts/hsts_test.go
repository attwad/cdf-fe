package hsts

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {}
	w := httptest.NewRecorder()
	NewHandler(http.HandlerFunc(h)).ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if got, want := w.Result().Header.Get("Strict-Transport-Security"), "max-age=31536000"; got != want {
		t.Errorf("header got=%s, want=%s", got, want)
	}
}
