package policy

import (
	"fmt"
	"strings"
)

// EvalResult holds the outcome of evaluating a policy against a set of paths.
type EvalResult struct {
	Environment string
	Violations  []string
	Passed      bool
}

// Evaluator runs policy checks against snapshot paths for a given environment.
type Evaluator struct {
	checker *Checker
}

// NewEvaluator creates an Evaluator using the provided Checker.
func NewEvaluator(c *Checker) *Evaluator {
	return &Evaluator{checker: c}
}

// Evaluate runs all policy rules against the given paths and returns an EvalResult.
func (e *Evaluator) Evaluate(env string, paths []string) EvalResult {
	violations := e.checker.Check(paths)
	var msgs []string
	for _, v := range violations {
		msgs = append(msgs, fmt.Sprintf("[%s] %s: %s", strings.ToUpper(string(v.Rule.Type)), v.Rule.Pattern, v.Message))
	}
	return EvalResult{
		Environment: env,
		Violations:  msgs,
		Passed:      len(msgs) == 0,
	}
}

// EvaluateAll evaluates multiple environments and returns one result per environment.
func (e *Evaluator) EvaluateAll(envPaths map[string][]string) []EvalResult {
	results := make([]EvalResult, 0, len(envPaths))
	for env, paths := range envPaths {
		results = append(results, e.Evaluate(env, paths))
	}
	return results
}
