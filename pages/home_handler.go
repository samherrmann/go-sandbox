package pages

import "net/http"

func NewHomeHandler(redirectPaths string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Process exact match only.
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, redirectPaths, http.StatusMovedPermanently)
	})
}
