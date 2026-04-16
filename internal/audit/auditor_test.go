package audit_test

import (
	"os"
	"testing"

	"github.com/vaultwatch/internal/audit"
	"github.com/vaultwatch/internal/snapshot"
)

func setupManager(t *testing.T) *snapshot.Manager {
	t.Helper()
	dir := t.TempDir()
	return snapshot.NewManager(dir)
}

func TestAudit_BasicDiff(t *testing.T) {
	m := setupManager(t)

	m.Save("prod", &snapshot.Snapshot{Paths: []string{"secret/a", "secret/b"}})
	m.Save("staging", &snapshot.Snapshot{Paths: []string{"secret/b", "secret/c"}})

	a := audit.NewAuditor(m)
	report, err := a.Audit("prod", "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(report.Result.OnlyInA) != 1 || report.Result.OnlyInA[0] != "secret/a" {
		t.Errorf("expected OnlyInA=[secret/a], got %v", report.Result.OnlyInA)
	}
	if len(report.Result.OnlyInB) != 1 || report.Result.OnlyInB[0] != "secret/c" {
		t.Errorf("expected OnlyInB=[secret/c], got %v", report.Result.OnlyInB)
	}
	if len(report.Result.InBoth) != 1 || report.Result.InBoth[0] != "secret/b" {
		t.Errorf("expected InBoth=[secret/b], got %v", report.Result.InBoth)
	}
}

func TestAudit_MissingEnvironment(t *testing.T) {
	m := setupManager(t)
	a := audit.NewAuditor(m)

	_, err := a.Audit("ghost", "phantom")
	if err == nil {
		t.Fatal("expected error for missing environments")
	}
}

func TestAudit_ReportMetadata(t *testing.T) {
	m := setupManager(t)
	m.Save("dev", &snapshot.Snapshot{Paths: []string{"secret/x"}})
	m.Save("prod", &snapshot.Snapshot{Paths: []string{"secret/x"}})

	a := audit.NewAuditor(m)
	report, err := a.Audit("dev", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.EnvironmentA != "dev" || report.EnvironmentB != "prod" {
		t.Errorf("unexpected environment labels: %q %q", report.EnvironmentA, report.EnvironmentB)
	}
	if report.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	_ = os.Getenv("CI") // suppress unused import lint
}
