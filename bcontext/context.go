package bcontext

import (
	"context"
	"sync"
	"time"
)

type defaultContext struct {
	orig context.Context
	kv   map[string]any

	timer  context.Context
	cancel context.CancelFunc

	sync.RWMutex
}

func New() Context {
	return &defaultContext{
		orig: nil,
		kv:   make(map[string]any),
	}
}

func NewWithCtx(ctx context.Context) Context {
	out := &defaultContext{
		orig: ctx,
		kv:   make(map[string]any),
	}
	if deadline, exist := ctx.Deadline(); exist {
		out.timer, out.cancel = context.WithDeadline(out, deadline)
	}
	return out
}

// WithTimeout override context timeout
func (ctx *defaultContext) WithTimeout(timeout time.Duration) {
	ctx.Lock()
	defer ctx.Unlock()
	origTimer, origCancel := ctx.timer, ctx.cancel
	// reset timeout
	ctx.timer, ctx.cancel = context.WithTimeout(context.Background(), timeout)
	if origTimer.Err() == nil {
		// release original timer
		origCancel()
	}
	return
}

// Cancel trigger context timeout early
func (ctx *defaultContext) Cancel() {
	if ctx.cancel != nil {
		ctx.cancel()
	}
}

// GetOrigCtx get original context
// returns 'false' if the original context does not exist
func (ctx *defaultContext) GetOrigCtx() (context.Context, bool) {
	if ctx.orig != nil {
		return ctx.orig, true
	}
	return nil, false
}

func (ctx *defaultContext) Set(key string, value any) {
	ctx.Lock()
	ctx.kv[key] = value
	ctx.Unlock()
}

func (ctx *defaultContext) Get(key string) (value any, exists bool) {
	ctx.RLock()
	value, exists = ctx.kv[key]
	ctx.RUnlock()
	return
}

func (ctx *defaultContext) Deadline() (deadline time.Time, ok bool) {
	ctx.RLock()
	defer ctx.RUnlock()
	if ctx.timer != nil {
		deadline, ok = ctx.timer.Deadline()
		return
	}
	return
}

func (ctx *defaultContext) Done() <-chan struct{} {
	ctx.RLock()
	defer ctx.RUnlock()
	if ctx.timer != nil {
		return ctx.timer.Done()
	}
	return nil
}

func (ctx *defaultContext) Err() error {
	ctx.RLock()
	defer ctx.RUnlock()
	if ctx.timer != nil {
		return ctx.timer.Err()
	}
	return context.Canceled
}

func (ctx *defaultContext) Value(key any) any {
	ctx.RLock()
	defer ctx.RUnlock()
	if strKey, ok := key.(string); ok {
		if val, exist := ctx.Get(strKey); exist {
			return val
		}
	}
	if ctx.orig != nil {
		return ctx.orig.Value(key)
	}
	return nil
}
