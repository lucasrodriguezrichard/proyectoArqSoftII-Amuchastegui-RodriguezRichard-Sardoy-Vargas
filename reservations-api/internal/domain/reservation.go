package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationStatus string

const (
	StatusPending   ReservationStatus = "pending"
	StatusConfirmed ReservationStatus = "confirmed"
	StatusCancelled ReservationStatus = "cancelled"
)

type Reservation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    uint               `bson:"userId"         json:"userId"`
	Date      time.Time          `bson:"date"           json:"date"`      // fecha y hora de la reserva
	PartySize int                `bson:"partySize"      json:"partySize"` // cantidad de personas
	TableID   string             `bson:"tableId"        json:"tableId"`   // mesa asignada (puede ser "" en pending)
	Status    ReservationStatus  `bson:"status"         json:"status"`
	Notes     string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt time.Time          `bson:"createdAt"      json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"      json:"updatedAt"`
}
