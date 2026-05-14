package env

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// AccessLevel represents the permission level for an environment.
type AccessLevel int

const (
	AccessNone  AccessLevel = iota
	AccessRead              // read-only
	AccessWrite             // read + write
	AccessAdmin             // read + write + manage
)

func (a AccessLevel) String() string {
	switch a {
	case AccessRead:
		return "read"
	case AccessWrite:
		return "write"
	case AccessAdmin:
		return "admin"
	default:
		return "none"
	}
}

// ParseAccessLevel converts a string to an AccessLevel.
func ParseAccessLevel(s string) (AccessLevel, error) {
	switch strings.ToLower(s) {
	case "read":
		return AccessRead, nil
	case "write":
		return AccessWrite, nil
	case "admin":
		return AccessAdmin, nil
	default:
		return AccessNone, fmt.Errorf("unknown access level: %q", s)
	}
}

// AccessPolicy maps principal identifiers to their access level for an environment.
type AccessPolicy struct {
	Environment string
	Grants      map[string]AccessLevel // principal -> level
}

// NewAccessPolicy creates an empty policy for the given environment.
func NewAccessPolicy(env string) *AccessPolicy {
	return &AccessPolicy{Environment: env, Grants: make(map[string]AccessLevel)}
}

// Grant sets the access level for a principal.
func (p *AccessPolicy) Grant(principal string, level AccessLevel) error {
	if strings.TrimSpace(principal) == "" {
		return errors.New("principal must not be empty")
	}
	p.Grants[principal] = level
	return nil
}

// Revoke removes access for a principal.
func (p *AccessPolicy) Revoke(principal string) {
	delete(p.Grants, principal)
}

// Check returns the access level for a principal (AccessNone if not found).
func (p *AccessPolicy) Check(principal string) AccessLevel {
	return p.Grants[principal]
}

// Principals returns all principals in deterministic order.
func (p *AccessPolicy) Principals() []string {
	keys := make([]string, 0, len(p.Grants))
	for k := range p.Grants {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
