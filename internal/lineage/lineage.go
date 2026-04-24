// Package lineage tracks the origin and ancestry of secret paths across
// environments, enabling users to trace how paths have evolved over time.
package lineage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single lineage record for a secret path.
type Entry struct {
	Path        string    `json:"path"`
	Environment string    `json:"environment"`
	SeenAt      time.Time `json:"seen_at"`
	Source      string    `json:"source"` // snapshot label or "live"
}

// Record groups all entries observed for a given path.
type Record struct {
	Path    string  `json:"path"`
	History []Entry `json:"history"`
}

// Store persists lineage records to disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("lineage: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) filePath(env string) string {
	return filepath.Join(s.dir, env+".json")
}

// Record appends entries for the given environment.
func (s *Store) Record(env string, paths []string, source string) error {
	existing, _ := s.Load(env)
	now := time.Now().UTC()
	for _, p := range paths {
		existing = append(existing, Entry{
			Path:        p,
			Environment: env,
			SeenAt:      now,
			Source:      source,
		})
	}
	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return fmt.Errorf("lineage: marshal: %w", err)
	}
	return os.WriteFile(s.filePath(env), data, 0o644)
}

// Load returns all entries recorded for env.
func (s *Store) Load(env string) ([]Entry, error) {
	data, err := os.ReadFile(s.filePath(env))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("lineage: read: %w", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("lineage: unmarshal: %w", err)
	}
	return entries, nil
}
