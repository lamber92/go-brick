package bslice

import (
	"sort"

	"github.com/lamber92/go-brick/btype"
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

// RemoveDuplicates remove duplicate items
func RemoveDuplicates[T comparable](src []T) []T {
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

// SortNumbers sort numeric slice
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

// descStringSlice attaches the methods of Interface to []string, sorting in increasing order.
type ascStringSlice[T btype.String] []T

func (x ascStringSlice[T]) Len() int           { return len(x) }
func (x ascStringSlice[T]) Less(i, j int) bool { return x[i] < x[j] }
func (x ascStringSlice[T]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// descStringSlice attaches the methods of Interface to []string, sorting in decreasing order.
type descStringSlice[T btype.String] []T

func (x descStringSlice[T]) Len() int           { return len(x) }
func (x descStringSlice[T]) Less(i, j int) bool { return x[i] > x[j] }
func (x descStringSlice[T]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// SortStrings sort string slice
func SortStrings[T btype.String](src []T, desc ...bool) []T {
	if len(src) == 0 {
		return src
	}
	if desc != nil && desc[0] {
		sort.Sort(descStringSlice[T](src))
	} else {
		sort.Sort(ascStringSlice[T](src))
	}
	return src
}
