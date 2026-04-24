// Package index provides path indexing and reverse-lookup capabilities
// for Vault secret paths across environments.
package index

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Entry represents a single indexed path with its associated environment.
type Entry struct {
	Environment string
	Path        string
}

// Index maps path segments to the environments that contain them,
// enabling fast reverse-lookup and cross-environment queries.
type Index struct {
	mu      sync.RWMutex
	byEnv   map[string][]string
	byPath  map[string][]string
}

// New returns an empty Index.
func New() *Index {
	return &Index{
		byEnv:  make(map[string][]string),
		byPath: make(map[string][]string),
	}
}

// Add inserts a path for the given environment into the index.
func (idx *Index) Add(env, path string) {
	if env == "" || path == "" {
		return
	}
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.byEnv[env] = append(idx.byEnv[env], path)
	idx.byPath[path] = appendUnique(idx.byPath[path], env)
}

// PathsForEnv returns all indexed paths for the given environment, sorted.
func (idx *Index) PathsForEnv(env string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	paths := make([]string, len(idx.byEnv[env]))
	copy(paths, idx.byEnv[env])
	sort.Strings(paths)
	return paths
}

// EnvsForPath returns all environments that contain the given path, sorted.
func (idx *Index) EnvsForPath(path string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	envs := make([]string, len(idx.byPath[path]))
	copy(envs, idx.byPath[path])
	sort.Strings(envs)
	return envs
}

// Search returns all entries whose path contains the given substring.
func (idx *Index) Search(query string) []Entry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	var results []Entry
	for path, envs := range idx.byPath {
		if strings.Contains(path, query) {
			for _, env := range envs {
				results = append(results, Entry{Environment: env, Path: path})
			}
		}
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Environment != results[j].Environment {
			return results[i].Environment < results[j].Environment
		}
		return results[i].Path < results[j].Path
	})
	return results
}

// Environments returns all known environments in the index, sorted.
func (idx *Index) Environments() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	envs := make([]string, 0, len(idx.byEnv))
	for env := range idx.byEnv {
		envs = append(envs, env)
	}
	sort.Strings(envs)
	return envs
}

// Stats returns a summary string of index contents.
func (idx *Index) Stats() string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return fmt.Sprintf("environments=%d unique_paths=%d", len(idx.byEnv), len(idx.byPath))
}

func appendUnique(slice []string, val string) []string {
	for _, v := range slice {
		if v == val {
			return slice
		}
	}
	return append(slice, val)
}
