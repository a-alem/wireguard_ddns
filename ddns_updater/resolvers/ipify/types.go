package ipify

type Config struct {
	ResolverUrl    string
	APIToken       string
	ResponseFormat string
}

type ResolveResponse struct {
	IP string `json:"ip"`
}
