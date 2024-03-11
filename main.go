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

	todoPage, err := pages.NewToDo(http.DefaultServeMux, "/todo", logger)
	if err != nil {
		return err
	}

	// Redirect home page to todo page.
	http.Handle(
		"/",
		http.RedirectHandler(todoPage.Path, http.StatusMovedPermanently),
	)

	addr := ":8080"
	fmt.Printf("Listening on %v...\n", addr)
	return http.ListenAndServe(addr, nil)
}
