package btrace_test

import (
	"context"
	"go-brick/bconst"
	"go-brick/bcontext"
	"go-brick/btrace"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceID(t *testing.T) {
	traceID := btrace.GenTraceID()
	t.Log(traceID)
	assert.Equal(t, "", btrace.GetTraceID(context.Background()))

	ctx1 := context.WithValue(context.Background(), bconst.KeyTraceID, traceID)
	ctx2 := bcontext.NewWithCtx(ctx1)
	ctx3 := bcontext.New().(bcontext.Context)
	ctx3.Set(bconst.KeyTraceID, traceID)
	assert.Equal(t, traceID, btrace.GetTraceID(ctx1))
	assert.Equal(t, traceID, btrace.GetTraceID(ctx2))
	assert.Equal(t, traceID, btrace.GetTraceID(ctx3))
}
