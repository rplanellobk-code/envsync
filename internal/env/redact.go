package env

import (
	"regexp"
	"strings"
)

// RedactOptions controls how sensitive values are masked.
type RedactOptions struct {
	// MaskChar is the character used to fill the masked region. Defaults to "*".
	MaskChar string
	// VisibleSuffix is the number of trailing characters to leave visible. Defaults to 0.
	VisibleSuffix int
}

var defaultRedactOpts = RedactOptions{MaskChar: "*", VisibleSuffix: 0}

// sensitivePattern matches common secret-like key names.
var sensitivePattern = regexp.MustCompile(
	`(?i)(password|passwd|secret|token|api[_-]?key|auth|private[_-]?key|credential|access[_-]?key|signing[_-]?key)`,
)

// IsSensitiveKey reports whether a key name looks like it holds a secret value.
func IsSensitiveKey(key string) bool {
	return sensitivePattern.MatchString(key)
}

// RedactValue masks a secret value according to opts.
// If opts is nil, the default options are used.
func RedactValue(value string, opts *RedactOptions) string {
	if opts == nil {
		opts = &defaultRedactOpts
	}
	mask := opts.MaskChar
	if mask == "" {
		mask = "*"
	}
	if opts.VisibleSuffix <= 0 || opts.VisibleSuffix >= len(value) {
		return strings.Repeat(mask, 8)
	}
	suffix := value[len(value)-opts.VisibleSuffix:]
	return strings.Repeat(mask, 8) + suffix
}

// Redact returns a copy of env where every value whose key is considered
// sensitive has been replaced with a masked string.
func Redact(env map[string]string, opts *RedactOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if IsSensitiveKey(k) {
			out[k] = RedactValue(v, opts)
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactKeys returns a copy of env where only the explicitly listed keys are
// redacted, regardless of their names.
func RedactKeys(env map[string]string, keys []string, opts *RedactOptions) map[string]string {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		if _, ok := set[k]; ok {
			out[k] = RedactValue(v, opts)
		} else {
			out[k] = v
		}
	}
	return out
}
