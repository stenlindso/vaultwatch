// Package classify provides path classification by sensitivity level.
package classify

import (
	"regexp"
	"strings"
)

// Level represents the sensitivity level of a secret path.
type Level int

const (
	LevelPublic   Level = iota // no special handling required
	LevelInternal              // internal use, limited exposure
	LevelSecret                // sensitive, should be audited
	LevelCritical              // highly sensitive, requires strict controls
)

func (l Level) String() string {
	switch l {
	case LevelPublic:
		return "public"
	case LevelInternal:
		return "internal"
	case LevelSecret:
		return "secret"
	case LevelCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// Rule maps a pattern to a classification level.
type Rule struct {
	Pattern string `json:"pattern"`
	Level   Level  `json:"level"`
	regexp  *regexp.Regexp
}

// Result holds the classification outcome for a single path.
type Result struct {
	Path  string
	Level Level
	Rule  string // matched rule pattern, empty if default
}

// Classifier assigns sensitivity levels to secret paths.
type Classifier struct {
	rules []Rule
}

// New creates a Classifier from the provided rules.
// Rules are evaluated in order; the first match wins.
func New(rules []Rule) (*Classifier, error) {
	compiled := make([]Rule, 0, len(rules))
	for _, r := range rules {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, err
		}
		r.regexp = re
		compiled = append(compiled, r)
	}
	return &Classifier{rules: compiled}, nil
}

// Classify returns the sensitivity level for a single path.
func (c *Classifier) Classify(path string) Result {
	norm := strings.ToLower(path)
	for _, r := range c.rules {
		if r.regexp.MatchString(norm) {
			return Result{Path: path, Level: r.Level, Rule: r.Pattern}
		}
	}
	return Result{Path: path, Level: LevelPublic}
}

// ClassifyAll classifies a slice of paths and returns results grouped by level.
func (c *Classifier) ClassifyAll(paths []string) map[Level][]Result {
	out := make(map[Level][]Result)
	for _, p := range paths {
		r := c.Classify(p)
		out[r.Level] = append(out[r.Level], r)
	}
	return out
}
