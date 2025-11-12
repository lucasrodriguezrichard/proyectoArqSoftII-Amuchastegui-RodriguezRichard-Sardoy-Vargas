package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher handles publishing messages to RabbitMQ
type RabbitMQPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queue    string
}

// EventMessage represents the message format for RabbitMQ
type EventMessage struct {
	Operation  string    `json:"operation"`  // create, update, delete
	EntityID   string    `json:"entity_id"`  // reservation ID
	EntityType string    `json:"entity_type"` // always "reservation"
	Timestamp  time.Time `json:"timestamp"`
}

// NewRabbitMQPublisher creates a new RabbitMQ publisher
func NewRabbitMQPublisher(uri, exchange, queue string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange
	err = channel.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue
	_, err = channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	err = channel.QueueBind(
		queue,                // queue name
		"reservation.*",      // routing key
		exchange,             // exchange
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	log.Printf("RabbitMQ publisher connected to exchange: %s, queue: %s", exchange, queue)

	return &RabbitMQPublisher{
		conn:     conn,
		channel:  channel,
		exchange: exchange,
		queue:    queue,
	}, nil
}

// Publish sends a message to RabbitMQ
func (p *RabbitMQPublisher) Publish(operation, entityID string) error {
	msg := EventMessage{
		Operation:  operation,
		EntityID:   entityID,
		EntityType: "reservation",
		Timestamp:  time.Now(),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	routingKey := fmt.Sprintf("reservation.%s", operation)

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange,  // exchange
		routingKey,  // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to RabbitMQ: %s %s", operation, entityID)
	return nil
}

// Close closes the RabbitMQ connection
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}
