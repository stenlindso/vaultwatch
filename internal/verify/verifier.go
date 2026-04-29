// Package verify provides path existence verification across environments,
// checking whether expected secret paths are present or missing in a snapshot.
package verify

import (
	"fmt"
	"sort"
	"time"
)

// Result holds the outcome of a verification run for a single environment.
type Result struct {
	Environment string
	Timestamp   time.Time
	Present     []string
	Missing     []string
}

// HasMissing returns true if any expected paths were not found.
func (r Result) HasMissing() bool {
	return len(r.Missing) > 0
}

// Summary returns a short human-readable summary line.
func (r Result) Summary() string {
	return fmt.Sprintf("env=%s present=%d missing=%d", r.Environment, len(r.Present), len(r.Missing))
}

// Verifier checks a set of expected paths against a snapshot of actual paths.
type Verifier struct {
	expected []string
}

// New creates a Verifier with the given expected paths.
func New(expected []string) *Verifier {
	norm := make([]string, len(expected))
	copy(norm, expected)
	sort.Strings(norm)
	return &Verifier{expected: norm}
}

// Check compares the expected paths against the provided actual paths for the
// named environment and returns a Result.
func (v *Verifier) Check(env string, actual []string) Result {
	actualSet := make(map[string]struct{}, len(actual))
	for _, p := range actual {
		actualSet[p] = struct{}{}
	}

	var present, missing []string
	for _, p := range v.expected {
		if _, ok := actualSet[p]; ok {
			present = append(present, p)
		} else {
			missing = append(missing, p)
		}
	}

	return Result{
		Environment: env,
		Timestamp:   time.Now().UTC(),
		Present:     present,
		Missing:     missing,
	}
}

// CheckAll runs Check for every entry in the envPaths map and returns all results.
func (v *Verifier) CheckAll(envPaths map[string][]string) []Result {
	envs := make([]string, 0, len(envPaths))
	for e := range envPaths {
		envs = append(envs, e)
	}
	sort.Strings(envs)

	results := make([]Result, 0, len(envs))
	for _, e := range envs {
		results = append(results, v.Check(e, envPaths[e]))
	}
	return results
}
