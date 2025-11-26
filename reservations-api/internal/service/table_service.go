package service

import (
	"context"
	"fmt"

	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/domain"
	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TableService handles business logic for tables
type TableService struct {
	tableRepo repository.TableRepository
	publisher *RabbitMQPublisher
}

// NewTableService creates a new table service
func NewTableService(tableRepo repository.TableRepository, publisher *RabbitMQPublisher) *TableService {
	return &TableService{
		tableRepo: tableRepo,
		publisher: publisher,
	}
}

// CreateTableInput represents input for creating a table
type CreateTableInput struct {
	TableNumber int    `json:"table_number" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required,min=1,max=50"`
	MealType    string `json:"meal_type" binding:"required,oneof=breakfast lunch dinner event"`
}

// UpdateTableInput represents input for updating a table
type UpdateTableInput struct {
	TableNumber int    `json:"table_number" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required,min=1,max=50"`
	MealType    string `json:"meal_type" binding:"required,oneof=breakfast lunch dinner event"`
}

// CreateTable creates a new table
func (s *TableService) CreateTable(ctx context.Context, input CreateTableInput) (*domain.Table, error) {
	// Check if a table with the same number and meal type already exists
	existing, err := s.tableRepo.GetByTableNumberAndMealType(ctx, input.TableNumber, input.MealType)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("table number %d already exists for meal type %s", input.TableNumber, input.MealType)
	}

	table := &domain.Table{
		TableNumber: input.TableNumber,
		Capacity:    input.Capacity,
		MealType:    input.MealType,
	}

	if err := s.tableRepo.Create(ctx, table); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Publish table creation event to RabbitMQ
	if s.publisher != nil {
		if err := s.publisher.PublishWithType("create", table.ID.Hex(), "table", tableMetadata(table)); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to publish table create event: %v\n", err)
		}
	}

	return table, nil
}

// GetTableByID retrieves a table by ID
func (s *TableService) GetTableByID(ctx context.Context, id string) (*domain.Table, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid table ID: %w", err)
	}

	table, err := s.tableRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetAllTables retrieves all tables
func (s *TableService) GetAllTables(ctx context.Context) ([]domain.Table, error) {
	tables, err := s.tableRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Return empty array instead of nil if no tables
	if tables == nil {
		return []domain.Table{}, nil
	}

	return tables, nil
}

// GetTablesByMealType retrieves tables by meal type
func (s *TableService) GetTablesByMealType(ctx context.Context, mealType string) ([]domain.Table, error) {
	tables, err := s.tableRepo.GetByMealType(ctx, mealType)
	if err != nil {
		return nil, err
	}

	// Return empty array instead of nil if no tables
	if tables == nil {
		return []domain.Table{}, nil
	}

	return tables, nil
}

// UpdateTable updates an existing table
func (s *TableService) UpdateTable(ctx context.Context, id string, input UpdateTableInput) (*domain.Table, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid table ID: %w", err)
	}

	// Check if table exists
	table, err := s.tableRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Check if updating would create a duplicate (different table with same number and meal type)
	existing, err := s.tableRepo.GetByTableNumberAndMealType(ctx, input.TableNumber, input.MealType)
	if err == nil && existing != nil && existing.ID != objectID {
		return nil, fmt.Errorf("table number %d already exists for meal type %s", input.TableNumber, input.MealType)
	}

	table.TableNumber = input.TableNumber
	table.Capacity = input.Capacity
	table.MealType = input.MealType

	if err := s.tableRepo.Update(ctx, objectID, table); err != nil {
		return nil, fmt.Errorf("failed to update table: %w", err)
	}

	// Publish table update event to RabbitMQ
	if s.publisher != nil {
		if err := s.publisher.PublishWithType("update", objectID.Hex(), "table", tableMetadata(table)); err != nil {
			fmt.Printf("Warning: failed to publish table update event: %v\n", err)
		}
	}

	return table, nil
}

// DeleteTable deletes a table
func (s *TableService) DeleteTable(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid table ID: %w", err)
	}

	// fetch table info for metadata before deletion
	table, err := s.tableRepo.GetByID(ctx, objectID)
	if err != nil {
		return err
	}

	if err := s.tableRepo.Delete(ctx, objectID); err != nil {
		return err
	}

	// Publish table deletion event to RabbitMQ
	if s.publisher != nil {
		if err := s.publisher.PublishWithType("delete", objectID.Hex(), "table", tableMetadata(table)); err != nil {
			fmt.Printf("Warning: failed to publish table delete event: %v\n", err)
		}
	}

	return nil
}

func tableMetadata(table *domain.Table) map[string]interface{} {
	if table == nil {
		return nil
	}
	return map[string]interface{}{
		"table_number": table.TableNumber,
		"meal_type":    table.MealType,
		"capacity":     table.Capacity,
	}
}

// GetTableCapacity returns the capacity for a specific table number and meal type
func (s *TableService) GetTableCapacity(ctx context.Context, tableNumber int, mealType string) (int, error) {
	table, err := s.tableRepo.GetByTableNumberAndMealType(ctx, tableNumber, mealType)
	if err != nil {
		return 0, err
	}

	return table.Capacity, nil
}
