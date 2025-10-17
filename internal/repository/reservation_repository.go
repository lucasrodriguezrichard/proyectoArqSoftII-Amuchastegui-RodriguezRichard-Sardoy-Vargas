package repository

import (
	"restaurant-system/internal/dao"
	"restaurant-system/internal/domain"
	"time"
)

// ReservationRepository defines the interface for reservation operations
type ReservationRepository interface {
	Create(reservation *domain.Reservation) error
	GetByID(id string) (*domain.Reservation, error)
	GetByDateRange(startDate, endDate time.Time) ([]*domain.Reservation, error)
	UpdateStatus(id string, status domain.ReservationStatus) error
	Delete(id string) error
}

// reservationRepository implements ReservationRepository
type reservationRepository struct {
	dao *dao.ReservationDAO
}

// NewReservationRepository creates a new reservation repository
func NewReservationRepository(dao *dao.ReservationDAO) ReservationRepository {
	return &reservationRepository{dao: dao}
}

func (r *reservationRepository) Create(reservation *domain.Reservation) error {
	return r.dao.Create(reservation)
}

func (r *reservationRepository) GetByID(id string) (*domain.Reservation, error) {
	return r.dao.GetByID(id)
}

func (r *reservationRepository) GetByDateRange(startDate, endDate time.Time) ([]*domain.Reservation, error) {
	return r.dao.GetByDateRange(startDate, endDate)
}

func (r *reservationRepository) UpdateStatus(id string, status domain.ReservationStatus) error {
	return r.dao.UpdateStatus(id, status)
}

func (r *reservationRepository) Delete(id string) error {
	return r.dao.Delete(id)
}
