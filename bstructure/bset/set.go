package bset

import "go-brick/btype"

// GetFromSlice convert the incoming slice into a mapping table consistent with the data type of the slice element.
func GetFromSlice[T btype.Number | ~string](source []T) map[T]struct{} {
	r := make(map[T]struct{}, len(source))
	for _, v := range source {
		r[v] = struct{}{}
	}
	return r
}
