package bmap

import "go-brick/btype"

// FromSlice Convert the incoming slice into a mapping table consistent with the data type of the slice element.
func FromSlice[T btype.Number | ~string](source []T) map[T]struct{} {
	r := make(map[T]struct{}, len(source))
	for _, v := range source {
		r[v] = struct{}{}
	}
	return r
}

// Keys Get the key slice of the mapping table.
// the returned results are unordered
func Keys[Tk btype.Number | ~string, Tv any](source map[Tk]Tv) []Tk {
	r := make([]Tk, 0, len(source))
	for k := range source {
		r = append(r, k)
	}
	return r
}

// Values Get the value slice of the mapping table.
// the returned results are unordered
func Values(source map[any]any) []any {
	r := make([]any, 0, len(source))
	for _, v := range source {
		r = append(r, v)
	}
	return r
}
