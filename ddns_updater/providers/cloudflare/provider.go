package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/a-alem/wireguard_ddns/ddns_updater/internal"
)

type Provider struct {
	client  *Client
	zoneID  string
	proxied *bool
}

func New(cfg Config) *Provider {
	return &Provider{
		client:  NewClient(cfg.BaseUrl, cfg.APIToken),
		zoneID:  cfg.ZoneID,
		proxied: cfg.Proxied,
	}
}

func (p *Provider) UpdateRecord(ctx context.Context, record internal.Record) error {
	if err := p.validateRecord(record); err != nil {
		log.Println("error occurred while validating DNS record prior to updating it")
		return err
	}

	path := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", p.zoneID, record.Id)

	reqBody := UpdateDNSRecordRequest{
		Content: record.Value,
		TTL:     record.TTL,
		Type:    string(record.Type),
		Proxied: p.proxied,
	}

	req, err := p.client.newRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		log.Println("error occurred while building DNS record update request")
		return err
	}

	resp, err := p.client.do(req)
	if err != nil {
		log.Printf("error occurred while calling cloudflare DNS update API: %s", err.Error())
		return err
	}

	defer func() {
		if readerErr := resp.Body.Close(); readerErr != nil {
			err = errors.Join(err, fmt.Errorf("close response body: %s", readerErr.Error()))
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error occurred while reading response")
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println("error, cloudflare api returned unsuccessful response")
		var errorParsed CloudflareApiError
		if err := json.Unmarshal(respBody, &errorParsed); err != nil {
			return fmt.Errorf("error unmarshalling cloudflare error response: %s", err.Error())
		}
		return fmt.Errorf("cloudflare api returned unsuccessful response: %+v", errorParsed)
	}

	return nil
}

func (p *Provider) validateRecord(record internal.Record) error {
	if strings.TrimSpace(p.zoneID) == "" {
		return errors.New("zoneID must have a non whitespace only or empty string value")
	}

	if strings.TrimSpace(record.Name) == "" {
		return errors.New("record name must have a non whitespace only or empty string value")
	}

	if parsedIP := net.ParseIP(record.Value); parsedIP == nil {
		return errors.New("record value must be a valid IPv4 address for A type record")
	}

	if record.TTL < 30 {
		// More here: https://developers.cloudflare.com/api/resources/dns/subresources/records/methods/update#(resource)%20dns.records%20%3E%20(model)%20a_record%20%3E%20(schema)%20%3E%20(property)%20ttl
		return errors.New("record ttl must be a at least 30 seconds in accordance with cloudflare DNS API docs")
	}

	return nil
}
