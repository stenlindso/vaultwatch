package redact_test

import (
	"testing"

	"github.com/your-org/vaultwatch/internal/redact"
)

func TestApply_NoRules_PassesThrough(t *testing.T) {
	r, err := redact.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	paths := []string{"secret/prod/db", "secret/prod/api"}
	got := r.Apply(paths)
	if len(got) != len(paths) {
		t.Fatalf("expected %d paths, got %d", len(paths), len(got))
	}
}

func TestApply_SubstringRule_Omits(t *testing.T) {
	rules := []redact.Rule{
		{Pattern: "password", Replacement: ""},
	}
	r, _ := redact.New(rules)
	paths := []string{"secret/prod/password", "secret/prod/api-key"}
	got := r.Apply(paths)
	if len(got) != 1 || got[0] != "secret/prod/api-key" {
		t.Fatalf("expected only api-key path, got %v", got)
	}
}

func TestApply_SubstringRule_Replaces(t *testing.T) {
	rules := []redact.Rule{
		{Pattern: "prod", Replacement: "***"},
	}
	r, _ := redact.New(rules)
	paths := []string{"secret/prod/db"}
	got := r.Apply(paths)
	if len(got) != 1 || got[0] != "secret/***/db" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestApply_RegexRule_Omits(t *testing.T) {
	rules := []redact.Rule{
		{Pattern: `secret/prod/.*key.*`, Regex: true, Replacement: ""},
	}
	r, _ := redact.New(rules)
	paths := []string{"secret/prod/api-key", "secret/prod/db", "secret/prod/ssh-key"}
	got := r.Apply(paths)
	if len(got) != 1 || got[0] != "secret/prod/db" {
		t.Fatalf("expected only db path, got %v", got)
	}
}

func TestApply_RegexRule_Replaces(t *testing.T) {
	rules := []redact.Rule{
		{Pattern: `(token|secret)/([^/]+)`, Regex: true, Replacement: "$1/[REDACTED]"},
	}
	r, _ := redact.New(rules)
	paths := []string{"token/myvalue"}
	got := r.Apply(paths)
	if len(got) != 1 || got[0] != "token/[REDACTED]" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestNew_InvalidRegex_ReturnsError(t *testing.T) {
	rules := []redact.Rule{
		{Pattern: `[invalid`, Regex: true},
	}
	_, err := redact.New(rules)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestApply_MultipleRules_FirstMatchWins(t *testing.T) {
	rules := []redact.Rule{
		{Pattern: "prod", Replacement: "[ENV]"},
		{Pattern: "prod", Replacement: "[OTHER]"},
	}
	r, _ := redact.New(rules)
	paths := []string{"secret/prod/db"}
	got := r.Apply(paths)
	if len(got) != 1 || got[0] != "secret/[ENV]/db" {
		t.Fatalf("expected first rule to win, got %v", got)
	}
}
