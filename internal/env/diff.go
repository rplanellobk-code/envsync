package env

import "sort"

// ChangeType represents the type of change between two env maps.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
)

// Change represents a single difference between two env maps.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Diff computes the differences between a local and remote env map.
// It returns a slice of Change entries describing keys that were added,
// removed, or updated going from base to target.
func Diff(base, target map[string]string) []Change {
	var changes []Change

	// Find added and updated keys
	for key, targetVal := range target {
		baseVal, exists := base[key]
		if !exists {
			changes = append(changes, Change{
				Key:      key,
				Type:     ChangeAdded,
				NewValue: targetVal,
			})
		} else if baseVal != targetVal {
			changes = append(changes, Change{
				Key:      key,
				Type:     ChangeUpdated,
				OldValue: baseVal,
				NewValue: targetVal,
			})
		}
	}

	// Find removed keys
	for key, baseVal := range base {
		if _, exists := target[key]; !exists {
			changes = append(changes, Change{
				Key:      key,
				Type:     ChangeRemoved,
				OldValue: baseVal,
			})
		}
	}

	// Sort for deterministic output
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Key != changes[j].Key {
			return changes[i].Key < changes[j].Key
		}
		return changes[i].Type < changes[j].Type
	})

	return changes
}

// Apply merges changes into a base env map, returning a new map.
// ChangeAdded and ChangeUpdated set the key; ChangeRemoved deletes it.
func Apply(base map[string]string, changes []Change) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}
	for _, c := range changes {
		switch c.Type {
		case ChangeAdded, ChangeUpdated:
			result[c.Key] = c.NewValue
		case ChangeRemoved:
			delete(result, c.Key)
		}
	}
	return result
}
