package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReservationRepository defines the interface for reservation persistence
type ReservationRepository interface {
	Create(ctx context.Context, reservation *domain.Reservation) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Reservation, error)
	GetByUserID(ctx context.Context, userID string) ([]domain.Reservation, error)
	Update(ctx context.Context, id primitive.ObjectID, reservation *domain.Reservation) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetReservedTableNumbers(ctx context.Context, date string, mealType string) ([]int, error)
}

// MongoReservationRepository implements ReservationRepository using MongoDB
type MongoReservationRepository struct {
	collection *mongo.Collection
}

// NewMongoReservationRepository creates a new MongoDB reservation repository
func NewMongoReservationRepository(collection *mongo.Collection) *MongoReservationRepository {
	return &MongoReservationRepository{
		collection: collection,
	}
}

// Create inserts a new reservation
func (r *MongoReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	result, err := r.collection.InsertOne(ctx, reservation)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	// Set the ID from the inserted document
	reservation.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a reservation by ID
func (r *MongoReservationRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error) {
	var reservation domain.Reservation

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&reservation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("reservation not found")
		}
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}

	return &reservation, nil
}

// GetAll retrieves all reservations with pagination
func (r *MongoReservationRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Reservation, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "date_time", Value: -1}}) // Sort by date descending

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}
	defer cursor.Close(ctx)

	var reservations []domain.Reservation
	if err := cursor.All(ctx, &reservations); err != nil {
		return nil, fmt.Errorf("failed to decode reservations: %w", err)
	}

	return reservations, nil
}

// GetByUserID retrieves all reservations for a specific user
func (r *MongoReservationRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Reservation, error) {
	filter := bson.M{"owner_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "date_time", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reservations: %w", err)
	}
	defer cursor.Close(ctx)

	var reservations []domain.Reservation
	if err := cursor.All(ctx, &reservations); err != nil {
		return nil, fmt.Errorf("failed to decode reservations: %w", err)
	}

	return reservations, nil
}

// Update updates an existing reservation
func (r *MongoReservationRepository) Update(ctx context.Context, id primitive.ObjectID, reservation *domain.Reservation) error {
	reservation.UpdatedAt = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": reservation}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("reservation not found")
	}

	return nil
}

// Delete removes a reservation
func (r *MongoReservationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("reservation not found")
	}

	return nil
}

// GetReservedTableNumbers returns table numbers that are reserved for a given date and meal type
func (r *MongoReservationRepository) GetReservedTableNumbers(ctx context.Context, date string, mealType string) ([]int, error) {
	// Parse date to get start and end of day
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Query for reservations on that date with that meal type
	// Only count non-cancelled reservations
	filter := bson.M{
		"meal_type": mealType,
		"date_time": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
		"status": bson.M{
			"$ne": domain.StatusCancelled, // Exclude cancelled reservations
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query reservations: %w", err)
	}
	defer cursor.Close(ctx)

	var reservations []domain.Reservation
	if err := cursor.All(ctx, &reservations); err != nil {
		return nil, fmt.Errorf("failed to decode reservations: %w", err)
	}

	// Extract table numbers
	tableNumbers := make([]int, 0, len(reservations))
	for _, res := range reservations {
		tableNumbers = append(tableNumbers, res.TableNumber)
	}

	return tableNumbers, nil
}
