package env

import (
	"errors"
	"fmt"
	"regexp"
)

// SchemaField describes a single expected environment variable.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern the value must match
}

// Schema holds a set of field definitions for an environment.
type Schema struct {
	Fields []SchemaField
}

// ValidationError holds all violations found during schema validation.
type ValidationError struct {
	Violations []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("schema validation failed with %d violation(s): %v", len(e.Violations), e.Violations)
}

// Validate checks the provided env map against the schema.
// It returns a *ValidationError listing all violations, or nil if valid.
func (s *Schema) Validate(env map[string]string) error {
	var violations []string

	for _, field := range s.Fields {
		val, ok := env[field.Key]
		if !ok || val == "" {
			if field.Required {
				violations = append(violations, fmt.Sprintf("missing required key: %s", field.Key))
			}
			continue
		}

		if field.Pattern != "" {
			matched, err := regexp.MatchString(field.Pattern, val)
			if err != nil {
				violations = append(violations, fmt.Sprintf("invalid pattern for key %s: %v", field.Key, err))
				continue
			}
			if !matched {
				violations = append(violations, fmt.Sprintf("key %s value %q does not match pattern %q", field.Key, val, field.Pattern))
			}
		}
	}

	if len(violations) > 0 {
		return &ValidationError{Violations: violations}
	}
	return nil
}

// ParseSchema builds a Schema from a slice of SchemaField definitions.
// Returns an error if any field key is empty.
func ParseSchema(fields []SchemaField) (*Schema, error) {
	for _, f := range fields {
		if f.Key == "" {
			return nil, errors.New("schema field key must not be empty")
		}
	}
	return &Schema{Fields: fields}, nil
}
