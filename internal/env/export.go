package env

import (
	"fmt"
	"strings"
)

// ExportFormat represents the output format for exported env vars.
type ExportFormat int

const (
	// FormatDotenv outputs KEY=VALUE pairs (default .env format).
	FormatDotenv ExportFormat = iota
	// FormatExport outputs shell export statements.
	FormatExport
	// FormatDocker outputs --env KEY=VALUE flags for docker run.
	FormatDocker
)

// ExportOptions configures the Export function.
type ExportOptions struct {
	Format ExportFormat
	// OmitKeys excludes specific keys from the output.
	OmitKeys []string
}

// Export renders an env map to a string in the requested format.
// Keys are emitted in deterministic (sorted) order.
func Export(vars map[string]string, opts ExportOptions) string {
	omit := make(map[string]struct{}, len(opts.OmitKeys))
	for _, k := range opts.OmitKeys {
		omit[k] = struct{}{}
	}

	keys := sortedKeys(vars)
	lines := make([]string, 0, len(keys))

	for _, k := range keys {
		if _, skip := omit[k]; skip {
			continue
		}
		v := vars[k]
		switch opts.Format {
		case FormatExport:
			lines = append(lines, fmt.Sprintf("export %s=%s", k, quoteIfNeeded(v)))
		case FormatDocker:
			lines = append(lines, fmt.Sprintf("--env %s=%s", k, v))
		default:
			lines = append(lines, fmt.Sprintf("%s=%s", k, quoteIfNeeded(v)))
		}
	}

	return strings.Join(lines, "\n")
}

// quoteIfNeeded wraps the value in double quotes if it contains spaces or
// special shell characters that could cause issues when sourced.
func quoteIfNeeded(v string) string {
	if v == "" {
		return `""`
	}
	specials := " \t\n#$&|;<>(){}!*?[]"
	if strings.ContainsAny(v, specials) {
		escaped := strings.ReplaceAll(v, `"`, `\"`)
		return fmt.Sprintf(`"%s"`, escaped)
	}
	return v
}
