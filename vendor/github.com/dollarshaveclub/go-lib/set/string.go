package set

// StringSet is a typesafe set of strings
// Wrap in a mutex for thread-safety
type StringSet struct {
	Set map[string]bool
}

func (s StringSet) Items() []string {
	items := []string{}
	for k := range s.Set {
		items = append(items, k)
	}
	return items
}

func (s *StringSet) Add(k string) {
	s.Set[k] = true
}

func (s *StringSet) Remove(k string) {
	delete(s.Set, k)
}

func (s StringSet) Contains(k string) bool {
	_, ok := s.Set[k]
	return ok
}

func (s StringSet) IsEqual(other StringSet) bool {
	for val := range s.Set {
		if !other.Contains(val) {
			return false
		}
	}
	return true
}

func (s StringSet) IsSubset(other StringSet) bool {
	for val := range s.Set {
		if !other.Contains(val) {
			return false
		}
	}
	return true
}

func (s StringSet) IsSuperset(other StringSet) bool {
	return other.IsSubset(s)
}

func (s StringSet) Intersection(other StringSet) StringSet {
	intersection := NewStringSet([]string{})
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

func (s StringSet) Difference(other StringSet) StringSet {
	diff := NewStringSet([]string{})
	for val := range s.Set {
		if !other.Contains(val) {
			diff.Add(val)
		}
	}
	return diff
}

func (s StringSet) SymmetricDifference(other StringSet) StringSet {
	diffA := s.Difference(other)
	diffB := other.Difference(s)
	return diffA.Union(diffB)
}

func (s StringSet) Union(other StringSet) StringSet {
	union := NewStringSet([]string{})
	for val := range s.Set {
		union.Add(val)
	}
	for val := range other.Set {
		union.Add(val)
	}
	return union
}

func NewStringSet(c []string) StringSet {
	return StringSet{
		Set: mapFromStringSlice(c),
	}
}

func mapFromStringSlice(s []string) map[string]bool {
	m := make(map[string]bool)
	for _, val := range s {
		m[val] = true
	}
	return m
}
