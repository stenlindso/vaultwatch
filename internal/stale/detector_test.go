package stale

import (
	"strings"
	"testing"
	"time"
)

func baseTime() time.Time {
	return time.Now().Add(-45 * 24 * time.Hour)
}

func TestDetect_NoStalePaths(t *testing.T) {
	current := Snapshot{
		Environment: "prod",
		Paths:       []string{"secret/app/db", "secret/app/api"},
		CapturedAt:  time.Now(),
	}
	reference := Snapshot{
		Environment: "prod",
		Paths:       []string{"secret/app/db", "secret/app/api"},
		CapturedAt:  baseTime(),
	}

	result := Detect(current, reference, DefaultOptions())
	if result.HasStale() {
		t.Errorf("expected no stale paths, got %d", len(result.Entries))
	}
}

func TestDetect_StalePath(t *testing.T) {
	current := Snapshot{
		Environment: "prod",
		Paths:       []string{"secret/app/db"},
		CapturedAt:  time.Now(),
	}
	reference := Snapshot{
		Environment: "prod",
		Paths:       []string{"secret/app/db", "secret/app/old-key"},
		CapturedAt:  baseTime(),
	}

	result := Detect(current, reference, DefaultOptions())
	if !result.HasStale() {
		t.Fatal("expected stale paths")
	}
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 stale entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Path != "secret/app/old-key" {
		t.Errorf("unexpected stale path: %s", result.Entries[0].Path)
	}
}

func TestDetect_BelowThreshold(t *testing.T) {
	current := Snapshot{
		Environment: "staging",
		Paths:       []string{"secret/app/db"},
		CapturedAt:  time.Now(),
	}
	// Reference is only 5 days old — below the 30-day threshold.
	reference := Snapshot{
		Environment: "staging",
		Paths:       []string{"secret/app/db", "secret/app/temp"},
		CapturedAt:  time.Now().Add(-5 * 24 * time.Hour),
	}

	result := Detect(current, reference, DefaultOptions())
	if result.HasStale() {
		t.Errorf("expected no stale paths below threshold, got %d", len(result.Entries))
	}
}

func TestDetect_SortedEntries(t *testing.T) {
	current := Snapshot{Environment: "prod", Paths: []string{}, CapturedAt: time.Now()}
	reference := Snapshot{
		Environment: "prod",
		Paths:       []string{"secret/z", "secret/a", "secret/m"},
		CapturedAt:  baseTime(),
	}

	result := Detect(current, reference, DefaultOptions())
	if len(result.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result.Entries))
	}
	for i := 1; i < len(result.Entries); i++ {
		if result.Entries[i].Path < result.Entries[i-1].Path {
			t.Errorf("entries not sorted at index %d", i)
		}
	}
}

func TestSummary_NoStale(t *testing.T) {
	r := Result{Environment: "dev", Entries: nil}
	s := Summary(r)
	if !strings.Contains(s, "no stale") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithStale(t *testing.T) {
	r := Result{
		Environment: "prod",
		Entries: []Entry{
			{Path: "secret/old", Environment: "prod", AgeDays: 45},
		},
	}
	s := Summary(r)
	if !strings.Contains(s, "1 stale") {
		t.Errorf("unexpected summary: %s", s)
	}
}
