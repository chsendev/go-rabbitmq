package mq

import (
	"github.com/ChsenDev/go-rabbitmq/config"
	"github.com/ChsenDev/go-rabbitmq/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"sync"
)

var connectionFactory = &cachingConnectionFactory{abstractConnectionFactory: &abstractConnectionFactory{}}

type ConnectionFactory interface {
	CreateConnection() (*Connection, error)
}

type abstractConnectionFactory struct {
	sync.Mutex
}

func (a *abstractConnectionFactory) CreateConnection() (*Connection, error) {
	a.Lock()
	defer a.Unlock()
	conn, err := amqp.Dial(config.Conf.Url)
	if err != nil {
		log.Error("CreateConnection error", zap.Error(err))
		return nil, err
	}
	return &Connection{Connection: conn, channelQueue: NewLinkedQueue[*amqp.Channel]()}, nil
}

type cachingConnectionFactory struct {
	*abstractConnectionFactory
	sync.Mutex
	conn *Connection
}

func (c *cachingConnectionFactory) CreateConnection() (*Connection, error) {
	c.Lock()
	defer c.Unlock()
	if c.conn == nil || c.conn.IsClosed() {
		var err error
		c.conn, err = c.abstractConnectionFactory.CreateConnection()
		if err != nil {
			return nil, err
		}
	}
	return c.conn, nil
}
