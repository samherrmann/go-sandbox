package httputil

import "net/http"

// SeeOther writes a 303 response. The primary use case for this response is to
// implement the PRG (Post/Redirect/Get) pattern.
// https://en.wikipedia.org/wiki/Post/Redirect/Get
func SeeOther(w http.ResponseWriter, path string) {
	w.Header().Add("Location", path)
	w.WriteHeader(http.StatusSeeOther)
}

// ExactPathHandler returns an http.Handler that only calls the next handler if
// the request path matches the given path exactly. Otherwise, it calls
// http.NotFound.
func ExactPathHandler(path string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
