package rmq

import "C"
import (
	"context"
	"github.com/chsendev/go-rabbitmq/binding"
	"github.com/chsendev/go-rabbitmq/log"
	"github.com/chsendev/go-rabbitmq/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"math"
	"time"
)

type HandlerFunc func(ctx *Context) error
type HandlersChain []HandlerFunc

var abortIndex = math.MaxInt16

type Context struct {
	msg         *mq.Message
	ctx         context.Context
	handlers    HandlersChain
	index       int
	hasAck      bool
	handlerErrs []error
}

func (c *Context) GetHeader(key string) any {
	return c.msg.Headers[key]
}

func (c *Context) GetMessageId() string {
	return c.msg.MessageId
}

func (c *Context) GetRawMessage() *amqp.Delivery {
	return c.msg.Delivery
}

func (c *Context) ShouldBind(obj any) error {
	b := binding.Default(c.msg.ContentType)
	err := b.BindBody(c.msg.Body, obj)
	// 反序列化失败，直接ack掉消息，避免死循环
	if err != nil {
		log.Error("Due to message deserialization failure, it has been automatically acked", zap.Any("error", err))
		return c.Ack()
	}
	return nil
}

func (c *Context) Ack() error {
	if c.hasAck {
		return nil
	}
	if err := c.msg.Ack(); err != nil {
		return err
	}
	c.hasAck = true
	return nil
}

func (c *Context) Nack() error {
	if c.hasAck {
		return nil
	}
	if err := c.msg.Nack(); err != nil {
		return err
	}
	c.hasAck = true
	return nil
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *Context) reset() {
	c.msg = nil
	c.ctx = nil
	c.handlers = nil
	c.index = -1
	c.hasAck = false
}

func (c *Context) Next() {
	c.index++

	for c.index < len(c.handlers) {
		if err := c.handlers[c.index](c); err != nil {
			log.Error(err.Error())
			c.handlerErrs = append(c.handlerErrs, err)
		}
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}
