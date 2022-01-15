package config

import (
	"encoding/json"
	"os"

	"github.com/samherrmann/go-sanbox/relpath"
)

type Config struct {
	FooPath relpath.Path `json:"fooPath"`
	Bars    []BarConfig  `json:"bars"`
}

type BarConfig struct {
	Path relpath.Path `json:"path"`
}

func ReadFile(path string) (*Config, error) {
	config := &Config{}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := relpath.Lock(path); err != nil {
		return nil, err
	}
	defer relpath.Unlock()
	if err := json.Unmarshal(bytes, config); err != nil {
		return nil, err
	}
	return config, nil
}

func WriteFile(path string, config *Config) error {
	if err := relpath.Lock(path); err != nil {
		return err
	}
	defer relpath.Unlock()
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, os.ModePerm)
}
