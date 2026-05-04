// Package shadow provides path shadowing detection across Vault environments.
// A path is considered "shadowed" when it exists in a lower-priority environment
// but is overridden or absent in a higher-priority environment.
package shadow

import "sort"

// Result holds the shadowing analysis for a set of environments.
type Result struct {
	// Shadowed maps a path to the list of environments where it is shadowed.
	Shadowed map[string][]string
	// Unique maps a path to the single environment that owns it exclusively.
	Unique map[string]string
	Environments []string
}

// HasShadows returns true if any shadowed paths were detected.
func (r Result) HasShadows() bool {
	return len(r.Shadowed) > 0
}

// Analyze detects shadowed paths given a priority-ordered list of environments
// and their path sets. Environments earlier in the priority slice have higher
// priority. A path is shadowed in env[i] if it also appears in env[j] where j < i.
func Analyze(priority []string, snapshots map[string][]string) Result {
	result := Result{
		Shadowed:     make(map[string][]string),
		Unique:       make(map[string]string),
		Environments: priority,
	}

	// Build a map from path -> set of environments containing it.
	pathEnvs := make(map[string][]string)
	for _, env := range priority {
		for _, p := range snapshots[env] {
			pathEnvs[p] = append(pathEnvs[p], env)
		}
	}

	// Determine priority rank for each environment.
	rank := make(map[string]int, len(priority))
	for i, env := range priority {
		rank[env] = i
	}

	for path, envs := range pathEnvs {
		if len(envs) == 1 {
			result.Unique[path] = envs[0]
			continue
		}
		// Sort by priority rank.
		sort.Slice(envs, func(i, j int) bool {
			return rank[envs[i]] < rank[envs[j]]
		})
		// All envs except the highest-priority one are shadowed.
		for _, env := range envs[1:] {
			result.Shadowed[path] = append(result.Shadowed[path], env)
		}
	}

	return result
}
