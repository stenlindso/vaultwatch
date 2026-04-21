package alert

import (
	"bytes"
	"strings"
	"testing"
)

func TestDispatch_NoAlert(t *testing.T) {
	called := false
	h := func(a *Alert) { called = true }
	d := NewDispatcher(DefaultThreshold(), h)
	d.Dispatch("staging", 1)
	if called {
		t.Error("handler should not be called below threshold")
	}
}

func TestDispatch_WarningTriggersHandler(t *testing.T) {
	var received *Alert
	h := func(a *Alert) { received = a }
	d := NewDispatcher(DefaultThreshold(), h)
	d.Dispatch("staging", 8)
	if received == nil {
		t.Fatal("expected handler to be called")
	}
	if received.Severity != SeverityWarning {
		t.Errorf("expected warning, got %s", received.Severity)
	}
}

func TestDispatch_MultipleHandlers(t *testing.T) {
	count := 0
	h := func(a *Alert) { count++ }
	d := NewDispatcher(DefaultThreshold(), h, h, h)
	d.Dispatch("prod", 25)
	if count != 3 {
		t.Errorf("expected 3 handler calls, got %d", count)
	}
}

func TestWriterHandler_Output(t *testing.T) {
	var buf bytes.Buffer
	h := WriterHandler(&buf)
	h(&Alert{
		Environment: "production",
		Severity:    SeverityCritical,
		Message:     "30 paths changed",
		PathCount:   30,
	})
	out := buf.String()
	if !strings.Contains(out, "production") {
		t.Errorf("expected env in output, got: %s", out)
	}
	if !strings.Contains(out, "critical") {
		t.Errorf("expected severity in output, got: %s", out)
	}
}
