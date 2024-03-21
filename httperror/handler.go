package httperror

import (
	"log/slog"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

func LogHandler(h Handler, l *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h.ServeHTTP(w, r); err != nil {
			l.Error(err.Error())
		}
	})
}
