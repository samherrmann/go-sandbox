package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if err := app(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func app() error {

	config, err := readConfig("config.json")
	if err != nil {
		return err
	}

	query := `{ __type(name: "Order") { fields { name } } }`

	res, err := sendGraphQL(config.Shop, config.AccessToken, query)
	if err != nil {
		return fmt.Errorf("sending GraphQL request: %w", err)
	}

	body, err := io.ReadAll(res)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}
	defer res.Close()

	fmt.Println(string(body))

	return nil
}

func sendGraphQL(shop string, accessToken string, query string) (io.ReadCloser, error) {
	url := fmt.Sprintf(
		"https://%s.myshopify.com/admin/api/unstable/graphql.json",
		shop,
	)

	payload := fmt.Sprintf(`{"query": %q}`, query)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Access-Token", accessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send HTTP request: %w", err)
	}

	if res.StatusCode > 399 {
		return nil, fmt.Errorf(res.Status)
	}

	return res.Body, nil
}
