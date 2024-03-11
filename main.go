package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/samherrmann/go-sandbox/pages"
)

func main() {
	if err := app(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func app() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	router, err := newRouter(logger)
	if err != nil {
		return err
	}

	addr := ":8080"
	fmt.Printf("Listening on %v...\n", addr)
	return http.ListenAndServe(addr, router)
}

func newRouter(logger *slog.Logger) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	// Create handlers.
	todoPath := "/todo"
	todoHandler, err := pages.NewTodoHandler(todoPath, logger)
	if err != nil {
		return nil, err
	}
	homeHandler := pages.NewHomeHandler(todoPath)

	// Register handlers in mux.
	//
	// Trailing-slash redirection:
	// https://pkg.go.dev/net/http#hdr-Trailing_slash_redirection
	mux.Handle(todoPath+"/", http.StripPrefix(todoPath, todoHandler))
	mux.Handle("/", homeHandler)

	return mux, nil
}
