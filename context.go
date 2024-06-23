package rmq

import (
	"context"
	"github.com/cscoder0/go-rabbitmq/binding"
	"github.com/cscoder0/go-rabbitmq/mq"
	"time"
)

type HandlerFunc func(ctx *Context) error
type HandlersChain []HandlerFunc

type Context struct {
	msg      *mq.Message
	Ctx      context.Context
	Handlers HandlersChain
}

func (c *Context) ShouldBind(obj any) error {
	b := binding.Default(c.msg.ContentType)
	return b.BindBody(c.msg.Body, obj)
}

func (c *Context) Ack() error {
	return c.msg.Ack()
}
func (c *Context) Nack() error {
	return c.msg.Nack()
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.Ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.Ctx.Done()
}

func (c *Context) Err() error {
	return c.Ctx.Err()
}

func (c *Context) Value(key any) any {
	return c.Ctx.Value(key)
}
