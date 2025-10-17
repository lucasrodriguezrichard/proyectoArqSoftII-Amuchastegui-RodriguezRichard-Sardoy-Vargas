package domain

import (
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusServed    OrderStatus = "served"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID         string  `json:"id" db:"id"`
	OrderID    string  `json:"order_id" db:"order_id"`
	MenuItemID string  `json:"menu_item_id" db:"menu_item_id"`
	Quantity   int     `json:"quantity" db:"quantity"`
	Price      float64 `json:"price" db:"price"`
	Notes      string  `json:"notes" db:"notes"`
}

// Order represents a food order
type Order struct {
	ID            string      `json:"id" db:"id"`
	ReservationID string      `json:"reservation_id" db:"reservation_id"`
	TableID       string      `json:"table_id" db:"table_id"`
	CustomerName  string      `json:"customer_name" db:"customer_name"`
	Status        OrderStatus `json:"status" db:"status"`
	Items         []OrderItem `json:"items" db:"-"`
	Subtotal      float64     `json:"subtotal" db:"subtotal"`
	Tax           float64     `json:"tax" db:"tax"`
	Tip           float64     `json:"tip" db:"tip"`
	Total         float64     `json:"total" db:"total"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}
