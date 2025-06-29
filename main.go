package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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
			"query": "name:" + orderNumber,
		},
	}

	res, err := graphql.Send(config.Shop, config.AccessToken, req)
	if err != nil {
		return fmt.Errorf("sending GraphQL request: %w", err)
	}
	defer res.Close()

	orderTpl, err := template.ParseFiles("invoice.tpl.html")
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		fmt.Sprintf("invoice-%v.html", orderNumber),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(res)
	if err != nil {
		return err
	}

	m := map[string]any{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	return orderTpl.Execute(file, m)
}
