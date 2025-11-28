package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock implementations

type mockReservationRepository struct {
	createFunc                  func(ctx context.Context, reservation *domain.Reservation) error
	getByIDFunc                 func(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error)
	getAllFunc                  func(ctx context.Context, limit, offset int) ([]domain.Reservation, error)
	getByUserIDFunc             func(ctx context.Context, userID string) ([]domain.Reservation, error)
	updateFunc                  func(ctx context.Context, id primitive.ObjectID, reservation *domain.Reservation) error
	deleteFunc                  func(ctx context.Context, id primitive.ObjectID) error
	getReservedTableNumbersFunc func(ctx context.Context, date, mealType string) ([]int, error)
}

func (m *mockReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, reservation)
	}
	reservation.ID = primitive.NewObjectID()
	return nil
}

func (m *mockReservationRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockReservationRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Reservation, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, limit, offset)
	}
	return []domain.Reservation{}, nil
}

func (m *mockReservationRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Reservation, error) {
	if m.getByUserIDFunc != nil {
		return m.getByUserIDFunc(ctx, userID)
	}
	return []domain.Reservation{}, nil
}

func (m *mockReservationRepository) Update(ctx context.Context, id primitive.ObjectID, reservation *domain.Reservation) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, reservation)
	}
	return nil
}

func (m *mockReservationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockReservationRepository) GetReservedTableNumbers(ctx context.Context, date, mealType string) ([]int, error) {
	if m.getReservedTableNumbersFunc != nil {
		return m.getReservedTableNumbersFunc(ctx, date, mealType)
	}
	return []int{}, nil
}

type mockTableRepository struct {
	createFunc                      func(ctx context.Context, table *domain.Table) error
	getByIDFunc                     func(ctx context.Context, id primitive.ObjectID) (*domain.Table, error)
	getAllFunc                      func(ctx context.Context) ([]domain.Table, error)
	getByMealTypeFunc               func(ctx context.Context, mealType string) ([]domain.Table, error)
	getByTableNumberAndMealTypeFunc func(ctx context.Context, tableNumber int, mealType string) (*domain.Table, error)
	updateFunc                      func(ctx context.Context, id primitive.ObjectID, table *domain.Table) error
	deleteFunc                      func(ctx context.Context, id primitive.ObjectID) error
}

func (m *mockTableRepository) Create(ctx context.Context, table *domain.Table) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, table)
	}
	table.ID = primitive.NewObjectID()
	return nil
}

func (m *mockTableRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Table, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockTableRepository) GetAll(ctx context.Context) ([]domain.Table, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return []domain.Table{}, nil
}

func (m *mockTableRepository) GetByMealType(ctx context.Context, mealType string) ([]domain.Table, error) {
	if m.getByMealTypeFunc != nil {
		return m.getByMealTypeFunc(ctx, mealType)
	}
	return []domain.Table{}, nil
}

func (m *mockTableRepository) GetByTableNumberAndMealType(ctx context.Context, tableNumber int, mealType string) (*domain.Table, error) {
	if m.getByTableNumberAndMealTypeFunc != nil {
		return m.getByTableNumberAndMealTypeFunc(ctx, tableNumber, mealType)
	}
	return nil, errors.New("not found")
}

func (m *mockTableRepository) Update(ctx context.Context, id primitive.ObjectID, table *domain.Table) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, table)
	}
	return nil
}

func (m *mockTableRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

// Helper functions to create test dependencies

func createMockUserClient(statusCode int) *UserClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			w.Write([]byte(`{"id":1,"username":"testuser","email":"test@example.com"}`))
		}
	}))

	return NewUserClient(server.URL)
}

// Tests

func TestReservationService_CreateReservation_Success(t *testing.T) {
	mockRepo := &mockReservationRepository{
		getReservedTableNumbersFunc: func(ctx context.Context, date, mealType string) ([]int, error) {
			return []int{2, 3}, nil // Table 1 is available
		},
		createFunc: func(ctx context.Context, reservation *domain.Reservation) error {
			reservation.ID = primitive.NewObjectID()
			return nil
		},
	}

	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil) // nil is handled gracefully

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	futureDate := time.Now().Add(24 * time.Hour)
	req := domain.CreateReservationRequest{
		OwnerID:     "1",
		TableNumber: 1,
		Guests:      4,
		DateTime:    futureDate,
		MealType:    domain.MealTypeLunch,
	}

	reservation, err := service.CreateReservation(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if reservation.ID.IsZero() {
		t.Error("expected reservation ID to be set")
	}

	if reservation.OwnerID != "1" {
		t.Errorf("expected owner_id '1', got %s", reservation.OwnerID)
	}

	if reservation.TableNumber != 1 {
		t.Errorf("expected table_number 1, got %d", reservation.TableNumber)
	}

	if reservation.Status != domain.StatusPending {
		t.Errorf("expected status 'pending', got %s", reservation.Status)
	}

	if reservation.TotalPrice <= 0 {
		t.Errorf("expected total_price > 0, got %f", reservation.TotalPrice)
	}
}

func TestReservationService_CreateReservation_UserValidationFails(t *testing.T) {
	mockRepo := &mockReservationRepository{}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusNotFound)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	futureDate := time.Now().Add(24 * time.Hour)
	req := domain.CreateReservationRequest{
		OwnerID:     "999",
		TableNumber: 1,
		Guests:      4,
		DateTime:    futureDate,
		MealType:    domain.MealTypeLunch,
	}

	_, err := service.CreateReservation(context.Background(), req)

	if err == nil {
		t.Fatal("expected error for invalid user, got nil")
	}

	if !contains(err.Error(), "user validation failed") {
		t.Errorf("expected 'user validation failed' error, got %v", err)
	}
}

func TestReservationService_CreateReservation_TableAlreadyReserved(t *testing.T) {
	mockRepo := &mockReservationRepository{
		getReservedTableNumbersFunc: func(ctx context.Context, date, mealType string) ([]int, error) {
			return []int{1, 2, 3}, nil // Table 1 is reserved
		},
	}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	futureDate := time.Now().Add(24 * time.Hour)
	req := domain.CreateReservationRequest{
		OwnerID:     "1",
		TableNumber: 1, // This table is already reserved
		Guests:      4,
		DateTime:    futureDate,
		MealType:    domain.MealTypeLunch,
	}

	_, err := service.CreateReservation(context.Background(), req)

	if err == nil {
		t.Fatal("expected error for reserved table, got nil")
	}

	if !contains(err.Error(), "already reserved") {
		t.Errorf("expected 'already reserved' error, got %v", err)
	}
}

func TestReservationService_GetReservation_Success(t *testing.T) {
	reservationID := primitive.NewObjectID()
	expectedReservation := &domain.Reservation{
		ID:          reservationID,
		OwnerID:     "1",
		TableNumber: 1,
		Guests:      4,
		Status:      domain.StatusPending,
	}

	mockRepo := &mockReservationRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error) {
			if id == reservationID {
				return expectedReservation, nil
			}
			return nil, errors.New("not found")
		},
	}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	reservation, err := service.GetReservation(context.Background(), reservationID.Hex())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if reservation.ID != reservationID {
		t.Errorf("expected reservation ID %s, got %s", reservationID.Hex(), reservation.ID.Hex())
	}

	if reservation.OwnerID != "1" {
		t.Errorf("expected owner_id '1', got %s", reservation.OwnerID)
	}
}

func TestReservationService_GetReservation_InvalidID(t *testing.T) {
	mockRepo := &mockReservationRepository{}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	_, err := service.GetReservation(context.Background(), "invalid-id")

	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}

	if !contains(err.Error(), "invalid reservation ID") {
		t.Errorf("expected 'invalid reservation ID' error, got %v", err)
	}
}

func TestReservationService_ConfirmReservation_Success(t *testing.T) {
	reservationID := primitive.NewObjectID()
	futureDate := time.Now().Add(24 * time.Hour)
	existingReservation := &domain.Reservation{
		ID:          reservationID,
		OwnerID:     "1",
		TableNumber: 1,
		Guests:      4,
		DateTime:    futureDate,
		MealType:    domain.MealTypeLunch,
		Status:      domain.StatusPending,
		TotalPrice:  100.0,
	}

	mockRepo := &mockReservationRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error) {
			if id == reservationID {
				return existingReservation, nil
			}
			return nil, errors.New("not found")
		},
		updateFunc: func(ctx context.Context, id primitive.ObjectID, reservation *domain.Reservation) error {
			return nil
		},
	}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	req := domain.ConfirmReservationRequest{
		ConfirmationNotes: "Confirmed by admin",
	}

	reservation, err := service.ConfirmReservation(context.Background(), reservationID.Hex(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if reservation.Status != domain.StatusConfirmed {
		t.Errorf("expected status 'confirmed', got %s", reservation.Status)
	}
}

func TestReservationService_ConfirmReservation_AlreadyConfirmed(t *testing.T) {
	reservationID := primitive.NewObjectID()
	futureDate := time.Now().Add(24 * time.Hour)
	existingReservation := &domain.Reservation{
		ID:          reservationID,
		OwnerID:     "1",
		TableNumber: 1,
		Guests:      4,
		DateTime:    futureDate,
		MealType:    domain.MealTypeLunch,
		Status:      domain.StatusConfirmed, // Already confirmed
		TotalPrice:  100.0,
	}

	mockRepo := &mockReservationRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error) {
			return existingReservation, nil
		},
	}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	req := domain.ConfirmReservationRequest{}

	_, err := service.ConfirmReservation(context.Background(), reservationID.Hex(), req)

	if err == nil {
		t.Fatal("expected error for already confirmed reservation, got nil")
	}

	if !contains(err.Error(), "already confirmed") {
		t.Errorf("expected 'already confirmed' error, got %v", err)
	}
}

func TestReservationService_CancelReservation_Success(t *testing.T) {
	reservationID := primitive.NewObjectID()
	futureDate := time.Now().Add(24 * time.Hour)
	existingReservation := &domain.Reservation{
		ID:          reservationID,
		OwnerID:     "1",
		TableNumber: 1,
		Guests:      4,
		DateTime:    futureDate,
		MealType:    domain.MealTypeLunch,
		Status:      domain.StatusPending,
		TotalPrice:  100.0,
	}

	mockRepo := &mockReservationRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error) {
			return existingReservation, nil
		},
		updateFunc: func(ctx context.Context, id primitive.ObjectID, reservation *domain.Reservation) error {
			return nil
		},
	}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	reservation, err := service.CancelReservation(context.Background(), reservationID.Hex())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if reservation.Status != domain.StatusCancelled {
		t.Errorf("expected status 'cancelled', got %s", reservation.Status)
	}
}

func TestReservationService_CancelReservation_NotPending(t *testing.T) {
	reservationID := primitive.NewObjectID()
	futureDate := time.Now().Add(24 * time.Hour)
	existingReservation := &domain.Reservation{
		ID:          reservationID,
		OwnerID:     "1",
		TableNumber: 1,
		Guests:      4,
		DateTime:    futureDate,
		MealType:    domain.MealTypeLunch,
		Status:      domain.StatusConfirmed, // Not pending
		TotalPrice:  100.0,
	}

	mockRepo := &mockReservationRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error) {
			return existingReservation, nil
		},
	}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	_, err := service.CancelReservation(context.Background(), reservationID.Hex())

	if err == nil {
		t.Fatal("expected error for non-pending reservation, got nil")
	}

	if !contains(err.Error(), "only pending reservations can be cancelled") {
		t.Errorf("expected 'only pending reservations can be cancelled' error, got %v", err)
	}
}

func TestReservationService_DeleteReservation_Success(t *testing.T) {
	reservationID := primitive.NewObjectID()

	mockRepo := &mockReservationRepository{
		deleteFunc: func(ctx context.Context, id primitive.ObjectID) error {
			if id == reservationID {
				return nil
			}
			return errors.New("not found")
		},
	}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	err := service.DeleteReservation(context.Background(), reservationID.Hex())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestReservationService_GetAvailableTables_Success(t *testing.T) {
	mockRepo := &mockReservationRepository{
		getReservedTableNumbersFunc: func(ctx context.Context, date, mealType string) ([]int, error) {
			return []int{2, 3}, nil // Tables 2 and 3 are reserved
		},
	}

	mockTableRepo := &mockTableRepository{
		getByMealTypeFunc: func(ctx context.Context, mealType string) ([]domain.Table, error) {
			return []domain.Table{
				{TableNumber: 1, Capacity: 4, MealType: domain.MealTypeLunch},
				{TableNumber: 2, Capacity: 6, MealType: domain.MealTypeLunch},
				{TableNumber: 3, Capacity: 2, MealType: domain.MealTypeLunch},
				{TableNumber: 4, Capacity: 8, MealType: domain.MealTypeLunch},
			}, nil
		},
	}

	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	date := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	tables, err := service.GetAvailableTables(context.Background(), date, domain.MealTypeLunch)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should have tables 1 and 4 available (2 and 3 are reserved)
	if len(tables) != 2 {
		t.Errorf("expected 2 available tables, got %d", len(tables))
	}

	// Verify table numbers
	tableNumbers := make(map[int]bool)
	for _, table := range tables {
		tableNumbers[table.TableNumber] = true
	}

	if !tableNumbers[1] {
		t.Error("expected table 1 to be available")
	}

	if !tableNumbers[4] {
		t.Error("expected table 4 to be available")
	}

	if tableNumbers[2] || tableNumbers[3] {
		t.Error("tables 2 and 3 should not be available")
	}
}

func TestReservationService_GetAllReservations_Success(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)
	expectedReservations := []domain.Reservation{
		{
			ID:          primitive.NewObjectID(),
			OwnerID:     "1",
			TableNumber: 1,
			Guests:      4,
			DateTime:    futureDate,
			MealType:    domain.MealTypeLunch,
			Status:      domain.StatusPending,
		},
		{
			ID:          primitive.NewObjectID(),
			OwnerID:     "2",
			TableNumber: 2,
			Guests:      2,
			DateTime:    futureDate,
			MealType:    domain.MealTypeDinner,
			Status:      domain.StatusConfirmed,
		},
	}

	mockRepo := &mockReservationRepository{
		getAllFunc: func(ctx context.Context, limit, offset int) ([]domain.Reservation, error) {
			return expectedReservations, nil
		},
	}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	reservations, err := service.GetAllReservations(context.Background(), 10, 0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(reservations) != 2 {
		t.Errorf("expected 2 reservations, got %d", len(reservations))
	}
}

func TestReservationService_GetUserReservations_Success(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)
	expectedReservations := []domain.Reservation{
		{
			ID:          primitive.NewObjectID(),
			OwnerID:     "1",
			TableNumber: 1,
			Guests:      4,
			DateTime:    futureDate,
			MealType:    domain.MealTypeLunch,
			Status:      domain.StatusPending,
		},
	}

	mockRepo := &mockReservationRepository{
		getByUserIDFunc: func(ctx context.Context, userID string) ([]domain.Reservation, error) {
			if userID == "1" {
				return expectedReservations, nil
			}
			return []domain.Reservation{}, nil
		},
	}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	reservations, err := service.GetUserReservations(context.Background(), "1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(reservations) != 1 {
		t.Errorf("expected 1 reservation, got %d", len(reservations))
	}

	if reservations[0].OwnerID != "1" {
		t.Errorf("expected owner_id '1', got %s", reservations[0].OwnerID)
	}
}

func TestReservationService_GetUserReservations_UserNotFound(t *testing.T) {
	mockRepo := &mockReservationRepository{}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusNotFound)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	_, err := service.GetUserReservations(context.Background(), "999")

	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}

	if !contains(err.Error(), "user validation failed") {
		t.Errorf("expected 'user validation failed' error, got %v", err)
	}
}

func TestReservationService_UpdateReservation_Success(t *testing.T) {
	reservationID := primitive.NewObjectID()
	futureDate := time.Now().Add(24 * time.Hour)
	existingReservation := &domain.Reservation{
		ID:          reservationID,
		OwnerID:     "1",
		TableNumber: 1,
		Guests:      4,
		DateTime:    futureDate,
		MealType:    domain.MealTypeLunch,
		Status:      domain.StatusPending,
		TotalPrice:  100.0,
	}

	newGuests := 6
	updateReq := domain.UpdateReservationRequest{
		Guests: &newGuests,
	}

	mockRepo := &mockReservationRepository{
		getByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*domain.Reservation, error) {
			if id == reservationID {
				return existingReservation, nil
			}
			return nil, errors.New("not found")
		},
		updateFunc: func(ctx context.Context, id primitive.ObjectID, reservation *domain.Reservation) error {
			return nil
		},
	}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	reservation, err := service.UpdateReservation(context.Background(), reservationID.Hex(), updateReq)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if reservation.Guests != 6 {
		t.Errorf("expected guests to be updated to 6, got %d", reservation.Guests)
	}
}

func TestReservationService_UpdateReservation_InvalidID(t *testing.T) {
	mockRepo := &mockReservationRepository{}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	newGuests := 6
	updateReq := domain.UpdateReservationRequest{
		Guests: &newGuests,
	}

	_, err := service.UpdateReservation(context.Background(), "invalid-id", updateReq)

	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}

	if !contains(err.Error(), "invalid reservation ID") {
		t.Errorf("expected 'invalid reservation ID' error, got %v", err)
	}
}

func TestReservationService_DeleteReservation_InvalidID(t *testing.T) {
	mockRepo := &mockReservationRepository{}
	mockTableRepo := &mockTableRepository{}
	mockUserClient := createMockUserClient(http.StatusOK)
	mockPublisher := (*RabbitMQPublisher)(nil)

	service := NewReservationService(mockRepo, mockUserClient, mockPublisher, mockTableRepo)

	err := service.DeleteReservation(context.Background(), "invalid-id")

	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}

	if !contains(err.Error(), "invalid reservation ID") {
		t.Errorf("expected 'invalid reservation ID' error, got %v", err)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
