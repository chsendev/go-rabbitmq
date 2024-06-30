package mq

import (
	"github.com/ChsenDev/go-rabbitmq/config"
)

var DefaultOpts = []Option{new(qosOpt)}

type Opt func(*Client) error

type Option interface {
	Do() Opt
	Enabled() bool
}

type qosOpt struct {
}

func (c *qosOpt) Enabled() bool {
	return config.Conf.Listener.Prefetch > 0
}

func (c *qosOpt) Do() Opt {
	return func(client *Client) error {
		return client.Channel.Qos(config.Conf.Listener.Prefetch, 0, false)
	}
}
