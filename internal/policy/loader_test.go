package policy

import (
	"encoding/json"
	"os"
	"testing"
)

func writeTempPolicy(t *testing.T, rules []Rule) string {
	t.Helper()
	f, err := os.CreateTemp("", "policy-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(PolicyFile{Rules: rules}); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestLoadFromFile_Valid(t *testing.T) {
	rules := []Rule{
		{Pattern: `^secret/prod`, Deny: true},
		{Pattern: `^secret/prod/tls`, Required: true},
	}
	path := writeTempPolicy(t, rules)
	defer os.Remove(path)

	loaded, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded) != 2 {
		t.Errorf("expected 2 rules, got %d", len(loaded))
	}
}

func TestLoadFromFile_Missing(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/policy.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFromFile_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp("", "bad-policy-*.json")
	f.WriteString("not json")
	f.Close()
	defer os.Remove(f.Name())

	_, err := LoadFromFile(f.Name())
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDefaultRules_NotEmpty(t *testing.T) {
	rules := DefaultRules()
	if len(rules) == 0 {
		t.Error("expected at least one default rule")
	}
}
