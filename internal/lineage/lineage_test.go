package lineage_test

import (
	"os"
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/lineage"
)

func newTempStore(t *testing.T) *lineage.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "lineage-test-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := lineage.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestRecord_And_Load(t *testing.T) {
	s := newTempStore(t)
	paths := []string{"secret/app/db", "secret/app/api"}
	if err := s.Record("staging", paths, "snapshot-v1"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, err := s.Load("staging")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Source != "snapshot-v1" {
		t.Errorf("unexpected source: %s", entries[0].Source)
	}
}

func TestLoad_Missing(t *testing.T) {
	s := newTempStore(t)
	entries, err := s.Load("nonexistent")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for missing env")
	}
}

func TestRecord_Appends(t *testing.T) {
	s := newTempStore(t)
	_ = s.Record("prod", []string{"secret/a"}, "snap-1")
	time.Sleep(2 * time.Millisecond)
	_ = s.Record("prod", []string{"secret/b"}, "snap-2")

	entries, err := s.Load("prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries after two records, got %d", len(entries))
	}
}
