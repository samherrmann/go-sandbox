package httputil

import "net/http"

// SeeOther writes a 303 response. The primary use case for this response is to
// implement the PRG (Post/Redirect/Get) pattern.
// https://en.wikipedia.org/wiki/Post/Redirect/Get
func SeeOther(w http.ResponseWriter, path string) {
	w.Header().Add("Location", path)
	w.WriteHeader(http.StatusSeeOther)
}
