package env

import "fmt"

// ScopedDiff computes a Diff restricted to the keys matched by the given Scope.
// Only changes involving keys within the scope are returned.
func ScopedDiff(scope *Scope, local, remote map[string]string) []DiffEntry {
	l := scope.Filter(local)
	r := scope.Filter(remote)
	return Diff(l, r)
}

// ScopedApply applies only the DiffEntries whose keys are matched by scope.
// Entries outside the scope are silently skipped.
func ScopedApply(scope *Scope, base map[string]string, entries []DiffEntry) map[string]string {
	filtered := make([]DiffEntry, 0, len(entries))
	for _, e := range entries {
		if scope.Match(e.Key) {
			filtered = append(filtered, e)
		}
	}
	return Apply(base, filtered)
}

// ScopedMerge merges local and remote maps considering only keys within scope,
// using the provided MergeStrategy. Keys outside the scope are taken from local
// unchanged.
func ScopedMerge(scope *Scope, local, remote map[string]string, strategy MergeStrategy) (map[string]string, error) {
	// Build scoped views.
	scopedLocal := scope.Filter(local)
	scopedRemote := scope.Filter(remote)

	merged, err := Merge(scopedLocal, scopedRemote, strategy)
	if err != nil {
		return nil, fmt.Errorf("scoped merge: %w", err)
	}

	// Start from full local, then overwrite with merged scoped keys.
	out := make(map[string]string, len(local))
	for k, v := range local {
		out[k] = v
	}
	for k, v := range merged {
		out[k] = v
	}
	return out, nil
}
