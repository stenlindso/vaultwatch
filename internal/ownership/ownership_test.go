package ownership_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultwatch/internal/ownership"
)

func newTempStore(t *testing.T) *ownership.Store {
	t.Helper()
	dir := t.TempDir()
	return ownership.NewStore(filepath.Join(dir, "ownership.json"))
}

func TestSet_And_Get(t *testing.T) {
	s := newTempStore(t)
	r := ownership.Record{Path: "secret/app/db", Owner: "alice", Team: "platform"}
	if err := s.Set(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := s.Get("secret/app/db")
	if !ok {
		t.Fatal("expected record to be present")
	}
	if got.Owner != "alice" || got.Team != "platform" {
		t.Errorf("unexpected record: %+v", got)
	}
}

func TestSet_EmptyPath_ReturnsError(t *testing.T) {
	s := newTempStore(t)
	err := s.Set(ownership.Record{Path: "", Owner: "alice"})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSet_EmptyOwner_ReturnsError(t *testing.T) {
	s := newTempStore(t)
	err := s.Set(ownership.Record{Path: "secret/x", Owner: ""})
	if err == nil {
		t.Fatal("expected error for empty owner")
	}
}

func TestUnowned_ReturnsUnregisteredPaths(t *testing.T) {
	s := newTempStore(t)
	_ = s.Set(ownership.Record{Path: "secret/known", Owner: "bob"})
	unowned := s.Unowned([]string{"secret/known", "secret/unknown"})
	if len(unowned) != 1 || unowned[0] != "secret/unknown" {
		t.Errorf("unexpected unowned: %v", unowned)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "ownership.json")

	s1 := ownership.NewStore(file)
	_ = s1.Set(ownership.Record{Path: "secret/svc/token", Owner: "carol", Contact: "carol@example.com"})
	if err := s1.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	s2 := ownership.NewStore(file)
	if err := s2.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}
	got, ok := s2.Get("secret/svc/token")
	if !ok {
		t.Fatal("record not found after reload")
	}
	if got.Contact != "carol@example.com" {
		t.Errorf("unexpected contact: %s", got.Contact)
	}
}

func TestLoad_MissingFile_NoError(t *testing.T) {
	s := ownership.NewStore(filepath.Join(t.TempDir(), "missing.json"))
	if err := s.Load(); err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(file, []byte("not-json"), 0o644)
	s := ownership.NewStore(file)
	if err := s.Load(); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestAll_SortedByPath(t *testing.T) {
	s := newTempStore(t)
	_ = s.Set(ownership.Record{Path: "z/path", Owner: "x"})
	_ = s.Set(ownership.Record{Path: "a/path", Owner: "y"})
	all := s.All()
	if len(all) != 2 || all[0].Path != "a/path" {
		t.Errorf("expected sorted order, got: %v", all)
	}
}
