package env

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Scope represents a named filter over environment variable keys,
// allowing operations to target subsets of variables by prefix or pattern.
type Scope struct {
	Name    string `json:"name"`
	Prefix  string `json:"prefix,omitempty"`
	Pattern string `json:"pattern,omitempty"`

	compiled *regexp.Regexp
}

// Match reports whether the given key falls within this scope.
func (s *Scope) Match(key string) bool {
	if s.Prefix != "" && !strings.HasPrefix(key, s.Prefix) {
		return false
	}
	if s.compiled != nil {
		return s.compiled.MatchString(key)
	}
	return true
}

// Filter returns a new map containing only the keys matched by this scope.
func (s *Scope) Filter(env map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range env {
		if s.Match(k) {
			out[k] = v
		}
	}
	return out
}

// compile pre-compiles the Pattern field if set.
func (s *Scope) compile() error {
	if s.Pattern == "" {
		return nil
	}
	re, err := regexp.Compile(s.Pattern)
	if err != nil {
		return fmt.Errorf("scope %q: invalid pattern: %w", s.Name, err)
	}
	s.compiled = re
	return nil
}

// marshalScope is used for JSON round-trips (compiled field excluded).
type marshalScope struct {
	Name    string `json:"name"`
	Prefix  string `json:"prefix,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (s Scope) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalScope{Name: s.Name, Prefix: s.Prefix, Pattern: s.Pattern})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Scope) UnmarshalJSON(data []byte) error {
	var m marshalScope
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	s.Name = m.Name
	s.Prefix = m.Prefix
	s.Pattern = m.Pattern
	return s.compile()
}
