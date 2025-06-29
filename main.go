package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/samherrmann/go-sandbox/graphql"
)

func main() {
	if err := app(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func app() error {

	if len(os.Args) < 2 {
		return fmt.Errorf("missing order number")
	}

	orderNumber := os.Args[1]
	if !strings.HasPrefix(orderNumber, "#") {
		orderNumber = "#" + orderNumber
	}

	config, err := readConfig("config.json")
	if err != nil {
		return err
	}

	query, err := os.ReadFile("query.graphql")
	if err != nil {
		return err
	}

	req := &graphql.Request{
		Query: string(query),
		Variables: map[string]string{
			"name": orderNumber,
		},
	}

	res, err := graphql.Send(config.Shop, config.AccessToken, req)
	if err != nil {
		return fmt.Errorf("sending GraphQL request: %w", err)
	}

	if _, err := io.Copy(os.Stdout, res); err != nil {
		return fmt.Errorf("reading response: %w", err)
	}
	defer res.Close()

	return nil
}
