package service

import (
	"context"
	"errors"
	"testing"

	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tests for TableService

func TestTableService_CreateTable_Success(t *testing.T) {
	mockTableRepo := &mockTableRepository{
		getByTableNumberAndMealTypeFunc: func(ctx context.Context, tableNumber int, mealType string) (*domain.Table, error) {
			return nil, errors.New("not found")
		},
		createFunc: func(ctx context.Context, table *domain.Table) error {
			table.ID = primitive.NewObjectID()
			return nil
		},
	}

	service := NewTableService(mockTableRepo, nil)

	input := CreateTableInput{
		TableNumber: 1,
		Capacity:    4,
		MealType:    "lunch",
	}

	table, err := service.CreateTable(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if table.TableNumber != 1 {
		t.Errorf("expected table number 1, got %d", table.TableNumber)
	}

	if table.Capacity != 4 {
		t.Errorf("expected capacity 4, got %d", table.Capacity)
	}

	if table.MealType != "lunch" {
		t.Errorf("expected meal type 'lunch', got %s", table.MealType)
	}
}

func TestTableService_CreateTable_AlreadyExists(t *testing.T) {
	existingTable := &domain.Table{
		ID:          primitive.NewObjectID(),
		TableNumber: 1,
		Capacity:    4,
		MealType:    "lunch",
	}

	mockTableRepo := &mockTableRepository{
		getByTableNumberAndMealTypeFunc: func(ctx context.Context, tableNumber int, mealType string) (*domain.Table, error) {
			if tableNumber == 1 && mealType == "lunch" {
				return existingTable, nil
			}
			return nil, errors.New("not found")
		},
	}

	service := NewTableService(mockTableRepo, nil)

	input := CreateTableInput{
		TableNumber: 1,
		Capacity:    4,
		MealType:    "lunch",
	}

	_, err := service.CreateTable(context.Background(), input)

	if err == nil {
		t.Fatal("expected error for duplicate table, got nil")
	}

	if !contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got %v", err)
	}
}

func TestTableService_GetTableByID_Success(t *testing.T) {
	tableID := primitive.NewObjectID()
	expectedTable := &domain.Table{
		ID:          tableID,
		TableNumber: 1,
		Capacity:    4,
		MealType:    "lunch",
	}

	mockTableRepo := &mockTableRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Table, error) {
			if id == tableID {
				return expectedTable, nil
			}
			return nil, errors.New("not found")
		},
	}

	service := NewTableService(mockTableRepo, nil)

	table, err := service.GetTableByID(context.Background(), tableID.Hex())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if table.ID != tableID {
		t.Errorf("expected table ID %s, got %s", tableID.Hex(), table.ID.Hex())
	}
}

func TestTableService_GetTableByID_InvalidID(t *testing.T) {
	mockTableRepo := &mockTableRepository{}
	service := NewTableService(mockTableRepo, nil)

	_, err := service.GetTableByID(context.Background(), "invalid-id")

	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}

	if !contains(err.Error(), "invalid table ID") {
		t.Errorf("expected 'invalid table ID' error, got %v", err)
	}
}

func TestTableService_GetAllTables_Success(t *testing.T) {
	expectedTables := []domain.Table{
		{ID: primitive.NewObjectID(), TableNumber: 1, Capacity: 4, MealType: "lunch"},
		{ID: primitive.NewObjectID(), TableNumber: 2, Capacity: 6, MealType: "dinner"},
	}

	mockTableRepo := &mockTableRepository{
		getAllFunc: func(ctx context.Context) ([]domain.Table, error) {
			return expectedTables, nil
		},
	}

	service := NewTableService(mockTableRepo, nil)

	tables, err := service.GetAllTables(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tables) != 2 {
		t.Errorf("expected 2 tables, got %d", len(tables))
	}
}

func TestTableService_GetAllTables_EmptyResult(t *testing.T) {
	mockTableRepo := &mockTableRepository{
		getAllFunc: func(ctx context.Context) ([]domain.Table, error) {
			return nil, nil
		},
	}

	service := NewTableService(mockTableRepo, nil)

	tables, err := service.GetAllTables(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tables == nil {
		t.Error("expected empty array, got nil")
	}

	if len(tables) != 0 {
		t.Errorf("expected 0 tables, got %d", len(tables))
	}
}

func TestTableService_GetTablesByMealType_Success(t *testing.T) {
	expectedTables := []domain.Table{
		{ID: primitive.NewObjectID(), TableNumber: 1, Capacity: 4, MealType: "lunch"},
		{ID: primitive.NewObjectID(), TableNumber: 2, Capacity: 6, MealType: "lunch"},
	}

	mockTableRepo := &mockTableRepository{
		getByMealTypeFunc: func(ctx context.Context, mealType string) ([]domain.Table, error) {
			if mealType == "lunch" {
				return expectedTables, nil
			}
			return []domain.Table{}, nil
		},
	}

	service := NewTableService(mockTableRepo, nil)

	tables, err := service.GetTablesByMealType(context.Background(), "lunch")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tables) != 2 {
		t.Errorf("expected 2 tables, got %d", len(tables))
	}
}

func TestTableService_UpdateTable_Success(t *testing.T) {
	tableID := primitive.NewObjectID()
	existingTable := &domain.Table{
		ID:          tableID,
		TableNumber: 1,
		Capacity:    4,
		MealType:    "lunch",
	}

	mockTableRepo := &mockTableRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Table, error) {
			if id == tableID {
				return existingTable, nil
			}
			return nil, errors.New("not found")
		},
		getByTableNumberAndMealTypeFunc: func(ctx context.Context, tableNumber int, mealType string) (*domain.Table, error) {
			return nil, errors.New("not found")
		},
		updateFunc: func(ctx context.Context, id primitive.ObjectID, table *domain.Table) error {
			return nil
		},
	}

	service := NewTableService(mockTableRepo, nil)

	input := UpdateTableInput{
		TableNumber: 1,
		Capacity:    6,
		MealType:    "lunch",
	}

	table, err := service.UpdateTable(context.Background(), tableID.Hex(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if table.Capacity != 6 {
		t.Errorf("expected capacity 6, got %d", table.Capacity)
	}
}

func TestTableService_UpdateTable_InvalidID(t *testing.T) {
	mockTableRepo := &mockTableRepository{}
	service := NewTableService(mockTableRepo, nil)

	input := UpdateTableInput{
		TableNumber: 1,
		Capacity:    6,
		MealType:    "lunch",
	}

	_, err := service.UpdateTable(context.Background(), "invalid-id", input)

	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}

	if !contains(err.Error(), "invalid table ID") {
		t.Errorf("expected 'invalid table ID' error, got %v", err)
	}
}

func TestTableService_UpdateTable_NotFound(t *testing.T) {
	mockTableRepo := &mockTableRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Table, error) {
			return nil, errors.New("not found")
		},
	}

	service := NewTableService(mockTableRepo, nil)

	input := UpdateTableInput{
		TableNumber: 1,
		Capacity:    6,
		MealType:    "lunch",
	}

	_, err := service.UpdateTable(context.Background(), primitive.NewObjectID().Hex(), input)

	if err == nil {
		t.Fatal("expected error for not found table, got nil")
	}
}

func TestTableService_UpdateTable_DuplicateTableNumber(t *testing.T) {
	tableID1 := primitive.NewObjectID()
	tableID2 := primitive.NewObjectID()

	existingTable1 := &domain.Table{
		ID:          tableID1,
		TableNumber: 1,
		Capacity:    4,
		MealType:    "lunch",
	}

	existingTable2 := &domain.Table{
		ID:          tableID2,
		TableNumber: 2,
		Capacity:    6,
		MealType:    "lunch",
	}

	mockTableRepo := &mockTableRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Table, error) {
			if id == tableID1 {
				return existingTable1, nil
			}
			return nil, errors.New("not found")
		},
		getByTableNumberAndMealTypeFunc: func(ctx context.Context, tableNumber int, mealType string) (*domain.Table, error) {
			if tableNumber == 2 && mealType == "lunch" {
				return existingTable2, nil
			}
			return nil, errors.New("not found")
		},
	}

	service := NewTableService(mockTableRepo, nil)

	input := UpdateTableInput{
		TableNumber: 2, // Trying to change to table 2, which already exists
		Capacity:    6,
		MealType:    "lunch",
	}

	_, err := service.UpdateTable(context.Background(), tableID1.Hex(), input)

	if err == nil {
		t.Fatal("expected error for duplicate table number, got nil")
	}

	if !contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got %v", err)
	}
}

func TestTableService_DeleteTable_Success(t *testing.T) {
	tableID := primitive.NewObjectID()
	existingTable := &domain.Table{
		ID:          tableID,
		TableNumber: 1,
		Capacity:    4,
		MealType:    "lunch",
	}

	mockTableRepo := &mockTableRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Table, error) {
			if id == tableID {
				return existingTable, nil
			}
			return nil, errors.New("not found")
		},
		deleteFunc: func(ctx context.Context, id primitive.ObjectID) error {
			return nil
		},
	}

	service := NewTableService(mockTableRepo, nil)

	err := service.DeleteTable(context.Background(), tableID.Hex())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestTableService_DeleteTable_InvalidID(t *testing.T) {
	mockTableRepo := &mockTableRepository{}
	service := NewTableService(mockTableRepo, nil)

	err := service.DeleteTable(context.Background(), "invalid-id")

	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}

	if !contains(err.Error(), "invalid table ID") {
		t.Errorf("expected 'invalid table ID' error, got %v", err)
	}
}

func TestTableService_DeleteTable_NotFound(t *testing.T) {
	mockTableRepo := &mockTableRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Table, error) {
			return nil, errors.New("not found")
		},
	}

	service := NewTableService(mockTableRepo, nil)

	err := service.DeleteTable(context.Background(), primitive.NewObjectID().Hex())

	if err == nil {
		t.Fatal("expected error for not found table, got nil")
	}
}

func TestTableService_GetTableCapacity_Success(t *testing.T) {
	expectedTable := &domain.Table{
		ID:          primitive.NewObjectID(),
		TableNumber: 1,
		Capacity:    4,
		MealType:    "lunch",
	}

	mockTableRepo := &mockTableRepository{
		getByTableNumberAndMealTypeFunc: func(ctx context.Context, tableNumber int, mealType string) (*domain.Table, error) {
			if tableNumber == 1 && mealType == "lunch" {
				return expectedTable, nil
			}
			return nil, errors.New("not found")
		},
	}

	service := NewTableService(mockTableRepo, nil)

	capacity, err := service.GetTableCapacity(context.Background(), 1, "lunch")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if capacity != 4 {
		t.Errorf("expected capacity 4, got %d", capacity)
	}
}

func TestTableService_GetTableCapacity_NotFound(t *testing.T) {
	mockTableRepo := &mockTableRepository{
		getByTableNumberAndMealTypeFunc: func(ctx context.Context, tableNumber int, mealType string) (*domain.Table, error) {
			return nil, errors.New("not found")
		},
	}

	service := NewTableService(mockTableRepo, nil)

	_, err := service.GetTableCapacity(context.Background(), 999, "lunch")

	if err == nil {
		t.Fatal("expected error for not found table, got nil")
	}
}
