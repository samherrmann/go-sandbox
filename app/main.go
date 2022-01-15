package main

import (
	"fmt"

	"github.com/samherrmann/go-sanbox/app/config"
)

func main() {
	configPath := "path/to/config.json"

	c, err := config.ReadFile(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Use the paths from the configuration.
	fmt.Println(c.FooPath.String())
	for _, b := range c.Bars {
		fmt.Println(b.Path.String())
	}

	if err := config.WriteFile(configPath, c); err != nil {
		fmt.Println(err)
	}
}
