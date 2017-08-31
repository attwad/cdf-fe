package gzip

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGzipHandler(t *testing.T) {
	tests := []struct {
		msg            string
		req            *http.Request
		acceptHeader   string
		wantStatus     int
		wantBody       string
		wantRespHeader string
	}{
		{
			msg:            "no accepting gzip",
			req:            httptest.NewRequest("GET", "/testurl", nil),
			acceptHeader:   "text/html",
			wantStatus:     200,
			wantBody:       "/testurl",
			wantRespHeader: "",
		}, {
			msg:            "accepting gzip",
			req:            httptest.NewRequest("GET", "/testurl", nil),
			acceptHeader:   "gzip,deflate",
			wantStatus:     200,
			wantBody:       "",
			wantRespHeader: "gzip",
		},
	}
	for _, test := range tests {
		test.req.Header.Set("Accept-Encoding", test.acceptHeader)
		// A handler that writes the request URL in its response body.
		origHandler := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", r.URL.String())
		}
		gzh := NewGZipHTTPHandler(http.HandlerFunc(origHandler))
		w := httptest.NewRecorder()
		gzh.ServeHTTP(w, test.req)
		resp := w.Result()
		if got, want := resp.StatusCode, test.wantStatus; got != want {
			t.Errorf("[%s] Status code got=%d, want=%d", test.msg, got, want)
		}
		if test.wantStatus != 200 {
			continue
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := string(b), test.wantBody; test.wantBody != "" && got != want {
			t.Errorf("[%s] body got=%s, want=%s", test.msg, got, want)
		}
		if got, want := resp.Header.Get("Content-Encoding"), test.wantRespHeader; got != want {
			t.Errorf("[%s] resp header got=%s, want=%s", test.msg, got, want)
		}
	}
}
