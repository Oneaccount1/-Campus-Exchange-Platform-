package messaging

import (
	"campus/internal/utils/logger"
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"strconv"
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

	logger.Debug("开始连接RabbitMQ",
		zap.String("url", url),
		zap.String("exchange", exchangeName),
		zap.String("queue", queueName))

	rmq.conn, err = amqp.Dial(url)
	if err != nil {
		logger.Error("RabbitMQ连接失败",
			zap.String("url", url),
			zap.Error(err))
		return nil, err
	}

	rmq.chanel, err = rmq.conn.Channel()
	if err != nil {
		logger.Error("RabbitMQ创建通道失败", zap.Error(err))
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
		logger.Error("RabbitMQ声明交换机失败",
			zap.String("exchange", exchangeName),
			zap.Error(err))
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
		logger.Error("RabbitMQ声明队列失败",
			zap.String("queue", queueName),
			zap.Error(err))
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
		logger.Error("RabbitMQ绑定队列到交换机失败",
			zap.String("queue", queueName),
			zap.String("exchange", exchangeName),
			zap.Error(err))
		return nil, err
	}

	logger.Info("RabbitMQ初始化成功",
		zap.String("exchange", exchangeName),
		zap.String("queue", queueName))
	return rmq, nil
}

// Close 关闭rabbitMQ连接
func (r *RabbitMQ) Close() error {
	if r.chanel != nil {
		r.chanel.Close()
	}
	if r.conn != nil {
		logger.Debug("关闭RabbitMQ连接",
			zap.String("exchange", r.exchangeName),
			zap.String("queue", r.queueName))
		return r.conn.Close()
	}
	return nil
}

// 发布消息

func (r *RabbitMQ) Publish(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		logger.Error("RabbitMQ消息序列化失败", zap.Error(err))
		return err
	}

	err = r.chanel.Publish(
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

	if err != nil {
		logger.Error("RabbitMQ发布消息失败",
			zap.String("exchange", r.exchangeName),
			zap.String("routingKey", r.routingKey),
			zap.Error(err))
	} else {
		logger.Debug("RabbitMQ发布消息成功",
			zap.String("exchange", r.exchangeName),
			zap.String("routingKey", r.routingKey),
			zap.Int("size", len(body)))
	}

	return err
}

// Consume 消费消息
func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	deliveries, err := r.chanel.Consume(
		r.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logger.Error("RabbitMQ启动消费者失败",
			zap.String("queue", r.queueName),
			zap.Error(err))
	} else {
		logger.Info("RabbitMQ启动消费者成功", zap.String("queue", r.queueName))
	}

	return deliveries, err
}

// GetUserQueue 根据用户ID获取队列名
func GetUserQueue(userID uint) string {
	return "user_queue_" + strconv.FormatUint(uint64(userID), 10)
}
