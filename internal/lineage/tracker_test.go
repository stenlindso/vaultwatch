package lineage_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultwatch/internal/lineage"
)

func TestBuildRecords_GroupsByPath(t *testing.T) {
	s := newTempStore(t)
	_ = s.Record("dev", []string{"secret/x", "secret/y"}, "snap-1")
	_ = s.Record("dev", []string{"secret/x"}, "snap-2")

	tracker := lineage.NewTracker(s)
	records, err := tracker.BuildRecords("dev")
	if err != nil {
		t.Fatalf("BuildRecords: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
	for _, r := range records {
		if r.Path == "secret/x" && len(r.History) != 2 {
			t.Errorf("secret/x: expected 2 history entries, got %d", len(r.History))
		}
	}
}

func TestFirstSeen_Found(t *testing.T) {
	s := newTempStore(t)
	_ = s.Record("qa", []string{"secret/cfg"}, "initial")

	tracker := lineage.NewTracker(s)
	entry, err := tracker.FirstSeen("qa", "secret/cfg")
	if err != nil {
		t.Fatalf("FirstSeen: %v", err)
	}
	if entry.Source != "initial" {
		t.Errorf("expected source 'initial', got %q", entry.Source)
	}
}

func TestFirstSeen_NotFound(t *testing.T) {
	s := newTempStore(t)
	tracker := lineage.NewTracker(s)
	_, err := tracker.FirstSeen("qa", "secret/missing")
	if err == nil {
		t.Fatal("expected error for missing path")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestFormatText_ContainsPath(t *testing.T) {
	s := newTempStore(t)
	_ = s.Record("prod", []string{"secret/db"}, "live")
	tracker := lineage.NewTracker(s)
	records, _ := tracker.BuildRecords("prod")

	var sb strings.Builder
	lineage.FormatText(&sb, "prod", records)
	out := sb.String()
	if !strings.Contains(out, "secret/db") {
		t.Errorf("output missing path: %s", out)
	}
	if !strings.Contains(out, "prod") {
		t.Errorf("output missing env name: %s", out)
	}
}
