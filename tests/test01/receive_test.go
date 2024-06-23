package test01

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"reflect"
	"testing"
	"time"
)

//func failOnError(err error, msg string) {
//	if err != nil {
//		log.Panicf("%s: %s", msg, err)
//	}
//}

func TestReceive(t *testing.T) {
	conn, err := amqp.Dial("amqp://test:123@192.168.31.207:5672//test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			time.Sleep(time.Millisecond * 200)
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func TestReceive2(t *testing.T) {
	conn, err := amqp.Dial("amqp://test:123@192.168.31.207:5672//test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			time.Sleep(time.Millisecond * 20)
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func test() {
	//mq(func(a any) {
	//	// 我希望，这里拿到的a已经是真实的类型
	//})
}

func TestName(t *testing.T) {
	data := ([]byte)(`1`)
	a := new(int)
	json.Unmarshal(data, a)
	fmt.Println(*a) //1
	//mq(func(ctx context.Context, a int) {
	//	fmt.Println(a)
	//})
	//mq(func(ctx context.Context, u User) {
	//	fmt.Println(u)
	//})
}

type User struct {
	Name string `json:"name"`
}

func mq(f func(ctx context.Context, any2 any)) {
	data := ([]byte)(`{"name":"abc"}`)

	doFun := reflect.TypeOf(f)

	in := make([]reflect.Value, doFun.NumIn())
	for i := range in {
		paramType := doFun.In(i)
		switch paramType.Kind() {
		case reflect.Int:
			var p int
			_ = json.Unmarshal(data, &p)
			in[i] = reflect.ValueOf(p)
		case reflect.Struct:
			paramValue := reflect.New(paramType).Elem()
			_ = json.Unmarshal(data, paramValue.Addr().Interface())
			in[i] = paramValue
		}
	}
	reflect.ValueOf(f).Call(in)
}
