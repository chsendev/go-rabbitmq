package mq

import amqp "github.com/rabbitmq/amqp091-go"

type rawMessage struct {
	*amqp.Delivery
}

type Message struct {
	*rawMessage
	acknowledge
}
