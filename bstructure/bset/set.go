package bset

// FromSlice convert the incoming slice into a mapping table consistent with the data type of the slice element.
func FromSlice[T comparable](source []T) map[T]struct{} {
	r := make(map[T]struct{}, len(source))
	for _, v := range source {
		r[v] = struct{}{}
	}
	return r
}

// ToSlice convert the incoming set's keys into a slice.
func ToSlice[T comparable](source map[T]struct{}) []T {
	r := make([]T, 0, len(source))
	for k := range source {
		r = append(r, k)
	}
	return r
}

// ToSafeSet convert the incoming set to a SafeSet
func ToSafeSet[T comparable](source map[T]struct{}) SafeSet[T] {
	r := NewSafeSet[T]()
	for k := range source {
		r.Add(k)
	}
	return r
}
