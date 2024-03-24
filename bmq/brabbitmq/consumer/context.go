package consumer

import (
	"go-brick/bcontext"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler func(*Context, []*amqp.Delivery, uint) error

type Context struct {
	bcontext.Context

	layer    int       // the current layer of the handler chain
	hdrChain []Handler // handlers chain
	hdrCount int       // the length of handlers chain
}

func NewContext(chain []Handler, withCtx ...bcontext.Context) *Context {
	res := &Context{
		hdrChain: chain,
		hdrCount: len(chain),
	}
	if len(withCtx) > 0 && withCtx[0] != nil {
		res.Context = withCtx[0]
	} else {
		res.Context = bcontext.New()
	}
	return res
}

func (c *Context) Handle(d []*amqp.Delivery, index uint) error {
	return c.hdrChain[0](c, d, index)
}

// Next call to the next handler in chain order
// this method prohibits concurrent calls
func (c *Context) Next(d []*amqp.Delivery, index uint) error {
	c.layer++
	if c.layer == c.hdrCount {
		return nil
	}
	return c.hdrChain[c.layer](c, d, index)
}
