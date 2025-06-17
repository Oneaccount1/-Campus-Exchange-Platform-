package messaging

import (
	"campus/internal/utils/logger"
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

// 定义RabbitMQ连接相关常量
const (
	// 重连间隔时间
	reconnectDelay = 5 * time.Second
	// 发布确认超时时间
	publishConfirmTimeout = 5 * time.Second
	// 最大重试次数
	maxReconnectAttempts = 10
)

// RabbitMQ 连接管理器
type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
	queueName    string
	routingKey   string
	url          string

	// 连接状态
	connected bool
	connMutex sync.RWMutex

	// 连接关闭通知通道
	connClose    chan *amqp.Error
	channelClose chan *amqp.Error

	// 停止信号
	done chan struct{}
}

// NewRabbitMQ 创建新的RabbitMQ连接
func NewRabbitMQ(url, exchangeName, queueName, routingKey string) (*RabbitMQ, error) {
	rmq := &RabbitMQ{
		url:          url,
		exchangeName: exchangeName,
		queueName:    queueName,
		routingKey:   routingKey,
		connected:    false,
		done:         make(chan struct{}),
	}

	// 建立初始连接
	if err := rmq.connect(); err != nil {
		return nil, err
	}

	// 启动连接监控
	go rmq.reconnectLoop()

	return rmq, nil
}

// connect 建立RabbitMQ连接
func (r *RabbitMQ) connect() error {
	var err error

	logger.Debug("开始连接RabbitMQ",
		zap.String("url", r.url),
		zap.String("exchange", r.exchangeName),
		zap.String("queue", r.queueName))

	// 建立连接
	r.conn, err = amqp.Dial(r.url)
	if err != nil {
		logger.Error("RabbitMQ连接失败",
			zap.String("url", r.url),
			zap.Error(err))
		return err
	}

	// 监听连接关闭
	r.connClose = make(chan *amqp.Error)
	r.conn.NotifyClose(r.connClose)

	// 创建通道
	r.channel, err = r.conn.Channel()
	if err != nil {
		logger.Error("RabbitMQ创建通道失败", zap.Error(err))
		r.conn.Close()
		return err
	}

	// 监听通道关闭
	r.channelClose = make(chan *amqp.Error)
	r.channel.NotifyClose(r.channelClose)

	// 声明交换机
	err = r.channel.ExchangeDeclare(
		r.exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("RabbitMQ声明交换机失败",
			zap.String("exchange", r.exchangeName),
			zap.Error(err))
		r.conn.Close()
		return err
	}

	// 声明队列
	_, err = r.channel.QueueDeclare(
		r.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("RabbitMQ声明队列失败",
			zap.String("queue", r.queueName),
			zap.Error(err))
		r.conn.Close()
		return err
	}

	// 绑定队列到交换机
	err = r.channel.QueueBind(
		r.queueName,
		r.routingKey,
		r.exchangeName,
		false,
		nil,
	)
	if err != nil {
		logger.Error("RabbitMQ绑定队列到交换机失败",
			zap.String("queue", r.queueName),
			zap.String("exchange", r.exchangeName),
			zap.Error(err))
		r.conn.Close()
		return err
	}

	// 设置QoS，限制预取数量
	err = r.channel.Qos(
		10,    // 预取计数
		0,     // 预取大小
		false, // 全局设置
	)
	if err != nil {
		logger.Error("RabbitMQ设置QoS失败", zap.Error(err))
		r.conn.Close()
		return err
	}

	// 标记为已连接
	r.connMutex.Lock()
	r.connected = true
	r.connMutex.Unlock()

	logger.Info("RabbitMQ初始化成功",
		zap.String("exchange", r.exchangeName),
		zap.String("queue", r.queueName))

	return nil
}

// reconnectLoop 监控连接状态并自动重连
func (r *RabbitMQ) reconnectLoop() {
	attempts := 0

	for {
		select {
		case <-r.done:
			// 收到关闭信号，退出重连循环
			return

		case err := <-r.connClose:
			// 连接关闭，尝试重连
			r.connMutex.Lock()
			r.connected = false
			r.connMutex.Unlock()

			logger.Warn("RabbitMQ连接已关闭，准备重连",
				zap.String("reason", err.Reason),
				zap.Error(err))

			// 重连逻辑
			for attempts < maxReconnectAttempts {
				attempts++
				logger.Info("尝试重新连接RabbitMQ",
					zap.Int("attempt", attempts),
					zap.Int("maxAttempts", maxReconnectAttempts))

				// 等待重连间隔
				time.Sleep(reconnectDelay)

				// 尝试重新连接
				if err := r.connect(); err == nil {
					// 重连成功
					logger.Info("RabbitMQ重连成功")
					attempts = 0
					break
				} else {
					logger.Error("RabbitMQ重连失败", zap.Error(err))
				}
			}

			if attempts >= maxReconnectAttempts {
				logger.Error("RabbitMQ重连失败次数过多，放弃重连")
				return
			}

		case err := <-r.channelClose:
			// 通道关闭，尝试重建通道
			r.connMutex.Lock()
			r.connected = false
			r.connMutex.Unlock()

			logger.Warn("RabbitMQ通道已关闭，准备重建",
				zap.String("reason", err.Reason),
				zap.Error(err))

			// 尝试重建通道
			if r.conn != nil && !r.conn.IsClosed() {
				var err error
				r.channel, err = r.conn.Channel()
				if err != nil {
					logger.Error("重建RabbitMQ通道失败", zap.Error(err))
				} else {
					// 重新设置通道监听
					r.channelClose = make(chan *amqp.Error)
					r.channel.NotifyClose(r.channelClose)

					// 重新设置QoS
					r.channel.Qos(10, 0, false)

					r.connMutex.Lock()
					r.connected = true
					r.connMutex.Unlock()

					logger.Info("RabbitMQ通道重建成功")
				}
			}
		}
	}
}

// Close 关闭RabbitMQ连接
func (r *RabbitMQ) Close() error {
	// 发送停止信号
	close(r.done)

	// 关闭通道和连接
	if r.channel != nil {
		r.channel.Close()
	}

	if r.conn != nil {
		logger.Debug("关闭RabbitMQ连接",
			zap.String("exchange", r.exchangeName),
			zap.String("queue", r.queueName))
		return r.conn.Close()
	}

	return nil
}

// IsConnected 检查连接状态
func (r *RabbitMQ) IsConnected() bool {
	r.connMutex.RLock()
	defer r.connMutex.RUnlock()
	return r.connected
}

// Publish 发布消息
func (r *RabbitMQ) Publish(data interface{}) error {
	// 检查连接状态
	if !r.IsConnected() {
		return errors.New("RabbitMQ连接不可用")
	}

	body, err := json.Marshal(data)
	if err != nil {
		logger.Error("RabbitMQ消息序列化失败", zap.Error(err))
		return err
	}

	err = r.channel.Publish(
		r.exchangeName,
		r.routingKey,
		false, // mandatory
		false, // immediate
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

		// 如果发布失败，标记连接为不可用，触发重连
		r.connMutex.Lock()
		r.connected = false
		r.connMutex.Unlock()
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
	// 检查连接状态
	if !r.IsConnected() {
		return nil, errors.New("RabbitMQ连接不可用")
	}

	// 设置消费者，不自动确认
	deliveries, err := r.channel.Consume(
		r.queueName,
		"",    // 消费者标签
		false, // 自动确认
		false, // 排他性
		false, // 不等待服务器确认
		false, // 参数
		nil,   // 参数
	)

	if err != nil {
		logger.Error("RabbitMQ启动消费者失败",
			zap.String("queue", r.queueName),
			zap.Error(err))
		return nil, err
	}

	logger.Info("RabbitMQ启动消费者成功", zap.String("queue", r.queueName))
	return deliveries, nil
}

// GetUserQueue 根据用户ID获取队列名
func GetUserQueue(userID uint) string {
	return "user_queue_" + strconv.FormatUint(uint64(userID), 10)
}

// CheckHealth 检查RabbitMQ连接健康状态
func (r *RabbitMQ) CheckHealth() error {
	if !r.IsConnected() {
		return errors.New("RabbitMQ连接不可用")
	}

	if r.conn == nil || r.conn.IsClosed() {
		return errors.New("RabbitMQ连接已关闭")
	}

	if r.channel == nil {
		return errors.New("RabbitMQ通道不可用")
	}

	return nil
}
