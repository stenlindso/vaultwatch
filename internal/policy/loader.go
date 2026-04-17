package policy

import (
	"encoding/json"
	"fmt"
	"os"
)

// PolicyFile represents the JSON structure of a policy file.
type PolicyFile struct {
	Rules []Rule `json:"rules"`
}

// LoadFromFile reads and parses a policy JSON file from disk.
func LoadFromFile(path string) ([]Rule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open policy file: %w", err)
	}
	defer f.Close()

	var pf PolicyFile
	if err := json.NewDecoder(f).Decode(&pf); err != nil {
		return nil, fmt.Errorf("decode policy file: %w", err)
	}
	return pf.Rules, nil
}

// DefaultRules returns a minimal set of sensible default rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Pattern:  `(?i)^secret/.*/temp/`,
			Deny:     true,
		},
		{
			Pattern:  `(?i)test|tmp|scratch`,
			Deny:     true,
		},
	}
}
