package summary

import (
	"testing"
	"time"
)

func makeInput(env string, paths, added, removed []string) Input {
	return Input{
		Environment: env,
		Paths:       paths,
		Added:       added,
		Removed:     removed,
		LastScanned: time.Now().UTC(),
	}
}

func TestSummarize_BasicAggregation(t *testing.T) {
	inputs := []Input{
		makeInput("prod", []string{"a", "b", "c"}, []string{"c"}, []string{}),
		makeInput("staging", []string{"a", "b"}, []string{}, []string{"b"}),
	}

	r := Summarize(inputs)

	if r.TotalPaths != 5 {
		t.Errorf("expected TotalPaths=5, got %d", r.TotalPaths)
	}
	if r.TotalAdded != 1 {
		t.Errorf("expected TotalAdded=1, got %d", r.TotalAdded)
	}
	if r.TotalRemoved != 1 {
		t.Errorf("expected TotalRemoved=1, got %d", r.TotalRemoved)
	}
}

func TestSummarize_SortedByEnvironment(t *testing.T) {
	inputs := []Input{
		makeInput("staging", []string{"x"}, nil, nil),
		makeInput("dev", []string{"y"}, nil, nil),
		makeInput("prod", []string{"z"}, nil, nil),
	}

	r := Summarize(inputs)

	order := []string{"dev", "prod", "staging"}
	for i, env := range order {
		if r.Environments[i].Environment != env {
			t.Errorf("index %d: expected %s, got %s", i, env, r.Environments[i].Environment)
		}
	}
}

func TestSummarize_Empty(t *testing.T) {
	r := Summarize(nil)

	if r.TotalPaths != 0 || r.TotalAdded != 0 || r.TotalRemoved != 0 {
		t.Error("expected all zero counts for empty input")
	}
	if r.HasChanges() {
		t.Error("expected HasChanges=false for empty input")
	}
}

func TestReport_HasChanges(t *testing.T) {
	inputs := []Input{
		makeInput("prod", []string{"a"}, []string{"a"}, nil),
	}
	r := Summarize(inputs)
	if !r.HasChanges() {
		t.Error("expected HasChanges=true")
	}
}

func TestReport_String(t *testing.T) {
	inputs := []Input{
		makeInput("prod", []string{"a", "b"}, []string{"b"}, []string{"c"}),
	}
	r := Summarize(inputs)
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string from Report.String()")
	}
}
