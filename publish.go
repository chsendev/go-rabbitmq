package rmq

import (
	"context"
	"github.com/cscoder0/go-rabbitmq/mq/publish"
)

func Publish(ctx context.Context, exchange, routingKey string, data any, opt ...publish.Opt) error {
	return publish.New(opt...).Publish(ctx, exchange, routingKey, data).Error()
}
