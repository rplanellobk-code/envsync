package env

import (
	"fmt"
	"strings"
)

// ImportFormat represents a supported import source format.
type ImportFormat string

const (
	FormatDotenv ImportFormat = "dotenv"
	FormatShell  ImportFormat = "shell"
	FormatDocker ImportFormat = "docker"
)

// ImportOptions configures how an import is performed.
type ImportOptions struct {
	Format      ImportFormat
	Overwrite   bool
	IgnoreKeys  []string
}

// Import parses raw input in the given format and merges it into dst.
// Keys listed in IgnoreKeys are skipped. If Overwrite is false, existing
// keys in dst are preserved.
func Import(dst map[string]string, input string, opts ImportOptions) (map[string]string, error) {
	parsed, err := parseFormatted(input, opts.Format)
	if err != nil {
		return nil, fmt.Errorf("import: %w", err)
	}

	ignore := make(map[string]struct{}, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignore[k] = struct{}{}
	}

	result := make(map[string]string, len(dst))
	for k, v := range dst {
		result[k] = v
	}

	for k, v := range parsed {
		if _, skip := ignore[k]; skip {
			continue
		}
		if _, exists := result[k]; exists && !opts.Overwrite {
			continue
		}
		result[k] = v
	}

	return result, nil
}

// parseFormatted dispatches to the appropriate parser based on format.
func parseFormatted(input string, format ImportFormat) (map[string]string, error) {
	switch format {
	case FormatDotenv, "":
		return Parse(input)
	case FormatShell:
		return parseShellExport(input)
	case FormatDocker:
		return parseDockerEnv(input)
	default:
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
}

// parseShellExport handles lines like: export KEY=VALUE
func parseShellExport(input string) (map[string]string, error) {
	lines := strings.Split(input, "\n")
	result := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[strings.TrimSpace(parts[0])] = strings.Trim(strings.TrimSpace(parts[1]), `"`)
	}
	return result, nil
}

// parseDockerEnv handles lines like: KEY=VALUE (no quoting, docker --env-file style)
func parseDockerEnv(input string) (map[string]string, error) {
	return Parse(input)
}
