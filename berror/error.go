package berror

import (
	"runtime"
)

const maxStackDepth = 32

// defaultError
// Provide built-in error status carrier
type defaultError struct {
	err    error // original error
	subErr Error // the error in the next layer that is wrapped

	status Status    // business information
	stack  []uintptr // stack information when this object(*defaultError) was created
}

// newWithSkip
// Create and return an error containing the stack trace.
// @offset: offset stack depth
func newWithSkip(err error, status Status, offset int) Error {
	return &defaultError{
		err:    err,
		status: status,
		stack:  callers(offset),
	}
}

// callers
// Get stack information ptr
func callers(skip ...int) []uintptr {
	var (
		pcs [maxStackDepth]uintptr
		n   = 3 // Because the call to this func has gone through 3 layers
	)
	if len(skip) > 0 {
		n += skip[0]
	}
	return pcs[:runtime.Callers(n, pcs[:])]
}
