package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"baseApi/config"
	"baseApi/logger"

	"github.com/streadway/amqp"
)

type RabbitMQPublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
}

var rabbitMQInstance *RabbitMQPublisher

/* InitRabbitMQ initializes RabbitMQ connection and sets up topic exchange */
func InitRabbitMQ(cfg *config.Config) error {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare topic exchange
	err = ch.ExchangeDeclare(
		cfg.RabbitMQExchange, // name
		"topic",              // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	rabbitMQInstance = &RabbitMQPublisher{
		connection: conn,
		channel:    ch,
		exchange:   cfg.RabbitMQExchange,
	}

	logger.Info("RabbitMQ initialized successfully with topic exchange:", cfg.RabbitMQExchange)
	return nil
}

/* GetRabbitMQPublisher returns the singleton RabbitMQ publisher instance */
func GetRabbitMQPublisher() *RabbitMQPublisher {
	return rabbitMQInstance
}

/* PublishJSON publishes JSON data to RabbitMQ topic exchange */
func (r *RabbitMQPublisher) PublishJSON(routingKey string, data interface{}) error {
	if r == nil || r.channel == nil {
		return fmt.Errorf("RabbitMQ publisher not initialized")
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Publish message
	err = r.channel.Publish(
		r.exchange,   // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	logger.Info(fmt.Sprintf("Published message to exchange '%s' with routing key '%s'", r.exchange, routingKey))
	return nil
}

/* PublishUserEvent publishes user-related events */
func (r *RabbitMQPublisher) PublishUserEvent(eventType string, userID uint, data interface{}) error {
	routingKey := fmt.Sprintf("user.%s", eventType)
	
	eventData := map[string]interface{}{
		"event_type": eventType,
		"user_id":    userID,
		"data":       data,
		"timestamp":  fmt.Sprintf("%d", getCurrentTimestamp()),
	}

	return r.PublishJSON(routingKey, eventData)
}

/* PublishSystemEvent publishes system-related events */
func (r *RabbitMQPublisher) PublishSystemEvent(eventType string, data interface{}) error {
	routingKey := fmt.Sprintf("system.%s", eventType)
	
	eventData := map[string]interface{}{
		"event_type": eventType,
		"data":       data,
		"timestamp":  fmt.Sprintf("%d", getCurrentTimestamp()),
	}

	return r.PublishJSON(routingKey, eventData)
}

/* Close closes RabbitMQ connection and channel */
func (r *RabbitMQPublisher) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.Printf("Error closing RabbitMQ channel: %v", err)
		}
	}
	if r.connection != nil {
		if err := r.connection.Close(); err != nil {
			log.Printf("Error closing RabbitMQ connection: %v", err)
		}
	}
	return nil
}

/* getCurrentTimestamp returns current unix timestamp */
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}