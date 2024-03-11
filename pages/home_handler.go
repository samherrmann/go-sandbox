package pages

import (
	"net/http"

	"github.com/samherrmann/go-sandbox/httputil"
)

func NewHomeHandler(path string, redirectPaths string) http.Handler {
	redirectHandler := http.RedirectHandler(
		redirectPaths,
		http.StatusMovedPermanently,
	)
	return httputil.ExactPathHandler(path, redirectHandler)
}
