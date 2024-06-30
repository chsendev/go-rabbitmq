package mq

import (
	"github.com/ChsenDev/go-rabbitmq/config"
	"github.com/ChsenDev/go-rabbitmq/log"
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
	if config.Conf.Retry != nil {
		wait := config.Conf.Retry.InitialInterval
		for i := 0; i < config.Conf.Retry.MaxAttempts; i++ {
			instance, err := createClient()
			if instance != nil {
				return instance, nil
			}
			if err != nil {
				log.Error(err.Error())
			}
			time.Sleep(wait)
			wait = wait * time.Duration(config.Conf.Retry.Multiplier)
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
