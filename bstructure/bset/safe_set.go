package bset

import (
	"go-brick/btype"
	"sync"
)

type SafeSet[T btype.Key] interface {
	// Add append elements
	Add(items ...T)
}

type defaultSafeSet[T btype.Key] struct {
	sync.RWMutex
	m map[T]struct{}
}

func NewSafeSet[T btype.Key](items ...T) SafeSet[T] {
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
