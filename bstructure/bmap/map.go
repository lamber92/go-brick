package bmap

import "go-brick/btype"

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
func Values[Tk btype.Number | ~string, Tv any](source map[Tk]Tv) []Tv {
	r := make([]Tv, 0, len(source))
	for _, v := range source {
		r = append(r, v)
	}
	return r
}
