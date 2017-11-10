package redirect

import "net/http"

const headerForward = "X-Forwarded-Proto"

// ToHTTPS redirects all HTTP connections that were marked as such by a GCP
// load balancer by looking at the X-Forwarded-Proto header.
func ToHTTPS(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(headerForward) == "http" {
			newURL := "https://" + r.URL.Host + r.URL.Path
			q := r.URL.Query().Encode()
			if len(q) > 0 {
				newURL += "?" + q
			}
			http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
			return
		}
		hf.ServeHTTP(w, r)
	})
}
