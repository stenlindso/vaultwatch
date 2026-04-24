// Package resolve provides path resolution utilities for mapping
// short or aliased Vault paths to their canonical full paths.
package resolve

import (
	"fmt"
	"strings"
)

// Alias maps a short name to a canonical Vault path prefix.
type Alias struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

// Resolver resolves aliased or relative Vault paths to canonical paths.
type Resolver struct {
	aliases map[string]string
}

// New creates a Resolver from a slice of Alias definitions.
func New(aliases []Alias) *Resolver {
	m := make(map[string]string, len(aliases))
	for _, a := range aliases {
		key := strings.TrimSpace(a.Name)
		if key != "" {
			m[key] = strings.TrimRight(a.Prefix, "/")
		}
	}
	return &Resolver{aliases: m}
}

// Resolve expands a path that may begin with an alias into its full form.
// If no alias matches, the path is returned unchanged.
func (r *Resolver) Resolve(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("resolve: empty path")
	}
	for alias, prefix := range r.aliases {
		target := alias + "/"
		if strings.HasPrefix(path, target) {
			remainder := strings.TrimPrefix(path, target)
			return prefix + "/" + remainder, nil
		}
		if path == alias {
			return prefix, nil
		}
	}
	return path, nil
}

// ResolveAll resolves a slice of paths, returning the first error encountered.
func (r *Resolver) ResolveAll(paths []string) ([]string, error) {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		resolved, err := r.Resolve(p)
		if err != nil {
			return nil, err
		}
		out = append(out, resolved)
	}
	return out, nil
}

// Aliases returns a copy of the registered alias map.
func (r *Resolver) Aliases() map[string]string {
	copy := make(map[string]string, len(r.aliases))
	for k, v := range r.aliases {
		copy[k] = v
	}
	return copy
}
