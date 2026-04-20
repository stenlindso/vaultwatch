package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var sampleTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func sampleRecords() []Record {
	return []Record{
		{Environment: "prod", Path: "secret/db", CapturedAt: sampleTime},
		{Environment: "staging", Path: "secret/api", CapturedAt: sampleTime},
	}
}

func TestExportJSON_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, sampleRecords(), FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []Record
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 records, got %d", len(out))
	}
	if out[0].Environment != "prod" {
		t.Errorf("expected prod, got %s", out[0].Environment)
	}
}

func TestExportCSV_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, sampleRecords(), FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 { // header + 2 records
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "environment") {
		t.Errorf("expected CSV header, got %q", lines[0])
	}
	if !strings.Contains(lines[1], "prod") {
		t.Errorf("expected prod in first record, got %q", lines[1])
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, sampleRecords(), Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExportJSON_EmptyRecords(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, []Record{}, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[]") {
		t.Errorf("expected empty JSON array, got %q", buf.String())
	}
}
