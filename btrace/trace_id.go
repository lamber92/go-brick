package btrace

import (
	"context"
	"encoding/hex"
	"go-brick/bcontext"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
)

const (
	KeyTraceID = "b_trace_id"
)

type TraceIDGenerator interface {
	// GenTraceID generate a Stack-ID
	GenTraceID() string
}

var traceIDGen = newUUIDV4Generator()

// ReplaceTraceIDGenerator overrides the default Stack-ID generator
func ReplaceTraceIDGenerator(gen TraceIDGenerator) {
	traceIDGen = gen
}

// uuidV4Generator uuid v4 version generator,
// used as Stack-ID generator
type uuidV4Generator struct{}

func newUUIDV4Generator() TraceIDGenerator {
	return &uuidV4Generator{}
}

// GenTraceID generate a Stack-ID
func (*uuidV4Generator) GenTraceID() string {
	return hex.EncodeToString(uuid.NewV4().Bytes())
}

// GenTraceID generate a Stack-ID
func GenTraceID() string {
	return traceIDGen.GenTraceID()
}

// SetTraceID set Stack-ID into context.
// if traceIDs is not passed in, a built-in new ID will be used.
func SetTraceID(ctx context.Context, traceIDs ...string) context.Context {
	var traceID string
	if len(traceIDs) > 0 {
		traceID = traceIDs[0]
	} else {
		traceID = traceIDGen.GenTraceID()
	}

	switch tmp := ctx.(type) {
	case bcontext.Context:
		return tmp.Set(KeyTraceID, traceID)
	default:
		return context.WithValue(ctx, KeyTraceID, traceID)
	}
}

// GetTraceID get Stack-ID from context
func GetTraceID(ctx context.Context) string {
	tmp := ctx.Value(KeyTraceID)
	traceID, ok := tmp.(string)
	if ok {
		return traceID
	}
	return cast.ToString(traceID)
}
