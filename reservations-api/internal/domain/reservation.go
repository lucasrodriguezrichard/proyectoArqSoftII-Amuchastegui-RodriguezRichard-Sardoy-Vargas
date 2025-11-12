package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Reservation statuses
const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)

// Meal types
const (
	MealTypeBreakfast = "breakfast"
	MealTypeLunch     = "lunch"
	MealTypeDinner    = "dinner"
	MealTypeEvent     = "event"
)

// Reservation represents a restaurant reservation
type Reservation struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID         string             `bson:"owner_id" json:"owner_id"`
	TableNumber     int                `bson:"table_number" json:"table_number"`
	Guests          int                `bson:"guests" json:"guests"`
	DateTime        time.Time          `bson:"date_time" json:"date_time"`
	MealType        string             `bson:"meal_type" json:"meal_type"`
	Status          string             `bson:"status" json:"status"`
	TotalPrice      float64            `bson:"total_price" json:"total_price"`
	SpecialRequests string             `bson:"special_requests,omitempty" json:"special_requests,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateReservationRequest DTO for creating a reservation
type CreateReservationRequest struct {
	OwnerID         string    `json:"owner_id" binding:"required"`
	TableNumber     int       `json:"table_number" binding:"required,min=1"`
	Guests          int       `json:"guests" binding:"required,min=1,max=20"`
	DateTime        time.Time `json:"date_time" binding:"required"`
	MealType        string    `json:"meal_type" binding:"required,oneof=breakfast lunch dinner event"`
	SpecialRequests string    `json:"special_requests,omitempty"`
}

// UpdateReservationRequest DTO for updating a reservation
type UpdateReservationRequest struct {
	TableNumber     *int       `json:"table_number,omitempty" binding:"omitempty,min=1"`
	Guests          *int       `json:"guests,omitempty" binding:"omitempty,min=1,max=20"`
	DateTime        *time.Time `json:"date_time,omitempty"`
	MealType        *string    `json:"meal_type,omitempty" binding:"omitempty,oneof=breakfast lunch dinner event"`
	SpecialRequests *string    `json:"special_requests,omitempty"`
	Status          *string    `json:"status,omitempty" binding:"omitempty,oneof=pending confirmed cancelled completed"`
}

// ConfirmReservationRequest DTO for confirming a reservation
type ConfirmReservationRequest struct {
	ConfirmationNotes string `json:"confirmation_notes,omitempty"`
}

// Validate checks if the reservation data is valid
func (r *Reservation) Validate() error {
	if r.OwnerID == "" {
		return errors.New("owner_id is required")
	}
	if r.TableNumber < 1 {
		return errors.New("table_number must be positive")
	}
	if r.Guests < 1 || r.Guests > 20 {
		return errors.New("guests must be between 1 and 20")
	}
	if r.DateTime.Before(time.Now()) {
		return errors.New("date_time must be in the future")
	}
	if !isValidMealType(r.MealType) {
		return errors.New("invalid meal_type")
	}
	if !isValidStatus(r.Status) {
		return errors.New("invalid status")
	}
	return nil
}

func isValidMealType(mt string) bool {
	switch mt {
	case MealTypeBreakfast, MealTypeLunch, MealTypeDinner, MealTypeEvent:
		return true
	}
	return false
}

func isValidStatus(s string) bool {
	switch s {
	case StatusPending, StatusConfirmed, StatusCancelled, StatusCompleted:
		return true
	}
	return false
}

// NewReservation creates a new reservation with default values
func NewReservation(req CreateReservationRequest) Reservation {
	now := time.Now()
	return Reservation{
		OwnerID:         req.OwnerID,
		TableNumber:     req.TableNumber,
		Guests:          req.Guests,
		DateTime:        req.DateTime,
		MealType:        req.MealType,
		Status:          StatusPending,
		TotalPrice:      0, // will be calculated
		SpecialRequests: req.SpecialRequests,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
