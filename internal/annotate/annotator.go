// Package annotate provides functionality for attaching and retrieving
// user-defined annotations (notes, labels, ownership) on Vault secret paths.
package annotate

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Annotation holds metadata attached to a secret path.
type Annotation struct {
	Path      string    `json:"path"`
	Env       string    `json:"env"`
	Note      string    `json:"note"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
}

// Store persists annotations to a JSON file keyed by env+path.
type Store struct {
	dir string
}

// NewStore returns a Store that persists data under dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) filePath(env string) string {
	return filepath.Join(s.dir, fmt.Sprintf("annotations_%s.json", env))
}

// Save writes the annotation for the given env, replacing any existing entry
// for the same path.
func (s *Store) Save(env string, a Annotation) error {
	if env == "" {
		return errors.New("annotate: env must not be empty")
	}
	if a.Path == "" {
		return errors.New("annotate: path must not be empty")
	}
	a.Env = env
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now().UTC()
	}

	existing, err := s.Load(env)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	updated := make(map[string]Annotation)
	for _, ann := range existing {
		updated[ann.Path] = ann
	}
	updated[a.Path] = a

	records := make([]Annotation, 0, len(updated))
	for _, v := range updated {
		records = append(records, v)
	}

	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath(env), data, 0o644)
}

// Load returns all annotations for the given env.
func (s *Store) Load(env string) ([]Annotation, error) {
	data, err := os.ReadFile(s.filePath(env))
	if err != nil {
		return nil, err
	}
	var records []Annotation
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("annotate: invalid JSON in store: %w", err)
	}
	return records, nil
}

// Get returns the annotation for a specific path in env, or an error if not found.
func (s *Store) Get(env, path string) (Annotation, error) {
	all, err := s.Load(env)
	if err != nil {
		return Annotation{}, err
	}
	for _, a := range all {
		if a.Path == path {
			return a, nil
		}
	}
	return Annotation{}, fmt.Errorf("annotate: no annotation found for path %q in env %q", path, env)
}
