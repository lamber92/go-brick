package producer

import (
	"context"
	"go-brick/btrace"
	"go-brick/internal/bufferpool"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap/zapcore"
)

type TraceFunc func(ctx context.Context, data *amqp.Publishing, since time.Duration)

func defaultTraceFunc(ctx context.Context, data *amqp.Publishing, since time.Duration) {
	btrace.AppendMDIntoCtx(ctx, newTraceMD(data.Type, string(data.Body), since.Milliseconds()))
}

type trace struct {
	module btrace.Module
	typ    string
	body   string
	cost   int64
}

func newTraceMD(typ, body string, cost int64) btrace.Metadata {
	return &trace{
		module: "rabbitmq-producer",
		typ:    typ,
		body:   body,
		cost:   cost,
	}
}

func (t *trace) Module() btrace.Module {
	return t.module
}

func (t *trace) String() string {
	buff := bufferpool.Get()
	buff.AppendString("module: ")
	buff.AppendString(string(t.module))
	buff.AppendByte(',')
	buff.AppendByte(' ')
	buff.AppendString("type: ")
	buff.AppendString(t.typ)
	buff.AppendByte(',')
	buff.AppendByte(' ')
	buff.AppendString("body: ")
	buff.AppendString(t.body)
	buff.AppendByte(',')
	buff.AppendByte(' ')
	buff.AppendString("cost: ")
	buff.AppendInt(t.cost)
	out := buff.String()
	buff.Free()
	return out
}

func (t *trace) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if t == nil {
		return nil
	}
	enc.AddString("module", string(t.module))
	enc.AddString("type", t.typ)
	enc.AddString("body", t.body)
	enc.AddInt64("cost", t.cost)
	return nil
}
