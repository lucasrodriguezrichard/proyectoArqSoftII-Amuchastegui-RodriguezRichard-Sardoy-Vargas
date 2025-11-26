package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TableRepository defines the interface for table persistence
type TableRepository interface {
	Create(ctx context.Context, table *domain.Table) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Table, error)
	GetAll(ctx context.Context) ([]domain.Table, error)
	GetByMealType(ctx context.Context, mealType string) ([]domain.Table, error)
	GetByTableNumberAndMealType(ctx context.Context, tableNumber int, mealType string) (*domain.Table, error)
	Update(ctx context.Context, id primitive.ObjectID, table *domain.Table) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// MongoTableRepository implements TableRepository using MongoDB
type MongoTableRepository struct {
	collection *mongo.Collection
}

// NewMongoTableRepository creates a new MongoDB table repository
func NewMongoTableRepository(collection *mongo.Collection) *MongoTableRepository {
	return &MongoTableRepository{
		collection: collection,
	}
}

// Create inserts a new table
func (r *MongoTableRepository) Create(ctx context.Context, table *domain.Table) error {
	table.CreatedAt = time.Now()
	table.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, table)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Set the ID from the inserted document
	table.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a table by ID
func (r *MongoTableRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Table, error) {
	var table domain.Table
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&table)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("table not found")
		}
		return nil, fmt.Errorf("failed to get table: %w", err)
	}
	return &table, nil
}

// GetAll retrieves all tables
func (r *MongoTableRepository) GetAll(ctx context.Context) ([]domain.Table, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	defer cursor.Close(ctx)

	var tables []domain.Table
	if err := cursor.All(ctx, &tables); err != nil {
		return nil, fmt.Errorf("failed to decode tables: %w", err)
	}

	return tables, nil
}

// GetByMealType retrieves all tables for a specific meal type
func (r *MongoTableRepository) GetByMealType(ctx context.Context, mealType string) ([]domain.Table, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"meal_type": mealType})
	if err != nil {
		return nil, fmt.Errorf("failed to get tables by meal type: %w", err)
	}
	defer cursor.Close(ctx)

	var tables []domain.Table
	if err := cursor.All(ctx, &tables); err != nil {
		return nil, fmt.Errorf("failed to decode tables: %w", err)
	}

	return tables, nil
}

// GetByTableNumberAndMealType retrieves a table by table number and meal type
func (r *MongoTableRepository) GetByTableNumberAndMealType(ctx context.Context, tableNumber int, mealType string) (*domain.Table, error) {
	var table domain.Table
	err := r.collection.FindOne(ctx, bson.M{
		"table_number": tableNumber,
		"meal_type":    mealType,
	}).Decode(&table)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("table not found")
		}
		return nil, fmt.Errorf("failed to get table: %w", err)
	}
	return &table, nil
}

// Update updates an existing table
func (r *MongoTableRepository) Update(ctx context.Context, id primitive.ObjectID, table *domain.Table) error {
	table.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"table_number": table.TableNumber,
			"capacity":     table.Capacity,
			"meal_type":    table.MealType,
			"updated_at":   table.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update table: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("table not found")
	}

	return nil
}

// Delete removes a table
func (r *MongoTableRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("table not found")
	}

	return nil
}
