package resolve

import (
	"encoding/json"
	"fmt"
	"os"
)

// aliasFile is the expected on-disk structure for an alias config file.
type aliasFile struct {
	Aliases []Alias `json:"aliases"`
}

// LoadAliases reads alias definitions from a JSON file and returns a Resolver.
// The file must contain a top-level "aliases" array.
func LoadAliases(path string) (*Resolver, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("resolve: open aliases file: %w", err)
	}
	defer f.Close()

	var af aliasFile
	if err := json.NewDecoder(f).Decode(&af); err != nil {
		return nil, fmt.Errorf("resolve: decode aliases file: %w", err)
	}
	return New(af.Aliases), nil
}

// DefaultAliases returns a Resolver pre-loaded with common Vault path aliases
// used as a fallback when no alias file is present.
func DefaultAliases() *Resolver {
	return New([]Alias{
		{Name: "secret", Prefix: "secret"},
		{Name: "kv", Prefix: "kv/data"},
	})
}
