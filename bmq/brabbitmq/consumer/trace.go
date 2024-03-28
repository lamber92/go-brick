package consumer

import (
	"go-brick/btrace"
	"go-brick/internal/bufferpool"
	"go-brick/internal/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap/zapcore"
)

type TraceFunc func(ctx *Context, err error, ds []*amqp.Delivery, since time.Duration, idx uint)

func defaultTraceFunc(ctx *Context, err error, ds []*amqp.Delivery, since time.Duration, idx uint) {
	btrace.AppendMDIntoCtx(ctx, newTraceMD(err, ds, since.Milliseconds(), idx))
}

func newTraceMD(err error, ds []*amqp.Delivery, cost int64, idx uint) btrace.Metadata {
	res := &trace{
		module:     "rabbitmq-consumer",
		cost:       cost,
		err:        err,
		consumerId: idx,
	}
	res.parseDeliveryList(ds)
	return res
}

type trace struct {
	module     btrace.Module
	messages   traceItemList
	consumerId uint
	cost       int64
	err        error
}

func (t *trace) parseDeliveryList(ds []*amqp.Delivery) {
	t.messages = make(traceItemList, 0, len(ds))
	for _, d := range ds {
		t.messages = append(t.messages, traceItem{
			MessageId:   d.MessageId,
			Timestamp:   d.Timestamp,
			Type:        d.Type,
			AppId:       d.AppId,
			ConsumerTag: d.ConsumerTag,
			RoutingKey:  d.RoutingKey,
			Body:        d.Body,
		})
	}
}

func (t *trace) Module() btrace.Module {
	return t.module
}

func (t *trace) String() string {
	buff := bufferpool.Get()
	buff.AppendString("module: ")
	buff.AppendString(string(t.module))
	buff.AppendString(" | messages: ")
	tmp, _ := json.MarshalToString(t.messages)
	buff.AppendString(tmp)
	buff.AppendString(" | consumer_id: ")
	buff.AppendUint(uint64(t.consumerId))
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
	_ = enc.AddArray("messages", t.messages)
	enc.AddUint("consumer_id", t.consumerId)
	enc.AddInt64("cost", t.cost)
	if t.err == nil {
		enc.AddString("err", "")
	} else {
		enc.AddString("err", t.err.Error())
	}
	return nil
}

type traceItemList []traceItem

func (l traceItemList) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, v := range l {
		if err := enc.AppendObject(v); err != nil {
			return err
		}
	}
	return nil
}

type traceItem struct {
	MessageId string    // application use - message identifier
	Timestamp time.Time // application use - message timestamp
	Type      string    // application use - message type name
	AppId     string    // application use - creating application id
	// Valid only with Channel.Consume
	ConsumerTag string
	RoutingKey  string // basic.publish routing key
	Body        []byte
}

func (i traceItem) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("message_id", i.MessageId)
	enc.AddInt64("timestamp", i.Timestamp.UnixMilli())
	enc.AddString("type", i.Type)
	enc.AddString("app_id", i.AppId)
	enc.AddString("consumer_tag", i.ConsumerTag)
	enc.AddString("routing_key", i.RoutingKey)
	enc.AddString("body", string(i.Body))
	return nil
}
