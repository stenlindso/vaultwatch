// Package filter provides path filtering utilities for Vault secret paths,
// supporting prefix, suffix, and regex-based inclusion/exclusion rules.
package filter

import (
	"regexp"
	"strings"
)

// Rule defines a single filter rule applied to a secret path.
type Rule struct {
	Type    string `json:"type"`    // "prefix", "suffix", or "regex"
	Pattern string `json:"pattern"` // the pattern to match
	Exclude bool   `json:"exclude"` // if true, matching paths are excluded
}

// Filter applies a set of rules to a list of paths.
type Filter struct {
	rules []Rule
}

// New creates a new Filter with the given rules.
func New(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Apply returns only the paths that pass all filter rules.
// Exclusion rules take priority over inclusion rules.
func (f *Filter) Apply(paths []string) ([]string, error) {
	var result []string
	for _, path := range paths {
		include, err := f.evaluate(path)
		if err != nil {
			return nil, err
		}
		if include {
			result = append(result, path)
		}
	}
	return result, nil
}

// evaluate returns true if the path should be included after applying all rules.
func (f *Filter) evaluate(path string) (bool, error) {
	for _, rule := range f.rules {
		matched, err := matches(rule, path)
		if err != nil {
			return false, err
		}
		if matched && rule.Exclude {
			return false, nil
		}
		if matched && !rule.Exclude {
			return true, nil
		}
	}
	// No inclusion rules defined — include everything by default.
	hasInclusion := false
	for _, rule := range f.rules {
		if !rule.Exclude {
			hasInclusion = true
			break
		}
	}
	return !hasInclusion, nil
}

func matches(rule Rule, path string) (bool, error) {
	switch rule.Type {
	case "prefix":
		return strings.HasPrefix(path, rule.Pattern), nil
	case "suffix":
		return strings.HasSuffix(path, rule.Pattern), nil
	case "regex":
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return false, err
		}
		return re.MatchString(path), nil
	default:
		return false, nil
	}
}
