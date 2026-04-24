// Package score computes a health score for a Vault environment based on
// drift, policy violations, and staleness of snapshots.
package score

import "time"

// Grade represents a letter-grade summary of an environment's health.
type Grade string

const (
	GradeA Grade = "A"
	GradeB Grade = "B"
	GradeC Grade = "C"
	GradeD Grade = "D"
	GradeF Grade = "F"
)

// Input holds the raw metrics used to compute a score.
type Input struct {
	Environment      string
	TotalPaths       int
	DriftedPaths     int
	PolicyViolations int
	SnapshotAge      time.Duration
	MaxSnapshotAge   time.Duration // treated as 24h if zero
}

// Result is the computed score for a single environment.
type Result struct {
	Environment string
	Score       int // 0–100
	Grade       Grade
	Reasons     []string
}

// Compute derives a 0–100 score and letter grade from the supplied Input.
func Compute(in Input) Result {
	if in.MaxSnapshotAge == 0 {
		in.MaxSnapshotAge = 24 * time.Hour
	}

	score := 100
	var reasons []string

	// Penalise drift: up to -40 points proportionally.
	if in.TotalPaths > 0 && in.DriftedPaths > 0 {
		ratio := float64(in.DriftedPaths) / float64(in.TotalPaths)
		penalty := int(ratio * 40)
		if penalty > 40 {
			penalty = 40
		}
		score -= penalty
		reasons = append(reasons, "drifted paths detected")
	}

	// Penalise policy violations: -10 per violation, capped at -40.
	if in.PolicyViolations > 0 {
		penalty := in.PolicyViolations * 10
		if penalty > 40 {
			penalty = 40
		}
		score -= penalty
		reasons = append(reasons, "policy violations present")
	}

	// Penalise stale snapshot: up to -20 points.
	if in.SnapshotAge > in.MaxSnapshotAge {
		ratio := float64(in.SnapshotAge) / float64(in.MaxSnapshotAge)
		if ratio > 2 {
			ratio = 2
		}
		penalty := int((ratio - 1) * 20)
		if penalty > 20 {
			penalty = 20
		}
		score -= penalty
		reasons = append(reasons, "snapshot is stale")
	}

	if score < 0 {
		score = 0
	}

	return Result{
		Environment: in.Environment,
		Score:       score,
		Grade:       toGrade(score),
		Reasons:     reasons,
	}
}

func toGrade(score int) Grade {
	switch {
	case score >= 90:
		return GradeA
	case score >= 75:
		return GradeB
	case score >= 60:
		return GradeC
	case score >= 40:
		return GradeD
	default:
		return GradeF
	}
}
