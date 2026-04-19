package cloudflare

type Config struct {
	BaseUrl  string
	APIToken string
	ZoneID   string
	Proxied  *bool `yaml:"proxied,omitempty"`
}

type UpdateDNSRecordRequest struct {
	Name           string `json:"name"`
	TTL            int64  `json:"ttl"`
	Type           string `json:"type"`
	Comment        string `json:"comment,omitempty"`
	Content        string `json:"content,omitempty"`
	Proxied        *bool  `json:"proxied,omitempty"`
	PrivateRouting *bool  `json:"private_routing,omitempty"`
}

type CloudflareApiGenericResponse[T any] struct {
	Errors   []CloudflareApiError   `json:"errors"`
	Messages []CloudflareApiMessage `json:"messages"`
	Success  *bool                  `json:"success"`
	Result   T                      `json:"result"`
}

type CloudflareApiResultDnsRecordResponse struct {
}

type CloudflareApiMessage struct {
	Code             int64               `json:"code"`
	Message          string              `json:"message"`
	DocumentationUrl string              `json:"documentation_url"`
	Source           CloudflareApiSource `json:"source"`
}

type CloudflareApiSource struct {
	Pointer string `json:"pointer"`
}

type CloudflareApiError struct {
	Code             int64               `json:"code"`
	Message          string              `json:"message"`
	DocumentationUrl string              `json:"documentation_url"`
	Source           CloudflareApiSource `json:"source"`
}
