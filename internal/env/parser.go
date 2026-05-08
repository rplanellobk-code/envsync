// Package env provides utilities for parsing and serializing .env files.
package env

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Map represents a set of key-value pairs parsed from a .env file.
type Map map[string]string

// Parse reads key=value pairs from r, skipping blank lines and comments.
// Lines beginning with '#' are treated as comments.
func Parse(r io.Reader) (Map, error) {
	result := make(Map)
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("env: invalid syntax on line %d: %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("env: empty key on line %d", lineNum)
		}

		// Strip optional surrounding quotes from value.
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env: scanner error: %w", err)
	}

	return result, nil
}

// Serialize writes the Map to w in sorted key=value format.
func Serialize(m Map, w io.Writer) error {
	keys := sortedKeys(m)
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, m[k]); err != nil {
			return fmt.Errorf("env: write error: %w", err)
		}
	}
	return nil
}

func sortedKeys(m Map) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Simple insertion sort — env files are typically small.
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
