package rollup

import (
	"testing"
)

func TestRollup_BasicGrouping(t *testing.T) {
	envPaths := map[string][]string{
		"prod": {"secret/app/db", "secret/app/cache", "secret/infra/net"},
		"staging": {"secret/app/db", "secret/infra/net"},
	}

	result := Rollup(envPaths, 2)

	if result.TotalEnvs != 2 {
		t.Errorf("expected 2 total envs, got %d", result.TotalEnvs)
	}

	if len(result.Entries) == 0 {
		t.Fatal("expected non-empty entries")
	}

	prefixMap := make(map[string]Entry)
	for _, e := range result.Entries {
		prefixMap[e.Prefix] = e
	}

	appEntry, ok := prefixMap["secret/app"]
	if !ok {
		t.Fatal("expected entry for secret/app")
	}
	if appEntry.Count != 2 {
		t.Errorf("expected secret/app in 2 envs, got %d", appEntry.Count)
	}
}

func TestRollup_Coverage(t *testing.T) {
	envPaths := map[string][]string{
		"prod":    {"kv/shared/key"},
		"staging": {"kv/shared/key"},
		"dev":     {"kv/other/key"},
	}

	result := Rollup(envPaths, 1)

	for _, e := range result.Entries {
		if e.Prefix == "kv" {
			if e.Coverage != 1.0 {
				t.Errorf("expected full coverage for kv, got %.2f", e.Coverage)
			}
			return
		}
	}
	t.Error("expected entry for prefix 'kv'")
}

func TestRollup_EmptyInput(t *testing.T) {
	result := Rollup(map[string][]string{}, 1)
	if result.TotalEnvs != 0 {
		t.Errorf("expected 0 total envs")
	}
	if len(result.Entries) != 0 {
		t.Errorf("expected empty entries")
	}
}

func TestRollup_DepthOne(t *testing.T) {
	envPaths := map[string][]string{
		"prod": {"secret/app/db", "secret/infra/net", "kv/data/x"},
	}

	result := Rollup(envPaths, 1)

	prefixes := make(map[string]bool)
	for _, e := range result.Entries {
		prefixes[e.Prefix] = true
	}

	if !prefixes["secret"] {
		t.Error("expected prefix 'secret'")
	}
	if !prefixes["kv"] {
		t.Error("expected prefix 'kv'")
	}
}

func TestExtractPrefix(t *testing.T) {
	cases := []struct {
		path     string
		depth    int
		expected string
	}{
		{"secret/app/db", 1, "secret"},
		{"secret/app/db", 2, "secret/app"},
		{"/secret/app/db/", 2, "secret/app"},
		{"single", 2, "single"},
		{"", 1, ""},
	}

	for _, tc := range cases {
		got := extractPrefix(tc.path, tc.depth)
		if got != tc.expected {
			t.Errorf("extractPrefix(%q, %d) = %q, want %q", tc.path, tc.depth, got, tc.expected)
		}
	}
}
