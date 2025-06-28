package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

	req := &GraphQLRequest{
		Query: string(query),
		Variables: map[string]string{
			"name": orderNumber,
		},
	}

	res, err := sendGraphQL(config.Shop, config.AccessToken, req)
	if err != nil {
		return fmt.Errorf("sending GraphQL request: %w", err)
	}

	if _, err := io.Copy(os.Stdout, res); err != nil {
		return fmt.Errorf("reading response: %w", err)
	}
	defer res.Close()

	return nil
}

func sendGraphQL(shop string, accessToken string, req *GraphQLRequest) (io.ReadCloser, error) {
	url := fmt.Sprintf(
		"https://%s.myshopify.com/admin/api/unstable/graphql.json",
		shop,
	)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Shopify-Access-Token", accessToken)

	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send HTTP request: %w", err)
	}

	if res.StatusCode > 399 {
		return nil, fmt.Errorf(res.Status)
	}

	return res.Body, nil
}

type GraphQLRequest struct {
	Query     string `json:"query"`
	Variables any    `json:"variables"`
}
