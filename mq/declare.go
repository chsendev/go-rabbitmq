package mq

func Queue(queueName string) error {
	client, err := CreateClient()
	if err != nil {
		return err
	}
	defer client.Close()
	return queue(client, queueName)
}

func queue(client *Client, queue string) error {
	_, err := client.Channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func Exchange(exchangeName string, exchangeType ExchangeType) error {
	client, err := CreateClient()
	if err != nil {
		return err
	}
	defer client.Close()
	return exchange(client, exchangeName, exchangeType)
}

func exchange(client *Client, exchange string, exchangeType ExchangeType) error {
	return client.Channel.ExchangeDeclare(
		exchange,
		string(exchangeType),
		true,
		false,
		false,
		false,
		nil,
	)
}

func BindingKey(exchange string, exchangeType ExchangeType, queue string, bindingKeyName ...string) error {
	client, err := CreateClient()
	if err != nil {
		return err
	}
	defer client.Close()
	return bindingKey(client, exchange, exchangeType, queue, bindingKeyName...)
}

func bindingKey(client *Client, exchange string, exchangeType ExchangeType, queue string, bindingKey ...string) error {
	for _, key := range bindingKey {
		err := client.Channel.QueueBind(
			queue,    // queue name
			key,      // routing key
			exchange, // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func Binding(exchangeName string, exchangeType ExchangeType, queueName string, bindingKeyName ...string) error {
	client, err := CreateClient()
	if err != nil {
		return err
	}
	defer client.Close()
	if err = queue(client, queueName); err != nil {
		return err
	}
	if err = exchange(client, exchangeName, exchangeType); err != nil {
		return err
	}
	if err = bindingKey(client, exchangeName, exchangeType, queueName, bindingKeyName...); err != nil {
		return err
	}
	return nil
}
