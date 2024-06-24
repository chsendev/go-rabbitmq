package rmq

import (
	"context"
	"fmt"
	"github.com/cscoder0/go-rabbitmq/config"
	"github.com/cscoder0/go-rabbitmq/mq"
	"github.com/cscoder0/go-rabbitmq/mq/publish"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"testing"
	"time"
)

func getConfig() *config.RabbitmqConfig {
	var conf config.RabbitmqConfig
	viper.SetConfigFile("./rabbitmq.yaml")
	if e := viper.ReadInConfig(); e != nil {
		panic(e)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
	return &conf
}

func TestName(t *testing.T) {
	engine := New(getConfig()).
		Binding("go-demo", mq.Direct, "demo.queue1", "q1").
		Binding("go-demo", mq.Direct, "demo.queue2", "q2").
		Use(f1).Use(f2)

	g := engine.Group(f3, f4)
	{
		g.Listen("demo.queue1", f5, f6)
	}
	//g2 := engine.Group(f7, f8)
	//{
	//	g2.Listen("demo.queue2", f9)
	//}

	fmt.Println("11233123")
	//Listen("demo.queue1", lis).
	//	Listen("demo.queue2", lis).Error()
	if engine.Error() != nil {
		panic(engine.Error())
	}
	time.Sleep(time.Hour)
}

func f1(ctx *Context) error {
	return nil
}
func f2(ctx *Context) error {
	return nil
}
func f3(ctx *Context) error {
	return nil
}
func f4(ctx *Context) error {
	return nil
}
func f5(ctx *Context) error {
	return nil
}
func f6(ctx *Context) error {
	return nil
}
func f7(ctx *Context) error {
	return nil
}
func f8(ctx *Context) error {
	return nil
}
func f9(ctx *Context) error {
	return nil
}

func TestPublish(t *testing.T) {
	m := make(map[string]any)
	m["name"] = "jack"
	m["age"] = "10"
	//config.Init(getConfig())
	Init(getConfig())
	//fmt.Println(mq.Publish(context.Background(), "go-demo", "q1", m))

	//publish.Pub(context.Background(), "go-demo", "q2", m)

	//ret := make(chan amqp.Return)
	////publish.New(publish.NotifyReturn(ret)).Publish()
	//opt := publish.NotifyReturn(ret, func(returns <-chan amqp.Return) {
	//	r := <-ret
	//	fmt.Println(r)
	//})
	//if err := publish.New().Publish(context.Background(), "go-demo", "q2", m).Error(); err != nil {
	//	fmt.Println(err)
	//}
	publish.ListenNotifyReturn(func(a *amqp.Return) {
		fmt.Println(a)
	})
	if err := publish.New().Publish(context.Background(), "go-demo", "q11", m).Error(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(1)
	if err := publish.New().Publish(context.Background(), "go-demo", "q12", m).Error(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(1)
	if err := publish.New().Publish(context.Background(), "go-demo", "q13", m).Error(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(1)
	if err := publish.New().Publish(context.Background(), "go-demo", "q14", m).Error(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(1)
	if err := publish.New().Publish(context.Background(), "go-demo", "q15", m).Error(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(1)
	if err := publish.New().Publish(context.Background(), "go-demo", "q16", m).Error(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(1)
	if err := publish.New().Publish(context.Background(), "go-demo", "q17", m).Error(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(1)
	if err := publish.New().Publish(context.Background(), "go-demo", "q18", m).Error(); err != nil {
		fmt.Println(err)
	}
	//if err := publish.New(publish.NotifyReturn()).Publish(context.Background(), "go-demo", "q3", m).Error(); err != nil {
	//	fmt.Println(err)
	//}
	//if err := publish.New(publish.NotifyReturn()).Publish(context.Background(), "go-demo", "q3", m).Error(); err != nil {
	//	fmt.Println(err)
	//}
	fmt.Println(123)

	time.Sleep(time.Hour)
	//r := <-ret
	//fmt.Println(r)
}

func TestName123(t *testing.T) {
	conn, err := amqp.Dial("amqp://test:123@192.168.31.207:5672//test")
	if err != nil {
		fmt.Println(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch.IsClosed())
	fmt.Println()
	for i := 0; i < 100; i++ {
		fmt.Println(ch.IsClosed())
	}
	fmt.Println(ch.IsClosed())
}

func TestQueue(t *testing.T) {
	q := mq.NewLinkedQueue[string]()
	q.Push("1")
	q.Push("2")
	q.Push("3")
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())

	fmt.Println(q)
}

func TestChan(t *testing.T) {
	var ch chan int

	go func() {
		for {
			select {
			case i, ok := <-ch:
				fmt.Println(i, ok)
				if !ok {
					ch = make(chan int)
				}
			}
		}
	}()

	ch <- 10
	time.Sleep(time.Second * 3)
	close(ch)
	time.Sleep(time.Second * 3)
	ch <- 20
	close(ch)
}

func TestChan2(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 10
	go func() {
		for {
			select {
			case i := <-ch:
				fmt.Println(i)
				ch = make(chan int, 1)
				ch <- 10
			}
		}
	}()
	time.Sleep(time.Hour)
}
