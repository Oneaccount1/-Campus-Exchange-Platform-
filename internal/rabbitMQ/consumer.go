package rabbitMQ

import (
	"campus/internal/modules/message/api"
	"campus/internal/utils/logger"
	"campus/internal/websocket"
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

// StartConsumer initializes and runs the message consumer.
// It should be run as a goroutine from the main application.
func StartConsumer(url string, wsManager *websocket.Manager) {
	logger.Info("Starting message consumer...")

	// Loop indefinitely to handle reconnects
	for {
		err := runConsumer(url, wsManager)
		if err != nil {
			logger.Error("Message consumer error. Reconnecting...", zap.Error(err))
			time.Sleep(reconnectDelay)
		}
	}
}

func runConsumer(url string, wsManager *websocket.Manager) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Ensure the topology (exchange, queue, binding) exists.
	if err := setupTopology(ch); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		messageQueue, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	logger.Info("Message consumer is waiting for messages.")

	// Process messages
	for d := range msgs {
		processMessage(d, wsManager)
	}

	return nil
}

func setupTopology(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		messageExchange, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		messageQueue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		messageQueue, messageRoutingKey, messageExchange, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func processMessage(d amqp.Delivery, wsManager *websocket.Manager) {
	var msgResponse api.MessageResponse
	if err := json.Unmarshal(d.Body, &msgResponse); err != nil {
		logger.Error("Failed to unmarshal message", zap.Error(err))
		d.Nack(false, false) // Nack and don't requeue, bad message format
		return
	}

	receiverID := msgResponse.ReceiverID
	logger.Debug("Processing message for user", zap.Uint("userID", receiverID))

	if wsManager.IsUserOnline(receiverID) {
		logger.Debug("User is online, attempting to send via WebSocket", zap.Uint("userID", receiverID))
		if ok := wsManager.SendMessage(receiverID, d.Body); ok {
			logger.Info("Successfully sent message via WebSocket", zap.Uint("userID", receiverID))
			d.Ack(false) // Acknowledge the message
		} else {
			logger.Warn("Failed to send message via WebSocket, user might have just disconnected.", zap.Uint("userID", receiverID))
			d.Nack(false, true) // Requeue the message
		}
	} else {
		logger.Debug("User is offline. Discarding message from queue since it's already in DB.", zap.Uint("userID", receiverID))
		// If user is offline, we still Ack the message to remove it from the queue.
		// The authoritative store is the database, which the user will query for history.
		// This prevents the queue from filling with messages for offline users.
		d.Ack(false)
	}
}
