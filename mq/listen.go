package mq

import (
	"fmt"
	"github.com/cscoder0/go-rabbitmq/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

func Listen(queue string) (<-chan *Message, error) {
	message := make(chan *Message)
	go func() {
		for {
			msg, client, err := listenMsg(queue)
			if err != nil {
				continue
			}
			closeCh := make(chan *amqp.Error)
		listenFor:
			for {
				fmt.Println("正在进入for循环")
				select {
				case d, ok := <-msg:
					fmt.Println("111112312321")
					if ok {
						log.Debug("Receive a message", zap.Any("data", d))
						raw := &rawMessage{Delivery: &d}
						message <- &Message{
							rawMessage:  raw,
							acknowledge: getAcknowledge(raw),
						}
					} else {
						log.Info("Receive a not ok message", zap.Any("data", d))
						break listenFor
					}
				case <-client.Channel.NotifyClose(closeCh):
					log.Warn("channel already close")
					break listenFor
				}
			}
			fmt.Println("1111111111111111111111111")
			fmt.Println("22222222222222222222222222")
			fmt.Println("Retrying...")
			log.Debug("Retrying...")
			time.Sleep(time.Second * 3)
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
