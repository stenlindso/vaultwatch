// Package normalize provides utilities for cleaning and standardizing
// Vault secret paths before comparison or storage.
package normalize

import (
	"path"
	"strings"
)

// Options controls normalization behavior.
type Options struct {
	// StripLeadingSlash removes a leading slash from each path.
	StripLeadingSlash bool
	// StripTrailingSlash removes a trailing slash from each path.
	StripTrailingSlash bool
	// Lowercase converts all path segments to lowercase.
	Lowercase bool
	// Clean applies path.Clean to collapse double slashes and dot segments.
	Clean bool
}

// DefaultOptions returns a sensible default normalization configuration.
func DefaultOptions() Options {
	return Options{
		StripLeadingSlash:  true,
		StripTrailingSlash: true,
		Lowercase:          false,
		Clean:              true,
	}
}

// Normalizer applies a fixed set of options to vault paths.
type Normalizer struct {
	opts Options
}

// New creates a Normalizer with the provided options.
func New(opts Options) *Normalizer {
	return &Normalizer{opts: opts}
}

// One normalizes a single path string.
func (n *Normalizer) One(p string) string {
	if n.opts.Clean {
		p = path.Clean(p)
	}
	if n.opts.Lowercase {
		p = strings.ToLower(p)
	}
	if n.opts.StripLeadingSlash {
		p = strings.TrimPrefix(p, "/")
	}
	if n.opts.StripTrailingSlash {
		p = strings.TrimSuffix(p, "/")
	}
	return p
}

// All normalizes a slice of paths, returning a new slice.
func (n *Normalizer) All(paths []string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		out[i] = n.One(p)
	}
	return out
}

// Deduplicate returns a slice with duplicate paths removed, preserving order.
func (n *Normalizer) Deduplicate(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		norm := n.One(p)
		if _, ok := seen[norm]; !ok {
			seen[norm] = struct{}{}
			out = append(out, norm)
		}
	}
	return out
}
