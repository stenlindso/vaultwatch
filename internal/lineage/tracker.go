package lineage

import (
	"fmt"
	"sort"
)

// Tracker builds per-path Record views from raw entries.
type Tracker struct {
	store *Store
}

// NewTracker wraps a Store with higher-level query helpers.
func NewTracker(store *Store) *Tracker {
	return &Tracker{store: store}
}

// BuildRecords loads all entries for env and groups them by path.
func (t *Tracker) BuildRecords(env string) ([]Record, error) {
	entries, err := t.store.Load(env)
	if err != nil {
		return nil, fmt.Errorf("tracker: load %s: %w", env, err)
	}

	index := make(map[string]*Record)
	for _, e := range entries {
		if _, ok := index[e.Path]; !ok {
			index[e.Path] = &Record{Path: e.Path}
		}
		index[e.Path].History = append(index[e.Path].History, e)
	}

	records := make([]Record, 0, len(index))
	for _, r := range index {
		sort.Slice(r.History, func(i, j int) bool {
			return r.History[i].SeenAt.Before(r.History[j].SeenAt)
		})
		records = append(records, *r)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Path < records[j].Path
	})
	return records, nil
}

// FirstSeen returns the earliest Entry for a path in the given environment,
// or an error if the path has no recorded history.
func (t *Tracker) FirstSeen(env, path string) (Entry, error) {
	records, err := t.BuildRecords(env)
	if err != nil {
		return Entry{}, err
	}
	for _, r := range records {
		if r.Path == path && len(r.History) > 0 {
			return r.History[0], nil
		}
	}
	return Entry{}, fmt.Errorf("tracker: path %q not found in %s", path, env)
}
