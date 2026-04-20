package export

import (
	"testing"
	"time"
)

func TestBuildRecordsFromPaths_Basic(t *testing.T) {
	paths := []string{"secret/a", "secret/b", "secret/c"}
	now := time.Now().UTC()
	records := BuildRecordsFromPaths("dev", paths, now)

	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}
	for _, r := range records {
		if r.Environment != "dev" {
			t.Errorf("expected env dev, got %s", r.Environment)
		}
		if r.CapturedAt != now {
			t.Errorf("unexpected captured_at: %v", r.CapturedAt)
		}
	}
	if records[0].Path != "secret/a" {
		t.Errorf("unexpected path order: %s", records[0].Path)
	}
}

func TestBuildRecordsFromPaths_Empty(t *testing.T) {
	records := BuildRecordsFromPaths("prod", []string{}, time.Now())
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}

func TestBuildRecordsFromPaths_PreservesOrder(t *testing.T) {
	paths := []string{"z/path", "a/path", "m/path"}
	records := BuildRecordsFromPaths("staging", paths, time.Now())
	for i, r := range records {
		if r.Path != paths[i] {
			t.Errorf("index %d: expected %s, got %s", i, paths[i], r.Path)
		}
	}
}
