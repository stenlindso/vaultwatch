package annotate_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vaultwatch/internal/annotate"
)

func newTempStore(t *testing.T) *annotate.Store {
	t.Helper()
	dir := t.TempDir()
	return annotate.NewStore(dir)
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	s := newTempStore(t)
	ann := annotate.Annotation{
		Path:  "secret/db/password",
		Note:  "primary db cred",
		Owner: "platform-team",
	}
	if err := s.Save("prod", ann); err != nil {
		t.Fatalf("Save: %v", err)
	}
	results, err := s.Load("prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(results))
	}
	if results[0].Path != ann.Path {
		t.Errorf("path mismatch: got %q", results[0].Path)
	}
	if results[0].Owner != ann.Owner {
		t.Errorf("owner mismatch: got %q", results[0].Owner)
	}
}

func TestSave_SetsCreatedAt(t *testing.T) {
	s := newTempStore(t)
	before := time.Now().UTC()
	if err := s.Save("staging", annotate.Annotation{Path: "secret/key"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	results, _ := s.Load("staging")
	if results[0].CreatedAt.Before(before) {
		t.Error("expected CreatedAt to be set to current time")
	}
}

func TestSave_OverwritesExistingPath(t *testing.T) {
	s := newTempStore(t)
	s.Save("dev", annotate.Annotation{Path: "secret/api", Note: "old"})
	s.Save("dev", annotate.Annotation{Path: "secret/api", Note: "new"})
	results, _ := s.Load("dev")
	if len(results) != 1 {
		t.Fatalf("expected 1 annotation after overwrite, got %d", len(results))
	}
	if results[0].Note != "new" {
		t.Errorf("expected note 'new', got %q", results[0].Note)
	}
}

func TestGet_Found(t *testing.T) {
	s := newTempStore(t)
	s.Save("prod", annotate.Annotation{Path: "secret/token", Owner: "sec-team"})
	ann, err := s.Get("prod", "secret/token")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if ann.Owner != "sec-team" {
		t.Errorf("expected owner 'sec-team', got %q", ann.Owner)
	}
}

func TestGet_NotFound(t *testing.T) {
	s := newTempStore(t)
	s.Save("prod", annotate.Annotation{Path: "secret/other"})
	_, err := s.Get("prod", "secret/missing")
	if err == nil {
		t.Error("expected error for missing path")
	}
}

func TestLoad_Missing(t *testing.T) {
	s := newTempStore(t)
	_, err := s.Load("nonexistent")
	if !os.IsNotExist(err) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestSave_EmptyEnv_ReturnsError(t *testing.T) {
	s := newTempStore(t)
	err := s.Save("", annotate.Annotation{Path: "secret/x"})
	if err == nil {
		t.Error("expected error for empty env")
	}
}

func TestSave_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "store")
	s := annotate.NewStore(dir)
	if err := s.Save("prod", annotate.Annotation{Path: "secret/x"}); err != nil {
		t.Fatalf("Save with nested dir: %v", err)
	}
}
