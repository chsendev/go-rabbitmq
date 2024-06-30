# go-rabbitmq

基于amqp的封装，除了一些基础功能（声明交换机、队列、发送消息、消费消息），还封装了一些高级功能：

- 生产者确认
- 消费者确认
- 消息多拦截器
- 延迟消息
- 统一消息处理
- 消费端断开重连
- 内置连接池
- ...

# 安装
```shell
go get github.com/ChsenDev/go-rabbitmq
```

# 基础功能

## 初始化

```go
rmq.Init("amqp://test:123@127.0.0.1:5672//test")
```

## 声明资源

**方式一（常用）**

```go
// 声明名称为go-demo类型为Direct的交换机
// 与demo.queue1队列进行绑定，binding key为q1
// 与demo.queue2队列进行绑定，binding key为q2、q3
engine := rmq.New().
Binding("go-demo", mq.Direct, "demo.queue1", "q1")
Binding("go-demo", mq.Direct, "demo.queue2", "q2", "q3")
```

**方式二**

```go
engine := rmq.New().
Exchange("test.direct", mq.Direct). // 声明交换机
Queue("test.queue1").  // 声明队列
BindingKey("k1", "k2") // 声明test.direct和test.queue1的binging key
```

1. 可只声明某个交换机或者某个队列
2. 声明BindingKey之前必须声明交换机和声明队列

## 监听消息

```go
engine.Listen("test.queue1", func (ctx *rmq.Context) error {
var u User
if err := ctx.ShouldBind(&u); err != nil {
return err
}
/// ...
})
```

## 获取Error

```go
if engine.Error() != nil {
panic(engine.Error())
}
```

## 发送消息

```go
m := make(map[string]any)
m["name"] = "jack"
m["age"] = "10"
if err := rmq.Publish(context.Background(), "go-demo", "q1", m).Error(); err != nil {
fmt.Println(err)
}
```

# 高级功能

## 消费者拉取消息限制

有时候需要限制每个消费者一次性拉取消息的条数，如果设置得当的话，可以起到能者多劳的效果（性能较好的消费者能够消费更多的消息）

```go
rmq.Init("amqp://test:123@127.0.0.1:5672//test", rmq.WithPrefetch(1))
```

## 设置消费者Ack策略

框架默认提供了三种策略：

- none：自动进行Ack，无论消息是否消费成功
- auto(默认)：如果Handlers处理成功，则进行Ack，否则Nack
- manual：完全由开发者手动调用ctx.Ack进行Ack

```go
rmq.Init("amqp://test:123@127.0.0.1:5672//test", rmq.WithAckMode(config.AcknowledgeModeAuto))
```

## 设置框架日志

框架使用zap作为日志的管理，默认的级别为Info，若需要调整可使用WithLogLevel进行设置

```go
rmq.Init("amqp://test:123@127.0.0.1:5672//test", rmq.WithLogLevel("debug"))
```

## 连接重试策略

当与rabbitmq通行的channel断开链接时，框架内部会进行重试。

initialInterval(默认1s)：失败后的初始等待时间
multiplier(默认2)：失败后下次的等待时长倍数，下次等待时长 = initialInterval * multiplier
maxAttempts(默认3)：最大重试次数

```go
rmq.Init("amqp://test:123@127.0.0.1:5672//test", rmq.WithRetry(time.Second, 2, 3))
```

## 生产者确认

1. Publish Confirm
   ```go
    if err := rmq.Publish(context.Background(), "go-demo", "q3", m, publish.WithConfirm(time.Second*3)); err != nil {
        fmt.Println(err)
    }
    ```
   若rabbitmq未在3秒内返回ack消息，则返回error
2. Publish Return
   ```go
    // 监听Publish Return的通道
	rmq.ListenNotifyReturn(func(msg *amqp.Return) {
		fmt.Println(msg)
	})
	if err := rmq.Publish(context.Background(), "go-demo123", "q1", m, publish.WithNotifyReturn()); err != nil {
		panic(err)
	}
    ```

## 消费者确认

框架默认提供了三种策略：none、auto、manual，如果设置manual，需要开发者手动调用ctx.Ack进行Ack

```go
 rmq.Init("amqp://test:123@127.0.0.1:5672//test", rmq.WithAckMode(config.AcknowledgeModeManual))
engine := rmq.New().Binding("go-demo", mq.Direct, "demo.queue1", "q1")
engine.Listen("demo.queue1", func (ctx *rmq.Context) error {
// ...
err := ctx.Ack() // 或ctx.Nack()
// ...
})

if engine.Error() != nil {
panic(engine.Error())
}
```

## 消息拦截器

有时候消息可能需要经过多个处理器(handler)进行处理，例如：自定义消息的重试次数、消息消费统一错误处理等
自定义消息的重试次数：
```go
type User struct {
	Name string
}

var messageDb = make(map[string]int)

// 前置检查
func preCheck(ctx *rmq.Context) error {
	num := messageDb[ctx.GetMessageId()]

	// 限制重试三次
	if num >= 3 {
		fmt.Println("达到重试次数")
		ctx.Abort()
		return ctx.Ack()
	}
	messageDb[ctx.GetMessageId()] = num + 1
	return nil
}

// 业务逻辑
func logic(ctx *rmq.Context) error {
	var u User
	if err := ctx.ShouldBind(&u); err != nil {
		return err
	}
	fmt.Println("获取到的消息：", u)
	// 手动触发异常
	return errors.New("手动触发的错误")
}

func TestRetry(t *testing.T) {
	rmq.Init("amqp://test:123@127.0.0.1:5672//test", rmq.WithAckMode(config.AcknowledgeModeAuto))
	engine := rmq.New().Binding("go-demo", mq.Direct, "demo.queue1", "q1")

	g1 := engine.Group(preCheck)
	{
		g1.Listen("demo.queue1", logic)
	}

	if engine.Error() != nil {
		panic(engine.Error())
	}

	time.Sleep(time.Hour)
}

func TestRetryPublish(t *testing.T) {
	u := &User{Name: "jack"}
	rmq.Init("amqp://test:123@127.0.0.1:5672//test", rmq.WithAckMode(config.AcknowledgeModeNone))
	err := rmq.Publish(context.Background(), "go-demo", "q1", u)
	fmt.Println(err)
}
```

## 延迟消息
延时消息：可以允许消息在指定时间后再发送给消费者，此功能需要额外的插件支持：
https://github.com/rabbitmq/rabbitmq-delayed-message-exchange
```go
func TestListenDelay(t *testing.T) {
	rmq.Init("amqp://test:123@127.0.0.1:5672//test", rmq.WithLogLevel("debug"))
	engine := rmq.New().BindingWithDelay("delay-demo", mq.Direct, "demo.queue1", "q1")

	engine.Listen("demo.queue1", func(ctx *rmq.Context) error {
		var s string
		if err := ctx.ShouldBind(&s); err != nil {
			return err
		}
		log.Println("receive message: ", s)
		return nil
	})

	if engine.Error() != nil {
		panic(engine.Error())
	}

	time.Sleep(time.Hour)
}

func TestPublishDelay(t *testing.T) {
	rmq.Init("amqp://test:123@127.0.0.1:5672//test", rmq.WithLogLevel("debug"))
	msg := "test delay message"
	err := rmq.Publish(context.Background(), "delay-demo", "q1", msg, publish.WithDelay(time.Second*10))
	if err != nil {
		panic(err)
	}
	log.Println("publish message: ", msg)
}
```
