package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Format represents the output format for exported data.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Record represents a single exported secret path entry.
type Record struct {
	Environment string    `json:"environment"`
	Path        string    `json:"path"`
	CapturedAt  time.Time `json:"captured_at"`
}

// Export writes records to the given writer in the specified format.
func Export(w io.Writer, records []Record, format Format) error {
	switch format {
	case FormatJSON:
		return exportJSON(w, records)
	case FormatCSV:
		return exportCSV(w, records)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func exportJSON(w io.Writer, records []Record) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func exportCSV(w io.Writer, records []Record) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"environment", "path", "captured_at"}); err != nil {
		return fmt.Errorf("writing csv header: %w", err)
	}
	for _, r := range records {
		if err := cw.Write([]string{r.Environment, r.Path, r.CapturedAt.UTC().Format(time.RFC3339)}); err != nil {
			return fmt.Errorf("writing csv record: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}
