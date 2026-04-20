// Package export provides functionality for exporting Vault secret path
// snapshots to external formats such as JSON and CSV.
//
// Use BuildRecords to convert persisted snapshot data from a Manager into
// a flat list of Records, then pass those records to Export to write the
// desired output format to any io.Writer.
//
// Supported formats: JSON, CSV.
package export
