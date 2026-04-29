package ancestry

import (
	"os"
	"testing"
)

func newTempStore(t *testing.T) *Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "ancestry-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestRecord_And_Load(t *testing.T) {
	s := newTempStore(t)
	if err := s.Record("prod", "secret/base", "secret/derived"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	g, err := s.Load("prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(g.Edges) != 1 {
		t.Fatalf("want 1 edge, got %d", len(g.Edges))
	}
	if g.Edges[0].Parent != "secret/base" || g.Edges[0].Child != "secret/derived" {
		t.Errorf("unexpected edge: %+v", g.Edges[0])
	}
}

func TestLoad_Missing(t *testing.T) {
	s := newTempStore(t)
	g, err := s.Load("nonexistent")
	if err != nil {
		t.Fatalf("expected no error for missing env, got %v", err)
	}
	if len(g.Edges) != 0 {
		t.Errorf("expected empty graph, got %d edges", len(g.Edges))
	}
}

func TestChildren_ReturnsDirectChildren(t *testing.T) {
	s := newTempStore(t)
	_ = s.Record("staging", "secret/root", "secret/child-a")
	_ = s.Record("staging", "secret/root", "secret/child-b")
	_ = s.Record("staging", "secret/other", "secret/child-c")

	children, err := s.Children("staging", "secret/root")
	if err != nil {
		t.Fatalf("Children: %v", err)
	}
	if len(children) != 2 {
		t.Fatalf("want 2 children, got %d", len(children))
	}
}

func TestOrphans_DetectsOrphanedPaths(t *testing.T) {
	s := newTempStore(t)
	_ = s.Record("dev", "secret/parent", "secret/orphan")
	_ = s.Record("dev", "secret/parent", "secret/alive")

	// secret/alive is still active; secret/orphan is not
	orphans, err := s.Orphans("dev", []string{"secret/alive", "secret/parent"})
	if err != nil {
		t.Fatalf("Orphans: %v", err)
	}
	if len(orphans) != 1 || orphans[0] != "secret/orphan" {
		t.Errorf("unexpected orphans: %v", orphans)
	}
}

func TestRecord_EmptyFields_ReturnsError(t *testing.T) {
	s := newTempStore(t)
	if err := s.Record("", "p", "c"); err == nil {
		t.Error("expected error for empty env")
	}
	if err := s.Record("env", "", "c"); err == nil {
		t.Error("expected error for empty parent")
	}
	if err := s.Record("env", "p", ""); err == nil {
		t.Error("expected error for empty child")
	}
}

func TestRecord_Appends(t *testing.T) {
	s := newTempStore(t)
	_ = s.Record("prod", "a", "b")
	_ = s.Record("prod", "b", "c")
	g, _ := s.Load("prod")
	if len(g.Edges) != 2 {
		t.Errorf("want 2 edges, got %d", len(g.Edges))
	}
}
