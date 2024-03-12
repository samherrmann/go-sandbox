package pages

import (
	"net/http"
)

func NewHomeHandler(path string, redirectPaths string) http.Handler {
	return http.RedirectHandler(
		redirectPaths,
		http.StatusMovedPermanently,
	)
}
