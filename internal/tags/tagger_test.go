package tags_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultwatch/internal/tags"
)

func newTempStore(t *testing.T) *tags.TagStore {
	t.Helper()
	dir := t.TempDir()
	store, err := tags.NewTagStore(dir)
	if err != nil {
		t.Fatalf("NewTagStore: %v", err)
	}
	return store
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	store := newTempStore(t)
	entries := []tags.TagEntry{
		{Path: "secret/app/db", Tags: []string{"critical", "pii"}},
		{Path: "secret/app/api", Tags: []string{"critical"}},
	}
	if err := store.Save("prod", entries); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := store.Load("prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(loaded))
	}
	if loaded[0].Path != entries[0].Path {
		t.Errorf("path mismatch: got %q", loaded[0].Path)
	}
}

func TestLoad_Missing(t *testing.T) {
	store := newTempStore(t)
	_, err := store.Load("staging")
	if err == nil {
		t.Fatal("expected error for missing env, got nil")
	}
}

func TestSave_EmptyEnv(t *testing.T) {
	store := newTempStore(t)
	err := store.Save("", []tags.TagEntry{})
	if err == nil {
		t.Fatal("expected error for empty env")
	}
}

func TestFilterByTag_Match(t *testing.T) {
	entries := []tags.TagEntry{
		{Path: "secret/a", Tags: []string{"pii"}},
		{Path: "secret/b", Tags: []string{"internal"}},
		{Path: "secret/c", Tags: []string{"pii", "critical"}},
	}
	paths := tags.FilterByTag(entries, "pii")
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
}

func TestFilterByTag_NoMatch(t *testing.T) {
	entries := []tags.TagEntry{
		{Path: "secret/x", Tags: []string{"other"}},
	}
	paths := tags.FilterByTag(entries, "pii")
	if len(paths) != 0 {
		t.Errorf("expected 0 paths, got %d", len(paths))
	}
}

func TestBuildEntries(t *testing.T) {
	paths := []string{"secret/a", "secret/b"}
	entries := tags.BuildEntries(paths, []string{"auto", "scanned"})
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if len(e.Tags) != 2 {
			t.Errorf("expected 2 tags on %q, got %d", e.Path, len(e.Tags))
		}
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	store, _ := tags.NewTagStore(dir)
	_ = os.WriteFile(filepath.Join(dir, "dev.tags.json"), []byte("not-json{"), 0o644)
	_, err := store.Load("dev")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
