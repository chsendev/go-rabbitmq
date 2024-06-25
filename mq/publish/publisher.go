package publish

import (
	"context"
	"github.com/cscoder0/go-rabbitmq/log"
	"github.com/cscoder0/go-rabbitmq/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"sync"
	"time"
)

var once sync.Once
var notifyReturn chan amqp.Return

type Opt func(*Publisher) error

type Publisher struct {
	Client    *mq.Client
	options   []Opt
	err       error
	confirm   bool
	mandatory bool
	waitMilli time.Duration
	// delayed 延迟消息时间
	//delayed time.Duration
	headers amqp.Table
}

func New(opt ...Opt) *Publisher {
	p := &Publisher{options: opt, headers: make(map[string]any)}
	return p
}

func (p *Publisher) returnErr(err error) *Publisher {
	p.err = err
	return p
}

func (p *Publisher) Error() error {
	if p.Client != nil {
		p.Client.Close()
	}
	return p.err
}

func (p *Publisher) doDefaultOpt() error {
	for _, option := range defaultOpts {
		if option.Enabled() {
			if err := option.Do()(p); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Publisher) Publish(ctx context.Context, exchange, routingKey string, data any) *Publisher {
	if p.err != nil {
		return p
	}
	var err error
	p.Client, err = mq.CreateClient()
	if err != nil {
		return p.returnErr(err)
	}
	defer p.Client.Close()
	if err = p.doDefaultOpt(); err != nil {
		return p.returnErr(err)
	}
	for _, opt := range p.options {
		if err = opt(p); err != nil {
			return p.returnErr(err)
		}
	}
	return p.returnErr(getPublish(p).Publish(ctx, exchange, routingKey, data))
}

func ListenNotifyReturn(callback func(*amqp.Return)) {
	once.Do(func() {
		notifyReturn = make(chan amqp.Return)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Error("ListenNotifyReturn err", zap.Any("error", err))
				}
			}()
			for {
				select {
				case returned, ok := <-notifyReturn:
					if !ok {
						log.Debug("channel closed, rebuilding return channel")
						notifyReturn = make(chan amqp.Return)
						continue
					}
					log.Info("Received a return message", zap.Any("data", returned))
					callback(&returned)
				}
			}
		}()
	})
}
