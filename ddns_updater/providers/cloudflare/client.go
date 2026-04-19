package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	BaseURL  string
	APIToken string
	HTTP     *http.Client
}

func NewClient(baseurl, apiToken string) *Client {
	return &Client{
		BaseURL:  baseurl,
		APIToken: apiToken,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) newRequest(ctx context.Context, method, endpointPath string, body any) (*http.Request, error) {

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

	req, err := http.NewRequest(method, endpointPath, bodyReader)
	if err != nil {
		log.Printf("error occurred while creating a new HTTP request: %s", err.Error())
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	log.Printf("%+v", req)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		log.Printf("error while performing request: %s", err.Error())
		return nil, err
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		log.Printf("error, cloudflare api returned non-successful response with status code: %d", resp.StatusCode)

		defer func() {
			if readerErr := resp.Body.Close(); readerErr != nil {
				err = errors.Join(err, fmt.Errorf("close response body: %s", readerErr.Error()))
			}
		}()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("error occurred while reading response")
			return nil, err
		}

		log.Println(string(respBody))

		var errorParsed CloudflareApiError
		if err := json.Unmarshal(respBody, &errorParsed); err != nil {
			return nil, fmt.Errorf("error unmarshalling cloudflare error response: %s", err.Error())
		}
		return nil, fmt.Errorf("cloudflare api returned unsuccessful response: %+v", errorParsed)
	}

	return resp, nil
}
