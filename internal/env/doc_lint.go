// Package env provides lint utilities for validating .env key-value maps
// against configurable rules.
//
// # Lint
//
// Lint runs a set of LintRule functions over an env map and returns all
// LintViolation values found. Built-in rules enforce:
//
//   - Keys are UPPER_SNAKE_CASE
//   - Values have no leading or trailing whitespace
//   - Values are non-empty
//
// Custom rules can be supplied to override or extend the defaults:
//
//	violations := env.Lint(myEnv, []env.LintRule{
//		{Name: "no-localhost", Message: "...", Check: func(k, v string) bool { ... }},
//	})
package env
