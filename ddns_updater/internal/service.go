package internal

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Service struct {
	Provider     Provider
	Resolver     IPResolver
	State        StateStore
	RecordConfig RecordConfig
}

func (s *Service) RunOnce(ctx context.Context) error {
	ip, err := s.Resolver.Resolve(ctx)
	if err != nil {
		return fmt.Errorf("resolve current public ip: %w", err)
	}

	currentIP := ip.String()

	prevState, err := s.State.Load(ctx)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	if prevState != nil && prevState.LastIP == currentIP {
		log.Printf("ip unchanged, skipping dns update: ip=%s", currentIP)
		return nil
	}

	record := Record{
		Zone:  s.RecordConfig.Zone,
		Name:  s.RecordConfig.Name,
		Type:  s.RecordConfig.Type,
		TTL:   s.RecordConfig.TTL,
		Id:    s.RecordConfig.Id,
		Value: currentIP,
	}

	if err := s.Provider.UpdateRecord(ctx, record); err != nil {
		return fmt.Errorf("update dns record: %w", err)
	}

	newState := &State{
		LastIP:    currentIP,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.State.Save(ctx, newState); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	log.Printf("dns record updated successfully: ip=%s", currentIP)
	return nil
}
