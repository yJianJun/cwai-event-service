package rabbitmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

const bcgConfirmrderQueue = "queue.workorder.confirmWorkOrder"
const BcgExchange = "message_router"

// Producer 结构体表示RabbitMQ生产者
type Producer struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	exchange string
}

// NewProducer 创建并返回一个新的Producer实例
func NewProducer(url, exchange string) (*Producer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel for producer: %w", err)
	}

	logger.Info(context.Background(), "got Channel, declaring Exchange (message_router)")
	if err = ch.ExchangeDeclare(
		exchange, // name of the exchange
		"topic",  // type
		true,     // durable
		false,    // delete when complete
		false,    // internal
		false,    // noWait
		nil,      // arguments
	); err != nil {
		return nil, fmt.Errorf("failed to exchange declare: %v", err)
	}

	Producer := &Producer{
		conn:     conn,
		ch:       ch,
		exchange: exchange,
	}

	return Producer, nil
}

// Publish 发送消息到 RabbitMQ
func (p *Producer) Publish(routingKey string, message string) error {
	// 构造消息体
	body := []byte(message)
	msg := amqp.Publishing{
		ContentType: "text/json",
		Body:        body,
	}

	if routingKey == "" {
		routingKey = bcgConfirmrderQueue
	}

	// 发布消息到交换器
	err := p.ch.Publish(
		p.exchange, // 交换器名称
		routingKey, // 路由键
		false,      // 如果为true，则消息是持久的
		false,      // 如果为true，则消息在队列中等待消费者时，如果队列中没有消费者会立即返回给生产者
		msg,
	)
	if err != nil {
		logger.Errorf(context.Background(), "failed to publish a message to %v: %w", err, routingKey)
	}

	logger.Infof(context.Background(), "success to Sent %s", message)
	return nil
}

// Close 关闭连接和通道
func (p *Producer) Close() error {
	if err := p.ch.Close(); err != nil {
		return err
	}
	if err := p.conn.Close(); err != nil {
		return err
	}
	return nil
}

// 在主程序或其他地方使用Producer类
/*
func main() {
	// 创建并启动生产者
	Producer, err := NewProducer("amqp://dev-p1:l4KOMrzL9q@42.123.120.77:35672/billing", BcgExchange)
	if err != nil {
		log.Fatalf("Failed to create Producer: %v", err)
	}
	defer Producer.Close()

	// 开始消费消息
	if err := Producer.Publish(bcgConfirmrderQueue, ""); err != nil {
		log.Fatalf("Failed to producer messages: %v", err)
	}
}
*/
