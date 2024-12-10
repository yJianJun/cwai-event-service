package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"

	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

const rabbitMQTag = "CWAI"
const rabbitMQConsumer = "cwai"
const BcgMsgQueue = "cwai.noticeResInstanceAMQP"

// Consumer 结构体表示RabbitMQ消费者
type Consumer struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queue     string
	tag       string
	handler   func(delivery amqp.Delivery)
	closed    bool
	mu        sync.Mutex
	serverUrl string
	queueName string
}

// NewConsumer 创建并返回一个新的Consumer实例
func NewConsumer(url, queueName string, handler func(delivery amqp.Delivery)) (*Consumer, error) {
	consumer := &Consumer{
		tag:       rabbitMQTag,
		handler:   handler,
		serverUrl: url,
		queueName: queueName,
	}

	if err := consumer.Init(); err != nil {
		return nil, err
	}
	consumer.MonitorConn()
	return consumer, nil
}

func (c *Consumer) Init() error {
	conn, err := amqp.Dial(c.serverUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	logger.Info(context.Background(), "got Channel, declaring Exchange (message_router)")
	if err = ch.ExchangeDeclare(
		"message_router", // name of the exchange
		"topic",          // type
		true,             // durable
		false,            // delete when complete
		false,            // internal
		false,            // noWait
		nil,              // arguments
	); err != nil {
		return fmt.Errorf("failed to exchange declare: %v", err)
	}

	q, err := ch.QueueDeclare(
		c.queueName, // name
		true,        // durable
		false,       // auto-delete
		false,       // exclusive
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare a queue: %w", err)
	}
	logger.Infof(context.Background(),
		"declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		q.Name, q.Messages, q.Consumers, c.queueName)

	if err = ch.QueueBind(
		c.queueName,      // name of the queue
		c.queueName,      // bindingKey
		"message_router", // sourceExchange
		false,            // noWait
		nil,              // arguments
	); err != nil {
		return fmt.Errorf("failed to queue Bind: %v", err)
	}
	c.ch = ch
	c.conn = conn
	c.queue = q.Name
	c.closed = false
	return nil
}

// Consume 开始消费消息
func (c *Consumer) Consume(consumerName string) error {
	msgs, err := c.ch.Consume(
		c.queue,      // queue
		consumerName, // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			logger.Infof(context.Background(), "consumed a massage successfully")
			// 调用用户提供的消息处理函数
			c.handler(d)
		}
	}()

	logger.Info(context.Background(), " [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

// Close 关闭连接和channel
func (c *Consumer) Close() error {
	if err := c.ch.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (c *Consumer) MonitorConn() {
	go func() {
		for {
			select {
			case <-time.After(time.Second * 5):
				if c.conn.IsClosed() {
					c.mu.Lock()
					c.closed = true
					c.mu.Unlock()

					log.Println("Connection closed, attempting to reconnect...")
					if err := c.reconnect(); err != nil {
						log.Println("Failed to reconnect:", err)
					}
				}
			}
		}
	}()
}

func (c *Consumer) reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		return nil
	}

	c.Init()
	return nil
}

// 示例消息处理函数
/*
func HandleDelivery(delivery amqp.Delivery) {
	logger.Infof("Received a message: %s", delivery.Body)
	var msg groupv1.MqMessage
	if err := json.Unmarshal(delivery.Body, &msg); err != nil {
		logger.Errorf("failed to unmarshal mq message: %v", err)
		return
	}
	if err := query.Q.WithContext(context.TODO()).MqMessage.Create(&model.MqMessage{
		MasterOrderID: msg.MasterOrderId,
		RawMessage:    string(delivery.Body),
	}); err != nil {
		logger.Errorf("failed to save mq message to db: %v", err)
		return
	}
	// 在此处添加实际的消息处理逻辑
	delivery.Ack(true) // 根据实际情况决定是否立即确认消息
}*/

// 在主程序或其他地方使用Consumer类
/*
func main() {
	// 定义消息处理函数
	handleMsg := func(delivery amqp.Delivery) {
		log.Printf("Received a message: %s", delivery.Body)
		// 处理接收到的消息...
		delivery.Ack(false)
	}

	// 创建并启动消费者
	consumer, err := NewConsumer("amqp://dev-p1:l4KOMrzL9q@42.123.120.77:35672/", "cwai.noticeResInstanceAMQP", handleMsg)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// 开始消费消息
	if err := consumer.Consume("CWAI"); err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	// 等待，直到收到终止信号（例如CTRL+C）
	select {}
}
*/
