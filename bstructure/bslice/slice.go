package bslice

import (
	"go-brick/btype"
	"sort"
)

// Join concat two slices
func Join[T any](first []T, second []T) []T {
	r := make([]T, 0, len(first)+len(second))
	r = append(r, first...)
	r = append(r, second...)
	return r
}

// Joins concat multiple slices
func Joins[T any](source ...[]T) []T {
	length := 0
	for _, v := range source {
		length += len(v)
	}
	r := make([]T, 0, length)
	for _, v := range source {
		r = append(r, v...)
	}
	return r
}

// RemoveDuplicates Remove duplicate items
func RemoveDuplicates[T btype.Number | ~string](source []T) []T {
	r := make([]T, 0, len(source))
	check := make(map[T]struct{})
	for _, v := range source {
		_, exist := check[v]
		if !exist {
			check[v] = struct{}{}
			r = append(r, v)
		}
	}
	return r
}

// SortNumber Sort numeric slice
func SortNumber[T btype.Number](source []T, desc ...bool) []T {
	if len(source) == 0 {
		return source
	}
	if desc != nil && desc[0] {
		sort.Slice(source, func(i, j int) bool {
			if source[i] > source[j] {
				return true
			}
			return false
		})
	} else {
		sort.Slice(source, func(i, j int) bool {
			if source[i] < source[j] {
				return true
			}
			return false
		})
	}
	return source
}

// descStringSlice attaches the methods of Interface to []string, sorting in increasing order.
type descStringSlice []string

func (x descStringSlice) Len() int           { return len(x) }
func (x descStringSlice) Less(i, j int) bool { return x[i] > x[j] }
func (x descStringSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// SortSting Sort string slice
func SortSting(source []string, desc ...bool) []string {
	if len(source) == 0 {
		return source
	}
	if desc != nil && desc[0] {
		sort.Sort(descStringSlice(source))
	} else {
		sort.Sort(sort.StringSlice(source))
	}
	return source
}
