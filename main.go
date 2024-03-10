package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/samherrmann/go-sandbox/pages"
	"github.com/samherrmann/go-sandbox/templates"
)

func main() {
	if err := app(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func app() error {
	logger := slog.New(&slog.TextHandler{})
	tpls := templates.New()

	homePage := pages.NewHome(logger, tpls)

	http.HandleFunc("GET /", homePage.GetToDos())
	http.HandleFunc("POST /", homePage.AddToDo())

	addr := ":8080"
	fmt.Printf("Listening on %v...\n", addr)
	return http.ListenAndServe(":8080", nil)
}
