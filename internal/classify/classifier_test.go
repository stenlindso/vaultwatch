package classify

import (
	"testing"
)

func defaultRules() []Rule {
	return []Rule{
		{Pattern: `(password|passwd|secret|token)`, Level: LevelCritical},
		{Pattern: `(key|cert|tls|ssl)`, Level: LevelSecret},
		{Pattern: `(config|setting)`, Level: LevelInternal},
	}
}

func TestClassify_MatchesCritical(t *testing.T) {
	c, err := New(defaultRules())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := c.Classify("prod/app/db-password")
	if r.Level != LevelCritical {
		t.Errorf("expected critical, got %s", r.Level)
	}
	if r.Rule == "" {
		t.Error("expected non-empty rule")
	}
}

func TestClassify_MatchesSecret(t *testing.T) {
	c, _ := New(defaultRules())
	r := c.Classify("prod/app/tls-cert")
	if r.Level != LevelSecret {
		t.Errorf("expected secret, got %s", r.Level)
	}
}

func TestClassify_DefaultsToPublic(t *testing.T) {
	c, _ := New(defaultRules())
	r := c.Classify("prod/app/version")
	if r.Level != LevelPublic {
		t.Errorf("expected public, got %s", r.Level)
	}
	if r.Rule != "" {
		t.Errorf("expected empty rule for default, got %q", r.Rule)
	}
}

func TestClassify_CaseInsensitive(t *testing.T) {
	c, _ := New(defaultRules())
	r := c.Classify("prod/app/DB-PASSWORD")
	if r.Level != LevelCritical {
		t.Errorf("expected critical for uppercase path, got %s", r.Level)
	}
}

func TestClassifyAll_GroupsByLevel(t *testing.T) {
	c, _ := New(defaultRules())
	paths := []string{
		"prod/app/db-password",
		"prod/app/api-token",
		"prod/app/tls-key",
		"prod/app/version",
	}
	results := c.ClassifyAll(paths)
	if len(results[LevelCritical]) != 2 {
		t.Errorf("expected 2 critical, got %d", len(results[LevelCritical]))
	}
	if len(results[LevelSecret]) != 1 {
		t.Errorf("expected 1 secret, got %d", len(results[LevelSecret]))
	}
	if len(results[LevelPublic]) != 1 {
		t.Errorf("expected 1 public, got %d", len(results[LevelPublic]))
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]Rule{{Pattern: `[invalid`, Level: LevelCritical}})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestLevelString(t *testing.T) {
	cases := map[Level]string{
		LevelPublic:   "public",
		LevelInternal: "internal",
		LevelSecret:   "secret",
		LevelCritical: "critical",
	}
	for l, want := range cases {
		if got := l.String(); got != want {
			t.Errorf("Level(%d).String() = %q, want %q", l, got, want)
		}
	}
}
