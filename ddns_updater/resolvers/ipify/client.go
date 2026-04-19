package ipify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	resolverURL string
	apiToken    string
	httpClient  *http.Client
}

func NewClient(resolverUrl, apiToken string) *Client {
	return &Client{
		resolverURL: resolverUrl,
		apiToken:    apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) newRequest(ctx context.Context, method, endpointPath, responseFormat string, body any) (*http.Request, error) {

	var bodyReader *bytes.Reader
	// Check if request has body
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			log.Printf("error occurred while marshalling HTTP request body: %s", err.Error())
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBytes)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, endpointPath+"?format="+responseFormat, bodyReader)
	if err != nil {
		log.Printf("error occurred while creating a new HTTP request: %s", err.Error())
		return nil, err
	}

	if strings.TrimSpace(c.apiToken) != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiToken)
	}

	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("error while performing request: %s", err.Error())
		return nil, err
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		log.Printf("error, ipify api returned non-successful response with status code: %d", resp.StatusCode)
		return nil, errors.New("error, cloudflare api returned non-successful response")
	}

	return resp, nil
}
