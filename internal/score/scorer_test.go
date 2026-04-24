package score

import (
	"testing"
	"time"
)

func TestCompute_PerfectScore(t *testing.T) {
	r := Compute(Input{
		Environment:      "prod",
		TotalPaths:       100,
		DriftedPaths:     0,
		PolicyViolations: 0,
		SnapshotAge:      1 * time.Hour,
		MaxSnapshotAge:   24 * time.Hour,
	})
	if r.Score != 100 {
		t.Errorf("expected 100, got %d", r.Score)
	}
	if r.Grade != GradeA {
		t.Errorf("expected A, got %s", r.Grade)
	}
	if len(r.Reasons) != 0 {
		t.Errorf("expected no reasons, got %v", r.Reasons)
	}
}

func TestCompute_DriftPenalty(t *testing.T) {
	r := Compute(Input{
		Environment:    "staging",
		TotalPaths:     100,
		DriftedPaths:   50, // 50% drift → -20
		MaxSnapshotAge: 24 * time.Hour,
	})
	if r.Score != 80 {
		t.Errorf("expected 80, got %d", r.Score)
	}
	if r.Grade != GradeB {
		t.Errorf("expected B, got %s", r.Grade)
	}
}

func TestCompute_PolicyViolationPenalty(t *testing.T) {
	r := Compute(Input{
		Environment:      "dev",
		TotalPaths:       50,
		PolicyViolations: 3, // -30
		MaxSnapshotAge:   24 * time.Hour,
	})
	if r.Score != 70 {
		t.Errorf("expected 70, got %d", r.Score)
	}
}

func TestCompute_StaleSnapshot(t *testing.T) {
	r := Compute(Input{
		Environment:    "prod",
		TotalPaths:     10,
		SnapshotAge:    48 * time.Hour, // 2× max → -20
		MaxSnapshotAge: 24 * time.Hour,
	})
	if r.Score != 80 {
		t.Errorf("expected 80, got %d", r.Score)
	}
	if r.Grade != GradeB {
		t.Errorf("expected B, got %s", r.Grade)
	}
}

func TestCompute_ScoreFloorIsZero(t *testing.T) {
	r := Compute(Input{
		Environment:      "chaos",
		TotalPaths:       10,
		DriftedPaths:     10,
		PolicyViolations: 10,
		SnapshotAge:      96 * time.Hour,
		MaxSnapshotAge:   24 * time.Hour,
	})
	if r.Score < 0 {
		t.Errorf("score should not be negative, got %d", r.Score)
	}
	if r.Grade != GradeF {
		t.Errorf("expected F, got %s", r.Grade)
	}
}

func TestCompute_DefaultMaxSnapshotAge(t *testing.T) {
	// MaxSnapshotAge zero → defaults to 24h; snapshot 12h old should not penalise.
	r := Compute(Input{
		Environment: "prod",
		SnapshotAge: 12 * time.Hour,
	})
	if r.Score != 100 {
		t.Errorf("expected 100, got %d", r.Score)
	}
}

func TestToGrade_Boundaries(t *testing.T) {
	cases := []struct {
		score int
		want  Grade
	}{
		{100, GradeA},
		{90, GradeA},
		{89, GradeB},
		{75, GradeB},
		{74, GradeC},
		{60, GradeC},
		{59, GradeD},
		{40, GradeD},
		{39, GradeF},
		{0, GradeF},
	}
	for _, c := range cases {
		if g := toGrade(c.score); g != c.want {
			t.Errorf("toGrade(%d) = %s, want %s", c.score, g, c.want)
		}
	}
}
