package rmq

import (
	"context"
	"github.com/chsendev/go-rabbitmq/mq/publish"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(ctx context.Context, exchange, routingKey string, data any, opt ...publish.Opt) error {
	return publish.New(opt...).Publish(ctx, exchange, routingKey, data).Error()
}

func ListenNotifyReturn(callback func(*amqp.Return)) {
	publish.ListenNotifyReturn(callback)
}
