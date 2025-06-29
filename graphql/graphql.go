package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Send(shop string, accessToken string, req *Request) (io.ReadCloser, error) {
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

type Request struct {
	Query     string `json:"query"`
	Variables any    `json:"variables"`
}
