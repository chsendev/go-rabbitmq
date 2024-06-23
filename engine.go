package rmq

import (
	"context"
	"fmt"
	"github.com/cscoder0/go-rabbitmq/config"
	"github.com/cscoder0/go-rabbitmq/log"
	"github.com/cscoder0/go-rabbitmq/mq"
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

func New(conf *config.RabbitmqConfig) *Engine {
	Init(conf)
	var h Engine
	h.poll.New = func() any {
		return &Context{Ctx: context.Background()}
	}
	h.handlerMap = make(map[string]HandlersChain)
	h.RouterGroup = &RouterGroup{Engine: &h}
	return &h
}

func (h *Engine) addHandlers(queue string, handlers HandlersChain) {
	if h.err != nil {
		return
	}
	_, has := h.handlerMap[queue]
	if has {
		panic(fmt.Errorf("queue %s already exists", queue))
	}
	h.handlerMap[queue] = handlers
	msg, err := mq.Listen(queue)
	if err != nil {
		h.err = err
	}
	go h.listen(msg, handlers)
}

func (h *Engine) listen(msg <-chan *mq.Message, handlers HandlersChain) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("handle panic", zap.Any("error", err))
		}
	}()

	for {
		select {
		case d, ok := <-msg:
			fmt.Println(ok)
			go h.handle(d, handlers)
		}
	}
}

func (h *Engine) handle(msg *mq.Message, chain HandlersChain) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("handle panic", zap.Any("error", err))
			if err := msg.Nack(); err != nil {
				log.Error("nack err", zap.Error(err))
			}
		} else {
			if err := msg.Ack(); err != nil {
				log.Error("ack err", zap.Error(err))
			}
		}
	}()
	log.Debug("Received a message", zap.Any("data", msg))
	c := h.poll.Get().(*Context)
	c.msg = msg
	for _, handle := range chain {
		err := handle(c) // todo error
		if err != nil {
			msg.Nack()
		}
	}
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
	err := mq.Binding(exchange, exchangeType, queue, bindingKey...)
	return h.returnErr(err)
}

func (h *Engine) Error() error {
	return h.err
}

func (h *Engine) returnErr(err error) *Engine {
	h.err = err
	return h
}
