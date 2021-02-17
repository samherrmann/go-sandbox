package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var embeddedFS embed.FS

func main() {
	serverRoot, err := fs.Sub(embeddedFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(serverRoot)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
