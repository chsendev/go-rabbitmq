package test03

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(n int) (res int, err error) {
	conn, err := amqp.Dial("amqp://test:123@192.168.31.207:5672//test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

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

	corrId := randomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			failOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

func TestClient(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	//n := bodyFrom(os.Args)
	n := 30

	log.Printf(" [x] Requesting fib(%d)", n)
	res, err := fibonacciRPC(n)
	failOnError(err, "Failed to handle RPC request")

	log.Printf(" [.] Got %d", res)
}

func bodyFrom(args []string) int {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	failOnError(err, "Failed to convert arg to integer")
	return n
}

func TestName(t *testing.T) {
	var domainRegex = regexp.MustCompile(`(https://|http://)?(.*?)/(.*)/products/`)
	fmt.Println(domainRegex.MatchString("https://www.burga.com/products/all-eyes-on-me-apple-watch-band?variant=42305230209199&size=38mm%20%2F%2040mm%20%2F%2041mm&color=Gold"))

	u, _ := url.Parse("https://www.burga.com/products/all-eyes-on-me-apple-watch-band?variant=42305230209199&size=38mm%20%2F%2040mm%20%2F%2041mm&color=Gold")
	fmt.Println(u.Scheme + "://" + u.Host)

	fmt.Println(path.Ext("https://uk.yeti.com/cdn/shop/files/Drinkware_Tumbler_20oz_Navy_Stud.png?v=1705679835&width=60"))
	//url.JoinPath("http")
	//fmt.Println(domainRegex)
}

func producer(ch chan<- int, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for i := 0; i < 1; i++ {
		fmt.Printf("Producer %d is producing %d\n", id, i)
		ch <- i
		time.Sleep(time.Millisecond * 100) // simulate some work
	}
}

func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Printf("Consumer is consuming %d\n", v)
		time.Sleep(time.Millisecond * 200) // simulate some work
	}
}

func TestGorouting(t *testing.T) {
	ch := make(chan int, 3)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go producer(ch, &wg, i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	consumer(ch)
}
