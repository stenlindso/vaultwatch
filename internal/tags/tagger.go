// Package tags provides functionality for tagging and filtering secret paths
// based on user-defined labels stored alongside snapshots.
package tags

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// TagStore manages path-to-tag mappings for a given environment.
type TagStore struct {
	dir string
}

// TagEntry maps a secret path to a set of tags.
type TagEntry struct {
	Path string   `json:"path"`
	Tags []string `json:"tags"`
}

// NewTagStore creates a TagStore that persists data under dir.
func NewTagStore(dir string) (*TagStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("tags: create dir: %w", err)
	}
	return &TagStore{dir: dir}, nil
}

func (s *TagStore) filePath(env string) string {
	return filepath.Join(s.dir, env+".tags.json")
}

// Save persists the tag entries for the given environment.
func (s *TagStore) Save(env string, entries []TagEntry) error {
	if env == "" {
		return fmt.Errorf("tags: env must not be empty")
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("tags: marshal: %w", err)
	}
	return os.WriteFile(s.filePath(env), data, 0o644)
}

// Load retrieves tag entries for the given environment.
func (s *TagStore) Load(env string) ([]TagEntry, error) {
	data, err := os.ReadFile(s.filePath(env))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("tags: no tag file for environment %q", env)
		}
		return nil, fmt.Errorf("tags: read: %w", err)
	}
	var entries []TagEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("tags: unmarshal: %w", err)
	}
	return entries, nil
}

// FilterByTag returns only those paths that carry the given tag.
func FilterByTag(entries []TagEntry, tag string) []string {
	var result []string
	for _, e := range entries {
		for _, t := range e.Tags {
			if t == tag {
				result = append(result, e.Path)
				break
			}
		}
	}
	sort.Strings(result)
	return result
}

// BuildEntries creates TagEntry records from a path list, assigning each the
// provided tags.
func BuildEntries(paths []string, tags []string) []TagEntry {
	entries := make([]TagEntry, 0, len(paths))
	for _, p := range paths {
		entries = append(entries, TagEntry{Path: p, Tags: tags})
	}
	return entries
}
