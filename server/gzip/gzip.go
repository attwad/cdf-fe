package gzip

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

type gzipHandler struct {
	h http.Handler
}

func (g *gzipHandler) sendGzip(w http.ResponseWriter, r *http.Request) error {
	// Get the content from the original handler.
	rr := httptest.NewRecorder()
	g.h.ServeHTTP(rr, r)
	resp := rr.Result()
	// Check the result, if okay, compress it else send the response back.
	if resp.StatusCode != 200 {
		log.Println("Status not okay when reading file from original handler")
		g.h.ServeHTTP(w, r)
		return nil
	}
	// Set required headers.
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-type", resp.Header.Get("Content-Type"))
	// Then copy the response but in gzip this time.
	zw := gzip.NewWriter(w)
	defer resp.Body.Close()
	defer zw.Close()
	if _, err := io.Copy(zw, resp.Body); err != nil {
		log.Println("Error copying response to gzip")
		return err
	}
	return nil
}

func (g *gzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accepts := strings.Split(r.Header.Get("Accept-Encoding"), ",")
	for _, accept := range accepts {
		if strings.TrimSpace(accept) == "gzip" {
			if err := g.sendGzip(w, r); err != nil {
				http.Error(w, "error compressing response", http.StatusInternalServerError)
			}
			return
		}
	}
	g.h.ServeHTTP(w, r)
}

// NewGZipHTTPHandler returns an HTTP handler that wraps the given handler.
func NewGZipHTTPHandler(orig http.Handler) http.Handler {
	return &gzipHandler{h: orig}
}
