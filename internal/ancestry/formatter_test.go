package ancestry

import (
	"strings"
	"testing"
	"time"
)

func makeGraph(edges ...Edge) Graph {
	return Graph{Edges: edges}
}

func TestFormatText_ContainsEnvName(t *testing.T) {
	g := makeGraph(Edge{Parent: "secret/a", Child: "secret/b", Env: "prod", RecordedAt: time.Now()})
	out := FormatText("prod", g)
	if !strings.Contains(out, "prod") {
		t.Errorf("expected env name in output, got:\n%s", out)
	}
}

func TestFormatText_ContainsEdge(t *testing.T) {
	g := makeGraph(Edge{Parent: "secret/parent", Child: "secret/child", Env: "dev", RecordedAt: time.Now()})
	out := FormatText("dev", g)
	if !strings.Contains(out, "secret/parent") || !strings.Contains(out, "secret/child") {
		t.Errorf("expected edge paths in output, got:\n%s", out)
	}
}

func TestFormatText_EmptyGraph(t *testing.T) {
	out := FormatText("staging", Graph{})
	if !strings.Contains(out, "no edges") {
		t.Errorf("expected empty message, got:\n%s", out)
	}
}

func TestFormatOrphans_ListsOrphans(t *testing.T) {
	out := FormatOrphans("prod", []string{"secret/stale", "secret/gone"})
	if !strings.Contains(out, "secret/stale") || !strings.Contains(out, "secret/gone") {
		t.Errorf("expected orphan paths in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Total orphans: 2") {
		t.Errorf("expected total count, got:\n%s", out)
	}
}

func TestFormatOrphans_Empty(t *testing.T) {
	out := FormatOrphans("dev", []string{})
	if !strings.Contains(out, "(none)") {
		t.Errorf("expected (none) for empty orphans, got:\n%s", out)
	}
}
