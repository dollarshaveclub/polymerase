package set

import (
	"testing"
)

var testIntegerSlice = []int{1, 2, 3}
var testIntegerSet = NewIntegerSet(testIntegerSlice)

func TestIntegerItems(t *testing.T) {
	items := testIntegerSet.Items()
	if len(items) != 3 {
		t.Fatalf("incorrect items returned")
	}
}

func TestIntegerAddRemove(t *testing.T) {
	testIntegerSet.Add(4)
	if !testIntegerSet.Contains(4) {
		t.Fatalf("should have contained 4")
	}
	testIntegerSet.Remove(4)
	if testIntegerSet.Contains(4) {
		t.Fatalf("should have not contained 4")
	}
}

func TestIntegerIsSubset(t *testing.T) {
	if !testIntegerSet.IsSubset(NewIntegerSet([]int{1, 2, 3, 4})) {
		t.Fatalf("should be subset")
	}
	if testIntegerSet.IsSubset(NewIntegerSet([]int{9, 10, 11})) {
		t.Fatalf("should not be subset")
	}
}

func TestIntegerIsSuperset(t *testing.T) {
	if !testIntegerSet.IsSuperset(NewIntegerSet([]int{1, 2})) {
		t.Fatalf("should be superset")
	}
	if testIntegerSet.IsSuperset(NewIntegerSet([]int{9, 10, 11})) {
		t.Fatalf("should not be superset")
	}
}

func TestIntegerDifference(t *testing.T) {
	diff := testIntegerSet.Difference(NewIntegerSet([]int{9, 10, 11}))
	if !diff.IsEqual(testIntegerSet) {
		t.Fatalf("incorrect difference")
	}
}

func TestIntegerIntersection(t *testing.T) {
	diff := testIntegerSet.Intersection(NewIntegerSet([]int{1, 9, 10}))
	if !diff.Contains(1) || len(diff.Set) != 1 {
		t.Fatalf("incorrect intersection")
	}
}

func TestIntegerUnion(t *testing.T) {
	add := []int{4, 5, 6}
	union := testIntegerSet.Union(NewIntegerSet(add))
	for _, v := range append(testIntegerSlice, add...) {
		if !union.Contains(v) {
			t.Fatalf("incorrect union")
		}
	}
}

func TestIntegerSymmetricDifference(t *testing.T) {
	otherSlice := []int{2, 3, 9}
	other := NewIntegerSet(otherSlice)
	sd := testIntegerSet.SymmetricDifference(other)
	for _, v := range []int{1, 9} {
		if !sd.Contains(v) {
			t.Fatalf("incorrectd symmetric difference")
		}
	}
}
