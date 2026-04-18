package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the top-level vaultwatch configuration.
type Config struct {
	Environments map[string]EnvConfig `json:"environments"`
	SnapshotDir  string               `json:"snapshot_dir"`
	PolicyFile   string               `json:"policy_file"`
}

// EnvConfig holds per-environment Vault connection settings.
type EnvConfig struct {
	Address   string `json:"address"`
	Token     string `json:"token"`
	MountPath string `json:"mount_path"`
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate performs basic sanity checks on the loaded config.
func (c *Config) validate() error {
	if len(c.Environments) == 0 {
		return fmt.Errorf("config: at least one environment must be defined")
	}
	for name, env := range c.Environments {
		if env.Address == "" {
			return fmt.Errorf("config: environment %q missing address", name)
		}
	}
	return nil
}

// EnvNames returns a sorted-stable slice of environment names.
func (c *Config) EnvNames() []string {
	names := make([]string, 0, len(c.Environments))
	for k := range c.Environments {
		names = append(names, k)
	}
	return names
}
