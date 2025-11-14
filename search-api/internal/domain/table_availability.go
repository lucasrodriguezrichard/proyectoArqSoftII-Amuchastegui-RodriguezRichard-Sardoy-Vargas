package domain

import (
	"fmt"
	"time"
)

// TableAvailability represents the MAIN ENTITY indexed in Solr
// It represents a table's availability for a specific date and meal type
// This is what gets indexed and searched in Solr
type TableAvailability struct {
	ID            string    `json:"id"`             // Format: "table-{meal_type}-{table_number}-{date}"
	TableNumber   int       `json:"table_number"`   // 1-10
	Capacity      int       `json:"capacity"`       // Number of people the table can hold
	MealType      string    `json:"meal_type"`      // breakfast, lunch, dinner, event
	Date          string    `json:"date"`           // Format: "2006-01-02"
	IsAvailable   bool      `json:"is_available"`   // true if available, false if reserved
	ReservationID string    `json:"reservation_id,omitempty"` // ID of reservation if taken
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GenerateTableAvailabilityID creates a unique ID for a table availability
// Format: table-{meal_type}-{table_number}-{YYYY-MM-DD}
func GenerateTableAvailabilityID(mealType string, tableNumber int, date string) string {
	return fmt.Sprintf("table-%s-%d-%s", mealType, tableNumber, date)
}

// NewTableAvailability creates a new table availability document
func NewTableAvailability(tableNumber int, capacity int, mealType string, date string) *TableAvailability {
	now := time.Now()
	return &TableAvailability{
		ID:          GenerateTableAvailabilityID(mealType, tableNumber, date),
		TableNumber: tableNumber,
		Capacity:    capacity,
		MealType:    mealType,
		Date:        date,
		IsAvailable: true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
