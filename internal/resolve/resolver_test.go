package resolve

import (
	"testing"
)

func makeResolver() *Resolver {
	return New([]Alias{
		{Name: "infra", Prefix: "secret/infrastructure"},
		{Name: "app", Prefix: "secret/applications"},
	})
}

func TestResolve_NoAlias(t *testing.T) {
	r := makeResolver()
	got, err := r.Resolve("secret/other/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/other/path" {
		t.Errorf("expected unchanged path, got %q", got)
	}
}

func TestResolve_WithAlias(t *testing.T) {
	r := makeResolver()
	got, err := r.Resolve("infra/networking/vpn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "secret/infrastructure/networking/vpn"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestResolve_ExactAlias(t *testing.T) {
	r := makeResolver()
	got, err := r.Resolve("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/applications" {
		t.Errorf("expected %q, got %q", "secret/applications", got)
	}
}

func TestResolve_EmptyPath(t *testing.T) {
	r := makeResolver()
	_, err := r.Resolve("")
	if err == nil {
		t.Error("expected error for empty path, got nil")
	}
}

func TestResolveAll_Mixed(t *testing.T) {
	r := makeResolver()
	input := []string{"infra/db", "secret/manual", "app/svc"}
	got, err := r.ResolveAll(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{
		"secret/infrastructure/db",
		"secret/manual",
		"secret/applications/svc",
	}
	for i, e := range expected {
		if got[i] != e {
			t.Errorf("index %d: expected %q, got %q", i, e, got[i])
		}
	}
}

func TestAliases_ReturnsCopy(t *testing.T) {
	r := makeResolver()
	a := r.Aliases()
	if len(a) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(a))
	}
	// Mutating the copy should not affect the resolver.
	delete(a, "infra")
	if _, ok := r.Aliases()["infra"]; !ok {
		t.Error("mutation of returned map affected resolver")
	}
}
