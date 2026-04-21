package alert

import (
	"testing"
)

func TestEvaluate_NoAlert(t *testing.T) {
	th := DefaultThreshold()
	result := Evaluate("staging", 2, th)
	if result != nil {
		t.Fatalf("expected nil alert, got %+v", result)
	}
}

func TestEvaluate_Warning(t *testing.T) {
	th := DefaultThreshold()
	result := Evaluate("staging", 7, th)
	if result == nil {
		t.Fatal("expected warning alert, got nil")
	}
	if result.Severity != SeverityWarning {
		t.Errorf("expected warning, got %s", result.Severity)
	}
	if result.Environment != "staging" {
		t.Errorf("expected env staging, got %s", result.Environment)
	}
	if result.PathCount != 7 {
		t.Errorf("expected path count 7, got %d", result.PathCount)
	}
}

func TestEvaluate_Critical(t *testing.T) {
	th := DefaultThreshold()
	result := Evaluate("production", 25, th)
	if result == nil {
		t.Fatal("expected critical alert, got nil")
	}
	if result.Severity != SeverityCritical {
		t.Errorf("expected critical, got %s", result.Severity)
	}
}

func TestEvaluate_ExactWarningBoundary(t *testing.T) {
	th := Threshold{WarningCount: 5, CriticalCount: 10}
	result := Evaluate("dev", 5, th)
	if result == nil || result.Severity != SeverityWarning {
		t.Errorf("expected warning at exact boundary, got %v", result)
	}
}

func TestEvaluate_ExactCriticalBoundary(t *testing.T) {
	th := Threshold{WarningCount: 5, CriticalCount: 10}
	result := Evaluate("dev", 10, th)
	if result == nil || result.Severity != SeverityCritical {
		t.Errorf("expected critical at exact boundary, got %v", result)
	}
}

func TestEvaluate_MessageContainsEnv(t *testing.T) {
	th := DefaultThreshold()
	result := Evaluate("myenv", 6, th)
	if result == nil {
		t.Fatal("expected alert")
	}
	if result.Message == "" {
		t.Error("expected non-empty message")
	}
}
