package resolve

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempAliases(t *testing.T, data any) string {
	t.Helper()
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	p := filepath.Join(t.TempDir(), "aliases.json")
	if err := os.WriteFile(p, b, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return p
}

func TestLoadAliases_Valid(t *testing.T) {
	p := writeTempAliases(t, map[string]any{
		"aliases": []map[string]string{
			{"name": "prod", "prefix": "secret/prod"},
			{"name": "staging", "prefix": "secret/staging"},
		},
	})
	r, err := LoadAliases(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := r.Resolve("prod/db")
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if got != "secret/prod/db" {
		t.Errorf("expected %q, got %q", "secret/prod/db", got)
	}
}

func TestLoadAliases_MissingFile(t *testing.T) {
	_, err := LoadAliases("/nonexistent/aliases.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadAliases_InvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(p, []byte("not json"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := LoadAliases(p)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestDefaultAliases_NotEmpty(t *testing.T) {
	r := DefaultAliases()
	aliases := r.Aliases()
	if len(aliases) == 0 {
		t.Error("expected non-empty default aliases")
	}
}

func TestDefaultAliases_KVPrefix(t *testing.T) {
	r := DefaultAliases()
	got, err := r.Resolve("kv/myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "kv/data/myapp/config" {
		t.Errorf("expected %q, got %q", "kv/data/myapp/config", got)
	}
}
