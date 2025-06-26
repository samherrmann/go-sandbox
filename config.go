package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Shop        string `json:"shop"`
	AccessToken string `json:"accessToken"`
}

func readConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}

	return &config, nil
}
