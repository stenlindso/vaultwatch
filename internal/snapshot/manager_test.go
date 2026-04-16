package snapshot_test

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/snapshot"
)

func TestManager_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	m, err := snapshot.NewManager(dir)
	if err != nil {
		t.Fatalf("NewManager() error: %v", err)
	}

	s := snapshot.New("prod", []string{"secret//db", "secret/prod/cache"})
	if err := m.SaveSnapshot(s); err != nil {
		t.Fatalf("SaveSnapshot() error: %v", err)
	}

	loaded"prod")
	if err != nil {
		t.Fatalf("LoadSnapshot() error: %v", err)
	}
	if loaded.Environment != "prod" {
		t.Errorf("expected prodEnvironment)
	}
}

func TestManager_ListEnvironments(t *testing.T) {
	dir := t.TempDir()
	m, err := snapshot.NewManager(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, env := range []string{"dev", "staging", "prod"} {
		if err := m.SaveSnapshot(snapshot.New(env, []string{})); err != nil {
			t.Fatalf("SaveSnapshot(%s) error: %v", env, err)
		}
	}

	envs, err := m.ListEnvironments()
	if err != nil {
		t.Fatalf("ListEnvironments() error: %v", err)
	}
	if len(envs) != 3 {
		t.Errorf("expected 3 environments, got %d", len(envs))
	}
}

func TestManager_LoadMissing(t *testing.T) {
	dir := t.TempDir()
	m, _ := snapshot.NewManager(dir)
	_, err := m.LoadSnapshot("ghost")
	if err == nil {
		t.Error("expected error loading non-existent snapshot")
	}
}
