package trend

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var base = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func pts(env string, counts ...int) []DataPoint {
	out := make([]DataPoint, len(counts))
	for i, c := range counts {
		out[i] = DataPoint{Env: env, Count: c, Timestamp: base.Add(time.Duration(i) * time.Hour)}
	}
	return out
}

func TestAnalyze_Growing(t *testing.T) {
	results := Analyze(pts("prod", 10, 15, 20))
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Direction != DirectionGrowing {
		t.Errorf("expected growing, got %s", r.Direction)
	}
	if r.Delta != 10 {
		t.Errorf("expected delta 10, got %d", r.Delta)
	}
	if r.Samples != 3 {
		t.Errorf("expected 3 samples, got %d", r.Samples)
	}
}

func TestAnalyze_Shrinking(t *testing.T) {
	results := Analyze(pts("staging", 30, 20, 10))
	if results[0].Direction != DirectionShrinking {
		t.Errorf("expected shrinking, got %s", results[0].Direction)
	}
	if results[0].Delta != -20 {
		t.Errorf("expected delta -20, got %d", results[0].Delta)
	}
}

func TestAnalyze_Stable(t *testing.T) {
	results := Analyze(pts("dev", 5, 5, 5))
	if results[0].Direction != DirectionStable {
		t.Errorf("expected stable, got %s", results[0].Direction)
	}
}

func TestAnalyze_SinglePoint(t *testing.T) {
	results := Analyze(pts("dev", 7))
	if results[0].Direction != DirectionStable {
		t.Errorf("single point should be stable")
	}
}

func TestAnalyze_MultipleEnvs_Sorted(t *testing.T) {
	all := append(pts("prod", 1, 3), pts("dev", 10, 8)...)
	results := Analyze(all)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Env != "dev" || results[1].Env != "prod" {
		t.Errorf("results not sorted by env name")
	}
}

func TestFormatText_ContainsEnv(t *testing.T) {
	results := Analyze(pts("prod", 5, 10))
	var buf bytes.Buffer
	FormatText(results, &buf)
	if !strings.Contains(buf.String(), "prod") {
		t.Error("expected env name in output")
	}
	if !strings.Contains(buf.String(), "growing") {
		t.Error("expected direction in output")
	}
}

func TestFormatText_Empty(t *testing.T) {
	var buf bytes.Buffer
	FormatText(nil, &buf)
	if !strings.Contains(buf.String(), "No trend data") {
		t.Error("expected empty message")
	}
}

func TestFormatOneLiner(t *testing.T) {
	results := Analyze(pts("prod", 10, 20))
	line := FormatOneLiner(results)
	if !strings.HasPrefix(line, "trends:") {
		t.Errorf("unexpected prefix: %s", line)
	}
	if !strings.Contains(line, "prod") {
		t.Error("expected env in one-liner")
	}
}
