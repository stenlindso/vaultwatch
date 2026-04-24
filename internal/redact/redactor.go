// Package redact provides utilities for redacting sensitive secret paths
// and values from output before display or export.
package redact

import (
	"regexp"
	"strings"
)

// Rule defines a single redaction rule matched against secret paths.
type Rule struct {
	// Pattern is a substring or regex pattern to match against a path.
	Pattern string `json:"pattern"`
	// Regex indicates whether Pattern should be treated as a regular expression.
	Regex bool `json:"regex"`
	// Replacement is the string used in place of the matched path segment.
	Replacement string `json:"replacement"`
}

// Redactor applies redaction rules to a list of secret paths.
type Redactor struct {
	rules   []Rule
	compiled []*regexp.Regexp
}

// New creates a Redactor from the provided rules. Returns an error if any
// regex rule fails to compile.
func New(rules []Rule) (*Redactor, error) {
	compiled := make([]*regexp.Regexp, len(rules))
	for i, r := range rules {
		if r.Regex {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, err
			}
			compiled[i] = re
		}
	}
	return &Redactor{rules: rules, compiled: compiled}, nil
}

// Apply returns a new slice of paths with sensitive entries replaced or
// removed according to the configured rules. Paths that match a rule with
// an empty Replacement are omitted entirely.
func (r *Redactor) Apply(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		result, redacted := r.redact(p)
		if redacted && result == "" {
			continue
		}
		out = append(out, result)
	}
	return out
}

// redact returns the (possibly replaced) path and whether any rule matched.
func (r *Redactor) redact(path string) (string, bool) {
	for i, rule := range r.rules {
		if rule.Regex {
			re := r.compiled[i]
			if re != nil && re.MatchString(path) {
				if rule.Replacement == "" {
					return "", true
				}
				return re.ReplaceAllString(path, rule.Replacement), true
			}
		} else {
			if strings.Contains(path, rule.Pattern) {
				if rule.Replacement == "" {
					return "", true
				}
				return strings.ReplaceAll(path, rule.Pattern, rule.Replacement), true
			}
		}
	}
	return path, false
}
