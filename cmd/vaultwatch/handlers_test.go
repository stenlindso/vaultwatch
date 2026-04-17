package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultwatch/internal/snapshot"
)

func TestDiffCommand_MissingSnapshots(t *testing.T) {
	dir := t.TempDir()
	flagSnapshotDir = dir
	flagOutputJSON = false

	cmd := diffCmd
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})

	err := runDiff(cmd, []string{"staging", "production"})
	if err == nil {
		t.Fatal("expected error for missing snapshots, got nil")
	}
}

func TestDiffCommand_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	flagSnapshotDir = dir
	flagOutputJSON = true

	mgr := snapshot.NewManager(dir)
	if err := mgr.Save("env1", []string{"secret/a", "secret/b"}); err != nil {
		t.Fatal(err)
	}
	if err := mgr.Save("env2", []string{"secret/b", "secret/c"}); err != nil {
		t.Fatal(err)
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runDiff(diffCmd, []string{"env1", "env2"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()

	if len(out) == 0 {
		t.Error("expected JSON output, got empty string")
	}
	if out[0] != '{' {
		t.Errorf("expected JSON object, got: %s", out[:min(20, len(out))])
	}
}

func TestSnapshotDir_Default(t *testing.T) {
	if flagSnapshotDir == "" {
		t.Error("snapshot dir should have a default value")
	}
	_ = filepath.Join(flagSnapshotDir, "test")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
