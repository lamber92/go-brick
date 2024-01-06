package bmap

// Keys Get the key slice of the mapping table.
// the returned results are unordered
func Keys[K comparable, V any](source map[K]V) []K {
	r := make([]K, 0, len(source))
	for k := range source {
		r = append(r, k)
	}
	return r
}

// Values Get the value slice of the mapping table.
// the returned results are unordered
func Values[K comparable, V any](source map[K]V) []V {
	r := make([]V, 0, len(source))
	for _, v := range source {
		r = append(r, v)
	}
	return r
}
