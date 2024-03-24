package producer

import (
	"go-brick/btrace"
	"go-brick/internal/bufferpool"

	"go.uber.org/zap/zapcore"
)

type trace struct {
	module btrace.Module
	typ    string
	body   string
	retry  uint
	cost   int64
}

func newTraceMD(typ, body string, retry uint, cost int64) btrace.Metadata {
	return &trace{
		module: "rabbitmq-producer",
		typ:    typ,
		body:   body,
		retry:  retry,
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
	buff.AppendString("retry: ")
	buff.AppendUint(uint64(t.retry))
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
	enc.AddUint("retry", t.retry)
	enc.AddInt64("cost", t.cost)
	return nil
}
