// Package watch provides continuous monitoring of Vault secret paths.
//
// A Watcher polls a Vault instance at a configurable interval, computes the
// difference between successive path listings, and emits Events on a channel
// whenever paths are added or removed.
//
// A Notifier consumes those Events and writes human-readable output to any
// io.Writer, making it easy to integrate with CLI output or log streams.
package watch
