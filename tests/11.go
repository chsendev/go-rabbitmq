package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

// MessageHandler 是一个函数类型，它接受一个消息内容（通常是某个结构体类型）并返回一个错误（如果有的话）
type MessageHandler func(ctx context.Context, msg interface{}) error

// GenericConsumer 是一个通用的 RabbitMQ 消费者函数
func GenericConsumer(ch *amqp.Channel, queueName string, handler MessageHandler, ack bool) {
	// 定义消费者标签，用于在取消消费时识别
	consumerTag := "generic_consumer"

	// 定义消息交付模式（手动或自动）
	//deliveryMode := amqp.Transient // 例如，使用 Transient 作为默认，但你可以根据需要更改为 Persistent
	//if ack {
	//	deliveryMode = amqp.Persistent
	//}

	// 声明队列，确保它存在
	_, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 定义消费者回调函数
	deliveries, err := ch.Consume(
		queueName,   // queue
		consumerTag, // consumer
		ack,         // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range deliveries {
			// 处理消息
			var msgBody interface{}                 // 这里可以是任何类型，例如特定的结构体
			err := json.Unmarshal(d.Body, &msgBody) // 假设消息体是 JSON 格式
			if err != nil {
				log.Printf("Failed to unmarshal message body: %v", err)
				if ack {
					d.Nack(false, true) // 不重新入队，标记为失败
				}
				continue
			}

			// 调用处理函数
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置超时时间
			defer cancel()
			if err := handler(ctx, msgBody); err != nil {
				log.Printf("Error processing message: %v", err)
				if ack {
					d.Nack(false, true) // 不重新入队，标记为失败
				}
			} else if ack {
				d.Ack(false) // 确认消息已成功处理
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 示例消息处理函数
func processMessage(ctx context.Context, msg MyMessage) error {
	// 处理 MyMessage 类型的消息
	fmt.Printf("Received a message: %v\n", msg)
	// 模拟处理时间
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(1 * time.Second):
		// 处理完成
	}
	return nil
}

// 示例消息结构体
type MyMessage struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

func main() {
	// ... 省略 RabbitMQ 连接和通道创建的代码 ...

	// 使用 GenericConsumer 函数启动消费者
	//GenericConsumer(nil, "my_queue", processMessage, true) // 第三个参数是消息处理函数，第四个参数表示是否启用手动确认
	//
	//// ... 等待程序
	gin.Default().Run()
}
