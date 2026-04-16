package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Manager handles storage and retrieval of snapshots in a directory.
type Manager struct {
	dir string
}

// NewManager creates a Manager that stores snapshots under dir.
func NewManager(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("creating snapshot dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// SaveSnapshot persists a snapshot, naming the file after the environment.
func (m *Manager) SaveSnapshot(s *Snapshot) error {
	path := m.pathFor(s.Environment)
	return Save(s, path)
}

// LoadSnapshot loads the snapshot for the given environment.
func (m *Manager) LoadSnapshot(env string) (*Snapshot, error) {
	return Load(m.pathFor(env))
}

// ListEnvironments returns all environments that have saved snapshots.
func (m *Manager) ListEnvironments() ([]string, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return nil, err
	}
	var envs []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			envs = append(envs, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	return envs, nil
}

func (m *Manager) pathFor(env string) string {
	return filepath.Join(m.dir, env+".json")
}
