package relpath_test

import (
	"encoding/json"
	"os"

	"github.com/samherrmann/go-sanbox/relpath"
)

func Example() {
	configPath := "path/to/config.json"
	config := &Config{}

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	if err := relpath.Lock(configPath); err != nil {
		panic(err)
	}
	defer relpath.Unlock()

	if err := json.Unmarshal(bytes, config); err != nil {
		panic(err)
	}
}

type Config struct {
	FooPath relpath.Path `json:"fooPath"`
	Bars    []BarConfig  `json:"bars"`
}

type BarConfig struct {
	Path relpath.Path `json:"path"`
}
