package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot holds a captured list of Vault secret paths for a given environment.
type Snapshot struct {
	Environment string    `json:"environment"`
	CapturedAt  time.Time `json:"captured_at"`
	Paths       []string  `json:"paths"`
}

// New creates a new Snapshot for the given environment and paths.
func New(env string, paths []string) *Snapshot {
	return &Snapshot{
		Environment: env,
		CapturedAt:  time.Now().UTC(),
		Paths:       paths,
	}
}

// Save writes the snapshot as JSON to the given file path.
func Save(s *Snapshot, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from a JSON file.
func Load(filePath string) (*Snapshot, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}
