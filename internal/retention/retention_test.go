package retention_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/retention"
)

func makeEntries(ages ...time.Duration) []retention.Entry {
	entries := make([]retention.Entry, len(ages))
	for i, age := range ages {
		entries[i] = retention.Entry{
			ID:        fmt.Sprintf("entry-%d", i),
			Env:       "production",
			CreatedAt: time.Now().Add(-age),
		}
	}
	return entries
}

import "fmt"

func TestApply_NoPolicy_RetainsAll(t *testing.T) {
	entries := makeEntries(1*time.Hour, 2*time.Hour, 3*time.Hour)
	p := retention.Policy{}
	r := retention.Apply(p, entries)
	if len(r.Retained) != 3 {
		t.Fatalf("expected 3 retained, got %d", len(r.Retained))
	}
	if len(r.Pruned) != 0 {
		t.Fatalf("expected 0 pruned, got %d", len(r.Pruned))
	}
}

func TestApply_MaxCount_PrunesOldest(t *testing.T) {
	entries := makeEntries(1*time.Hour, 2*time.Hour, 3*time.Hour, 4*time.Hour)
	p := retention.Policy{MaxCount: 2}
	r := retention.Apply(p, entries)
	if len(r.Retained) != 2 {
		t.Fatalf("expected 2 retained, got %d", len(r.Retained))
	}
	if len(r.Pruned) != 2 {
		t.Fatalf("expected 2 pruned, got %d", len(r.Pruned))
	}
}

func TestApply_MaxAge_PrunesStale(t *testing.T) {
	entries := makeEntries(1*time.Hour, 10*24*time.Hour, 40*24*time.Hour)
	p := retention.Policy{MaxAge: 30 * 24 * time.Hour}
	r := retention.Apply(p, entries)
	if len(r.Retained) != 2 {
		t.Fatalf("expected 2 retained, got %d", len(r.Retained))
	}
	if !r.HasPruned() {
		t.Fatal("expected HasPruned to be true")
	}
}

func TestApply_Combined_MaxCountWins(t *testing.T) {
	entries := makeEntries(1*time.Hour, 2*time.Hour, 3*time.Hour)
	p := retention.Policy{MaxAge: 30 * 24 * time.Hour, MaxCount: 1}
	r := retention.Apply(p, entries)
	if len(r.Retained) != 1 {
		t.Fatalf("expected 1 retained, got %d", len(r.Retained))
	}
}

func TestResult_Summary(t *testing.T) {
	r := retention.Result{
		Retained: []retention.Entry{{}, {}},
		Pruned:   []retention.Entry{{}},
	}
	s := r.Summary()
	if s != "retained=2 pruned=1" {
		t.Fatalf("unexpected summary: %s", s)
	}
}

func TestDefaultPolicy_NonZero(t *testing.T) {
	p := retention.DefaultPolicy()
	if p.MaxAge == 0 {
		t.Error("expected non-zero MaxAge")
	}
	if p.MaxCount == 0 {
		t.Error("expected non-zero MaxCount")
	}
}
