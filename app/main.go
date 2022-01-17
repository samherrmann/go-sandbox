package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func main() {
	configPath := "path/to/config.json"

	// Change working directory to the directory in which the config file is
	// located. That way all relative paths within the config file will be
	// evaluated correctly.
	if err := os.Chdir(filepath.Dir(configPath)); err != nil {
		panic(err)
	}

	bytes, err := os.ReadFile(filepath.Base(configPath))
	if err != nil {
		panic(err)
	}

	config := &Config{}
	if err := json.Unmarshal(bytes, config); err != nil {
		panic(err)
	}

	bytes, err = os.ReadFile(config.Path)
	if err != nil {
		panic(err)
	}

	println(string(bytes))
}

type Config struct {
	Path string `json:"path"`
}
