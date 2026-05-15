package env

import (
	"fmt"
	"regexp"
	"strings"
)

// LintRule defines a single rule applied to an env map.
type LintRule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// LintViolation records a rule violation for a specific key.
type LintViolation struct {
	Key     string
	Rule    string
	Message string
}

func (v LintViolation) Error() string {
	return fmt.Sprintf("key %q violated rule %q: %s", v.Key, v.Rule, v.Message)
}

var upperSnakeRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// DefaultRules returns the built-in lint rules.
func DefaultRules() []LintRule {
	return []LintRule{
		{
			Name:    "uppercase-key",
			Message: "key must be UPPER_SNAKE_CASE",
			Check: func(key, _ string) bool {
				return !upperSnakeRe.MatchString(key)
			},
		},
		{
			Name:    "no-trailing-space",
			Message: "value must not have leading or trailing whitespace",
			Check: func(_, value string) bool {
				return value != strings.TrimSpace(value)
			},
		},
		{
			Name:    "no-empty-value",
			Message: "value must not be empty",
			Check: func(_, value string) bool {
				return value == ""
			},
		},
	}
}

// Lint runs the provided rules against env and returns all violations.
// If rules is nil, DefaultRules() is used.
func Lint(env map[string]string, rules []LintRule) []LintViolation {
	if rules == nil {
		rules = DefaultRules()
	}
	var violations []LintViolation
	for _, key := range sortedKeys(env) {
		value := env[key]
		for _, rule := range rules {
			if rule.Check(key, value) {
				violations = append(violations, LintViolation{
					Key:     key,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}
	return violations
}
