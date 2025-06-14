package messaging

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

// RabbitMQ 连接管理器
type RabbitMQ struct {
	conn         *amqp.Connection
	chanel       *amqp.Channel
	exchangeName string
	queueName    string
	routingKey   string
	url          string
}

// 新建连接

func NewRabbitMQ(url, exchangeName, queueName, routingKey string) (*RabbitMQ, error) {
	rmq := &RabbitMQ{
		url:          url,
		exchangeName: exchangeName,
		queueName:    queueName,
		routingKey:   routingKey,
	}
	var err error

	rmq.conn, err = amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	rmq.chanel, err = rmq.conn.Channel()
	if err != nil {
		return nil, err
	}

	// 声明交换机
	err = rmq.chanel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// 声明队列
	_, err = rmq.chanel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// 绑定队列到交换机
	err = rmq.chanel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	log.Println("RabbitMQ初始化成功")
	return rmq, nil
}

// Close 关闭rabbitMQ连接
func (r *RabbitMQ) Close() error {
	if r.chanel != nil {
		r.chanel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// 发布消息

func (r *RabbitMQ) Publish(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.chanel.Publish(
		r.exchangeName,
		r.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

// Consume 消费消息
func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	return r.chanel.Consume(
		r.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

// GetUserQueue 获取用户专属队列
func GetUserQueue(userID uint) string {
	return "user_message_" + string(rune(userID))
}
