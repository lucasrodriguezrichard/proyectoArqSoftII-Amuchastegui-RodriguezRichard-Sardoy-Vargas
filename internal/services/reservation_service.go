package services

import (
	"fmt"
	"restaurant-system/internal/domain"
	"restaurant-system/internal/repository"
	"time"

	"github.com/google/uuid"
)

// ReservationService handles reservation business logic
type ReservationService struct {
	reservationRepo repository.ReservationRepository
}

// NewReservationService creates a new reservation service
func NewReservationService(reservationRepo repository.ReservationRepository) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
	}
}

// CreateReservation creates a new reservation
func (s *ReservationService) CreateReservation(req CreateReservationRequest) (*domain.Reservation, error) {
	// Validate reservation date is in the future
	if req.ReservationDate.Before(time.Now()) {
		return nil, fmt.Errorf("reservation date must be in the future")
	}

	// Validate party size
	if req.PartySize <= 0 {
		return nil, fmt.Errorf("party size must be greater than 0")
	}

	reservation := &domain.Reservation{
		ID:              uuid.New().String(),
		RestaurantID:    req.RestaurantID,
		TableID:         req.TableID,
		CustomerName:    req.CustomerName,
		CustomerPhone:   req.CustomerPhone,
		CustomerEmail:   req.CustomerEmail,
		ReservationDate: req.ReservationDate,
		MealType:        req.MealType,
		PartySize:       req.PartySize,
		Status:          domain.ReservationStatusPending,
		SpecialRequests: req.SpecialRequests,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err := s.reservationRepo.Create(reservation)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	return reservation, nil
}

// GetReservation retrieves a reservation by ID
func (s *ReservationService) GetReservation(id string) (*domain.Reservation, error) {
	return s.reservationRepo.GetByID(id)
}

// GetReservationsByDateRange retrieves reservations within a date range
func (s *ReservationService) GetReservationsByDateRange(startDate, endDate time.Time) ([]*domain.Reservation, error) {
	return s.reservationRepo.GetByDateRange(startDate, endDate)
}

// ConfirmReservation confirms a reservation
func (s *ReservationService) ConfirmReservation(id string) error {
	return s.reservationRepo.UpdateStatus(id, domain.ReservationStatusConfirmed)
}

// CancelReservation cancels a reservation
func (s *ReservationService) CancelReservation(id string) error {
	return s.reservationRepo.UpdateStatus(id, domain.ReservationStatusCancelled)
}

// CompleteReservation marks a reservation as completed
func (s *ReservationService) CompleteReservation(id string) error {
	return s.reservationRepo.UpdateStatus(id, domain.ReservationStatusCompleted)
}

// CreateReservationRequest represents the request to create a reservation
type CreateReservationRequest struct {
	RestaurantID    string          `json:"restaurant_id" binding:"required"`
	TableID         string          `json:"table_id" binding:"required"`
	CustomerName    string          `json:"customer_name" binding:"required"`
	CustomerPhone   string          `json:"customer_phone" binding:"required"`
	CustomerEmail   string          `json:"customer_email" binding:"required,email"`
	ReservationDate time.Time       `json:"reservation_date" binding:"required"`
	MealType        domain.MealType `json:"meal_type" binding:"required"`
	PartySize       int             `json:"party_size" binding:"required,min=1"`
	SpecialRequests string          `json:"special_requests"`
}
