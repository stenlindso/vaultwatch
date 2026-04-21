package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/youorg/vaultwatch/internal/history"
)

func newTempStore(t *testing.T) *history.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "history-test-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	store, err := history.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func TestRecord_And_List(t *testing.T) {
	store := newTempStore(t)

	e := history.Entry{
		ID:          "entry-001",
		Environment: "staging",
		Timestamp:   time.Now().UTC(),
		Added:       []string{"secret/foo"},
		Removed:     []string{},
		Violations:  0,
	}
	if err := store.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := store.List("staging")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].ID != "entry-001" {
		t.Errorf("expected ID entry-001, got %s", entries[0].ID)
	}
}

func TestList_Empty(t *testing.T) {
	store := newTempStore(t)
	entries, err := store.List("nonexistent")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestList_SortedByTimestamp(t *testing.T) {
	store := newTempStore(t)
	base := time.Now().UTC()

	for i, id := range []string{"c", "a", "b"} {
		err := store.Record(history.Entry{
			ID:          id,
			Environment: "prod",
			Timestamp:   base.Add(time.Duration(i) * time.Second),
		})
		if err != nil {
			t.Fatalf("Record %s: %v", id, err)
		}
	}

	entries, err := store.List("prod")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Timestamp.Before(entries[i-1].Timestamp) {
			t.Errorf("entries not sorted: index %d before %d", i, i-1)
		}
	}
}

func TestRecord_AutoID(t *testing.T) {
	store := newTempStore(t)
	e := history.Entry{
		Environment: "dev",
		Timestamp:   time.Now().UTC(),
	}
	if err := store.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, err := store.List("dev")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].ID == "" {
		t.Error("expected auto-generated ID, got empty string")
	}
}
