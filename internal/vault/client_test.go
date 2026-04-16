package vault

import (
	"testing"
)

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "test-token")

	_, err := NewClient(Config{EnvName: "test"})
	if err == nil {
		t.Fatal("expected error when VAULT_ADDR is missing, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient(Config{EnvName: "test"})
	if err == nil {
		t.Fatal("expected error when VAULT_TOKEN is missing, got nil")
	}
}

func TestNewClient_ExplicitConfig(t *testing.T) {
	cfg := Config{
		Address: "http://127.0.0.1:8200",
		Token:   "root",
		EnvName: "dev",
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.EnvName != "dev" {
		t.Errorf("expected EnvName 'dev', got %q", client.EnvName)
	}
}

func TestNewClient_EnvOverride(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://env-vault:8200")
	t.Setenv("VAULT_TOKEN", "env-token")

	client, err := NewClient(Config{EnvName: "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}
