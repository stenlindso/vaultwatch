// Package baseline provides functionality to establish and compare
// Vault secret path baselines across environments.
package baseline

import (
	"fmt"
	"time"

	"github.com/youorg/vaultwatch/internal/snapshot"
)

// Baseline represents a named reference point for a set of secret paths
// in a given environment.
type Baseline struct {
	Environment string    `json:"environment"`
	Label       string    `json:"label"`
	Paths       []string  `json:"paths"`
	CreatedAt   time.Time `json:"created_at"`
}

// Manager handles saving and loading baselines using the snapshot manager.
type Manager struct {
	snaps *snapshot.Manager
}

// NewManager returns a Manager backed by the given snapshot.Manager.
func NewManager(m *snapshot.Manager) *Manager {
	return &Manager{snaps: m}
}

// Save persists a baseline snapshot under a labelled key.
func (m *Manager) Save(env, label string, paths []string) (*Baseline, error) {
	if env == "" {
		return nil, fmt.Errorf("baseline: environment must not be empty")
	}
	if label == "" {
		return nil, fmt.Errorf("baseline: label must not be empty")
	}

	b := &Baseline{
		Environment: env,
		Label:       label,
		Paths:       paths,
		CreatedAt:   time.Now().UTC(),
	}

	key := baselineKey(env, label)
	if err := m.snaps.Save(key, paths); err != nil {
		return nil, fmt.Errorf("baseline: save %q: %w", key, err)
	}
	return b, nil
}

// Load retrieves a previously saved baseline.
func (m *Manager) Load(env, label string) ([]string, error) {
	key := baselineKey(env, label)
	paths, err := m.snaps.Load(key)
	if err != nil {
		return nil, fmt.Errorf("baseline: load %q: %w", key, err)
	}
	return paths, nil
}

// Diff returns paths added or removed relative to the saved baseline.
func (m *Manager) Diff(env, label string, current []string) (added, removed []string, err error) {
	base, err := m.Load(env, label)
	if err != nil {
		return nil, nil, err
	}

	baseSet := toSet(base)
	curSet := toSet(current)

	for p := range curSet {
		if !baseSet[p] {
			added = append(added, p)
		}
	}
	for p := range baseSet {
		if !curSet[p] {
			removed = append(removed, p)
		}
	}
	return added, removed, nil
}

func baselineKey(env, label string) string {
	return fmt.Sprintf("baseline_%s_%s", env, label)
}

func toSet(paths []string) map[string]bool {
	s := make(map[string]bool, len(paths))
	for _, p := range paths {
		s[p] = true
	}
	return s
}
