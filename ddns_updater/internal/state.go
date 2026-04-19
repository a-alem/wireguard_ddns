package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FileStateStore struct {
	Path string
}

func NewFileStateStore(path string) *FileStateStore {
	return &FileStateStore{
		Path: path,
	}
}

func (s *FileStateStore) Load(ctx context.Context) (*State, error) {
	_ = ctx

	data, err := os.ReadFile(s.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read state file: %w", err)
	}

	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("unmarshal state file: %w", err)
	}

	return &st, nil
}

func (s *FileStateStore) Save(ctx context.Context, st *State) error {
	_ = ctx

	if st == nil {
		return errors.New("state is nil")
	}

	dir := filepath.Dir(s.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}

	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}

	tmpPath := s.Path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
		return fmt.Errorf("write temp state file: %w", err)
	}

	if err := os.Rename(tmpPath, s.Path); err != nil {
		return fmt.Errorf("rename temp state file: %w", err)
	}

	return nil
}
