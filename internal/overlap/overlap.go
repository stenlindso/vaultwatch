// Package overlap identifies secret paths that appear in multiple environments,
// helping operators understand shared or duplicated configuration across envs.
package overlap

import "sort"

// Result holds the outcome of an overlap analysis across environments.
type Result struct {
	// SharedPaths maps each path to the sorted list of environments that contain it.
	SharedPaths map[string][]string
	// UniqueByEnv maps each environment to paths that appear only in that environment.
	UniqueByEnv map[string][]string
}

// HasOverlap returns true if any path appears in more than one environment.
func (r *Result) HasOverlap() bool {
	for _, envs := range r.SharedPaths {
		if len(envs) > 1 {
			return true
		}
	}
	return false
}

// Analyze computes path overlap across the provided environment snapshots.
// snapshots maps environment name -> list of secret paths.
func Analyze(snapshots map[string][]string) *Result {
	pathEnvs := make(map[string]map[string]struct{})

	for env, paths := range snapshots {
		for _, p := range paths {
			if pathEnvs[p] == nil {
				pathEnvs[p] = make(map[string]struct{})
			}
			pathEnvs[p][env] = struct{}{}
		}
	}

	shared := make(map[string][]string)
	for path, envSet := range pathEnvs {
		envList := make([]string, 0, len(envSet))
		for e := range envSet {
			envList = append(envList, e)
		}
		sort.Strings(envList)
		shared[path] = envList
	}

	uniqueByEnv := make(map[string][]string)
	for env := range snapshots {
		uniqueByEnv[env] = []string{}
	}
	for path, envList := range shared {
		if len(envList) == 1 {
			env := envList[0]
			uniqueByEnv[env] = append(uniqueByEnv[env], path)
		}
	}
	for env := range uniqueByEnv {
		sort.Strings(uniqueByEnv[env])
	}

	return &Result{
		SharedPaths: shared,
		UniqueByEnv: uniqueByEnv,
	}
}

// SharedAcrossAll returns paths that appear in every provided environment.
func SharedAcrossAll(snapshots map[string][]string) []string {
	if len(snapshots) == 0 {
		return nil
	}
	r := Analyze(snapshots)
	total := len(snapshots)
	var result []string
	for path, envs := range r.SharedPaths {
		if len(envs) == total {
			result = append(result, path)
		}
	}
	sort.Strings(result)
	return result
}
