package service

import (
	"testing"
	"time"
)

// Note: These tests focus on testing the EventMessage structure and basic functionality
// Full integration tests would require a running RabbitMQ instance

func TestEventMessage_Structure(t *testing.T) {
	now := time.Now()
	metadata := map[string]interface{}{
		"table_number": 1,
		"meal_type":    "lunch",
	}

	msg := EventMessage{
		Operation:  "create",
		EntityID:   "12345",
		EntityType: "reservation",
		Timestamp:  now,
		Metadata:   metadata,
	}

	if msg.Operation != "create" {
		t.Errorf("expected operation 'create', got %s", msg.Operation)
	}

	if msg.EntityID != "12345" {
		t.Errorf("expected entity_id '12345', got %s", msg.EntityID)
	}

	if msg.EntityType != "reservation" {
		t.Errorf("expected entity_type 'reservation', got %s", msg.EntityType)
	}

	if msg.Timestamp != now {
		t.Errorf("expected timestamp %v, got %v", now, msg.Timestamp)
	}

	if msg.Metadata == nil {
		t.Fatal("expected metadata to be set")
	}

	if msg.Metadata["table_number"] != 1 {
		t.Errorf("expected table_number 1, got %v", msg.Metadata["table_number"])
	}

	if msg.Metadata["meal_type"] != "lunch" {
		t.Errorf("expected meal_type 'lunch', got %v", msg.Metadata["meal_type"])
	}
}

func TestEventMessage_WithoutMetadata(t *testing.T) {
	msg := EventMessage{
		Operation:  "delete",
		EntityID:   "67890",
		EntityType: "table",
		Timestamp:  time.Now(),
		Metadata:   nil,
	}

	if msg.Metadata != nil {
		t.Error("expected nil metadata")
	}
}

func TestEventMessage_TableEvent(t *testing.T) {
	metadata := map[string]interface{}{
		"table_number": 5,
		"capacity":     8,
		"meal_type":    "dinner",
	}

	msg := EventMessage{
		Operation:  "update",
		EntityID:   "table-123",
		EntityType: "table",
		Timestamp:  time.Now(),
		Metadata:   metadata,
	}

	if msg.EntityType != "table" {
		t.Errorf("expected entity_type 'table', got %s", msg.EntityType)
	}

	if msg.Metadata["capacity"] != 8 {
		t.Errorf("expected capacity 8, got %v", msg.Metadata["capacity"])
	}
}

func TestEventMessage_ReservationEvent(t *testing.T) {
	msg := EventMessage{
		Operation:  "confirm",
		EntityID:   "res-456",
		EntityType: "reservation",
		Timestamp:  time.Now(),
		Metadata:   nil,
	}

	if msg.Operation != "confirm" {
		t.Errorf("expected operation 'confirm', got %s", msg.Operation)
	}

	if msg.EntityType != "reservation" {
		t.Errorf("expected entity_type 'reservation', got %s", msg.EntityType)
	}
}

// Test NewRabbitMQPublisher with invalid URI
func TestNewRabbitMQPublisher_InvalidURI(t *testing.T) {
	// This will fail after retries, which is expected behavior
	publisher, err := NewRabbitMQPublisher("amqp://invalid:5672", "test-exchange", "test-queue")

	if err == nil {
		t.Error("expected error for invalid URI, got nil")
		if publisher != nil {
			publisher.Close()
		}
	}
}
