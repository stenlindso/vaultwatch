package policy

import (
	"strings"
	"testing"
)

func makeChecker(rules []Rule) *Checker {
	return NewChecker(rules)
}

func TestEvaluate_NoViolations(t *testing.T) {
	c := makeChecker([]Rule{
		{Type: RuleRequired, Pattern: "secret/app"},
	})
	e := NewEvaluator(c)
	result := e.Evaluate("prod", []string{"secret/app", "secret/db"})
	if !result.Passed {
		t.Errorf("expected passed, got violations: %v", result.Violations)
	}
	if result.Environment != "prod" {
		t.Errorf("expected env prod, got %s", result.Environment)
	}
}

func TestEvaluate_WithViolation(t *testing.T) {
	c := makeChecker([]Rule{
		{Type: RuleDeny, Pattern: "secret/forbidden"},
	})
	e := NewEvaluator(c)
	result := e.Evaluate("staging", []string{"secret/app", "secret/forbidden"})
	if result.Passed {
		t.Error("expected failure due to deny rule")
	}
	if len(result.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(result.Violations))
	}
	if !strings.Contains(result.Violations[0], "DENY") {
		t.Errorf("expected DENY in violation message, got: %s", result.Violations[0])
	}
}

func TestEvaluateAll_MultipleEnvs(t *testing.T) {
	c := makeChecker([]Rule{
		{Type: RuleRequired, Pattern: "secret/common"},
	})
	e := NewEvaluator(c)
	envPaths := map[string][]string{
		"prod":    {"secret/common"},
		"staging": {"secret/other"},
	}
	results := e.EvaluateAll(envPaths)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	passed := 0
	for _, r := range results {
		if r.Passed {
			passed++
		}
	}
	if passed != 1 {
		t.Errorf("expected 1 passing env, got %d", passed)
	}
}
