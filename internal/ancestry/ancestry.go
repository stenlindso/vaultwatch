// Package ancestry tracks parent-child relationships between secret paths,
// enabling callers to walk derivation chains and detect orphaned secrets.
package ancestry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Edge represents a directed parent → child relationship.
type Edge struct {
	Parent    string    `json:"parent"`
	Child     string    `json:"child"`
	Env       string    `json:"env"`
	RecordedAt time.Time `json:"recorded_at"`
}

// Graph holds all recorded edges for a store.
type Graph struct {
	Edges []Edge `json:"edges"`
}

// Store persists ancestry edges to disk.
type Store struct {
	dir string
}

// NewStore returns a Store rooted at dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("ancestry: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) filePath(env string) string {
	return filepath.Join(s.dir, env+".json")
}

// Record appends a parent→child edge for the given environment.
func (s *Store) Record(env, parent, child string) error {
	if env == "" || parent == "" || child == "" {
		return fmt.Errorf("ancestry: env, parent, and child must not be empty")
	}
	g, _ := s.Load(env)
	g.Edges = append(g.Edges, Edge{
		Parent:     parent,
		Child:      child,
		Env:        env,
		RecordedAt: time.Now().UTC(),
	})
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return fmt.Errorf("ancestry: marshal: %w", err)
	}
	return os.WriteFile(s.filePath(env), data, 0o644)
}

// Load returns the Graph for env, or an empty Graph if none exists.
func (s *Store) Load(env string) (Graph, error) {
	data, err := os.ReadFile(s.filePath(env))
	if os.IsNotExist(err) {
		return Graph{}, nil
	}
	if err != nil {
		return Graph{}, fmt.Errorf("ancestry: read: %w", err)
	}
	var g Graph
	if err := json.Unmarshal(data, &g); err != nil {
		return Graph{}, fmt.Errorf("ancestry: unmarshal: %w", err)
	}
	return g, nil
}

// Children returns all direct children of parent in env.
func (s *Store) Children(env, parent string) ([]string, error) {
	g, err := s.Load(env)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range g.Edges {
		if e.Parent == parent {
			out = append(out, e.Child)
		}
	}
	sort.Strings(out)
	return out, nil
}

// Orphans returns paths that appear as children but never as parents,
// and are not present in the provided activePaths set.
func (s *Store) Orphans(env string, activePaths []string) ([]string, error) {
	g, err := s.Load(env)
	if err != nil {
		return nil, err
	}
	active := make(map[string]bool, len(activePaths))
	for _, p := range activePaths {
		active[p] = true
	}
	parents := make(map[string]bool)
	for _, e := range g.Edges {
		parents[e.Parent] = true
	}
	seen := make(map[string]bool)
	var out []string
	for _, e := range g.Edges {
		if !parents[e.Child] && !active[e.Child] && !seen[e.Child] {
			out = append(out, e.Child)
			seen[e.Child] = true
		}
	}
	sort.Strings(out)
	return out, nil
}
