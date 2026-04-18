package config

import (
	"encoding/json"
	"os"
	"testing"
)

func writeTempConfig(t *testing.T, cfg Config) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultwatch-cfg-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		t.Fatalf("encode config: %v", err)
	}
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	cfg := Config{
		Environments: map[string]EnvConfig{
			"prod": {Address: "https://vault.prod", Token: "tok", MountPath: "secret"},
		},
		SnapshotDir: "/tmp/snaps",
	}
	path := writeTempConfig(t, cfg)
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.Environments["prod"].Address != "https://vault.prod" {
		t.Errorf("expected prod address, got %q", loaded.Environments["prod"].Address)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/cfg.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad-*.json")
	f.WriteString("{not valid json")
	f.Close()
	_, err := Load(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoad_NoEnvironments(t *testing.T) {
	cfg := Config{SnapshotDir: "/tmp"}
	path := writeTempConfig(t, cfg)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty environments")
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	cfg := Config{
		Environments: map[string]EnvConfig{
			"staging": {Token: "tok"},
		},
	}
	path := writeTempConfig(t, cfg)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing address")
	}
}
