package internal

import (
	"context"
	"net"
	"time"
)

type AppConfig struct {
	Record    RecordConfig    `yaml:"record"`
	Providers ProvidersConfig `yaml:"providers"`
	Resolvers ResolversConfig `yaml:"resolvers"`
	State     StateConfig     `yaml:"state"`
}

type ProvidersConfig struct {
	Cloudflare CloudflareConfig `yaml:"cloudflare"`
}

type ResolversConfig struct {
	IPify IPifyResolverConfig `yaml:"ipify"`
}

type IPifyResolverConfig struct {
	URL            string `yaml:"url"`
	ResponseFormat string `yaml:"response_format"`
}

type StateConfig struct {
	Path string `yaml:"path"`
}

type RecordType string

const (
	RecordTypeA    RecordType = "A"
	RecordTypeAAAA RecordType = "AAAA"
)

type Record struct {
	Zone  string
	Name  string
	Id    string
	Type  RecordType
	TTL   int64
	Value string
}

// Provider is implemented by each DNS provider backend.
type Provider interface {
	UpdateRecord(ctx context.Context, record Record) error
}

// IPResolver resolves the current public IP address.
type IPResolver interface {
	Resolve(ctx context.Context) (net.IP, error)
}

// StateStore persists the last successfully applied IP.
// This avoids unnecessary DNS updates.
type StateStore interface {
	Load(ctx context.Context) (*State, error)
	Save(ctx context.Context, state *State) error
}

// State stores the last successful DDNS sync state.
type State struct {
	LastIP    string    `json:"last_ip"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Config is the runtime config needed by the sync service.
type RecordConfig struct {
	Zone string
	Name string
	Id   string
	Type RecordType
	TTL  int64
}

// Providers
type CloudflareConfig struct {
	APIToken string `yaml:"api_token"`
	ZoneID   string `yaml:"zone_id"`
	Proxied  *bool  `yaml:"proxied,omitempty"`
	BaseURL  string `yaml:"base_url,omitempty"`
}
