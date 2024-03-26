package producer

import (
	"context"
	"go-brick/btrace"
	"go-brick/internal/bufferpool"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap/zapcore"
)

type TraceFunc func(ctx context.Context, err error, data *amqp.Publishing, since time.Duration)

func defaultTraceFunc(ctx context.Context, err error, data *amqp.Publishing, since time.Duration) {
	btrace.AppendMDIntoCtx(ctx, newTraceMD(data.Type, err, string(data.Body), since.Milliseconds()))
}

type trace struct {
	module btrace.Module
	typ    string
	body   string
	cost   int64
	err    error
}

func newTraceMD(typ string, err error, body string, cost int64) btrace.Metadata {
	return &trace{
		module: "rabbitmq-producer",
		typ:    typ,
		body:   body,
		cost:   cost,
		err:    err,
	}
}

func (t *trace) Module() btrace.Module {
	return t.module
}

func (t *trace) String() string {
	buff := bufferpool.Get()
	buff.AppendString("module: ")
	buff.AppendString(string(t.module))
	buff.AppendString(" | type: ")
	buff.AppendString(t.typ)
	buff.AppendString(" | body: ")
	buff.AppendString(t.body)
	buff.AppendString(" | cost: ")
	buff.AppendInt(t.cost)
	if t.err != nil {
		buff.AppendString(" | err: ")
		buff.AppendString(t.err.Error())
	}
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
	if t.err == nil {
		enc.AddString("err", "")
	} else {
		enc.AddString("err", t.err.Error())
	}
	return nil
}
