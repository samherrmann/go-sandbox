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

	todoPage, err := pages.NewToDo(logger)
	if err != nil {
		return err
	}

	http.Handle("/", http.RedirectHandler("/todo", http.StatusMovedPermanently))
	http.HandleFunc("GET /todo", todoPage.Get())
	http.HandleFunc("POST /todo/add", todoPage.Add())
	http.HandleFunc("POST /todo/{id}/update", todoPage.Update())
	http.HandleFunc("POST /todo/{id}/delete", todoPage.Delete())

	addr := ":8080"
	fmt.Printf("Listening on %v...\n", addr)
	return http.ListenAndServe(addr, nil)
}
