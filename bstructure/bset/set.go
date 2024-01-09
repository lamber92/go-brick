package bset

// Clone deep copy set data
func Clone[T comparable](src map[T]struct{}) map[T]struct{} {
	r := make(map[T]struct{}, len(src))
	for k := range src {
		r[k] = struct{}{}
	}
	return r
}

// FromSlice convert the incoming slice into a mapping table consistent with the data type of the slice element.
func FromSlice[T comparable](src []T) map[T]struct{} {
	r := make(map[T]struct{}, len(src))
	for _, v := range src {
		r[v] = struct{}{}
	}
	return r
}

// ToSlice convert the incoming set's keys into a slice.
func ToSlice[T comparable](src map[T]struct{}) []T {
	r := make([]T, 0, len(src))
	for k := range src {
		r = append(r, k)
	}
	return r
}

// ToSafeSet convert the incoming set to a SafeSet
func ToSafeSet[T comparable](src map[T]struct{}) SafeSet[T] {
	return NewSafeSet[T](ToSlice(src)...)
}

// IntersectionSet return the intersection set of the source sets
// Example:
// A{1,2,3} ∩ B{2,3,4} = {2,3}
func IntersectionSet[T comparable](src ...map[T]struct{}) map[T]struct{} {
	if len(src) == 0 {
		return make(map[T]struct{})
	}
	r := Clone(src[0])
	for _, set := range src[1:] {
		temp := make(map[T]struct{})
		for k := range set {
			if _, ok := r[k]; ok {
				temp[k] = struct{}{}
			}
		}
		r = temp
	}
	return r
}

// UnionSet return the union set of the source sets
// Example:
// A{1,2,3} ∪ B{2,3,4} = {1,2,3,4}
func UnionSet[T comparable](src ...map[T]struct{}) map[T]struct{} {
	if len(src) == 0 {
		return make(map[T]struct{})
	}
	r := Clone(src[0])
	for _, set := range src[1:] {
		for k := range set {
			r[k] = struct{}{}
		}
	}
	return r
}

// ComplementSet return the relative complement set of B to A
// Example:
// B{1,2,3} \ A{2,3,4} = {1}
// B{2,3,4} \ A{1,2,3} = {4}
func ComplementSet[T comparable](A, B map[T]struct{}) map[T]struct{} {
	r := make(map[T]struct{})
	for k := range B {
		if _, ok := A[k]; !ok {
			r[k] = struct{}{}
		}
	}
	return r
}
