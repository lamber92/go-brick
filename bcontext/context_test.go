package bcontext_test

import (
	"context"
	"testing"
	"time"

	"github.com/lamber92/go-brick/bcontext"
	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	ctx := bcontext.New()

	key, value := "xxxx", "yyyy"

	v1, ok1 := ctx.Get(key)
	assert.Equal(t, false, ok1)
	assert.Equal(t, nil, v1)

	ctx.Set(key, value)
	v2, ok2 := ctx.Get(key)
	assert.Equal(t, true, ok2)
	assert.Equal(t, value, v2)
}

func TestNewCtxTimeout(t *testing.T) {
	ctx := bcontext.New()
	err1 := ctx.Err()
	assert.Equal(t, context.Canceled, err1)

	// test deadline time
	time1, ok1 := ctx.Deadline()
	assert.Equal(t, time.Time{}, time1)
	assert.Equal(t, false, ok1)

	begin := time.Now()
	sec := time.Second * 5

	ctx.WithTimeout(sec)
	time2, ok2 := ctx.Deadline()

	compareFunc2 := func() bool {
		time.Sleep(time.Second)
		return time2.Before(time.Now().Add(sec))
	}
	assert.Condition(t, compareFunc2)
	assert.Equal(t, true, ok2)

	// test timeout sec
	d3 := <-ctx.Done()
	compareFunc3 := func() bool {
		return time.Since(begin)/time.Second == 5
	}
	assert.Equal(t, struct{}{}, d3)
	assert.Condition(t, compareFunc3)

	// test done() again
	d4 := <-ctx.Done()
	assert.Condition(t, compareFunc3)
	assert.Equal(t, struct{}{}, d4)

	// test cancel
	begin = time.Now()
	ctx.WithTimeout(sec)
	go func() {
		time.Sleep(time.Second * 3)
		ctx.Cancel()
		ctx.Cancel() // duplicate call cancel()
	}()
	d5 := <-ctx.Done()
	compareFunc5 := func() bool {
		return time.Since(begin)/time.Second == 3
	}
	assert.Equal(t, struct{}{}, d5)
	assert.Condition(t, compareFunc5)
}

func TestNewWithCtxTimeout(t *testing.T) {
	begin := time.Now()
	sec := time.Second * 5

	orig := context.Background()
	orig, _ = context.WithTimeout(orig, sec)

	// test Err()
	ctx := bcontext.NewWithCtx(orig)
	err1 := ctx.Err()
	assert.Equal(t, nil, err1)

	// test deadline time
	time1, ok1 := ctx.Deadline()
	compareFunc1 := func() bool {
		time.Sleep(time.Second)
		return time1.Before(time.Now().Add(sec))
	}
	assert.Condition(t, compareFunc1)
	assert.Equal(t, true, ok1)

	// test timeout sec
	d3 := <-ctx.Done()
	compareFunc3 := func() bool {
		return time.Since(begin)/time.Second == 5
	}
	assert.Equal(t, struct{}{}, d3)
	assert.Condition(t, compareFunc3)

	// test done() again
	d4 := <-ctx.Done()
	assert.Condition(t, compareFunc3)
	assert.Equal(t, struct{}{}, d4)

	// test cancel
	begin = time.Now()
	ctx.WithTimeout(sec) // overwrite time out
	go func() {
		time.Sleep(time.Second * 3)
		ctx.Cancel()
		ctx.Cancel() // duplicate call cancel()
	}()
	d5 := <-ctx.Done()
	compareFunc5 := func() bool {
		return time.Since(begin)/time.Second == 3
	}
	assert.Equal(t, struct{}{}, d5)
	assert.Condition(t, compareFunc5)
}

func TestWithCancel(t *testing.T) {
	ctx := context.Background()
	ctx2, cancel2 := context.WithCancel(ctx)
	ctx3 := bcontext.NewWithCtx(ctx2)
	ctx3.WithCancel()

	time.Sleep(time.Second)
	cancel2()

	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

	go func() {
		t.Logf("[%s] cancel ctx after 3 seconds", time.Now().Format(time.RFC3339))
		time.Sleep(time.Second * 5)
		ctx3.Cancel()
	}()

	select {
	case <-timer.C:
		t.Logf("[%s] timeout", time.Now().Format(time.RFC3339))
	case <-ctx3.Done():
		t.Logf("[%s] bcontext cancel done", time.Now().Format(time.RFC3339))
	}
}

func TestCtxValue(t *testing.T) {
	ctx := context.Background()

	key, value := 123, "123"

	ctx = context.WithValue(ctx, key, value)
	ctx2 := bcontext.NewWithCtx(ctx)
	assert.Equal(t, value, ctx2.Value(key))

	// test original ctx
	ctx3, ok := ctx2.GetOrigCtx()
	assert.Equal(t, true, ok)
	assert.Equal(t, ctx, ctx3)
	assert.Equal(t, value, ctx3.Value(key))
}
