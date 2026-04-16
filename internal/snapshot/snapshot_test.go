package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultwatch/internal/snapshot"
)

func TestSaveAndLoad(t *testing.T) {
	paths := []string{"secret/app/db", "secret/app/api"}
	s := snapshot.New("staging", paths)

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(s, tmp); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := snapshot.Load(tmp)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Environment != "staging" {
		t.Errorf("expected environment staging, got %s", loaded.Environment)
	}
	if len(loaded.Paths) != len(paths) {
		t.Errorf("expected %d paths, got %d", len(paths), len(loaded.Paths))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(tmp, []byte("not json{"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := snapshot.Load(tmp)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
