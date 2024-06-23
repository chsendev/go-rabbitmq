package publish

import (
	"context"
	"encoding/json"
	"github.com/cscoder0/go-rabbitmq/log"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"strconv"
)

type Publish interface {
	Publish(ctx context.Context, exchange, routingKey string, data any) error
}

type defaultPublish struct {
	*Publisher
}

func (d *defaultPublish) Publish(ctx context.Context, exchange, routingKey string, data any) error {
	_, err := publish(ctx, d.Publisher, exchange, routingKey, data)
	return err
}

type publishWithConfirm struct {
	*Publisher
}

func (d *publishWithConfirm) Publish(ctx context.Context, exchange, routingKey string, data any) error {
	if err := d.Client.Channel.Confirm(false); err != nil {
		return err
	}
	confirmation, err := publish(ctx, d.Publisher, exchange, routingKey, data)
	ctx, cancel := context.WithTimeout(ctx, d.waitMilli)
	defer cancel()
	ack, err := confirmation.WaitContext(ctx)
	if err != nil {
		log.Debug("Publish confirm publish err", zap.Error(err))
	}
	log.Debug("Receive publish ack：" + strconv.FormatBool(ack))
	if !ack {
		return errors.New("Receive publish false ack")
	}
	return err
}

func publish(ctx context.Context, d *Publisher, exchange, routingKey string, data any) (*amqp.DeferredConfirmation, error) {
	j, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return d.Client.Channel.PublishWithDeferredConfirm(
		exchange,    // exchange
		routingKey,  // routing key
		d.mandatory, // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        j,
		})
}

func getPublish(p *Publisher) Publish {
	if p.confirm {
		return &publishWithConfirm{p}
	} else {
		return &defaultPublish{p}
	}
}
