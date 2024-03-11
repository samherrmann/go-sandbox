package pages

import "net/http"

func NewHomeHandler(redirectPaths string) http.Handler {
	return http.RedirectHandler(redirectPaths, http.StatusMovedPermanently)
}
