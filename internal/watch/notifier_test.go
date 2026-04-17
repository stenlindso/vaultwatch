package watch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeEvent(env string, added, removed []string) Event {
	return Event{Env: env, Added: added, Removed: removed, At: time.Now()}
}

func TestHandle_AddedPaths(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)
	n.Handle(makeEvent("staging", []string{"secret/new"}, nil))
	out := buf.String()
	if !strings.Contains(out, "secret/new") {
		t.Errorf("expected output to contain secret/new, got: %s", out)
	}
	if !strings.Contains(out, "added") {
		t.Errorf("expected output to contain 'added', got: %s", out)
	}
}

func TestHandle_RemovedPaths(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)
	n.Handle(makeEvent("prod", nil, []string{"secret/old"}))
	out := buf.String()
	if !strings.Contains(out, "secret/old") {
		t.Errorf("expected output to contain secret/old, got: %s", out)
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected output to contain 'removed', got: %s", out)
	}
}

func TestHandle_EnvInOutput(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)
	n.Handle(makeEvent("dev", []string{"secret/x"}, nil))
	if !strings.Contains(buf.String(), "env=dev") {
		t.Errorf("expected env=dev in output, got: %s", buf.String())
	}
}

func TestFormatEvent(t *testing.T) {
	ev := makeEvent("qa", []string{"a", "b"}, []string{"c"})
	s := FormatEvent(ev)
	if !strings.Contains(s, "env=qa") || !strings.Contains(s, "added=2") || !strings.Contains(s, "removed=1") {
		t.Errorf("unexpected format: %s", s)
	}
}
