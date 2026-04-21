package access_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/access"
)

func newTempStore(t *testing.T) *access.Store {
	t.Helper()
	dir := t.TempDir()
	store, err := access.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func sampleRecords(env string) []access.Record {
	now := time.Now().UTC()
	return []access.Record{
		{Environment: env, Path: "secret/alpha", Accessible: true, CheckedAt: now},
		{Environment: env, Path: "secret/beta", Accessible: false, CheckedAt: now},
		{Environment: env, Path: "secret/gamma", Accessible: true, CheckedAt: now},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	store := newTempStore(t)
	records := sampleRecords("staging")

	if err := store.Save("staging", records); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := store.Load("staging")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != len(records) {
		t.Errorf("expected %d records, got %d", len(records), len(loaded))
	}
}

func TestLoad_Missing(t *testing.T) {
	store := newTempStore(t)
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Error("expected error loading missing env, got nil")
	}
}

func TestSave_EmptyEnv(t *testing.T) {
	store := newTempStore(t)
	err := store.Save("", sampleRecords(""))
	if err == nil {
		t.Error("expected error for empty environment name")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	store, _ := access.NewStore(dir)
	badFile := filepath.Join(dir, "prod_access.json")
	_ = os.WriteFile(badFile, []byte("not-json{"), 0o644)
	_, err := store.Load("prod")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiff_GainedAndLost(t *testing.T) {
	now := time.Now().UTC()
	prev := []access.Record{
		{Path: "secret/alpha", Accessible: true, CheckedAt: now},
		{Path: "secret/beta", Accessible: true, CheckedAt: now},
	}
	curr := []access.Record{
		{Path: "secret/alpha", Accessible: true, CheckedAt: now},
		{Path: "secret/gamma", Accessible: true, CheckedAt: now},
	}
	report := access.Diff("prod", prev, curr)
	if report.Environment != "prod" {
		t.Errorf("expected env prod, got %s", report.Environment)
	}
	if len(report.Gained) != 1 || report.Gained[0] != "secret/gamma" {
		t.Errorf("unexpected gained: %v", report.Gained)
	}
	if len(report.Lost) != 1 || report.Lost[0] != "secret/beta" {
		t.Errorf("unexpected lost: %v", report.Lost)
	}
	if !report.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestDiff_NoChanges(t *testing.T) {
	records := sampleRecords("dev")
	report := access.Diff("dev", records, records)
	if report.HasChanges() {
		t.Errorf("expected no changes, got gained=%v lost=%v", report.Gained, report.Lost)
	}
}

func TestDiff_InaccessibleIgnored(t *testing.T) {
	now := time.Now().UTC()
	prev := []access.Record{
		{Path: "secret/hidden", Accessible: false, CheckedAt: now},
	}
	curr := []access.Record{}
	report := access.Diff("dev", prev, curr)
	if report.HasChanges() {
		t.Errorf("inaccessible paths should not appear in diff: %+v", report)
	}
}
