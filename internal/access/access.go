// Package access provides path-level access control tracking for Vault secrets.
// It records which environments have read access to specific secret paths
// and can detect when access patterns change between snapshots.
package access

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Record represents an access entry for a secret path in an environment.
type Record struct {
	Environment string    `json:"environment"`
	Path        string    `json:"path"`
	Accessible  bool      `json:"accessible"`
	CheckedAt   time.Time `json:"checked_at"`
}

// Report summarises access changes between two snapshots.
type Report struct {
	Environment string    `json:"environment"`
	Gained      []string  `json:"gained"`
	Lost        []string  `json:"lost"`
	GeneratedAt time.Time `json:"generated_at"`
}

// HasChanges returns true when any access has been gained or lost.
func (r Report) HasChanges() bool {
	return len(r.Gained) > 0 || len(r.Lost) > 0
}

// Store persists access records on disk as JSON files.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("access: create store dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes records for the given environment to disk.
func (s *Store) Save(env string, records []Record) error {
	if env == "" {
		return fmt.Errorf("access: environment must not be empty")
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("access: marshal records: %w", err)
	}
	return os.WriteFile(s.filePath(env), data, 0o644)
}

// Load reads access records for the given environment from disk.
func (s *Store) Load(env string) ([]Record, error) {
	data, err := os.ReadFile(s.filePath(env))
	if err != nil {
		return nil, fmt.Errorf("access: load %s: %w", env, err)
	}
	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("access: unmarshal %s: %w", env, err)
	}
	return records, nil
}

// Diff compares two sets of records and returns an access Report.
func Diff(env string, previous, current []Record) Report {
	prev := accessibleSet(previous)
	curr := accessibleSet(current)

	report := Report{
		Environment: env,
		GeneratedAt: time.Now().UTC(),
	}
	for p := range curr {
		if !prev[p] {
			report.Gained = append(report.Gained, p)
		}
	}
	for p := range prev {
		if !curr[p] {
			report.Lost = append(report.Lost, p)
		}
	}
	sort.Strings(report.Gained)
	sort.Strings(report.Lost)
	return report
}

func (s *Store) filePath(env string) string {
	return filepath.Join(s.dir, env+"_access.json")
}

func accessibleSet(records []Record) map[string]bool {
	out := make(map[string]bool, len(records))
	for _, r := range records {
		if r.Accessible {
			out[r.Path] = true
		}
	}
	return out
}
