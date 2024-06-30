package mq

import (
	"fmt"
	"github.com/ChsenDev/go-rabbitmq/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

func Listen(queue string, messageHandle func(msg *Message)) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("handle panic", zap.Any("error", err))
			}
		}()
		for {
			msg, client, err := listenMsg(queue)
			if err != nil {
				time.Sleep(time.Second * 3)
				log.Error("listen message error", zap.Any("error", err))
				continue
			}
			closeCh := make(chan *amqp.Error)
		listenFor:
			for {
				select {
				case d, ok := <-msg:
					if ok {
						log.Debug("Receive a message", zap.Any("data", d))
						raw := &rawMessage{Delivery: &d}
						messageHandle(&Message{
							rawMessage:  raw,
							acknowledge: getAcknowledge(raw),
						})
					} else {
						log.Info("Receive a not ok message", zap.Any("data", d))
						break listenFor
					}
				case e := <-client.Channel.NotifyClose(closeCh):
					fmt.Println(e)
					log.Warn("channel already close")
					break listenFor
				}
			}
			log.Debug("Retrying...")
			time.Sleep(time.Second * 3)
		}
	}()
}

func listenMsg(queue string) (<-chan amqp.Delivery, *Client, error) {
	client, err := CreateClient()
	if err != nil {
		return nil, nil, err
	}
	for _, opt := range DefaultOpts {
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
