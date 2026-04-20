package export

import (
	"fmt"
	"time"

	"github.com/yourusername/vaultwatch/internal/snapshot"
)

// BuildRecords converts snapshot data for one or more environments into
// a flat slice of Records suitable for export.
func BuildRecords(mgr *snapshot.Manager, envs []string) ([]Record, error) {
	var records []Record
	for _, env := range envs {
		snap, err := mgr.Load(env)
		if err != nil {
			return nil, fmt.Errorf("loading snapshot for %q: %w", env, err)
		}
		for _, path := range snap.Paths {
			records = append(records, Record{
				Environment: env,
				Path:        path,
				CapturedAt:  snap.CapturedAt,
			})
		}
	}
	if records == nil {
		records = []Record{}
	}
	return records, nil
}

// BuildRecordsFromPaths creates records directly from a path list without
// requiring a persisted snapshot — useful for live exports.
func BuildRecordsFromPaths(env string, paths []string, capturedAt time.Time) []Record {
	records := make([]Record, 0, len(paths))
	for _, p := range paths {
		records = append(records, Record{
			Environment: env,
			Path:        p,
			CapturedAt:  capturedAt,
		})
	}
	return records
}
