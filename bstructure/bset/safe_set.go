package bset

import (
	"sync"
)

// SafeSet Concurrency-safe Set
type SafeSet[T comparable] interface {
	// Add append elements into SafeSet
	Add(items ...T)
	// Delete remove elements from SafeSet
	Delete(items ...T)
	// Has determine whether an element exists in the SafeSet
	Has(item T) bool
	// Contains determine whether certain elements exist in the SafeSet
	// - if all elements exist or only some elements exist, return a new SafeSet with hit elements and a `true` flag.
	// - else return the `empty` SafeSet and a `false` flag.
	Contains(item ...T) (SafeSet[T], bool)
	// Clear remove all items
	Clear()
	// Len get size of SafeSet
	Len() int
	// IsEmpty check the size of SafeSet is 0 or not
	IsEmpty() bool
	// ToSlice convert the SafeSet-items into a new slice than return
	// sorting capabilities are not provided here,
	// because the current version of the generic type `comparable` is not compatible with the `<` and `>` operation.
	// if there is a sorting requirement, use the bslice.SortXYZ() method.
	ToSlice() []T
	// Clone deep copy a new SafeSet
	Clone() SafeSet[T]
	// IntersectionSet return the intersection set of the source sets
	// Example:
	// A{1,2,3} ∩ B{2,3,4} = {2,3}
	IntersectionSet(src ...SafeSet[T]) SafeSet[T]
	// UnionSet return the union set of the source sets
	// Example:
	// A{1,2,3} ∪ B{2,3,4} = {1,2,3,4}
	UnionSet(src ...SafeSet[T]) SafeSet[T]
	// ComplementSet return the relative complement set of B to A
	// Example:
	// B{1,2,3} \ A{2,3,4} = {1}
	// B{2,3,4} \ A{1,2,3} = {4}
	ComplementSet(B SafeSet[T]) SafeSet[T]
}

type defaultSafeSet[T comparable] struct {
	sync.RWMutex
	m map[T]struct{}
}

func NewSafeSet[T comparable](items ...T) SafeSet[T] {
	s := &defaultSafeSet[T]{
		m: make(map[T]struct{}, len(items)),
	}
	s.Add(items...)
	return s
}

func (s *defaultSafeSet[T]) Add(items ...T) {
	if items == nil {
		return
	}
	s.Lock()
	for _, v := range items {
		s.m[v] = struct{}{}
	}
	s.Unlock()
}

func (s *defaultSafeSet[T]) Delete(items ...T) {
	if len(items) == 0 {
		return
	}
	s.Lock()
	for _, v := range items {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *defaultSafeSet[T]) Has(item T) bool {
	s.RLock()
	_, ok := s.m[item]
	s.RUnlock()
	return ok
}

func (s *defaultSafeSet[T]) Contains(items ...T) (SafeSet[T], bool) {
	if len(items) == 0 {
		return nil, false
	}
	s.Lock()
	m := make(map[T]struct{})
	for _, v := range items {
		_, ok := s.m[v]
		if ok {
			m[v] = struct{}{}
		}
	}
	s.Unlock()
	if len(m) > 0 {
		return &defaultSafeSet[T]{m: m}, true
	}
	return &defaultSafeSet[T]{m: m}, false
}

func (s *defaultSafeSet[T]) Clear() {
	s.Lock()
	s.m = make(map[T]struct{})
	s.Unlock()
}

func (s *defaultSafeSet[T]) Len() int {
	s.RLock()
	length := len(s.m)
	s.RUnlock()
	return length
}

func (s *defaultSafeSet[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *defaultSafeSet[T]) ToSlice() []T {
	r := make([]T, 0, len(s.m))
	s.RLock()
	for k := range s.m {
		r = append(r, k)
	}
	s.RUnlock()
	return r
}

func (s *defaultSafeSet[T]) Clone() SafeSet[T] {
	return NewSafeSet(s.ToSlice()...)
}

func (s *defaultSafeSet[T]) IntersectionSet(src ...SafeSet[T]) SafeSet[T] {
	if len(src) == 0 {
		return s
	}
	r := s.Clone()
	for _, set := range src {
		sl := make([]T, 0, set.Len())
		for _, k := range set.ToSlice() {
			if r.Has(k) {
				sl = append(sl, k)
			}
		}
		r = NewSafeSet(sl...)
	}
	return r
}

func (s *defaultSafeSet[T]) UnionSet(src ...SafeSet[T]) SafeSet[T] {
	if len(src) == 0 {
		return NewSafeSet[T]()
	}
	r := s.Clone()
	for _, set := range src {
		sl := make([]T, 0, set.Len())
		for _, k := range set.ToSlice() {
			if !r.Has(k) {
				sl = append(sl, k)
			}
		}
		r.Add(sl...)
	}
	return r
}

func (s *defaultSafeSet[T]) ComplementSet(B SafeSet[T]) SafeSet[T] {
	if B == nil {
		return s
	}
	sl := make([]T, 0, B.Len())
	for _, k := range B.ToSlice() {
		if !s.Has(k) {
			sl = append(sl, k)
		}
	}
	return NewSafeSet(sl...)
}
