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
func Joins[T any](src ...[]T) []T {
	length := 0
	for _, v := range src {
		length += len(v)
	}
	r := make([]T, 0, length)
	for _, v := range src {
		r = append(r, v...)
	}
	return r
}

// Combine the contents of two-dimensional slices, reduce them to one-dimensional slices and return
func Combine[T any](src [][]T) []T {
	length := 0
	for _, v := range src {
		length += len(v)
	}
	r := make([]T, 0, length)
	for _, v := range src {
		r = append(r, v...)
	}
	return r
}

// RemoveDuplicates Remove duplicate items
func RemoveDuplicates[T btype.Number | ~string](src []T) []T {
	r := make([]T, 0, len(src))
	check := make(map[T]struct{})
	for _, v := range src {
		_, exist := check[v]
		if !exist {
			check[v] = struct{}{}
			r = append(r, v)
		}
	}
	return r
}

// SortNumbers Sort numeric slice
func SortNumbers[T btype.Number](src []T, desc ...bool) []T {
	if len(src) == 0 {
		return src
	}
	if desc != nil && desc[0] {
		sort.Slice(src, func(i, j int) bool {
			if src[i] > src[j] {
				return true
			}
			return false
		})
	} else {
		sort.Slice(src, func(i, j int) bool {
			if src[i] < src[j] {
				return true
			}
			return false
		})
	}
	return src
}

// descStringSlice attaches the methods of Interface to []string, sorting in decreasing order.
type descStringSlice []string

func (x descStringSlice) Len() int           { return len(x) }
func (x descStringSlice) Less(i, j int) bool { return x[i] > x[j] }
func (x descStringSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// SortStings Sort string slice
func SortStings(src []string, desc ...bool) []string {
	if len(src) == 0 {
		return src
	}
	if desc != nil && desc[0] {
		sort.Sort(descStringSlice(src))
	} else {
		sort.Sort(sort.StringSlice(src))
	}
	return src
}
