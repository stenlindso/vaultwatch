// Package ownership tracks team or service ownership of Vault secret paths.
package ownership

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Record associates a secret path with an owner and optional metadata.
type Record struct {
	Path    string `json:"path"`
	Owner   string `json:"owner"`
	Team    string `json:"team,omitempty"`
	Contact string `json:"contact,omitempty"`
}

// Store holds ownership records keyed by path.
type Store struct {
	file    string
	records map[string]Record
}

// NewStore creates a Store backed by the given file.
func NewStore(file string) *Store {
	return &Store{file: file, records: make(map[string]Record)}
}

// Set registers or updates the ownership record for a path.
func (s *Store) Set(r Record) error {
	if strings.TrimSpace(r.Path) == "" {
		return errors.New("ownership: path must not be empty")
	}
	if strings.TrimSpace(r.Owner) == "" {
		return errors.New("ownership: owner must not be empty")
	}
	s.records[r.Path] = r
	return nil
}

// Get returns the ownership record for a path, if present.
func (s *Store) Get(path string) (Record, bool) {
	r, ok := s.records[path]
	return r, ok
}

// Unowned returns paths from the provided list that have no ownership record.
func (s *Store) Unowned(paths []string) []string {
	var out []string
	for _, p := range paths {
		if _, ok := s.records[p]; !ok {
			out = append(out, p)
		}
	}
	return out
}

// All returns all records sorted by path.
func (s *Store) All() []Record {
	out := make([]Record, 0, len(s.records))
	for _, r := range s.records {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// Save persists the store to disk.
func (s *Store) Save() error {
	data, err := json.MarshalIndent(s.All(), "", "  ")
	if err != nil {
		return fmt.Errorf("ownership: marshal: %w", err)
	}
	return os.WriteFile(s.file, data, 0o644)
}

// Load reads records from disk into the store.
func (s *Store) Load() error {
	data, err := os.ReadFile(s.file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("ownership: read: %w", err)
	}
	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("ownership: unmarshal: %w", err)
	}
	for _, r := range records {
		s.records[r.Path] = r
	}
	return nil
}
