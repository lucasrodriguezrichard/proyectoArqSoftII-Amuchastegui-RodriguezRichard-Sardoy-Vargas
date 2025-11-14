package service

import (
	"context"
	"fmt"
	"log"

	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/domain"
	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReservationService defines the business logic for reservations
type ReservationService interface {
	CreateReservation(ctx context.Context, req domain.CreateReservationRequest) (*domain.Reservation, error)
	GetReservation(ctx context.Context, id string) (*domain.Reservation, error)
	GetAllReservations(ctx context.Context, limit, offset int) ([]domain.Reservation, error)
	GetUserReservations(ctx context.Context, userID string) ([]domain.Reservation, error)
	UpdateReservation(ctx context.Context, id string, req domain.UpdateReservationRequest) (*domain.Reservation, error)
	DeleteReservation(ctx context.Context, id string) error
	ConfirmReservation(ctx context.Context, id string, req domain.ConfirmReservationRequest) (*domain.Reservation, error)
	GetAvailableTables(ctx context.Context, date string, mealType string) ([]domain.TableConfig, error)
}

// reservationService implements ReservationService
type reservationService struct {
	repo        repository.ReservationRepository
	userClient  *UserClient
	rmqPublisher *RabbitMQPublisher
}

// NewReservationService creates a new reservation service
func NewReservationService(
	repo repository.ReservationRepository,
	userClient *UserClient,
	rmqPublisher *RabbitMQPublisher,
) ReservationService {
	return &reservationService{
		repo:        repo,
		userClient:  userClient,
		rmqPublisher: rmqPublisher,
	}
}

// CreateReservation creates a new reservation with validation and concurrent calculations
func (s *reservationService) CreateReservation(ctx context.Context, req domain.CreateReservationRequest) (*domain.Reservation, error) {
	// 1. Validate user exists via Users API
	if err := s.userClient.ValidateUser(req.OwnerID); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	// 2. Create reservation object
	reservation := domain.NewReservation(req)

	// 3. VALIDATE TABLE AVAILABILITY - Check if specific table is already reserved for this date/meal_type
	date := reservation.DateTime.Format("2006-01-02")
	reservedTables, err := s.repo.GetReservedTableNumbers(ctx, date, reservation.MealType)
	if err != nil {
		return nil, fmt.Errorf("failed to check table availability: %w", err)
	}

	// Check if requested table is in the reserved list
	for _, reservedNum := range reservedTables {
		if reservedNum == reservation.TableNumber {
			return nil, fmt.Errorf("table %d is already reserved for %s on %s", reservation.TableNumber, reservation.MealType, date)
		}
	}

	// 4. Perform concurrent calculations (pricing, discounts)
	calcResult, err := domain.CalculateReservationConcurrent(
		reservation.TableNumber,
		reservation.Guests,
		reservation.DateTime,
		reservation.MealType,
		reservation.OwnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("calculation failed: %w", err)
	}

	// 5. Check if table is available (additional validations like weekend hours)
	if !calcResult.Available {
		return nil, fmt.Errorf("reservation not available: %v", calcResult.Restrictions)
	}

	// 6. Set calculated price
	reservation.TotalPrice = calcResult.FinalPrice

	// 7. Validate reservation data
	if err := reservation.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 8. Save to database
	if err := s.repo.Create(ctx, &reservation); err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	// 9. Publish event to RabbitMQ (async notification)
	go func() {
		if err := s.rmqPublisher.Publish("create", reservation.ID.Hex()); err != nil {
			log.Printf("Warning: failed to publish create event: %v", err)
		}
	}()

	return &reservation, nil
}

// GetReservation retrieves a reservation by ID
func (s *reservationService) GetReservation(ctx context.Context, id string) (*domain.Reservation, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid reservation ID: %w", err)
	}

	reservation, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

// GetAllReservations retrieves all reservations with pagination
func (s *reservationService) GetAllReservations(ctx context.Context, limit, offset int) ([]domain.Reservation, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

// GetUserReservations retrieves all reservations for a specific user
func (s *reservationService) GetUserReservations(ctx context.Context, userID string) ([]domain.Reservation, error) {
	// Validate user exists
	if err := s.userClient.ValidateUser(userID); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	return s.repo.GetByUserID(ctx, userID)
}

// UpdateReservation updates an existing reservation
func (s *reservationService) UpdateReservation(ctx context.Context, id string, req domain.UpdateReservationRequest) (*domain.Reservation, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid reservation ID: %w", err)
	}

	// Get existing reservation
	reservation, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.TableNumber != nil {
		reservation.TableNumber = *req.TableNumber
	}
	if req.Guests != nil {
		reservation.Guests = *req.Guests
	}
	if req.DateTime != nil {
		reservation.DateTime = *req.DateTime
	}
	if req.MealType != nil {
		reservation.MealType = *req.MealType
	}
	if req.SpecialRequests != nil {
		reservation.SpecialRequests = *req.SpecialRequests
	}
	if req.Status != nil {
		reservation.Status = *req.Status
	}

	// Recalculate price if relevant fields changed
	if req.Guests != nil || req.DateTime != nil || req.MealType != nil {
		calcResult, err := domain.CalculateReservationConcurrent(
			reservation.TableNumber,
			reservation.Guests,
			reservation.DateTime,
			reservation.MealType,
			reservation.OwnerID,
		)
		if err != nil {
			return nil, fmt.Errorf("calculation failed: %w", err)
		}

		if !calcResult.Available {
			return nil, fmt.Errorf("reservation not available: %v", calcResult.Restrictions)
		}

		reservation.TotalPrice = calcResult.FinalPrice
	}

	// Validate
	if err := reservation.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update in database
	if err := s.repo.Update(ctx, objectID, reservation); err != nil {
		return nil, err
	}

	// Publish event to RabbitMQ
	go func() {
		if err := s.rmqPublisher.Publish("update", reservation.ID.Hex()); err != nil {
			log.Printf("Warning: failed to publish update event: %v", err)
		}
	}()

	return reservation, nil
}

// DeleteReservation deletes a reservation
func (s *reservationService) DeleteReservation(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid reservation ID: %w", err)
	}

	// Delete from database
	if err := s.repo.Delete(ctx, objectID); err != nil {
		return err
	}

	// Publish event to RabbitMQ
	go func() {
		if err := s.rmqPublisher.Publish("delete", id); err != nil {
			log.Printf("Warning: failed to publish delete event: %v", err)
		}
	}()

	return nil
}

// ConfirmReservation confirms a reservation with concurrent recalculation
func (s *reservationService) ConfirmReservation(ctx context.Context, id string, req domain.ConfirmReservationRequest) (*domain.Reservation, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid reservation ID: %w", err)
	}

	// Get existing reservation
	reservation, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Check if already confirmed
	if reservation.Status == domain.StatusConfirmed {
		return nil, fmt.Errorf("reservation already confirmed")
	}

	// Perform concurrent calculations again (may apply confirmation discount)
	calcResult, err := domain.CalculateReservationConcurrent(
		reservation.TableNumber,
		reservation.Guests,
		reservation.DateTime,
		reservation.MealType,
		reservation.OwnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("calculation failed: %w", err)
	}

	if !calcResult.Available {
		return nil, fmt.Errorf("reservation not available: %v", calcResult.Restrictions)
	}

	// Update status and price
	reservation.Status = domain.StatusConfirmed
	reservation.TotalPrice = calcResult.FinalPrice

	// Update in database
	if err := s.repo.Update(ctx, objectID, reservation); err != nil {
		return nil, err
	}

	// Publish event to RabbitMQ
	go func() {
		if err := s.rmqPublisher.Publish("confirm", reservation.ID.Hex()); err != nil {
			log.Printf("Warning: failed to publish confirm event: %v", err)
		}
	}()

	return reservation, nil
}

// GetAvailableTables returns available tables for a given date and meal type
func (s *reservationService) GetAvailableTables(ctx context.Context, date string, mealType string) ([]domain.TableConfig, error) {
	// Get all predefined tables for the meal type
	allTables := domain.GetTablesForMealType(mealType)

	// Get all reservations for the given date and meal type
	reservedTables, err := s.repo.GetReservedTableNumbers(ctx, date, mealType)
	if err != nil {
		return nil, fmt.Errorf("failed to get reserved tables: %w", err)
	}

	// Filter out reserved tables
	availableTables := []domain.TableConfig{}
	for _, table := range allTables {
		isReserved := false
		for _, reservedNum := range reservedTables {
			if table.TableNumber == reservedNum {
				isReserved = true
				break
			}
		}
		if !isReserved {
			availableTables = append(availableTables, table)
		}
	}

	return availableTables, nil
}
