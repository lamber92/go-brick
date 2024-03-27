package consumer

import (
	"fmt"
	"go-brick/bcontext"
	"go-brick/berror"
	"go-brick/blog/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	Handler        func(ctx *Context, deliveries []*amqp.Delivery, idx uint) error
	HandlerRecover func(ctx *Context, deliveries []*amqp.Delivery, idx uint, recover any) error
)

type Context struct {
	bcontext.Context

	layer      int       // the current layer of the handler chain
	hdrChain   []Handler // handlers chain
	hdrCount   int       // the length of handlers chain
	hdrRecover HandlerRecover
}

func newContext(chain []Handler, withCtx ...bcontext.Context) *Context {
	var ctx bcontext.Context
	if len(withCtx) > 0 && withCtx[0] != nil {
		ctx = withCtx[0]
	} else {
		ctx = bcontext.New()
	}
	res := &Context{
		Context:    ctx,
		hdrChain:   chain,
		hdrCount:   len(chain),
		hdrRecover: defaultHandlerRecover,
	}
	return res
}

func (c *Context) Handle(ds []*amqp.Delivery, index uint) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = c.hdrRecover(c, ds, index, r)
		}
	}()
	return c.hdrChain[0](c, ds, index)
}

// Next call to the next handler in chain order
// this method prohibits concurrent calls
func (c *Context) Next(ds []*amqp.Delivery, index uint) error {
	c.layer++
	if c.layer == c.hdrCount {
		return nil
	}
	return c.hdrChain[c.layer](c, ds, index)
}

func defaultHandlerRecover(ctx *Context, deliveries []*amqp.Delivery, idx uint, recover any) error {
	var err error
	switch tmp := recover.(type) {
	case string:
		err = berror.NewInternalError(nil, tmp)
	case error:
		err = tmp
	default:
		err = fmt.Errorf(fmt.Sprintf("%+v", recover))
	}

	var d *amqp.Delivery
	if len(deliveries) > 0 {
		d = deliveries[0]
	} else {
		d = &amqp.Delivery{}
	}
	logger.Infra.WithContext(ctx).WithError(err).
		With(logger.NewField().String("routine_key", d.RoutingKey)).
		With(logger.NewField().String("message_id", d.MessageId)).
		With(logger.NewField().String("body", string(d.Body))).
		Errorf("[rabbitmq-consumer][%d] delivery handlers recover", idx)
	return err
}
