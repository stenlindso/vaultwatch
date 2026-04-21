// Package drift detects configuration drift between a baseline and current
// vault secret paths, producing a structured drift report per environment.
package drift

import (
	"fmt"
	"time"
)

// Result holds the drift detection outcome for a single environment.
type Result struct {
	Environment string
	Label       string
	Added       []string
	Removed     []string
	DetectedAt  time.Time
}

// HasDrift returns true if there are any added or removed paths.
func (r Result) HasDrift() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0
}

// Summary returns a human-readable one-line summary of the drift result.
func (r Result) Summary() string {
	if !r.HasDrift() {
		return fmt.Sprintf("[%s/%s] no drift detected", r.Environment, r.Label)
	}
	return fmt.Sprintf("[%s/%s] drift detected: +%d added, -%d removed",
		r.Environment, r.Label, len(r.Added), len(r.Removed))
}

// Detector compares current paths against a labeled baseline snapshot.
type Detector struct {
	environment string
	label       string
}

// NewDetector creates a Detector for the given environment and baseline label.
func NewDetector(environment, label string) *Detector {
	return &Detector{environment: environment, label: label}
}

// Detect computes the drift between baseline paths and current paths.
// baseline and current are slices of secret path strings.
func (d *Detector) Detect(baseline, current []string) Result {
	baseSet := toSet(baseline)
	curSet := toSet(current)

	var added, removed []string

	for p := range curSet {
		if !baseSet[p] {
			added = append(added, p)
		}
	}
	for p := range baseSet {
		if !curSet[p] {
			removed = append(removed, p)
		}
	}

	return Result{
		Environment: d.environment,
		Label:       d.label,
		Added:       added,
		Removed:     removed,
		DetectedAt:  time.Now().UTC(),
	}
}

func toSet(paths []string) map[string]bool {
	s := make(map[string]bool, len(paths))
	for _, p := range paths {
		s[p] = true
	}
	return s
}
