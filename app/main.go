package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/samherrmann/go-sanbox/lib"
)

func main() {
	if err := lib.SomeFunc(); err != nil {
		fmt.Println(err)                               // foobar: invalid syntax
		fmt.Println(lib.IsFoobar(err))                 // true
		fmt.Println(errors.Is(err, strconv.ErrSyntax)) // true
	}
}
