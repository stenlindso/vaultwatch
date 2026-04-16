package vault

import (
	"context"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Lister recursively lists secret paths under a given prefix.
type Lister struct {
	client *vaultapi.Client
}

// NewLister creates a Lister from an existing vault client.
func NewLister(c *vaultapi.Client) *Lister {
	return &Lister{client: c}
}

// ListPaths returns all leaf secret paths under the given root path.
func (l *Lister) ListPaths(ctx context.Context, root string) ([]string, error) {
	root = strings.TrimSuffix(root, "/")
	var results []string
	if err := l.walk(ctx, root+"/", &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (l *Lister) walk(ctx context.Context, path string, out *[]string) error {
	secret, err := l.client.Logical().ListWithContext(ctx, toMetadataPath(path))
	if err != nil {
		return fmt.Errorf("listing %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil
	}
	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil
	}
	for _, k := range keys {
		key, _ := k.(string)
		full := path + key
		if strings.HasSuffix(key, "/") {
			if err := l.walk(ctx, full, out); err != nil {
				return err
			}
		} else {
			*out = append(*out, full)
		}
	}
	return nil
}

// toMetadataPath converts a KV v2 data path to its metadata equivalent.
func toMetadataPath(path string) string {
	return strings.Replace(path, "/data/", "/metadata/", 1)
}
