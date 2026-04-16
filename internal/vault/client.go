package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with environment metadata.
type Client struct {
	api     *vaultapi.Client
	EnvName string
}

// Config holds connection parameters for a Vault instance.
type Config struct {
	Address string
	Token   string
	EnvName string
}

// NewClient creates a new authenticated Vault client.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		cfg.Address = os.Getenv("VAULT_ADDR")
	}
	if cfg.Token == "" {
		cfg.Token = os.Getenv("VAULT_TOKEN")
	}
	if cfg.Address == "" {
		return nil, fmt.Errorf("vault address is required (set VAULT_ADDR or pass --address)")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("vault token is required (set VAULT_TOKEN or pass --token)")
	}

	vcfg := vaultapi.DefaultConfig()
	vcfg.Address = cfg.Address

	api, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}
	api.SetToken(cfg.Token)

	return &Client{api: api, EnvName: cfg.EnvName}, nil
}

// ListPaths recursively lists all secret paths under the given prefix (KV v2).
func (c *Client) ListPaths(mountPath, prefix string) ([]string, error) {
	var paths []string

	listPath := fmt.Sprintf("%s/metadata/%s", mountPath, prefix)
	secret, err := c.api.Logical().List(listPath)
	if err != nil {
		return nil, fmt.Errorf("listing %s: %w", listPath, err)
	}
	if secret == nil || secret.Data == nil {
		return paths, nil
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return paths, nil
	}

	for _, k := range keys {
		key, _ := k.(string)
		full := prefix + key
		if len(key) > 0 && key[len(key)-1] == '/' {
			sub, err := c.ListPaths(mountPath, full)
			if err != nil {
				return nil, err
			}
			paths = append(paths, sub...)
		} else {
			paths = append(paths, full)
		}
	}
	return paths, nil
}
