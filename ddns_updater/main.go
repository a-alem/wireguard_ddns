package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/a-alem/wireguard_ddns/ddns_updater/internal"
	"github.com/a-alem/wireguard_ddns/ddns_updater/providers/cloudflare"
	"github.com/a-alem/wireguard_ddns/ddns_updater/resolvers/ipify"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <config-file>", os.Args[0])
	}

	configPath := os.Args[1]

	cfg, err := internal.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	provider := cloudflare.New(cloudflare.Config{
		BaseUrl:  cfg.Providers.Cloudflare.BaseURL,
		APIToken: cfg.Providers.Cloudflare.APIToken,
		ZoneID:   cfg.Providers.Cloudflare.ZoneID,
		Proxied:  cfg.Providers.Cloudflare.Proxied,
	})

	resolver := ipify.New(ipify.Config{
		ResolverUrl:    cfg.Resolvers.IPify.URL,
		ResponseFormat: cfg.Resolvers.IPify.ResponseFormat,
	})

	stateStore := internal.NewFileStateStore(cfg.State.Path)

	service := &internal.Service{
		Provider: provider,
		Resolver: resolver,
		State:    stateStore,
		RecordConfig: internal.RecordConfig{
			Zone: cfg.Record.Zone,
			Name: cfg.Record.Name,
			Type: cfg.Record.Type,
			TTL:  cfg.Record.TTL,
			Id:   cfg.Record.Id,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := service.RunOnce(ctx); err != nil {
		log.Fatalf("ddns updater failed: %v", err)
	}

	log.Println("ddns updater completed successfully")
}
