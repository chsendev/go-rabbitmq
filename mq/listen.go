package mq

import (
	"github.com/cscoder0/go-rabbitmq/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

func Listen(queue string) (<-chan *Message, error) {
	msg, client, err := listenMsg(queue)
	if err != nil {
		return nil, err
	}
	message := make(chan *Message)
	go func() {
		for {
			select {
			case d, ok := <-msg:
				if ok {
					raw := &rawMessage{Delivery: &d}
					message <- &Message{
						rawMessage:  raw,
						acknowledge: getAcknowledge(raw),
					}
				} else {
					log.Info("Receive a not ok message", zap.Any("data", d))
				}
			case <-client.Channel.NotifyClose(make(chan *amqp.Error)):
				//fmt.Println(123213)
				log.Warn("channel already close")
				time.Sleep(time.Second * 3)
				reMsg, reClient, err := listenMsg(queue)
				if err != nil {
					log.Error("Re again listen failed", zap.Error(err))
				} else {
					msg = reMsg
					client = reClient
				}
			}
		}
	}()
	return message, nil
}

func listenMsg(queue string) (<-chan amqp.Delivery, *Client, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, nil, err
	}
	for _, opt := range defaultOpts {
		if opt.Enabled() {
			if err = opt.Do()(client); err != nil {
				return nil, nil, err
			}
		}
	}
	msg, err := client.Channel.Consume(
		queue,       // queue
		"",          // consumer
		isAutoAck(), // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	return msg, client, err
}
