package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	data = []byte(os.ExpandEnv(string(data)))

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *AppConfig) Validate() error {
	if strings.TrimSpace(c.Record.Zone) == "" {
		return errors.New("record.zone is required")
	}

	if strings.TrimSpace(c.Record.Name) == "" {
		return errors.New("record.name is required")
	}

	switch c.Record.Type {
	case RecordTypeA, RecordTypeAAAA:
	default:
		return fmt.Errorf("record.type must be one of: %s, %s", RecordTypeA, RecordTypeAAAA)
	}

	if c.Record.TTL <= 0 {
		return errors.New("record.ttl must be greater than 0")
	}

	if strings.TrimSpace(c.Providers.Cloudflare.APIToken) == "" {
		return errors.New("providers.cloudflare.api_token is required")
	}

	if strings.TrimSpace(c.Providers.Cloudflare.ZoneID) == "" {
		return errors.New("providers.cloudflare.zone_id is required")
	}

	if strings.TrimSpace(c.Providers.Cloudflare.BaseURL) == "" {
		c.Providers.Cloudflare.BaseURL = "https://api.cloudflare.com/client/v4"
	}

	if strings.TrimSpace(c.Resolvers.IPify.URL) == "" {
		return errors.New("resolvers.ipify.url is required")
	}

	switch strings.ToLower(strings.TrimSpace(c.Resolvers.IPify.ResponseFormat)) {
	case "json", "text", "":
		// allow empty and treat later as default
	default:
		return fmt.Errorf("resolvers.ipify.response_format must be either 'json' or 'text'")
	}

	if strings.TrimSpace(c.State.Path) == "" {
		return errors.New("state.path is required")
	}

	return nil
}
