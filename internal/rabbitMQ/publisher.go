package rabbitMQ

import (
	"campus/internal/utils/logger"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	reconnectDelay       = 5 * time.Second
	maxReconnectAttempts = 10
	messageExchange      = "message_exchange"
	messageQueue         = "messages.queue"
	messageRoutingKey    = "message.route"
)

// Publisher implements the services.RabbitMQPublisher interface.
// It is a RabbitMQ client designed for publishing messages in a reliable way,
// with auto-reconnect logic.
type Publisher struct {
	url       string
	conn      *amqp.Connection
	channel   *amqp.Channel
	connMutex sync.RWMutex
	done      chan struct{}
}

// NewPublisher creates a new Publisher instance.
// It establishes a connection and starts a monitoring goroutine for auto-reconnect.
func NewPublisher(url string) (*Publisher, error) {
	p := &Publisher{
		url:  url,
		done: make(chan struct{}),
	}

	if err := p.connect(); err != nil {
		return nil, err
	}

	go p.reconnectLoop()

	return p, nil
}

func (p *Publisher) connect() error {
	var err error
	p.conn, err = amqp.Dial(p.url)
	if err != nil {
		logger.Error("RabbitMQ publisher connection failed", zap.Error(err))
		return err
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		p.conn.Close()
		logger.Error("RabbitMQ publisher failed to open a channel", zap.Error(err))
		return err
	}

	// Declare the exchange
	err = p.channel.ExchangeDeclare(
		messageExchange,
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		p.conn.Close()
		logger.Error("RabbitMQ publisher failed to declare an exchange", zap.Error(err))
		return err
	}

	logger.Info("RabbitMQ Publisher connected and exchange declared")
	return nil
}

func (p *Publisher) reconnectLoop() {
	connClose := p.conn.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case <-p.done:
			return
		case err := <-connClose:
			if err != nil {
				logger.Warn("RabbitMQ publisher connection lost. Attempting to reconnect...", zap.Error(err))
			}

			for i := 0; i < maxReconnectAttempts; i++ {
				time.Sleep(reconnectDelay)
				if err := p.connect(); err == nil {
					logger.Info("RabbitMQ publisher reconnected successfully.")
					connClose = p.conn.NotifyClose(make(chan *amqp.Error))
					break
				}
				logger.Error("Failed to reconnect publisher", zap.Int("attempt", i+1))
			}
		}
	}
}

// Publish sends a message to the pre-configured exchange.
func (p *Publisher) Publish(body []byte, contentType string) error {
	p.connMutex.RLock()
	defer p.connMutex.RUnlock()

	if p.channel == nil {
		return fmt.Errorf("RabbitMQ channel is not available")

	}

	return p.channel.Publish(
		messageExchange,   // exchange
		messageRoutingKey, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// Close gracefully shuts down the publisher's connection.
func (p *Publisher) Close() {
	close(p.done)
	p.connMutex.Lock()
	defer p.connMutex.Unlock()
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	logger.Info("RabbitMQ publisher closed.")
}
