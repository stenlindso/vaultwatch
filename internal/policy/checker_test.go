package policy

import (
	"testing"
)

func TestCheck_DenyRule(t *testing.T) {
	rules := []Rule{
		{Pattern: `^secret/prod/.*`, Deny: true},
	}
	checker := NewChecker(rules)
	paths := []string{
		"secret/prod/db",
		"secret/staging/db",
	}
	violations := checker.Check(paths)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Path != "secret/prod/db" {
		t.Errorf("unexpected violation path: %s", violations[0].Path)
	}
}

func TestCheck_RequiredRule_Missing(t *testing.T) {
	rules := []Rule{
		{Pattern: `^secret/prod/tls`, Required: true},
	}
	checker := NewChecker(rules)
	paths := []string{"secret/prod/db"}
	violations := checker.Check(paths)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Path != "" {
		t.Errorf("expected empty path for required rule violation")
	}
}

func TestCheck_RequiredRule_Present(t *testing.T) {
	rules := []Rule{
		{Pattern: `^secret/prod/tls`, Required: true},
	}
	checker := NewChecker(rules)
	paths := []string{"secret/prod/tls", "secret/prod/db"}
	violations := checker.Check(paths)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_NoRules(t *testing.T) {
	checker := NewChecker(nil)
	violations := checker.Check([]string{"secret/prod/db"})
	if len(violations) != 0 {
		t.Errorf("expected no violations with no rules")
	}
}

func TestCheck_MultipleViolations(t *testing.T) {
	rules := []Rule{
		{Pattern: `^secret/prod/.*`, Deny: true},
	}
	checker := NewChecker(rules)
	paths := []string{
		"secret/prod/db",
		"secret/prod/api",
		"secret/staging/db",
	}
	violations := checker.Check(paths)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}

func TestSummary_NoViolations(t *testing.T) {
	s := Summary(nil)
	if s != "no policy violations found" {
		t.Errorf("unexpected summary: %s", s)
	}
}
