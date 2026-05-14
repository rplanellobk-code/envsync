package env

import (
	"fmt"
	"sort"
	"strings"
)

// TemplateResult holds the result of rendering a template against an env map.
type TemplateResult struct {
	Rendered  map[string]string
	Missing   []string
	Unused    []string
}

// RenderTemplate takes a template map (keys with optional default values encoded
// as "key=default" or just "key" for required) and an env map, and returns a
// TemplateResult. Required keys that are absent in env are collected in Missing.
func RenderTemplate(template map[string]string, env map[string]string) TemplateResult {
	rendered := make(map[string]string)
	var missing []string

	for key, defaultVal := range template {
		if val, ok := env[key]; ok {
			rendered[key] = val
		} else if defaultVal != "" {
			rendered[key] = defaultVal
		} else {
			missing = append(missing, key)
		}
	}

	var unused []string
	for key := range env {
		if _, ok := template[key]; !ok {
			unused = append(unused, key)
		}
	}

	sort.Strings(missing)
	sort.Strings(unused)

	return TemplateResult{
		Rendered: rendered,
		Missing:  missing,
		Unused:   unused,
	}
}

// ParseTemplate parses a .env-style template file where values are treated as
// defaults. Lines with no value ("KEY" or "KEY=") mark required keys.
func ParseTemplate(input string) (map[string]string, error) {
	tmpl := make(map[string]string)
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			// bare key — required, no default
			key := strings.TrimSpace(line)
			if key == "" {
				return nil, fmt.Errorf("empty key in template")
			}
			tmpl[key] = ""
			continue
		}
		key := strings.TrimSpace(line[:idx])
		if key == "" {
			return nil, fmt.Errorf("empty key in template")
		}
		tmpl[key] = strings.TrimSpace(line[idx+1:])
	}
	return tmpl, nil
}
