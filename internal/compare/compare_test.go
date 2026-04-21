package compare

import (
	"testing"
	"time"

	"github.com/vaultwatch/internal/snapshot"
)

func makeSnap(paths []string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		Paths:     paths,
		Timestamp: time.Now(),
	}
}

func TestBuildMatrix_TwoEnvs(t *testing.T) {
	envs := map[string]*snapshot.Snapshot{
		"staging": makeSnap([]string{"secret/a", "secret/b", "secret/c"}),
		"prod":    makeSnap([]string{"secret/b", "secret/c", "secret/d"}),
	}

	matrix := BuildMatrix(envs)

	if len(matrix.Environments) != 2 {
		t.Fatalf("expected 2 environments, got %d", len(matrix.Environments))
	}

	cell, ok := matrix.Cells[Key("staging", "prod")]
	if !ok {
		t.Fatal("expected cell for staging::prod")
	}

	if cell.Shared != 2 {
		t.Errorf("expected 2 shared paths, got %d", cell.Shared)
	}
	if len(cell.OnlyLeft) != 1 || cell.OnlyLeft[0] != "secret/a" {
		t.Errorf("unexpected OnlyLeft: %v", cell.OnlyLeft)
	}
	if len(cell.OnlyRight) != 1 || cell.OnlyRight[0] != "secret/d" {
		t.Errorf("unexpected OnlyRight: %v", cell.OnlyRight)
	}
}

func TestBuildMatrix_ThreeEnvs(t *testing.T) {
	envs := map[string]*snapshot.Snapshot{
		"dev":     makeSnap([]string{"secret/x"}),
		"staging": makeSnap([]string{"secret/x", "secret/y"}),
		"prod":    makeSnap([]string{"secret/x", "secret/y", "secret/z"}),
	}

	matrix := BuildMatrix(envs)

	// 3 envs => 3 pairs
	if len(matrix.Cells) != 3 {
		t.Errorf("expected 3 cells, got %d", len(matrix.Cells))
	}
}

func TestBuildMatrix_IdenticalEnvs(t *testing.T) {
	paths := []string{"secret/a", "secret/b"}
	envs := map[string]*snapshot.Snapshot{
		"alpha": makeSnap(paths),
		"beta":  makeSnap(paths),
	}

	matrix := BuildMatrix(envs)
	cell := matrix.Cells[Key("alpha", "beta")]

	if cell.Shared != 2 {
		t.Errorf("expected 2 shared, got %d", cell.Shared)
	}
	if len(cell.OnlyLeft) != 0 || len(cell.OnlyRight) != 0 {
		t.Errorf("expected no unique paths for identical envs")
	}
}

func TestKey_Canonical(t *testing.T) {
	if Key("a", "b") != Key("b", "a") {
		t.Error("Key should be symmetric")
	}
}
