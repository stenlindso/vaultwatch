package prune

import (
	"testing"
	"time"
)

var baseTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeEntry(env, label string, age time.Duration) Entry {
	return Entry{
		Environment: env,
		Label:       label,
		CreatedAt:   baseTime.Add(-age),
	}
}

func TestEvaluate_Empty(t *testing.T) {
	r := Evaluate(nil, Options{MaxAge: time.Hour}, baseTime)
	if r.HasPruned() {
		t.Fatal("expected no pruned entries for empty input")
	}
}

func TestEvaluate_AgePruning(t *testing.T) {
	entries := []Entry{
		makeEntry("prod", "snap-1", 48*time.Hour),
		makeEntry("prod", "snap-2", 2*time.Hour),
		makeEntry("prod", "snap-3", 30*time.Minute),
	}
	r := Evaluate(entries, Options{MaxAge: 24 * time.Hour}, baseTime)
	if len(r.Pruned) != 1 {
		t.Fatalf("expected 1 pruned, got %d", len(r.Pruned))
	}
	if r.Pruned[0].Label != "snap-1" {
		t.Errorf("expected snap-1 pruned, got %s", r.Pruned[0].Label)
	}
	if len(r.Retained) != 2 {
		t.Fatalf("expected 2 retained, got %d", len(r.Retained))
	}
}

func TestEvaluate_CountPruning(t *testing.T) {
	entries := []Entry{
		makeEntry("staging", "snap-old", 10*time.Hour),
		makeEntry("staging", "snap-mid", 5*time.Hour),
		makeEntry("staging", "snap-new", 1*time.Hour),
	}
	r := Evaluate(entries, Options{MaxCount: 2}, baseTime)
	if len(r.Pruned) != 1 {
		t.Fatalf("expected 1 pruned, got %d", len(r.Pruned))
	}
	if r.Pruned[0].Label != "snap-old" {
		t.Errorf("expected snap-old pruned, got %s", r.Pruned[0].Label)
	}
}

func TestEvaluate_CombinedRules(t *testing.T) {
	entries := []Entry{
		makeEntry("prod", "a", 72*time.Hour),
		makeEntry("prod", "b", 36*time.Hour),
		makeEntry("prod", "c", 1*time.Hour),
		makeEntry("dev", "x", 1*time.Hour),
	}
	opts := Options{MaxAge: 48 * time.Hour, MaxCount: 1}
	r := Evaluate(entries, opts, baseTime)
	// "a" pruned by age, "b" pruned by count (prod keeps only "c")
	if len(r.Pruned) != 2 {
		t.Fatalf("expected 2 pruned, got %d: %v", len(r.Pruned), r.Pruned)
	}
	if len(r.Retained) != 2 {
		t.Fatalf("expected 2 retained, got %d", len(r.Retained))
	}
}

func TestResult_Summary(t *testing.T) {
	r := Result{
		Pruned:   []Entry{{Environment: "prod", Label: "old"}},
		Retained: []Entry{{Environment: "prod", Label: "new"}, {Environment: "dev", Label: "snap"}},
	}
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestEvaluate_NoOptions(t *testing.T) {
	entries := []Entry{
		makeEntry("prod", "snap-1", 100*time.Hour),
		makeEntry("prod", "snap-2", 200*time.Hour),
	}
	r := Evaluate(entries, Options{}, baseTime)
	if r.HasPruned() {
		t.Error("expected no pruning when no options set")
	}
	if len(r.Retained) != 2 {
		t.Errorf("expected 2 retained, got %d", len(r.Retained))
	}
}
