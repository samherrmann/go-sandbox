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

	http.HandleFunc("GET /", homePage.GetToDos())
	http.HandleFunc("POST /", homePage.AddToDo())
	http.HandleFunc("POST /delete", homePage.RemoveToDo())

	addr := ":8080"
	fmt.Printf("Listening on %v...\n", addr)
	return http.ListenAndServe(":8080", nil)
}
