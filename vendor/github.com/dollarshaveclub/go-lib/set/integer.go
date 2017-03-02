package set

// IntegerSet is a typesafe set of ints
// Wrap in a mutex for thread-safety
type IntegerSet struct {
	Set map[int]bool
}

func (s IntegerSet) Items() []int {
	items := []int{}
	for k := range s.Set {
		items = append(items, k)
	}
	return items
}

func (s *IntegerSet) Add(k int) {
	s.Set[k] = true
}

func (s *IntegerSet) Remove(k int) {
	delete(s.Set, k)
}

func (s IntegerSet) Contains(k int) bool {
	_, ok := s.Set[k]
	return ok
}

func (s IntegerSet) IsEqual(other IntegerSet) bool {
	for val := range s.Set {
		if !other.Contains(val) {
			return false
		}
	}
	return true
}

func (s IntegerSet) IsSubset(other IntegerSet) bool {
	for val := range s.Set {
		if !other.Contains(val) {
			return false
		}
	}
	return true
}

func (s IntegerSet) IsSuperset(other IntegerSet) bool {
	return other.IsSubset(s)
}

func (s IntegerSet) Intersection(other IntegerSet) IntegerSet {
	intersection := NewIntegerSet([]int{})
	if len(other.Set) > len(s.Set) {
		for val := range s.Set {
			if other.Contains(val) {
				intersection.Add(val)
			}
		}
	} else {
		for val := range other.Set {
			if s.Contains(val) {
				intersection.Add(val)
			}
		}
	}
	return intersection
}

func (s IntegerSet) Difference(other IntegerSet) IntegerSet {
	diff := NewIntegerSet([]int{})
	for val := range s.Set {
		if !other.Contains(val) {
			diff.Add(val)
		}
	}
	return diff
}

func (s IntegerSet) SymmetricDifference(other IntegerSet) IntegerSet {
	diffA := s.Difference(other)
	diffB := other.Difference(s)
	return diffA.Union(diffB)
}

func (s IntegerSet) Union(other IntegerSet) IntegerSet {
	union := NewIntegerSet([]int{})
	for val := range s.Set {
		union.Add(val)
	}
	for val := range other.Set {
		union.Add(val)
	}
	return union
}

func NewIntegerSet(c []int) IntegerSet {
	return IntegerSet{
		Set: mapFromIntegerSlice(c),
	}
}

func mapFromIntegerSlice(s []int) map[int]bool {
	m := make(map[int]bool)
	for _, val := range s {
		m[val] = true
	}
	return m
}
