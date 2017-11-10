package redirect

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRedirectToHTTPS(t *testing.T) {
	handler1 := func(w http.ResponseWriter, r *http.Request) {}
	wrappedHandler := ToHTTPS(http.HandlerFunc(handler1))
	tests := []struct {
		msg        string
		reqURL     string
		wantStatus int
		respURL    string
	}{
		{
			msg:        "redirect with params",
			reqURL:     "http://www.domain.com/path?q=1&r=2",
			wantStatus: http.StatusTemporaryRedirect,
			respURL:    "https://www.domain.com/path?q=1&r=2",
		}, {
			msg:        "redirect no params",
			reqURL:     "http://www.domain.com",
			wantStatus: http.StatusTemporaryRedirect,
			respURL:    "https://www.domain.com",
		}, {
			msg:        "no redirection",
			reqURL:     "https://www.domain.com",
			wantStatus: http.StatusOK,
		},
	}
	for _, test := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", test.reqURL, nil)
		if !strings.HasPrefix(test.reqURL, "https") {
			r.Header.Set(headerForward, "http")
		}
		wrappedHandler.ServeHTTP(w, r)
		if got, want := w.Result().StatusCode, test.wantStatus; got != want {
			t.Errorf("got=%d but wanted wanted redirection with code=%d", got, want)
		}
		if test.respURL == "" {
			return
		}
		if got, want := w.HeaderMap["Location"][0], test.respURL; got != want {
			t.Errorf("got redirect location=%q, want=%q", got, want)
		}
	}
}
