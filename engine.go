package rmq

import (
	"context"
	"github.com/chsendev/go-rabbitmq/config"
	"github.com/chsendev/go-rabbitmq/log"
	"github.com/chsendev/go-rabbitmq/mq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
)

var declareQueue string
var declareExchange string
var declareExchangeType mq.ExchangeType

type Engine struct {
	*RouterGroup
	poll       sync.Pool
	err        error
	handlerMap map[string]HandlersChain
}

func New() *Engine {
	var h Engine
	h.poll.New = func() any {
		return &Context{}
	}
	h.handlerMap = make(map[string]HandlersChain)
	h.RouterGroup = &RouterGroup{Engine: &h}
	return &h
}

func (h *Engine) addHandlers(queue string, handlers HandlersChain) {
	if h.err != nil {
		return
	}
	if config.Conf.Listener.AcknowledgeMode == config.AcknowledgeModeAuto {
		handlers = append(handlers, AutoAck)
	}
	h.handlerMap[queue] = handlers
	mq.Listen(queue, func(msg *mq.Message) {
		h.handle(msg, handlers)
	})
}

func (h *Engine) handle(msg *mq.Message, chain HandlersChain) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("handle panic", zap.Any("error", err))
		}
	}()
	c := h.poll.Get().(*Context)
	c.reset()
	c.msg = msg
	c.ctx = context.Background()
	c.handlers = chain
	c.Next()
	h.poll.Put(c)
}

func (h *Engine) Queue(queue string) *Engine {
	if h.err != nil {
		return h
	}
	declareQueue = queue
	return h.returnErr(mq.Queue(queue))
}

func (h *Engine) Exchange(exchange string, exchangeType mq.ExchangeType) *Engine {
	if h.err != nil {
		return h
	}
	declareExchange = exchange
	declareExchangeType = exchangeType
	return h.returnErr(mq.Exchange(exchange, exchangeType))
}

func (h *Engine) BindingKey(bindingKey ...string) *Engine {
	if h.err != nil {
		return h
	}
	if declareQueue == "" || declareExchange == "" {
		return h.returnErr(errors.New("Missing queue and exchange, please call Queue() and Exchange() first"))
	}
	err := mq.BindingKey(declareExchange, declareExchangeType, declareQueue, bindingKey...)
	declareQueue = ""
	declareExchange = ""
	declareExchangeType = ""
	return h.returnErr(err)
}

func (h *Engine) Binding(exchange string, exchangeType mq.ExchangeType, queue string, bindingKey ...string) *Engine {
	if h.err != nil {
		return h
	}
	err := mq.Binding(exchange, exchangeType, queue, nil, bindingKey...)
	return h.returnErr(err)
}

func (h *Engine) BindingWithDelay(exchange string, exchangeType mq.ExchangeType, queue string, bindingKey ...string) *Engine {
	if h.err != nil {
		return h
	}
	err := mq.Binding(exchange, "x-delayed-message", queue, map[string]any{"x-delayed-type": string(exchangeType)}, bindingKey...)
	return h.returnErr(err)
}

// AutoAck auto ack
func AutoAck(ctx *Context) error {
	var err error
	if len(ctx.handlerErrs) > 0 {
		log.Debug("nack message", zap.Any("data", ctx.msg))
		if err = ctx.Nack(); err != nil {
			log.Error("nack err", zap.Error(err))
		}
	} else {
		log.Debug("ack message", zap.Any("data", ctx.msg))
		if err = ctx.Ack(); err != nil {
			log.Error("ack err", zap.Error(err))
		}
	}
	return err
}

func (h *Engine) Error() error {
	return h.err
}

func (h *Engine) returnErr(err error) *Engine {
	h.err = err
	return h
}
