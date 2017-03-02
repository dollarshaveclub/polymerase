package set

import (
	"testing"
)

var testStringSlice = []string{"a", "b", "c"}
var testStringSet = NewStringSet(testStringSlice)

func TestStringItems(t *testing.T) {
	items := testStringSet.Items()
	if len(items) != 3 {
		t.Fatalf("incorrect items returned")
	}
}

func TestStringAddRemove(t *testing.T) {
	testStringSet.Add("d")
	if !testStringSet.Contains("d") {
		t.Fatalf("should have contained d")
	}
	testStringSet.Remove("d")
	if testStringSet.Contains("d") {
		t.Fatalf("should have not contained d")
	}
}

func TestStringIsSubset(t *testing.T) {
	if !testStringSet.IsSubset(NewStringSet([]string{"a", "b", "c", "d"})) {
		t.Fatalf("should be subset")
	}
	if testStringSet.IsSubset(NewStringSet([]string{"x", "y", "z"})) {
		t.Fatalf("should not be subset")
	}
}

func TestStringIsSuperset(t *testing.T) {
	if !testStringSet.IsSuperset(NewStringSet([]string{"a", "b"})) {
		t.Fatalf("should be superset")
	}
	if testStringSet.IsSuperset(NewStringSet([]string{"x", "y", "z"})) {
		t.Fatalf("should not be superset")
	}
}

func TestStringDifference(t *testing.T) {
	diff := testStringSet.Difference(NewStringSet([]string{"x", "y", "z"}))
	if !diff.IsEqual(testStringSet) {
		t.Fatalf("incorrect difference")
	}
}

func TestStringIntersection(t *testing.T) {
	diff := testStringSet.Intersection(NewStringSet([]string{"a", "y", "z"}))
	if !diff.Contains("a") || len(diff.Set) != 1 {
		t.Fatalf("incorrect intersection")
	}
}

func TestStringUnion(t *testing.T) {
	add := []string{"d", "e", "f"}
	union := testStringSet.Union(NewStringSet(add))
	for _, v := range append(testStringSlice, add...) {
		if !union.Contains(v) {
			t.Fatalf("incorrect union")
		}
	}
}

func TestStringSymmetricDifference(t *testing.T) {
	otherSlice := []string{"b", "c", "y"}
	other := NewStringSet(otherSlice)
	sd := testStringSet.SymmetricDifference(other)
	for _, v := range []string{"a", "y"} {
		if !sd.Contains(v) {
			t.Fatalf("incorrectd symmetric difference")
		}
	}
}
