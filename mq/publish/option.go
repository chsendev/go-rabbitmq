package publish

import (
	"github.com/cscoder0/go-rabbitmq/config"
	"time"
)

var defaultOpts = []Option{new(confirmOpt), new(notifyReturnOpt)}

type Option interface {
	Do() Opt
	Enabled() bool
}

type confirmOpt struct {
}

func (c *confirmOpt) Enabled() bool {
	return config.Conf.Publisher.Confirm.Enabled
}

func (c *confirmOpt) Do() Opt {
	return func(p *Publisher) error {
		p.confirm = true
		p.waitMilli = time.Duration(config.Conf.Publisher.Confirm.WaitMilli) * time.Millisecond
		return nil
	}
}

type notifyReturnOpt struct {
}

func (c *notifyReturnOpt) Enabled() bool {
	return config.Conf.Publisher.NotifyReturn
}

func (c *notifyReturnOpt) Do() Opt {
	return func(p *Publisher) error {
		p.Client.Channel.NotifyReturn(notifyReturn)
		p.mandatory = true
		return nil
	}
}

func WithDelay(delayed time.Duration) Opt {
	return func(publisher *Publisher) error {
		publisher.headers["x-delay"] = delayed / time.Millisecond
		return nil
	}
}
