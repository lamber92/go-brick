// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package berrgroup provides synchronization, error propagation, and Context
// cancelation for groups of goroutines working on subtasks of a common task.
package berrgroup

import (
	"context"
	"fmt"
	"go-brick/bcontext"
	"go-brick/bpanic"
	"sync"
)

type token struct{}

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid, has no limit on the number of active goroutines,
// and does not ctx on error.
type Group struct {
	ctx bcontext.Context

	wg sync.WaitGroup

	sem chan token

	errOnce sync.Once
	err     error
}

func (g *Group) done() {
	if g.sem != nil {
		<-g.sem
	}
	g.wg.Done()
}

func (g *Group) handleError(err error) {
	g.errOnce.Do(func() {
		g.err = err
		if g.ctx != nil {
			g.ctx.Cancel()
		}
	})
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, bcontext.Context) {
	switch inner := ctx.(type) {
	case bcontext.Context:
		inner.WithCancel()
		return &Group{ctx: inner}, inner
	default:
		newCtx := bcontext.NewWithCtx(ctx)
		newCtx.WithCancel()
		return &Group{ctx: newCtx}, newCtx
	}
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.ctx != nil {
		g.ctx.Cancel()
	}
	return g.err
}

// Go calls the given function in a new goroutine.
// It blocks until the new goroutine can be added without the number of
// active goroutines in the group exceeding the configured limit.
//
// The first call to return a non-nil error cancels the group's context, if the
// group was created by calling WithContext. The error will be returned by Wait.
func (g *Group) Go(f func() error, hook ...func(error)) {
	if g.sem != nil {
		g.sem <- token{}
	}

	var handle func(error) = nil
	if len(hook) > 0 {
		handle = hook[0]
	}

	g.wg.Add(1)
	go func() {
		defer bpanic.Recover(func(err error) {
			if err != nil {
				g.handleError(err)
				handle(err)
			}
		})
		defer g.done()

		if err := f(); err != nil {
			g.handleError(err)
		}
	}()
}

// TryGo calls the given function in a new goroutine only if the number of
// active goroutines in the group is currently below the configured limit.
//
// The return value reports whether the goroutine was started.
func (g *Group) TryGo(f func() error, hook ...func(error)) bool {
	if g.sem != nil {
		select {
		case g.sem <- token{}:
			// Note: this allows barging iff channels in general allow barging.
		default:
			return false
		}
	}

	var handle func(error) = nil
	if len(hook) > 0 {
		handle = hook[0]
	}

	g.wg.Add(1)
	go func() {
		defer bpanic.Recover(func(err error) {
			if err != nil {
				g.handleError(err)
				handle(err)
			}
		})
		defer g.done()

		if err := f(); err != nil {
			g.handleError(err)
		}
	}()
	return true
}

// SetLimit limits the number of active goroutines in this group to at most n.
// A negative value indicates no limit.
//
// Any subsequent call to the Go method will block until it can add an active
// goroutine without exceeding the configured limit.
//
// The limit must not be modified while any goroutines in the group are active.
func (g *Group) SetLimit(n int) {
	if n < 0 {
		g.sem = nil
		return
	}
	if len(g.sem) != 0 {
		panic(fmt.Errorf("errgroup: modify limit while %v goroutines in the group are still active", len(g.sem)))
	}
	g.sem = make(chan token, n)
}
