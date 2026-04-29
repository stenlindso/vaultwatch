package schedule

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRegister_Valid(t *testing.T) {
	r := New()
	if err := r.Register("prod", time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.List()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.List()))
	}
}

func TestRegister_EmptyEnv(t *testing.T) {
	r := New()
	if err := r.Register("", time.Minute); err == nil {
		t.Fatal("expected error for empty env")
	}
}

func TestRegister_ZeroInterval(t *testing.T) {
	r := New()
	if err := r.Register("prod", 0); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestDue_ReturnsDueEntries(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	r := New()
	r.now = fixedNow(base)

	_ = r.Register("prod", time.Minute)

	// Advance clock past NextRun
	r.now = fixedNow(base.Add(2 * time.Minute))
	due := r.Due()
	if len(due) != 1 {
		t.Fatalf("expected 1 due entry, got %d", len(due))
	}
	if due[0].Environment != "prod" {
		t.Errorf("expected prod, got %s", due[0].Environment)
	}
}

func TestDue_NoneReady(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	r := New()
	r.now = fixedNow(base)
	_ = r.Register("prod", 10*time.Minute)

	due := r.Due()
	if len(due) != 0 {
		t.Fatalf("expected 0 due entries, got %d", len(due))
	}
}

func TestAdvance_UpdatesTimes(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	r := New()
	r.now = fixedNow(base)
	_ = r.Register("staging", time.Hour)

	r.now = fixedNow(base.Add(time.Hour))
	if err := r.Advance("staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e := r.entries["staging"]
	if e.LastRun.IsZero() {
		t.Error("LastRun should not be zero")
	}
	if !e.NextRun.After(e.LastRun) {
		t.Error("NextRun should be after LastRun")
	}
}

func TestAdvance_UnknownEnv(t *testing.T) {
	r := New()
	if err := r.Advance("ghost"); err == nil {
		t.Fatal("expected error for unknown env")
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	r := New()
	_ = r.Register("dev", time.Minute)
	r.Remove("dev")
	if len(r.List()) != 0 {
		t.Error("expected empty list after remove")
	}
}
