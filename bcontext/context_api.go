package bcontext

import (
	"context"
	"time"
)

// Context extension interface of context.Context
// Context's methods may be called by multiple goroutines simultaneously.
type Context interface {
	// WithTimeout override context timeout
	WithTimeout(timeout time.Duration)
	// Cancel trigger context timeout early
	Cancel()
	// GetOrigCtx get original context
	// returns 'false' if the original context does not exist
	GetOrigCtx() (context.Context, bool)
	// Set store key-value pairs
	Set(key string, value any) Context
	// Get fetch the stored value by key
	Get(key string) (value any, exists bool)

	/*
	   The following methods are consistent with the
	   official context.Context definition
	*/
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}
