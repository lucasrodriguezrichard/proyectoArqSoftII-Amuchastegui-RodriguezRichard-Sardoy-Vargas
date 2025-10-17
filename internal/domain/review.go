package domain

import (
	"time"
)

// Review represents a customer review
type Review struct {
	ID             string    `json:"id" db:"id"`
	ReservationID  string    `json:"reservation_id" db:"reservation_id"`
	CustomerName   string    `json:"customer_name" db:"customer_name"`
	Rating         int       `json:"rating" db:"rating"` // 1-5 stars
	Title          string    `json:"title" db:"title"`
	Comment        string    `json:"comment" db:"comment"`
	FoodRating     int       `json:"food_rating" db:"food_rating"`         // 1-5 stars
	ServiceRating  int       `json:"service_rating" db:"service_rating"`   // 1-5 stars
	AmbienceRating int       `json:"ambience_rating" db:"ambience_rating"` // 1-5 stars
	IsAnonymous    bool      `json:"is_anonymous" db:"is_anonymous"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Ticket represents a payment ticket/bill
type Ticket struct {
	ID            string      `json:"id" db:"id"`
	OrderID       string      `json:"order_id" db:"order_id"`
	ReservationID string      `json:"reservation_id" db:"reservation_id"`
	TableNumber   int         `json:"table_number" db:"table_number"`
	CustomerName  string      `json:"customer_name" db:"customer_name"`
	Items         []OrderItem `json:"items" db:"-"`
	Subtotal      float64     `json:"subtotal" db:"subtotal"`
	Tax           float64     `json:"tax" db:"tax"`
	Tip           float64     `json:"tip" db:"tip"`
	Total         float64     `json:"total" db:"total"`
	PaymentMethod string      `json:"payment_method" db:"payment_method"`
	PaidAt        time.Time   `json:"paid_at" db:"paid_at"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
}
