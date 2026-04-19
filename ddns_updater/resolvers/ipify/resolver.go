package ipify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type Resolver struct {
	resolverUrl    string
	client         *Client
	responseFormat string
}

func New(cfg Config) *Resolver {
	return &Resolver{
		responseFormat: cfg.ResponseFormat,
		resolverUrl:    cfg.ResolverUrl,
		client:         NewClient(cfg.ResolverUrl, cfg.APIToken),
	}
}

func (r *Resolver) Resolve(ctx context.Context) (net.IP, error) {
	// Resolve and return IP
	req, err := r.client.newRequest(ctx, http.MethodGet, r.resolverUrl, r.responseFormat, nil)
	if err != nil {
		log.Println("error occurred while building IP resolution request")
		return nil, err
	}

	resp, err := r.client.do(req)
	if err != nil {
		log.Println("error occurred while calling ipify IP resolution API")
		return nil, err
	}

	defer func() {
		if readerErr := resp.Body.Close(); readerErr != nil {
			err = errors.Join(err, fmt.Errorf("close response body: %s", readerErr.Error()))
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error occurred while reading response from ipify")
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println("error, ipify api returned unsuccessful response")
		return nil, fmt.Errorf("cloudflare api returned unsuccessful response: %s", string(respBody))
	}

	var parsedResp ResolveResponse
	if err := json.Unmarshal(respBody, &parsedResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling ipify response: %s", err.Error())
	}

	ip := net.ParseIP(parsedResp.IP)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse IP from response: %s", parsedResp.IP)
	}

	return ip, nil
}
