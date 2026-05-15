// Package env provides utilities for managing environment variable files.
//
// # Import
//
// The Import function merges environment variables from an external source
// into an existing map. It supports three input formats:
//
//   - FormatDotenv  — standard KEY=VALUE pairs (default)
//   - FormatShell   — shell-style `export KEY=VALUE` lines
//   - FormatDocker  — docker --env-file compatible KEY=VALUE pairs
//
// ImportOptions controls behaviour:
//
//   - Overwrite: when true, imported values replace existing keys in dst.
//   - IgnoreKeys: a list of keys that will never be imported.
//
// Example:
//
//	result, err := env.Import(current, rawInput, env.ImportOptions{
//		Format:    env.FormatShell,
//		Overwrite: false,
//		IgnoreKeys: []string{"AWS_SECRET_ACCESS_KEY"},
//	})
package env
