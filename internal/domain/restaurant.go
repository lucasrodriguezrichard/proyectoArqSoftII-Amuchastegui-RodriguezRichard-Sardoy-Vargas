package domain

import (
	"time"
)

// Restaurant represents a restaurant entity
type Restaurant struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Address     string    `json:"address" db:"address"`
	Phone       string    `json:"phone" db:"phone"`
	Email       string    `json:"email" db:"email"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Table represents a table in the restaurant
type Table struct {
	ID           string `json:"id" db:"id"`
	RestaurantID string `json:"restaurant_id" db:"restaurant_id"`
	Number       int    `json:"number" db:"number"`
	Capacity     int    `json:"capacity" db:"capacity"`
	IsAvailable  bool   `json:"is_available" db:"is_available"`
}

// MenuItem represents a menu item
type MenuItem struct {
	ID           string  `json:"id" db:"id"`
	RestaurantID string  `json:"restaurant_id" db:"restaurant_id"`
	Name         string  `json:"name" db:"name"`
	Description  string  `json:"description" db:"description"`
	Price        float64 `json:"price" db:"price"`
	Category     string  `json:"category" db:"category"`
	IsAvailable  bool    `json:"is_available" db:"is_available"`
}
