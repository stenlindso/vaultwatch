// Package history tracks audit run history and provides trend analysis
// across multiple snapshot comparisons over time.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry represents a single recorded audit run.
type Entry struct {
	ID          string    `json:"id"`
	Environment string    `json:"environment"`
	Timestamp   time.Time `json:"timestamp"`
	Added       []string  `json:"added"`
	Removed     []string  `json:"removed"`
	Violations  int       `json:"violations"`
}

// Store persists and retrieves audit history entries.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("history: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Record saves an Entry to disk.
func (s *Store) Record(e Entry) error {
	if e.ID == "" {
		e.ID = fmt.Sprintf("%d", e.Timestamp.UnixNano())
	}
	path := filepath.Join(s.dir, e.Environment)
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("history: create env dir: %w", err)
	}
	file := filepath.Join(path, e.ID+".json")
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	return os.WriteFile(file, data, 0o644)
}

// List returns all entries for the given environment, sorted by timestamp ascending.
func (s *Store) List(environment string) ([]Entry, error) {
	path := filepath.Join(s.dir, environment)
	entries, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: read dir: %w", err)
	}
	var results []Entry
	for _, de := range entries {
		if de.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(path, de.Name()))
		if err != nil {
			return nil, fmt.Errorf("history: read file %s: %w", de.Name(), err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: unmarshal %s: %w", de.Name(), err)
		}
		results = append(results, e)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.Before(results[j].Timestamp)
	})
	return results, nil
}
