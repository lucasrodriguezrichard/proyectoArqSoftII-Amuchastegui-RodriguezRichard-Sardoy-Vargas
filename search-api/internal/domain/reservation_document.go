package domain

import "time"

// ReservationDocument represents the Solr document for a reservation
type ReservationDocument struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	TableNumber int       `json:"table_number"`
	Guests      int       `json:"guests"`
	DateTime    time.Time `json:"date_time"`
	MealType    string    `json:"meal_type"`
	Status      string    `json:"status"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
