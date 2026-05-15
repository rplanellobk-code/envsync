package env

import (
	"strings"
	"testing"
)

func TestCompareIdentical(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}

	r := Compare(a, b, "staging", "production")

	if !r.IsIdentical() {
		t.Fatal("expected identical environments")
	}
	if len(r.Common) != 2 {
		t.Fatalf("expected 2 common keys, got %d", len(r.Common))
	}
}

func TestCompareOnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "yes"}
	b := map[string]string{"FOO": "bar"}

	r := Compare(a, b, "staging", "production")

	if len(r.OnlyInA) != 1 {
		t.Fatalf("expected 1 key only in A, got %d", len(r.OnlyInA))
	}
	if r.OnlyInA["ONLY_A"] != "yes" {
		t.Fatalf("unexpected value for ONLY_A: %q", r.OnlyInA["ONLY_A"])
	}
	if r.IsIdentical() {
		t.Fatal("expected non-identical")
	}
}

func TestCompareOnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "yes"}

	r := Compare(a, b, "staging", "production")

	if len(r.OnlyInB) != 1 {
		t.Fatalf("expected 1 key only in B, got %d", len(r.OnlyInB))
	}
}

func TestCompareDifferentValues(t *testing.T) {
	a := map[string]string{"FOO": "alpha", "BAR": "same"}
	b := map[string]string{"FOO": "beta", "BAR": "same"}

	r := Compare(a, b, "staging", "production")

	if len(r.Different) != 1 {
		t.Fatalf("expected 1 differing key, got %d", len(r.Different))
	}
	pair, ok := r.Different["FOO"]
	if !ok {
		t.Fatal("expected FOO in Different")
	}
	if pair[0] != "alpha" || pair[1] != "beta" {
		t.Fatalf("unexpected pair: %v", pair)
	}
	keys := r.DifferentKeys()
	if len(keys) != 1 || keys[0] != "FOO" {
		t.Fatalf("unexpected DifferentKeys: %v", keys)
	}
}

func TestCompareSummaryContainsNames(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}

	r := Compare(a, b, "staging", "production")
	summary := r.Summary()

	if !strings.Contains(summary, "staging") {
		t.Error("summary missing environment A name")
	}
	if !strings.Contains(summary, "production") {
		t.Error("summary missing environment B name")
	}
}

func TestCompareBothEmpty(t *testing.T) {
	r := Compare(map[string]string{}, map[string]string{}, "a", "b")
	if !r.IsIdentical() {
		t.Fatal("two empty envs should be identical")
	}
}
