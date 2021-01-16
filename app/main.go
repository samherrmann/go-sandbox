package main

import (
	"errors"
	"log"

	"github.com/samherrmann/go-sanbox/mylib"
)

func main() {
	err := mylib.DoSomething()
	var uaErr UnauthorizedError
	if errors.As(err, &uaErr) {
		log.Fatalf("%v", err)
	}
	if err != nil {
		log.Fatalf("an unknown error occured: %v", err)
	}
}

// UnauthorizedError ...
type UnauthorizedError interface {
	Unauthorized() bool
	Error() string
}
