package mq

import (
	"github.com/ChsenDev/go-rabbitmq/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"sync"
)

const (
	defaultChannelCacheSize = 25
)

type Connection struct {
	*amqp.Connection
	sync.Mutex
	channelQueue *LinkedQueue[*amqp.Channel]
}

func (c *Connection) CreateChannel() (*amqp.Channel, error) {
	c.Lock()
	defer c.Unlock()
	ch := findOpenChannel(c.channelQueue)
	var err error
	if ch == nil {
		ch, err = c.Channel()
		if err != nil {
			return nil, err
		}
		c.channelQueue.Push(ch)
	}
	return ch, err
}

func findOpenChannel(channelQueue *LinkedQueue[*amqp.Channel]) *amqp.Channel {
	var ch *amqp.Channel
	for channelQueue.GetCapacity() > 0 {
		ch = channelQueue.Pop()
		if ch != nil && !ch.IsClosed() {
			break
		}
	}
	return ch
}

func (c *Connection) logicalClose(channel *amqp.Channel) {
	if !channel.IsClosed() {
		if c.channelQueue.GetCapacity() > defaultChannelCacheSize {
			c.physicalClose(channel)
		} else {
			c.channelQueue.Push(channel)
		}
	}
}

func (c *Connection) physicalClose(channel *amqp.Channel) {
	// 如果生产者未得到确认消息，应该等待生产者收到消息再关闭
	// 由于目前生产者确认机制使用的是同步确认，则无需关注这个
	err := channel.Close()
	if err != nil {
		log.Error("physicalClose err", zap.Error(err))
	}
}
