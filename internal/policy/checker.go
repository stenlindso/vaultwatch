package policy

import (
	"regexp"
	"strings"
)

// Rule defines a policy rule for secret path validation.
type Rule struct {
	Pattern  string `json:"pattern"`
	Required bool   `json:"required"`
	Deny     bool   `json:"deny"`
}

// Violation represents a policy violation found during a check.
type Violation struct {
	Path    string
	Rule    Rule
	Message string
}

// Checker evaluates secret paths against a set of policy rules.
type Checker struct {
	rules []Rule
}

// NewChecker creates a Checker with the given rules.
func NewChecker(rules []Rule) *Checker {
	return &Checker{rules: rules}
}

// Check evaluates paths against all rules and returns any violations.
func (c *Checker) Check(paths []string) []Violation {
	var violations []Violation

	for _, rule := range c.rules {
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			continue
		}

		if rule.Deny {
			for _, p := range paths {
				if re.MatchString(p) {
					violations = append(violations, Violation{
						Path:    p,
						Rule:    rule,
						Message: "path matches deny rule: " + rule.Pattern,
					})
				}
			}
		}

		if rule.Required {
			matched := false
			for _, p := range paths {
				if re.MatchString(p) {
					matched = true
					break
				}
			}
			if !matched {
				violations = append(violations, Violation{
					Path:    "",
					Rule:    rule,
					Message: "required pattern not found: " + rule.Pattern,
				})
			}
		}
	}

	return violations
}

// Summary returns a human-readable summary of violations.
func Summary(violations []Violation) string {
	if len(violations) == 0 {
		return "no policy violations found"
	}
	var sb strings.Builder
	for _, v := range violations {
		sb.WriteString("VIOLATION: " + v.Message)
		if v.Path != "" {
			sb.WriteString(" [path: " + v.Path + "]")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
