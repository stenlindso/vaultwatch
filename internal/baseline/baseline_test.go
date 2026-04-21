package baseline_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/youorg/vaultwatch/internal/baseline"
	"github.com/youorg/vaultwatch/internal/snapshot"
)

func newManager(t *testing.T) *baseline.Manager {
	t.Helper()
	dir := t.TempDir()
	sm, err := snapshot.NewManager(dir)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return baseline.NewManager(sm)
}

func TestSaveAndLoad(t *testing.T) {
	m := newManager(t)
	paths := []string{"secret/a", "secret/b", "secret/c"}

	_, err := m.Save("prod", "v1", paths)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := m.Load("prod", "v1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	sort.Strings(got)
	sort.Strings(paths)
	if len(got) != len(paths) {
		t.Fatalf("expected %d paths, got %d", len(paths), len(got))
	}
	for i := range paths {
		if got[i] != paths[i] {
			t.Errorf("path[%d]: want %q, got %q", i, paths[i], got[i])
		}
	}
}

func TestSave_EmptyEnv(t *testing.T) {
	m := newManager(t)
	_, err := m.Save("", "v1", []string{"secret/a"})
	if err == nil {
		t.Fatal("expected error for empty environment")
	}
}

func TestSave_EmptyLabel(t *testing.T) {
	m := newManager(t)
	_, err := m.Save("prod", "", []string{"secret/a"})
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestDiff_AddedAndRemoved(t *testing.T) {
	m := newManager(t)
	base := []string{"secret/a", "secret/b", "secret/c"}
	_, err := m.Save("staging", "v1", base)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	current := []string{"secret/b", "secret/c", "secret/d"}
	added, removed, err := m.Diff("staging", "v1", current)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}

	if len(added) != 1 || added[0] != "secret/d" {
		t.Errorf("added: want [secret/d], got %v", added)
	}
	if len(removed) != 1 || removed[0] != "secret/a" {
		t.Errorf("removed: want [secret/a], got %v", removed)
	}
}

func TestDiff_MissingBaseline(t *testing.T) {
	m := newManager(t)
	_, _, err := m.Diff("prod", "nonexistent", []string{"secret/x"})
	if err == nil {
		t.Fatal("expected error for missing baseline")
	}
}

func TestDiff_NoChanges(t *testing.T) {
	m := newManager(t)
	paths := []string{"secret/a", "secret/b"}
	_, _ = m.Save("dev", "v1", paths)

	added, removed, err := m.Diff("dev", "v1", paths)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no changes, got added=%v removed=%v", added, removed)
	}
}

// Ensure unused import does not break compilation.
var _ = filepath.Join
var _ = os.TempDir
