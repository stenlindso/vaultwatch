// Package quota provides path count quota tracking and enforcement across environments.
package quota

import (
	"fmt"
	"sort"
)

// Policy defines quota limits for secret path counts.
type Policy struct {
	MaxPaths    int // hard limit; 0 means unlimited
	WarnAt      int // soft warning threshold; 0 means no warning
}

// Status represents the quota evaluation result for a single environment.
type Status struct {
	Environment string
	PathCount   int
	Limit       int
	WarnAt      int
	Exceeded    bool
	Warning     bool
}

// Report holds quota statuses for all evaluated environments.
type Report struct {
	Statuses []Status
}

// HasExceeded returns true if any environment exceeded its hard limit.
func (r *Report) HasExceeded() bool {
	for _, s := range r.Statuses {
		if s.Exceeded {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any environment crossed the warning threshold.
func (r *Report) HasWarnings() bool {
	for _, s := range r.Statuses {
		if s.Warning {
			return true
		}
	}
	return false
}

// Evaluate checks each environment's path count against the given policy.
// envPaths maps environment name to its list of secret paths.
func Evaluate(envPaths map[string][]string, policy Policy) *Report {
	envs := make([]string, 0, len(envPaths))
	for env := range envPaths {
		envs = append(envs, env)
	}
	sort.Strings(envs)

	statuses := make([]Status, 0, len(envs))
	for _, env := range envs {
		paths := envPaths[env]
		count := len(paths)
		s := Status{
			Environment: env,
			PathCount:   count,
			Limit:       policy.MaxPaths,
			WarnAt:      policy.WarnAt,
		}
		if policy.MaxPaths > 0 && count > policy.MaxPaths {
			s.Exceeded = true
		}
		if policy.WarnAt > 0 && count >= policy.WarnAt && !s.Exceeded {
			s.Warning = true
		}
		statuses = append(statuses, s)
	}
	return &Report{Statuses: statuses}
}

// Summary returns a human-readable one-line summary of the report.
func Summary(r *Report) string {
	exceeded := 0
	warnings := 0
	for _, s := range r.Statuses {
		if s.Exceeded {
			exceeded++
		} else if s.Warning {
			warnings++
		}
	}
	if exceeded == 0 && warnings == 0 {
		return fmt.Sprintf("quota OK: %d environment(s) within limits", len(r.Statuses))
	}
	return fmt.Sprintf("quota: %d exceeded, %d warning(s) across %d environment(s)", exceeded, warnings, len(r.Statuses))
}
