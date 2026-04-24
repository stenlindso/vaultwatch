// Package rollup aggregates secret path statistics across multiple environments
// into a consolidated view, grouping paths by common prefixes or segments.
package rollup

import (
	"sort"
	"strings"
)

// Entry represents a rolled-up prefix group with environment coverage.
type Entry struct {
	Prefix   string
	Envs     []string
	Count    int
	Coverage float64 // fraction of envs that contain this prefix
}

// Result holds all rolled-up entries for a set of environments.
type Result struct {
	TotalEnvs int
	Entries   []Entry
}

// Rollup computes prefix-level aggregation from a map of env -> paths.
// depth controls how many path segments to group by (1 = top-level only).
func Rollup(envPaths map[string][]string, depth int) Result {
	if depth < 1 {
		depth = 1
	}

	// prefix -> set of envs
	prefixEnvs := make(map[string]map[string]struct{})

	for env, paths := range envPaths {
		for _, p := range paths {
			prefix := extractPrefix(p, depth)
			if prefix == "" {
				continue
			}
			if prefixEnvs[prefix] == nil {
				prefixEnvs[prefix] = make(map[string]struct{})
			}
			prefixEnvs[prefix][env] = struct{}{}
		}
	}

	totalEnvs := len(envPaths)
	entries := make([]Entry, 0, len(prefixEnvs))

	for prefix, envSet := range prefixEnvs {
		envList := make([]string, 0, len(envSet))
		for e := range envSet {
			envList = append(envList, e)
		}
		sort.Strings(envList)

		coverage := 0.0
		if totalEnvs > 0 {
			coverage = float64(len(envList)) / float64(totalEnvs)
		}

		entries = append(entries, Entry{
			Prefix:   prefix,
			Envs:     envList,
			Count:    len(envList),
			Coverage: coverage,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Prefix < entries[j].Prefix
	})

	return Result{TotalEnvs: totalEnvs, Entries: entries}
}

// extractPrefix returns the first `depth` segments of a slash-separated path.
func extractPrefix(path string, depth int) string {
	path = strings.Trim(path, "/")
	parts := strings.SplitN(path, "/", depth+1)
	if len(parts) == 0 {
		return ""
	}
	if len(parts) > depth {
		parts = parts[:depth]
	}
	return strings.Join(parts, "/")
}
