package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"context"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventMessage struct {
	Operation  string                 `json:"operation"`
	EntityID   string                 `json:"entity_id"`
	EntityType string                 `json:"entity_type"`
	Timestamp  time.Time              `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type Consumer struct {
	uri      string
	exchange string
	queue    string
	sync     *service.SyncService
}

func NewConsumer(uri, exchange, queue string, sync *service.SyncService) *Consumer {
	return &Consumer{uri: uri, exchange: exchange, queue: queue, sync: sync}
}

func (c *Consumer) Run(ctx context.Context) error {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error
	for attempt := 1; attempt <= 10; attempt++ {
		conn, err = amqp.Dial(c.uri)
		if err == nil {
			break
		}
		wait := time.Duration(attempt) * 500 * time.Millisecond
		log.Printf("rabbitmq consumer connect failed (attempt %d): %v; retrying in %s", attempt, err, wait)
		time.Sleep(wait)
	}
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}
	defer conn.Close()

	ch, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}
	defer ch.Close()

	// Ensure queue exists and bindings for both reservation and table events
	_, err = ch.QueueDeclare(c.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}
	if err := ch.QueueBind(c.queue, "reservation.*", c.exchange, false, nil); err != nil {
		return fmt.Errorf("queue bind reservation: %w", err)
	}
	if err := ch.QueueBind(c.queue, "table.*", c.exchange, false, nil); err != nil {
		return fmt.Errorf("queue bind table: %w", err)
	}

	msgs, err := ch.Consume(c.queue, "search-api", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case m := <-msgs:
			var evt EventMessage
			if err := json.Unmarshal(m.Body, &evt); err != nil {
				log.Printf("bad message: %v", err)
				_ = m.Nack(false, false)
				continue
			}
			// Validate event
			if evt.EntityID == "" || evt.Operation == "" {
				log.Printf("invalid event: missing entityID or operation")
				_ = m.Nack(false, false)
				continue
			}

			// Handle different entity types
			var handleErr error
			if evt.EntityType == "reservation" {
				handleErr = c.sync.HandleEvent(ctx, evt.Operation, evt.EntityID)
			} else if evt.EntityType == "table" {
				handleErr = c.sync.HandleTableEvent(ctx, evt.Operation, evt.EntityID, evt.Metadata)
			} else {
				log.Printf("unknown entity type: %s", evt.EntityType)
				_ = m.Nack(false, false)
				continue
			}

			if handleErr != nil {
				log.Printf("sync error: %v", handleErr)
				_ = m.Nack(false, true)
				continue
			}
			_ = m.Ack(false)
		}
	}
}
