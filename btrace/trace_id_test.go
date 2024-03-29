package btrace_test

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/lamber92/go-brick/bcontext"
	"github.com/lamber92/go-brick/btrace"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestTraceID(t *testing.T) {
	traceID := btrace.GenTraceID()
	t.Logf("current tracing-id is : %s", traceID)
	assert.Equal(t, "", btrace.GetTraceID(context.Background()))

	ctx1 := context.WithValue(context.Background(), btrace.KeyTraceID, traceID)
	ctx2 := bcontext.NewWithCtx(ctx1)
	ctx3 := bcontext.New().(bcontext.Context).Set(btrace.KeyTraceID, traceID)
	assert.Equal(t, traceID, btrace.GetTraceID(ctx1))
	assert.Equal(t, traceID, btrace.GetTraceID(ctx2))
	assert.Equal(t, traceID, btrace.GetTraceID(ctx3))

	ctx4 := btrace.SetTraceID(ctx3)
	assert.NotEqual(t, traceID, btrace.GetTraceID(ctx4))
	ctx5 := btrace.SetTraceID(ctx4, traceID)
	assert.Equal(t, traceID, btrace.GetTraceID(ctx5))
}

type myGenerator struct{}

// GenTraceID generate a Stack-ID
func (*myGenerator) GenTraceID() string {
	return "test_" + hex.EncodeToString(uuid.NewV1().Bytes())
}

func TestReplaceTraceID(t *testing.T) {
	btrace.ReplaceTraceIDGenerator(&myGenerator{})
	// test generate
	traceID := btrace.GenTraceID()
	assert.Condition(t, func() bool {
		if strings.Index(traceID, "test_") == 0 {
			return true
		}
		return false
	})
	// test set/get
	ctx := btrace.SetTraceID(bcontext.New(), traceID)
	assert.Equal(t, traceID, btrace.GetTraceID(ctx))
}
