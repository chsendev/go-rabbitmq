package mq

import (
	"github.com/cscoder0/go-rabbitmq/config"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Client struct {
	Conn    *Connection
	Channel *amqp.Channel
}

func newClient(conn *Connection, channel *amqp.Channel) *Client {
	return &Client{Conn: conn, Channel: channel}
}

func (i *Client) Close() {
	if i.Conn != nil && i.Channel != nil {
		i.Conn.logicalClose(i.Channel)
	}
}

func CreateClient() (*Client, error) {
	if config.Conf.Publisher.Retry != nil && config.Conf.Publisher.Retry.Enabled {
		wait := config.Conf.Publisher.Retry.InitialInterval
		for i := 0; i < config.Conf.Publisher.Retry.MaxAttempts; i++ {
			instance, err := createClient()
			if err != nil {
				return nil, err
			}
			if instance != nil {
				return instance, nil
			}
			time.Sleep(time.Second * time.Duration(wait))
			wait = wait * config.Conf.Publisher.Retry.Multiplier
		}
		return nil, errors.New("Connect failed")
	} else {
		return createClient()
	}

}

func createClient() (*Client, error) {
	conn, err := connectionFactory.CreateConnection()
	if err != nil {
		return nil, err
	}
	channel, err := conn.CreateChannel()
	if err != nil {
		return nil, err
	}
	return newClient(conn, channel), nil
}
