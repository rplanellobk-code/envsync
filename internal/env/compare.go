package env

import (
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the result of comparing two environments.
type CompareResult struct {
	EnvironmentA string
	EnvironmentB string
	OnlyInA      map[string]string
	OnlyInB      map[string]string
	Different    map[string][2]string // key -> [valueA, valueB]
	Common       map[string]string
}

// Summary returns a human-readable summary of the comparison.
func (r *CompareResult) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Comparing %s <-> %s\n", r.EnvironmentA, r.EnvironmentB)
	fmt.Fprintf(&sb, "  Common keys:       %d\n", len(r.Common))
	fmt.Fprintf(&sb, "  Only in %-10s %d\n", r.EnvironmentA+":", len(r.OnlyInA))
	fmt.Fprintf(&sb, "  Only in %-10s %d\n", r.EnvironmentB+":", len(r.OnlyInB))
	fmt.Fprintf(&sb, "  Value differences: %d\n", len(r.Different))
	return sb.String()
}

// DifferentKeys returns a sorted list of keys whose values differ.
func (r *CompareResult) DifferentKeys() []string {
	keys := make([]string, 0, len(r.Different))
	for k := range r.Different {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// IsIdentical returns true when both environments have the same keys and values.
func (r *CompareResult) IsIdentical() bool {
	return len(r.OnlyInA) == 0 && len(r.OnlyInB) == 0 && len(r.Different) == 0
}

// Compare compares two env maps and returns a CompareResult.
func Compare(envA, envB map[string]string, nameA, nameB string) *CompareResult {
	result := &CompareResult{
		EnvironmentA: nameA,
		EnvironmentB: nameB,
		OnlyInA:      make(map[string]string),
		OnlyInB:      make(map[string]string),
		Different:    make(map[string][2]string),
		Common:       make(map[string]string),
	}

	for k, va := range envA {
		if vb, ok := envB[k]; ok {
			if va == vb {
				result.Common[k] = va
			} else {
				result.Different[k] = [2]string{va, vb}
			}
		} else {
				result.OnlyInA[k] = va
			}
		}

	for k, vb := range envB {
		if _, ok := envA[k]; !ok {
				result.OnlyInB[k] = vb
			}
		}

	return result
}
