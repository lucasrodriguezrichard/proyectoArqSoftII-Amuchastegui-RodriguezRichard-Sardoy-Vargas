package domain

import (
	"time"
)

// ReservationStatus represents the status of a reservation
type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusCompleted ReservationStatus = "completed"
)

// MealType represents the type of meal
type MealType string

const (
	MealTypeBreakfast MealType = "breakfast"
	MealTypeLunch     MealType = "lunch"
	MealTypeSnack     MealType = "snack"
	MealTypeDinner    MealType = "dinner"
	MealTypePrivate   MealType = "private_event"
	MealTypeOther     MealType = "other"
)

// Reservation represents a table reservation
type Reservation struct {
	ID              string            `json:"id" db:"id"`
	RestaurantID    string            `json:"restaurant_id" db:"restaurant_id"`
	TableID         string            `json:"table_id" db:"table_id"`
	CustomerName    string            `json:"customer_name" db:"customer_name"`
	CustomerPhone   string            `json:"customer_phone" db:"customer_phone"`
	CustomerEmail   string            `json:"customer_email" db:"customer_email"`
	ReservationDate time.Time         `json:"reservation_date" db:"reservation_date"`
	MealType        MealType          `json:"meal_type" db:"meal_type"`
	PartySize       int               `json:"party_size" db:"party_size"`
	Status          ReservationStatus `json:"status" db:"status"`
	SpecialRequests string            `json:"special_requests" db:"special_requests"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}
