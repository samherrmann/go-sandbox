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
	logger := slog.New(&slog.TextHandler{})

	homePage, err := pages.NewHome(logger)
	if err != nil {
		return err
	}

	http.Handle("/", http.RedirectHandler("/todo", http.StatusMovedPermanently))
	http.HandleFunc("GET /todo", homePage.GetToDos())
	http.HandleFunc("POST /todo/add", homePage.AddToDo())
	http.HandleFunc("POST /todo/{id}/delete", homePage.RemoveToDo())

	addr := ":8080"
	fmt.Printf("Listening on %v...\n", addr)
	return http.ListenAndServe(addr, nil)
}
